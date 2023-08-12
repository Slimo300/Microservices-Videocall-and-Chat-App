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

	t := &threadSafeWriter{Conn: conn, Mutex: sync.Mutex{}, ticker: time.NewTicker(PING_INTERVAL), closeChan: make(chan struct{})}

	t.SetPongHandler(func(string) error {
		log.Println("PONG")
		return nil
	})

	go func() {
		defer log.Println("Out of goroutine")
		for {
			select {
			case <-t.ticker.C:
				t.Lock()
				log.Println("PING")
				if err := t.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("Error pinging websocket")
					return
				}
				t.Unlock()
			case <-t.closeChan:
				t.ticker.Stop()
				return
			}
		}
	}()

	return t
}

func (t *threadSafeWriter) Close() {
	t.closeChan <- struct{}{}
	t.Conn.Close()
}
