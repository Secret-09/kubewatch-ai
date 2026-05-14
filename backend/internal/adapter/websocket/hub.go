package websocket

import (
    "context"
    "sync"

    "github.com/gorilla/websocket"
)

type Hub struct {
    clients    map[*Client]struct{}
    Register   chan *Client
    unregister chan *Client
    broadcast  chan []byte
    lock       sync.RWMutex
}

type Client struct {
    conn *websocket.Conn
    send chan []byte
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]struct{}),
        Register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan []byte, 32),
    }
}

func (h *Hub) Run(ctx context.Context) {
    for {
        select {
        case client := <-h.Register:
            h.lock.Lock()
            h.clients[client] = struct{}{}
            h.lock.Unlock()
            go client.writePump()
        case client := <-h.unregister:
            h.lock.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.lock.Unlock()
        case message := <-h.broadcast:
            h.lock.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.lock.RUnlock()
        case <-ctx.Done():
            h.lock.RLock()
            for client := range h.clients {
                close(client.send)
                _ = client.conn.Close()
            }
            h.lock.RUnlock()
            return
        }
    }
}

func (h *Hub) Broadcast(payload []byte) {
    select {
    case h.broadcast <- payload:
    default:
    }
}

func NewClient(conn *websocket.Conn) *Client {
    return &Client{conn: conn, send: make(chan []byte, 16)}
}

func (c *Client) writePump() {
    defer c.conn.Close()
    for message := range c.send {
        _ = c.conn.WriteMessage(websocket.TextMessage, message)
    }
}
