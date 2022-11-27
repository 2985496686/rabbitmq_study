package workMode

import (
	"github.com/streadway/amqp"
	"log"
	"middleware/rabbitmq"
	"os"
	"strconv"
	"strings"
	"time"
)

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func send() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.1.207:5672")
	defer conn.Close()
	rabbitmq.ErrorHandle(err, "rabbitmq connect fail")

	//创建通道
	ch, err := conn.Channel()
	rabbitmq.ErrorHandle(err, "Failed to create channel")

	//声明队列
	q, err := ch.QueueDeclare(
		"queue2",
		false,
		false,
		false,
		false,
		nil,
	)
	err = ch.Qos(
		2,    //prefetchCount
		0,    //prefetchSize
		true, //global
	)
	for i := 0; i < 10; i++ {
		body := "message" + strconv.Itoa(i)
		//push
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, //持久交互
				ContentType:  "test/plain",
				Body:         []byte(body),
			},
		)
		rabbitmq.ErrorHandle(err, "Failed to push message")
		log.Printf("[x] send %s\n", body)
		time.Sleep(1 * time.Second)
	}
}
