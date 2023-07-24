package webrtc

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const PING_INTERVAL = 55 * time.Second

// Helper to make Gorilla Websockets threadsafe
type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
	ticker    *time.Ticker
	closeChan chan struct{}
}

func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()
	t.ticker.Reset(PING_INTERVAL)
	return t.Conn.WriteJSON(v)
}

func newThreadSafeWriter(conn *websocket.Conn) *threadSafeWriter {

	ticker := time.NewTicker(PING_INTERVAL)
	defer ticker.Stop()
	ws := &threadSafeWriter{Conn: conn, Mutex: sync.Mutex{}, ticker: ticker, closeChan: make(chan struct{})}

	go func() {
		select {
		case <-ticker.C:
			ws.Lock()
			if err := ws.WriteMessage(websocket.PingMessage, []byte("keepAlive")); err != nil {
				log.Println("Error pinging websocket")
				return
			}
			ws.Unlock()
		case <-ws.closeChan:
			return
		}
	}()

	return ws
}

func (t *threadSafeWriter) Close() {
	t.closeChan <- struct{}{}
	t.Conn.Close()
}
