package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/gowebappdev/trace"
	"github.com/stretchr/objx"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
	avatar  Avatar
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("join new client")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("leave client")
		case msg := <-r.forward:
			r.tracer.Trace("get a new messge: ", msg.Message)
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("send to client")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("failed to send the msg")
				}
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("failed to get cookie:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
