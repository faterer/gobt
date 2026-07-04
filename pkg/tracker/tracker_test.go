package tracker

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gobt/pkg/bencode"
	"gobt/pkg/torrent"
)

func TestBuildAnnounceURL(t *testing.T) {
	infoHash := make([]byte, 20)
	peerID := make([]byte, 20)
	for i := 0; i < 20; i++ {
		infoHash[i] = byte(i)
		peerID[i] = byte(20 - i)
	}

	urlStr, err := BuildAnnounceURL("http://tracker.example.com:6969/announce", AnnounceRequest{
		InfoHash:   infoHash,
		PeerID:     peerID,
		Port:       6881,
		Uploaded:   123,
		Downloaded: 456,
		Left:       789,
		Event:      "started",
		Compact:    true,
		NumWant:    50,
	})
	if err != nil {
		t.Fatalf("BuildAnnounceURL returned error: %v", err)
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("url.Parse failed: %v", err)
	}

	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		t.Fatalf("url.ParseQuery failed: %v", err)
	}

	if got := query.Get("port"); got != "6881" {
		t.Fatalf("expected port 6881, got %s", got)
	}
	if got := query.Get("uploaded"); got != "123" {
		t.Fatalf("expected uploaded 123, got %s", got)
	}
	if got := query.Get("downloaded"); got != "456" {
		t.Fatalf("expected downloaded 456, got %s", got)
	}
	if got := query.Get("left"); got != "789" {
		t.Fatalf("expected left 789, got %s", got)
	}
	if got := query.Get("event"); got != "started" {
		t.Fatalf("expected event started, got %s", got)
	}
	if got := query.Get("compact"); got != "1" {
		t.Fatalf("expected compact 1, got %s", got)
	}

	if got := []byte(query.Get("info_hash")); len(got) != 20 {
		t.Fatalf("expected 20-byte info_hash, got %d", len(got))
	} else {
		for i := 0; i < 20; i++ {
			if got[i] != byte(i) {
				t.Fatalf("info_hash byte %d mismatch: expected %d, got %d", i, byte(i), got[i])
			}
		}
	}

	if got := []byte(query.Get("peer_id")); len(got) != 20 {
		t.Fatalf("expected 20-byte peer_id, got %d", len(got))
	} else {
		for i := 0; i < 20; i++ {
			if got[i] != byte(20-i) {
				t.Fatalf("peer_id byte %d mismatch: expected %d, got %d", i, byte(20-i), got[i])
			}
		}
	}
}

func TestParseAnnounceResponseCompactPeers(t *testing.T) {
	responseDict := map[string]interface{}{
		"interval":  int64(1800),
		"complete":  int64(12),
		"incomplete": int64(3),
		"peers": []byte{
			1, 2, 3, 4, 0x1a, 0xe1,
			5, 6, 7, 8, 0x1a, 0xe2,
		},
	}

	encoded, err := bencode.Encode(responseDict)
	if err != nil {
		t.Fatalf("bencode.Encode failed: %v", err)
	}

	resp, err := ParseAnnounceResponse(encoded)
	if err != nil {
		t.Fatalf("ParseAnnounceResponse returned error: %v", err)
	}

	if resp.Interval != 1800 {
		t.Fatalf("expected interval 1800, got %d", resp.Interval)
	}
	if resp.Complete != 12 {
		t.Fatalf("expected complete 12, got %d", resp.Complete)
	}
	if resp.Incomplete != 3 {
		t.Fatalf("expected incomplete 3, got %d", resp.Incomplete)
	}
	if len(resp.Peers) != 2 {
		t.Fatalf("expected 2 peers, got %d", len(resp.Peers))
	}

	if resp.Peers[0].IP != "1.2.3.4" || resp.Peers[0].Port != 6881 {
		t.Fatalf("unexpected first peer: %+v", resp.Peers[0])
	}
	if resp.Peers[1].IP != "5.6.7.8" || resp.Peers[1].Port != 6882 {
		t.Fatalf("unexpected second peer: %+v", resp.Peers[1])
	}
	if len(resp.RawPeersCompact) != 12 {
		t.Fatalf("expected raw compact peers length 12, got %d", len(resp.RawPeersCompact))
	}
}

