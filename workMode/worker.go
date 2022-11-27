package workMode

import (
	"bytes"
	"github.com/streadway/amqp"
	"log"
	"middleware/rabbitmq"
	"sync"
	"time"
)

func receive() {
	var wg sync.WaitGroup
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

	message, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to receive message")
	wg.Add(1)
	go func() {
		for m := range message {
			log.Printf("receive message %s\n", string(m.Body))
			count := bytes.Count(m.Body, []byte("."))
			time.Sleep(time.Duration(count) * time.Second)
			log.Printf("Done\n")
		}
		wg.Done()
	}()
	wg.Wait()
}
