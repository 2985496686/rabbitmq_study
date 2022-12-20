package dlx

import (
	"github.com/streadway/amqp"
	"log"
	"middleware/rabbitmq"
	"strconv"
	"time"
)

func comConsumer() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.101.25:5672")
	defer conn.Close()
	rabbitmq.ErrorHandle(err, "rabbitmq connect fail")

	ch, err := conn.Channel()
	rabbitmq.ErrorHandle(err, "Failed to get Channel")

	//func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table)
	msgs, err := ch.Consume(
		"common_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	err = ch.Qos(
		1,
		0,
		true,
	)
	for msg := range msgs {
		n, err := strconv.Atoi(string(msg.Body))
		rabbitmq.ErrorHandle(err, "Failed convert string to integer")
		log.Printf("[comConsumer] start work ,it will take %d seconds\n", n)
		time.Sleep(time.Second * time.Duration(n))
		msg.Ack(false)

		log.Printf("[comConsumer] complete work\n")
	}
}
