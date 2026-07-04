package tracker

import (
	"crypto/rand"
	"errors"
	"fmt"
	"gobt/pkg/bencode"
	"gobt/pkg/torrent"
	"gobt/pkg/utils"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort   = 6881
	defaultScheme = "http"
	peerIDPrefix  = "-GP420-"
)

// Peer describes a peer returned by a tracker.
type Peer struct {
	IP      string
	Port    int
	PeerID  string
	Compact bool
}

// AnnounceRequest holds the parameters sent to a tracker.
type AnnounceRequest struct {
	InfoHash   []byte
	PeerID     []byte
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string
	Compact    bool
	NoPeerID   bool
	NumWant    int
	Key        string
}

// AnnounceResponse represents the parsed tracker response.
type AnnounceResponse struct {
	Interval        int
	MinInterval     int
	TrackerID       string
	Complete        int
	Incomplete      int
	Downloaded      int
	WarningMessage  string
	FailureReason   string
	Peers           []Peer
	RawPeersCompact []byte
}

// TrackerError indicates a tracker failure response.
type TrackerError struct {
	Reason string
}

func (e *TrackerError) Error() string {
	return fmt.Sprintf("tracker error: %s", e.Reason)
}

// Client performs HTTP tracker announces.
type Client struct {
	HTTPClient *http.Client
	PeerID     []byte
	Port       int
	UserAgent  string
	Compact    bool
}

// NewClient creates a tracker client with sensible defaults.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		PeerID:     GeneratePeerID(),
		Port:       defaultPort,
		UserAgent:   fmt.Sprintf("gobt/%s", utils.Version()),
		Compact:    true,
	}
}

// Announce announces a torrent to the given tracker URL.
func (c *Client) Announce(trackerURL string, t *torrent.TorrentInfo, event string) (*AnnounceResponse, error) {
	if t == nil {
		return nil, errors.New("torrent cannot be nil")
	}

	infoHash, err := t.InfoHashBytes()
	if err != nil {
		return nil, err
	}

	req := AnnounceRequest{
		InfoHash:   infoHash,
		PeerID:     c.PeerID,
		Port:       c.Port,
		Uploaded:   0,
		Downloaded: 0,
		Left:       t.TotalSize(),
		Event:      event,
		Compact:    c.Compact,
		NumWant:    50,
	}

	urlStr, err := BuildAnnounceURL(trackerURL, req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	if c.UserAgent != "" {
		httpReq.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tracker returned status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return ParseAnnounceResponse(body)
}

// BuildAnnounceURL builds a tracker announce URL with encoded query parameters.
func BuildAnnounceURL(baseURL string, req AnnounceRequest) (string, error) {
	if len(req.InfoHash) != 20 {
		return "", fmt.Errorf("info_hash must be 20 bytes, got %d", len(req.InfoHash))
	}
	if len(req.PeerID) != 20 {
		return "", fmt.Errorf("peer_id must be 20 bytes, got %d", len(req.PeerID))
	}
	if req.Port <= 0 || req.Port > 65535 {
		return "", fmt.Errorf("invalid port %d", req.Port)
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" {
		parsed.Scheme = defaultScheme
	}

	query := parsed.Query()
	query.Set("info_hash", string(req.InfoHash))
	query.Set("peer_id", string(req.PeerID))
	query.Set("port", strconv.Itoa(req.Port))
	query.Set("uploaded", strconv.FormatInt(req.Uploaded, 10))
	query.Set("downloaded", strconv.FormatInt(req.Downloaded, 10))
	query.Set("left", strconv.FormatInt(req.Left, 10))
	query.Set("compact", boolToInt(req.Compact))
	if req.Event != "" {
		query.Set("event", req.Event)
	}
	if req.NoPeerID {
		query.Set("no_peer_id", "1")
	}
	if req.NumWant != 0 {
		query.Set("numwant", strconv.Itoa(req.NumWant))
	}
	if req.Key != "" {
		query.Set("key", req.Key)
	}

	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

// ParseAnnounceResponse parses a tracker response encoded with bencode.
func ParseAnnounceResponse(data []byte) (*AnnounceResponse, error) {
	decoder := bencode.NewDecoder(strings.NewReader(string(data)))
	decoded, err := decoder.Decode()
	if err != nil {
		return nil, err
	}

	dict, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected tracker response dict, got %T", decoded)
	}

	resp := &AnnounceResponse{}
	if reason, ok := dict["failure reason"].(string); ok && reason != "" {
		resp.FailureReason = reason
		return resp, &TrackerError{Reason: reason}
	}

	resp.Interval = toInt(dict["interval"])
	resp.MinInterval = toInt(dict["min interval"])
	resp.TrackerID = toString(dict["tracker id"])
	resp.Complete = toInt(dict["complete"])
	resp.Incomplete = toInt(dict["incomplete"])
	resp.Downloaded = toInt(dict["downloaded"])
	resp.WarningMessage = toString(dict["warning message"])

	if peersValue, exists := dict["peers"]; exists {
		peers, raw, err := parsePeers(peersValue)
		if err != nil {
			return nil, err
		}
		resp.Peers = peers
		resp.RawPeersCompact = raw
	}

	return resp, nil
}

// GeneratePeerID creates a 20-byte peer id using the project version.
func GeneratePeerID() []byte {
	peerID := make([]byte, 20)
	copy(peerID, []byte(peerIDPrefix))

	remaining := peerID[len(peerIDPrefix):]
	if _, err := rand.Read(remaining); err != nil {
		for i := range remaining {
			remaining[i] = byte('0' + i%10)
		}
		return peerID
	}

	for i, b := range remaining {
		remaining[i] = alnumByte(b)
	}
	return peerID
}

func parsePeers(value interface{}) ([]Peer, []byte, error) {
	switch v := value.(type) {
	case string:
		return parseCompactPeers([]byte(v))
	case []byte:
		return parseCompactPeers(v)
	case []interface{}:
		return parsePeerList(v)
	default:
		return nil, nil, fmt.Errorf("unsupported peers type: %T", value)
	}
}

func parseCompactPeers(data []byte) ([]Peer, []byte, error) {
	if len(data)%6 != 0 {
		return nil, nil, fmt.Errorf("compact peers length must be a multiple of 6, got %d", len(data))
	}

	peers := make([]Peer, 0, len(data)/6)
	for i := 0; i < len(data); i += 6 {
		ip := net.IPv4(data[i], data[i+1], data[i+2], data[i+3]).String()
		port := int(data[i+4])<<8 | int(data[i+5])
		peers = append(peers, Peer{IP: ip, Port: port, Compact: true})
	}

	return peers, data, nil
}

func parsePeerList(items []interface{}) ([]Peer, []byte, error) {
	peers := make([]Peer, 0, len(items))
	for _, item := range items {
		entry, ok := item.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("expected peer dictionary, got %T", item)
		}

		ip := toString(entry["ip"])
		port := toInt(entry["port"])
		if ip == "" || port == 0 {
			return nil, nil, fmt.Errorf("peer entry missing ip or port")
		}

		peer := Peer{IP: ip, Port: port}
		if peerID, ok := entry["peer id"].(string); ok {
			peer.PeerID = peerID
		}
		peers = append(peers, peer)
	}

	return peers, nil, nil
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int64:
		return int(val)
	case int:
		return val
	case uint64:
		return int(val)
	default:
		return 0
	}
}

func toString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func boolToInt(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func alnumByte(b byte) byte {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return chars[int(b)%len(chars)]
}
