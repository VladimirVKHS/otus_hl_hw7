package otusrabbit

import (
	"github.com/streadway/amqp"
	"os"
	"strconv"
)

var Ch *amqp.Channel
var Queue amqp.Queue

func Init() {
	host, _ := os.LookupEnv("RABBIT_HOST")
	port, _ := os.LookupEnv("RABBIT_PORT")
	wsExchange, _ := os.LookupEnv("RABBIT_WS_EXCHANGE")
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")
	if err != nil {
		panic(err)
	}

	Ch, err = conn.Channel()
	if err != nil {
		panic(err)
	}

	err = Ch.ExchangeDeclare(
		wsExchange, // name
		"direct",   // type
		false,      // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		panic(err)
	}

	Queue, err = Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}
}

func BindUserId(userId int) error {
	return Ch.QueueBind(
		Queue.Name,           // queue name
		strconv.Itoa(userId), // routing key
		"ws",                 // exchange
		false,
		nil,
	)
}

func UnbindUserId(userId int) error {
	return Ch.QueueUnbind(
		Queue.Name,           // queue name
		strconv.Itoa(userId), // routing key
		"ws",                 // exchange
		nil,
	)
}

func Close() {
	Ch.Close()
}
