package dlx

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"middleware/rabbitmq"
	"strconv"
	"time"
)

const (
	DLXExchangeName    = "dlx_exchange"
	CommonExchangeName = "common_exchange"
	MessageDDL         = 6000
	CommonRouterKey    = "common_router_key"
	DLXRouterKey       = "dlx_router_key"
)

func producer() {
	//连接rabbitmq
	//url格式: amqp://账号:密码@RabbitMQ地址:端口/
	conn, err := amqp.Dial("amqp://user:111111@192.168.101.25:5672")
	defer conn.Close()
	rabbitmq.ErrorHandle(err, "rabbitmq connect fail")

	ch, err := conn.Channel()
	rabbitmq.ErrorHandle(err, "Failed to get Channel")

	//声明一个普通交换机
	//func (ch *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table)
	err = ch.ExchangeDeclare(
		CommonExchangeName,
		amqp.ExchangeDirect,
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to declare exchange")

	//声明死信交换机
	//func (ch *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table)
	err = ch.ExchangeDeclare(
		DLXExchangeName,
		amqp.ExchangeDirect,
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to declare exchange")

	//通过设置参数Table来指定死信交换机和消息的过期时间
	table := make(amqp.Table)
	table["x-dead-letter-exchange"] = DLXExchangeName
	table["x-message-ttl"] = MessageDDL
	table["x-dead-letter-routing-key"] = DLXRouterKey
	//声明一个普通队列，并绑定到普通交换机上
	//func (ch *Channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error)
	commonQueue, err := ch.QueueDeclare(
		"common_queue",
		false,
		false,
		false,
		false,
		table,
	)
	rabbitmq.ErrorHandle(err, "Failed to declare queue")

	//func (ch *Channel) QueueBind(name, key, exchange string, noWait bool, args Table)
	err = ch.QueueBind(
		commonQueue.Name,
		CommonRouterKey,
		CommonExchangeName,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to bind queue")

	//声明一个死信队列(就是将一个普通队列绑定到死信交换机上)
	dlxQueue, err := ch.QueueDeclare(
		"dlx_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to declare queue")

	//func (ch *Channel) QueueBind(name, key, exchange string, noWait bool, args Table)
	err = ch.QueueBind(
		dlxQueue.Name,
		DLXRouterKey,
		DLXExchangeName,
		false,
		nil,
	)
	rabbitmq.ErrorHandle(err, "Failed to bind queue")

	//发布10个任务
	//func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg Publishing)
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		n := rand.Intn(10)
		log.Printf("Publish %d seconds work\n", n)
		err = ch.Publish(
			CommonExchangeName,
			CommonRouterKey,
			false,
			false,
			amqp.Publishing{
				Type: "text/plain",
				Body: []byte(strconv.Itoa(n)),
			})
		rabbitmq.ErrorHandle(err, "Failed to publish work")
	}
}
