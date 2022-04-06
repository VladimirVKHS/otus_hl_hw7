package otuswsserver

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"log"
	"net/http"
	jwt_helper "otus_sn_wsserver/internal/helpers/jwt"
	"otus_sn_wsserver/internal/logger"
	"otus_sn_wsserver/internal/otusrabbit"
	"strconv"
)

var Server *socketio.Server

// Easier to get running with CORS. Thanks for help @Vindexus and @erkie
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func Init() {
	Server = socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	Server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected")
		return nil
	})

	Server.OnEvent("/", "authorization", func(s socketio.Conn, msg string) {
		tokenString := msg
		userId, err := jwt_helper.GetUserIdFromToken(tokenString)
		if err != nil || userId == 0 {
			logger.Error("user id not recognized: " + err.Error())
			s.Close()
			fmt.Println("invalid user id")
		}
		s.SetContext(userId)
		roomName := strconv.Itoa(userId)
		s.Join(roomName)
		fmt.Println("authorized user_id: ", userId)
		otusrabbit.BindUserId(userId)
	})

	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		ctx := s.Context()
		if ctx == nil {
			return
		}
		userId := ctx.(int)
		roomName := strconv.Itoa(userId)
		if roomName == "" {
			return
		}
		roomLen := Server.RoomLen("/", roomName)
		if roomLen == 0 {
			otusrabbit.UnbindUserId(userId)
		}
		fmt.Println("closed", reason, roomName, roomLen)
	})

	go func() {
		if err := Server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
}

func Close() {
	Server.Close()
}