func TestParseAnnounceResponseDictionaryPeers(t *testing.T) {
	responseDict := map[string]interface{}{
		"interval": int64(900),
		"peers": []interface{}{
			map[string]interface{}{"ip": "9.8.7.6", "port": int64(51413), "peer id": "peer-1"},
			map[string]interface{}{"ip": "1.1.1.1", "port": int64(6881)},
		},
	}

	encoded, err := bencode.Encode(responseDict)
	if err != nil {
		t.Fatalf("bencode.Encode failed: %v", err)
	}

	resp, err := ParseAnnounceResponse(encoded)
	if err != nil {
		t.Fatalf("ParseAnnounceResponse returned error: %v", err)
	}

	if len(resp.Peers) != 2 {
		t.Fatalf("expected 2 peers, got %d", len(resp.Peers))
	}
	if resp.Peers[0].IP != "9.8.7.6" || resp.Peers[0].Port != 51413 || resp.Peers[0].PeerID != "peer-1" {
		t.Fatalf("unexpected first peer: %+v", resp.Peers[0])
	}
	if resp.Peers[1].IP != "1.1.1.1" || resp.Peers[1].Port != 6881 {
		t.Fatalf("unexpected second peer: %+v", resp.Peers[1])
	}
}

func TestClientAnnounce(t *testing.T) {
	var sawUserAgent string
	var sawQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawUserAgent = r.Header.Get("User-Agent")
		sawQuery = r.URL.Query()

		responseDict := map[string]interface{}{
			"interval": int64(1200),
			"peers": []byte{
				10, 0, 0, 1, 0x1a, 0xe1,
			},
		}
		encoded, err := bencode.Encode(responseDict)
		if err != nil {
			t.Fatalf("bencode.Encode failed: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(encoded)
	}))
	defer server.Close()

	torrentInfo := &torrent.TorrentInfo{
		Announce: server.URL + "/announce",
		Info: torrent.InfoDict{
			Name:        "example.txt",
			PieceLength: 4,
			Pieces:      []byte("12345678901234567890"),
			Length:      12,
		},
	}

	client := &Client{
		HTTPClient: server.Client(),
		PeerID:     []byte("-GP420-test-peer-id-"),
		Port:       6881,
		UserAgent:  "gobt/test",
		Compact:    true,
	}

	resp, err := client.Announce(server.URL+"/announce", torrentInfo, "started")
	if err != nil {
		t.Fatalf("Announce returned error: %v", err)
	}

	if sawUserAgent != "gobt/test" {
		t.Fatalf("expected user-agent gobt/test, got %q", sawUserAgent)
	}
	if sawQuery.Get("event") != "started" {
		t.Fatalf("expected event started, got %q", sawQuery.Get("event"))
	}
	if sawQuery.Get("port") != "6881" {
		t.Fatalf("expected port 6881, got %q", sawQuery.Get("port"))
	}
	if sawQuery.Get("left") != "12" {
		t.Fatalf("expected left 12, got %q", sawQuery.Get("left"))
	}
	if resp.Interval != 1200 {
		t.Fatalf("expected interval 1200, got %d", resp.Interval)
	}
	if len(resp.Peers) != 1 {
		t.Fatalf("expected 1 peer, got %d", len(resp.Peers))
	}
	if resp.Peers[0].IP != "10.0.0.1" || resp.Peers[0].Port != 6881 {
		t.Fatalf("unexpected peer: %+v", resp.Peers[0])
	}
}

func TestParseAnnounceResponseFailure(t *testing.T) {
	responseDict := map[string]interface{}{
		"failure reason": "invalid info_hash",
	}
	encoded, err := bencode.Encode(responseDict)
	if err != nil {
		t.Fatalf("bencode.Encode failed: %v", err)
	}

	resp, err := ParseAnnounceResponse(encoded)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var trackerErr *TrackerError
	if !errors.As(err, &trackerErr) {
		t.Fatalf("expected TrackerError, got %T", err)
	}
	if resp == nil || resp.FailureReason != "invalid info_hash" {
		t.Fatalf("expected failure reason to be preserved, got %+v", resp)
	}
}
