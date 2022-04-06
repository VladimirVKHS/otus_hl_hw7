package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	jwt_helper "otus_sn_wsserver/internal/helpers/jwt"
	live_update "otus_sn_wsserver/internal/live-update"
	"otus_sn_wsserver/internal/logger"
	"otus_sn_wsserver/internal/otusrabbit"
	"otus_sn_wsserver/internal/otuswsserver"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Config not found: .env!")
	}

	logger.Init()
	defer logger.Close()
	otusrabbit.Init()
	defer otusrabbit.Close()
	otuswsserver.Init()
	defer otuswsserver.Close()
	jwt_helper.Init()
	live_update.Init()

	http.Handle("/socket.io/", otuswsserver.Server)
	wsPort, _ := os.LookupEnv("WS_PORT")
	fmt.Println("WS server started, port: " + wsPort)
	log.Fatal(http.ListenAndServe(":"+wsPort, nil))

}
