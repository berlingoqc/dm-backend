package main

import (
	"github.com/berlingoqc/dm-backend/dm"
)

func main() {
	ws := dm.Run()
	<-ws.ChannelStop
}
