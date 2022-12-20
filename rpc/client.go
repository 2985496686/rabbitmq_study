package rpc

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"middleware/rabbitmq"
	"strconv"
)

func randInt(max, min int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes = append(bytes, byte(randInt(200, 0)))
	}
	return string(bytes)
}

func rpcClient() (result int, err error) {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.101.25:5672")
	defer conn.Close()
	rabbitmq.ErrorHandle(err, "rabbitmq connect fail")

	//创建通道
	ch, err := conn.Channel()
	rabbitmq.ErrorHandle(err, "Failed to create channel")

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to declare queue")
	corrId := randomString(32)

	n := randInt(30, 10)
	log.Printf("[client] publish %d\n", n)
	err = ch.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	rabbitmq.ErrorHandle(err, "Failed to publish")

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to consume")
	for msg := range msgs {
		if msg.CorrelationId == corrId {
			result, err = strconv.Atoi(string(msg.Body))
		}
		log.Printf("[client] %d -> result:%d\n", n, result)
		break
	}
	return
}
