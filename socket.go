package projectx

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

const (
	RecordSep = 0x1E
)

type WebSocketHandshakeResponse struct {
	NegotiateVersion    int           `json:"negotiateVersion"`
	ConnectionId        string        `json:"connectionId"`
	ConnectionToken     string        `json:"connectionToken"`
	AvailableTransports []interface{} `json:"availableTransports"`
}

// WebSocket handshake (minimal, validates Sec-WebSocket-Accept)
func dialWebsocket(u *url.URL) (net.Conn, error) {
	// TLS dial
	conn, err := tls.Dial("tcp", u.Host+":443", &tls.Config{})
	if err != nil {
		return nil, err
	}

	// Sec-WebSocket-Key
	key := make([]byte, 16)
	_, _ = rand.Read(key)
	wsKey := base64.StdEncoding.EncodeToString(key)

	// Request
	req := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\n\r\n", u.RequestURI(), u.Host, wsKey)
	if _, err := conn.Write([]byte(req)); err != nil {
		conn.Close()
		return nil, err
	}

	// Read HTTP response header (read until double CRLF)
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return nil, err
	}
	resp := string(buf[:n])
	if resp == "" || (len(resp) >= 12 && resp[:12] != "HTTP/1.1 101") && !contains101(resp) {
		conn.Close()
		return nil, fmt.Errorf("websocket upgrade failed: %s", resp)
	}

	// Validate Sec-WebSocket-Accept
	expected := computeAccept(wsKey)
	if !contains(resp, expected) {
		conn.Close()
		return nil, fmt.Errorf("invalid Sec-WebSocket-Accept")
	}
	return conn, nil
}

func contains101(s string) bool   { return bytes.Contains([]byte(s), []byte("101 Switching Protocols")) }
func contains(s, sub string) bool { return bytes.Contains([]byte(s), []byte(sub)) }

func computeAccept(key string) string {
	const GUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.Sum([]byte(key + GUID))
	return base64.StdEncoding.EncodeToString(h[:])
}

// send a masked text frame
func WriteTextFrame(conn net.Conn, payload []byte) error {
	frame := buildClientFrame(payload)
	_, err := conn.Write(frame)
	return err
}

func buildClientFrame(payload []byte) []byte {
	// FIN + text
	b := []byte{0x81}
	l := len(payload)
	if l < 126 {
		b = append(b, byte(0x80|byte(l))) // mask bit set
	} else {
		// minimal support for <65535
		b = append(b, 0x80|126, byte(l>>8), byte(l&0xff))
	}

	// mask
	mask := make([]byte, 4)
	_, _ = rand.Read(mask)
	b = append(b, mask...)

	// masked payload
	m := make([]byte, l)
	for i := 0; i < l; i++ {
		m[i] = payload[i] ^ mask[i%4]
	}
	return append(b, m...)
}

// readText reads a FULL text message (handles fragmentation)
func ReadFullText(conn net.Conn) ([]byte, error) {
	var message []byte
	started := false

	for {
		// --- read first 2 bytes ---
		h := make([]byte, 2)
		if _, err := io.ReadFull(conn, h); err != nil {
			log.Println("failed to io.readfull")
			return nil, err
		}

		fin := (h[0] & 0x80) != 0
		op := h[0] & 0x0F

		// --- opcode handling ---
		switch op {
		case 0x1: // text frame (start)
			if started {
				return nil, fmt.Errorf("unexpected new text frame")
			}
			started = true

		case 0x0: // continuation
			if !started {
				return nil, fmt.Errorf("continuation without start")
			}

		case 0x8: // close
			return nil, io.EOF

		case 0x9: // ping (ignore or pong)
			// you should send pong here if you want to be correct
			continue

		case 0xA: // pong
			continue

		default:
			return nil, fmt.Errorf("unsupported opcode: %d", op)
		}

		// --- payload length ---
		lenByte := int(h[1] & 0x7F)
		var length int

		if lenByte < 126 {
			length = lenByte
		} else if lenByte == 126 {
			ext := make([]byte, 2)
			if _, err := io.ReadFull(conn, ext); err != nil {
				return nil, err
			}
			length = int(ext[0])<<8 | int(ext[1])
		} else {
			ext := make([]byte, 8)
			if _, err := io.ReadFull(conn, ext); err != nil {
				return nil, err
			}
			length = int(ext[6])<<8 | int(ext[7]) // small messages assumption
		}

		// --- read payload ---
		payload := make([]byte, length)
		if _, err := io.ReadFull(conn, payload); err != nil {
			return nil, err
		}

		message = append(message, payload...)

		// --- end of message ---
		if fin {
			return message, nil
		}
	}
}

// build websocket url
func buildWSURL(base, connToken, token string) *url.URL {
	u, _ := url.Parse(base)
	q := u.Query()
	q.Set("id", connToken)
	q.Set("access_token", token)
	q.Set("transport", "webSockets")
	q.Set("negotiateVersion", "1")
	u.RawQuery = q.Encode()
	return u
}

// negotiate the connection token
func negotiate(negotiateUrl, token string) (string, error) {

	// Negotiate
	req, err := http.NewRequest("POST", negotiateUrl+"&access_token="+token, nil)
	if err != nil {
		log.Fatalf("An Error Occured: %v", err)
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	var resMap WebSocketHandshakeResponse

	err = json.Unmarshal([]byte(body), &resMap)
	if err != nil {
		log.Fatalln("Error unmarshalling json:", err)
		return "", err
	}

	return resMap.ConnectionToken, nil
}

func GetWsConn(negotiateUrl, wsUrl, token string) (net.Conn, error) {
	// Get Market Connection
	connToken, err := negotiate(negotiateUrl, token)
	if err != nil {
		panic(err)
	}

	// Step 2: Build websocket URL using connectionToken
	wsMarket := buildWSURL(wsUrl, connToken, token)

	// Step 3: Dial + WebSocket handshake
	conn, err := dialWebsocket(wsMarket)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	hs := fmt.Sprintf(`{"protocol":"json","version":1}`) + string(byte(RecordSep))
	WriteTextFrame(conn, []byte(hs))

	// read server ack (simple read; assumes first frame fits)
	msg, err := ReadFullText(conn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("projectx server handshake ack: %q\n", msg)
	return conn, nil
}
