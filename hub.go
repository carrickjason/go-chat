package main

type Hub struct {
	// Registered chat clients
	clients map[*Client]bool

	// received messages from clients
	broadcast chan []byte

	// reqeusts to register clients
	register chan *Client

	// reqeusts to unregister clients
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (hub *Hub) run() {
	// infinite loop?
	for {
		// Select blocks until any channel is prepared to send
		select {
		// hub.register channel sends to client
		// So if register channel gets a client, we add that client to the clients map
		case client := <-hub.register:
			hub.clients[client] = true
		// and remove when client is added to unregister channel (if exists)
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		case message := <-hub.broadcast:
			// loop over registered clients and send message
			for client := range hub.clients {
				select {
				case client.send <- message:
				// Does this run if there is no message or if client.send is not ready to recieve?
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}
