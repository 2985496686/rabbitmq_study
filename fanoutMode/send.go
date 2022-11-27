package fanoutMode

import (
	"github.com/streadway/amqp"
	"middleware/rabbitmq"
	"strconv"
)

func send() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.1.207:5672")
	defer conn.Close()
	rabbitmq.ErrorHandle(err, "rabbitmq connect fail")

	//创建通道
	ch, err := conn.Channel()
	rabbitmq.ErrorHandle(err, "Failed to create channel")

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	for i := 0; i < 10; i++ {
		rabbitmq.ErrorHandle(err, "Failed to declare exchange")
		body := "message" + strconv.Itoa(i)
		err = ch.Publish(
			"logs",
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		rabbitmq.ErrorHandle(err, "Failed to publish")
	}
}
