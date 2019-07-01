package rpcproxy

// RPCWS ...
type RPCWS struct{}

// ActiveClient ...
func (r *RPCWS) ActiveClient() []*ActiveWSClient {
	var c []*ActiveWSClient
	for _, v := range activeClient {
		c = append(c, v)
	}
	return c
}

// ReConnect ...
func (r *RPCWS) ReConnect(namespace string) string {
	if active, ok := activeClient[namespace]; ok {
		if !active.Connected {
			if err := StartWebSocketClient(active.Endpoint); err != nil {
				panic(err)
			}
		} else {
			panic("already connected")
		}
	}
	return "OK"
}

// Disconnect ...
func (r *RPCWS) Disconnect(namespace string) string {
	if active, ok := activeClient[namespace]; ok {
		if active.Connected && active.Conn != nil {
			if err := active.Conn.Close(); err != nil {
				panic(err)
			}
		}
	}
	return "OK"
}
