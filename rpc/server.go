package rpc

import (
	"github.com/streadway/amqp"
	"log"
	"middleware/rabbitmq"
	"strconv"
)

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func rpcServer() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.101.25:5672")
	defer conn.Close()
	rabbitmq.ErrorHandle(err, "rabbitmq connect fail")

	//创建通道
	ch, err := conn.Channel()
	rabbitmq.ErrorHandle(err, "Failed to create channel")

	rpcQueue, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "rpcQueue declare failed")

	err = ch.Qos(
		1,
		0,
		true,
	)

	rabbitmq.ErrorHandle(err, "Failed to set Qos")

	msgs, err := ch.Consume(
		rpcQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	wait := make(chan bool)
	go func() {
		for msg := range msgs {
			n, err := strconv.Atoi(string(msg.Body))
			log.Printf("[server] get %d from client", n)
			rabbitmq.ErrorHandle(err, "Failed to convert string to integer")
			result := fib(n)
			log.Printf("[server] publish %d", result)
			err = ch.Publish(
				"",
				msg.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: msg.CorrelationId,
					Body:          []byte(strconv.Itoa(result)),
				})
			rabbitmq.ErrorHandle(err, "Failed to publish message")
			err = msg.Ack(false)
			rabbitmq.ErrorHandle(err, "Failed to ack message")
		}
	}()
	<-wait
}
