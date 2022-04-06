package live_update

import (
	"encoding/json"
	"fmt"
	"otus_sn_wsserver/internal/otusrabbit"
	"otus_sn_wsserver/internal/otuswsserver"
	"strconv"
)

type WsData struct {
	UserId int
	Event  string
	Data   interface{}
}

func Init() {
	msgs, err := otusrabbit.Ch.Consume(
		otusrabbit.Queue.Name, // queue
		"",                    // consumer
		true,                  // auto ack
		false,                 // exclusive
		false,                 // no local
		false,                 // no wait
		nil,                   // args
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for d := range msgs {
			wsData := &WsData{}
			if err := json.Unmarshal(d.Body, wsData); err != nil {
				continue
			}
			roomName := strconv.Itoa(wsData.UserId)
			fmt.Println("DATA", wsData)
			otuswsserver.Server.BroadcastToRoom(
				"/",
				roomName,
				wsData.Event,
				wsData.Data,
			)
		}
	}()
}
