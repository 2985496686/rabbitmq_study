package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func ErrorHandle(err error, mess string) {
	if err != nil {
		log.Fatalf("%s: %s", err, mess)
	}
}

func ReceiveMessage() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.1.207:5672")
	defer conn.Close()
	ErrorHandle(err, "rabbitmq connect fail")

	//创建信道
	ch, err := conn.Channel()
	defer ch.Close()
	ErrorHandle(err, "Filed to create channel")

	//声明要操作的队列
	//若声明队列的参数不符合创建队列的条件，将会返回一个异常，并且关闭channel
	queue, err := ch.QueueDeclare(
		"queue1", //name :队列的名字，可以为空，为空是会主动生成唯一的name，并在并在队列的name字段返回
		false,    //durable: 是否持久化
		false,    //autoDelete: 自动删除
		false,    //exclusive: 是否独占，独占队列只能由声明它们的连接访问，若其它channel对该队列进行操作会返回error
		false,
		nil, //args :初始化参数
	)
	ErrorHandle(err, "Fail")
	message, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	ErrorHandle(err, "Failed to receive message")
	for m := range message {
		fmt.Println(string(m.Body))
	}
}

func PublishMessage() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.1.207:5672")
	defer conn.Close()
	ErrorHandle(err, "rabbitmq connect fail")

	//创建信道
	ch, err := conn.Channel()
	defer ch.Close()
	ErrorHandle(err, "Filed to create channel")

	//声明要操作的队列
	//若声明队列的参数不符合创建队列的条件，将会返回一个异常，并且关闭channel
	queue, err := ch.QueueDeclare(
		"queue1", //name :队列的名字，可以为空，为空是会主动生成唯一的name，并在并在队列的name字段返回
		false,    //durable: 是否持久化
		false,    //autoDelete: 自动删除
		false,    //exclusive: 是否独占，独占队列只能由声明它们的连接访问，若其它channel对该队列进行操作会返回error
		false,
		nil, //args :初始化参数
	)
	ErrorHandle(err, "Fail")

	//要发送的消息
	body := "hello world"

	err = ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	ErrorHandle(err, "Failed to push message")
}
