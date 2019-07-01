package rpcproxy

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// WSMessage ...
type WSMessage struct {
	Namespace string  `json:"domain"`
	Data      RPCCall `json:"data"`
}

// WSMessageTrapper ...
var WSMessageTrapper func(WSMessage)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make([]*websocket.Conn, 0)

var clientMessageChannel = make(chan WSMessage)

// SendMessageWS ...
func SendMessageWS(domain string, namespace string, event string, data interface{}) {
	msg := WSMessage{
		Namespace: domain,
		Data: RPCCall{
			Jsonrpc: "2.0",
			ID:      "qwer",
			Method:  namespace + "." + event,
			Result:  []interface{}{data},
		},
	}
	clientMessageChannel <- msg
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println(err.Error())
		return
	}
	clients = append(clients, c)
	c.WriteJSON(WSMessage{Data: RPCCall{}})
}

func clientMessageHandler() {
	for {
		v, ok := <-clientMessageChannel
		if ok == false {
			println("FALSE returning from clientMessageChannel")
			break
		}

		go WSMessageTrapper(v)

		println("BROADCAST message from", v.Namespace)
		for _, client := range clients {
			if err := client.WriteJSON(v); err != nil {
				println("ERROR seding to ", client.RemoteAddr())
			}
		}
	}
}

// StartWebSocketClient ...
func StartWebSocketClient(info RPCHandlerEndpoint) {
	u := url.URL{Scheme: "ws", Host: info.URL, Path: "/jsonrpc"}

	println("connecting to ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		println("ERROR connecting ws ", err.Error())
	}

	go func() {
		defer func() {
			println("DEFFFERRRING")
		}()
		for {
			if c == nil {
				return
			}
			_, msg, err := c.ReadMessage()
			if err != nil {
				println("ERROR reading ", err.Error())
				return
			}
			rpcCall := RPCCall{}
			if err := json.Unmarshal(msg, &rpcCall); err != nil {
				println("FAILED to unmarshall RPCCALL event ", err.Error())
			}

			clientMessageChannel <- WSMessage{
				Namespace: info.Namespace,
				Data:      rpcCall,
			}
		}
	}()
}

func init() {
	go clientMessageHandler()
}
