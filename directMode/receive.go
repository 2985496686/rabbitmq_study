package directMode

import (
	"github.com/streadway/amqp"
	"log"
	"middleware/rabbitmq"
)

func receiveLog1() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.1.207:5672")
	defer conn.Close()
	rabbitmq.ErrorHandle(err, "rabbitmq connect fail")

	//创建通道
	ch, err := conn.Channel()
	rabbitmq.ErrorHandle(err, "Failed to create channel")

	//声明队列
	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	//将队列绑定到交换机上
	err = ch.QueueBind(
		queue.Name,
		"log1", //router key
		"direct.logs",
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Filed to bind queue")

	message, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to consume")
	for m := range message {
		log.Printf("[log1] consum %s\n", []byte(m.Body))
	}
}

func receiveLog2() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.1.207:5672")
	defer conn.Close()
	rabbitmq.ErrorHandle(err, "rabbitmq connect fail")

	//创建通道
	ch, err := conn.Channel()
	rabbitmq.ErrorHandle(err, "Failed to create channel")

	//声明队列
	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	//将队列绑定到交换机上
	err = ch.QueueBind(
		queue.Name,
		"log2", //router key
		"direct.logs",
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Filed to bind queue")

	message, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to consume")
	for m := range message {
		log.Printf("[log2] consum %s\n", []byte(m.Body))
	}
}
