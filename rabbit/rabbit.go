package rabbit

import (
	"audit_engine/tool"
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
	ch   *amqp.Channel
}

func (mq *MQ) Close() {
	mq.conn.Close()
	mq.ch.Close()
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
	mq.ch, err = mq.conn.Channel()
	tool.FatalLog(err, "failed to open a channel")
}

//队列创建
func (mq *MQ) Create(qn string) amqp.Queue {
	durable, autoDelete := true, false
	if qn == "" {
		durable = false
		autoDelete = true
	}
	q, err := mq.ch.QueueDeclare(
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
func (mq *MQ) Consume(qn string, fn func([]byte) bool, noAck bool) {
	//consume resister
	msgs, err := mq.ch.Consume(
		qn,
		"audit-engine",
		false,
		false,
		false,
		false,
		nil,
	)
	tool.FatalLog(err, "failed to register a consumer")

	//set qos
	err = mq.ch.Qos(5, 0, false)
	tool.FatalLog(err, "failed to set channel qos")
	//consume work
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			log.Printf("==> [%s] task start...", qn)
			success := fn(d.Body)
			log.Printf("<== [%s] task done, result: [%v]!!", qn, success)
			log.Println("[*] waiting for message. To exit press CTRL+C")

			if success && !noAck {
				d.Ack(false)
			}
		}
	}()
	log.Println("[*] waiting for message. To exit press CTRL+C")
	<-forever
}

//队列发布消息
func (mq *MQ) Publish(qn string, data []byte, n int) {
	//publish msg
	msg := amqp.Publishing{
		Body:        data,
		ContentType: "text/plain",
	}

	for i := 0; i < n; i++ {
		err := mq.ch.Publish("", qn, false, false, msg)
		tool.FatalLog(err, "failed to publish a message")
		log.Println("send message finish")
	}
}
