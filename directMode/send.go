package directMode

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

	//声明交换器
	err = ch.ExchangeDeclare(
		"direct.logs",
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Filed to declare exchange")

	//发送消息
	for i := 0; i < 10; i++ {
		body := "message" + strconv.Itoa(i)
		err = ch.Publish(
			"direct.logs",
			"log1",
			false,
			false,
			amqp.Publishing{
				Body:        []byte(body),
				ContentType: "text/plain",
			},
		)
	}
	rabbitmq.ErrorHandle(err, "Failed to publish message")
}
