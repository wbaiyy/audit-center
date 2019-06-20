package rabbit

import (
	"audit-center/tool"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Config struct {
	Host  string
	Port  int
	User  string
	Pass  string
	Vhost string
}

type MQ struct {
	conn *amqp.Connection
	//ch   *amqp.Channel
	Channels  map[string]*amqp.Channel
}

func (mq *MQ) Close() {
	for _, ch := range mq.Channels {
		ch.Close();
	}
	mq.conn.Close()
	//mq.ch.Close()
}

//队列初始化
func (mq *MQ) Init(mqcf Config) {
	var err error
	//conn
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", mqcf.User, mqcf.Pass, mqcf.Host, mqcf.Port, mqcf.Vhost)
	log.Println(url)
	mq.conn, err = amqp.Dial(url)
	tool.FatalLog(err, "failed to connect to RabbitMQ")

	//channel
	//mq.ch, err = mq.conn.Channel()
	//tool.FatalLog(err, "failed to open a channel")
	//
	mq.Channels = make(map[string]*amqp.Channel)
}

func (mq *MQ) GetChannel(queueName string ) *amqp.Channel{
	channel, err := mq.conn.Channel()
	tool.FatalLog(err, "failed to open a channel")
	mq.Channels[queueName] = channel

	return channel
}

//队列创建
func (mq *MQ) Create(qn string) amqp.Queue {
	durable, autoDelete := true, false
	if qn == "" {
		durable = false
		autoDelete = true
	}

	q, err := mq.Channels[qn].QueueDeclare(
		qn,
		durable,
		autoDelete,
		false,
		false,
		nil,
	)

	tool.FatalLog(err, "failed to declare queue")
	return q
}

//队列消费程序绑定
func (mq *MQ) Consume(qn string)<-chan amqp.Delivery {
	//set qos
	err := mq.Channels[qn].Qos(5, 0, false)
	tool.FatalLog(err, "failed to set channel qos")

	//consume resister
	msgs, err := mq.Channels[qn].Consume(
		qn,
		"audit-center",
		false,
		false,
		false,
		false,
		nil,
	)
	tool.FatalLog(err, "failed to register a consumer")

	return msgs
}

//队列发布消息
func (mq *MQ) Publish(qn string, data []byte, n int) {
	//publish msg
	msg := amqp.Publishing{
		Body:        data,
		ContentType: "text/plain",
	}
	channel, isExist := mq.Channels[qn]
	if !isExist {
		channel = mq.GetChannel(qn)
	}
	log.Println(fmt.Sprintf("==>[%s] send result: %s", qn, string(data)))

	for i := 0; i < n; i++ {
		err := channel.Publish("", qn, false, false, msg)
		tool.FatalLog(err, "failed to publish a message")
		//log.Println("send message finish")
	}
}
