package webrtc

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const PING_INTERVAL = 50 * time.Second

// Helper to make Gorilla Websockets threadsafe
type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
	ticker *time.Ticker
}

func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()
	t.ticker.Reset(PING_INTERVAL)
	return t.Conn.WriteJSON(v)
}

func newThreadSafeWriter(conn *websocket.Conn) *threadSafeWriter {

	ticker := time.NewTicker(PING_INTERVAL)
	ws := &threadSafeWriter{Conn: conn, Mutex: sync.Mutex{}, ticker: ticker}

	go func() {
		for range ticker.C {
			ws.Lock()
			if err := ws.WriteMessage(websocket.PingMessage, []byte("keepAlive")); err != nil {
				log.Println("Error pinging websocket")
				return
			}
			ws.Unlock()
		}
		ticker.Stop()
	}()

	return ws
}
