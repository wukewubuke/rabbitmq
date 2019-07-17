package RabbitMq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//url格式 amqp://user:pass@IP:5672/VirtualHost
const MQURL = "amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc"

type RabbitMq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//
	QueueName string
	//交换机
	Exchange string
	//key
	Key   string
	MqUrl string
}

//创建结构体实例
func NewRabbitMq(queueName, exchange, key string) *RabbitMq {
	rabbitmq := &RabbitMq{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MqUrl:     MQURL,
	}
	var err error
	//创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

func (r *RabbitMq) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMq) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

/*
简单模式只用传入queueName即可
*/

//简单模式step1:创建简单模式下rabbitmq实例
func NewRabbitMqSimple(queueName string) *RabbitMq {
	return NewRabbitMq(queueName, "", "")
}

//简单模式step2:生产代码
func (r *RabbitMq) PublishSimple(message string) {
	//1 申请队列，如果队列不存在则申请队列，如果存在则跳过创建,保证队列存在消息能在队列中
	_, err := r.channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性(不常用，消息只对自己账号可见)
		false,
		//是否阻塞
		false,
		//额外属性
		nil)
	if err != nil {
		fmt.Printf("queue declare failed, error:%+v\n", err)
		return
	}

	//2发送消息到队列中
	r.channel.Publish(
		//交换机
		r.Exchange, //简单模式为空
		//队列名称
		r.QueueName,
		//mandatory如果为true会根据exchange类型和routekey规则，
		//如果无法找到符合条件的队列就把消息返还给发送者
		false,
		//immediate如果为true，当exchage发送消息到队列后，
		//如果发现队列上没有绑定消费者就把消息返还给发送者
		false,
		//发送的信息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

//简单模式step3:消费消息
func (r *RabbitMq) ConsumeSimple() {
	//1 申请队列，如果队列不存在则申请队列，如果存在则跳过创建,保证队列存在消息能在队列中
	_, err := r.channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性(不常用，消息只对自己账号可见)
		false,
		//是否阻塞false是阻塞
		false,
		//额外属性
		nil)
	if err != nil {
		fmt.Printf("queue declare failed, error:%+v\n", err)
		return
	}

	//2接收消息
	msgs, err := r.channel.Consume(
		//队列名称
		r.QueueName,
		//用来区分多个消费者，为空不区分
		"",
		//autoAck是否自动应答，如果把消息消费了，主动告诉rabbitmq服务器
		true,
		//是否具有排他性,只能查看到我自己账号创建的消息，和pulish对应
		false,
		//noLocal如果设置为true，不能将同一个connection发送的消息，传递给同一个connection中的其他消费者
		false,
		//是否阻塞,false阻塞
		false,
		//额外参数
		nil)

	if err != nil {
		fmt.Printf("channel.Consume failed, error:%+v\n", err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//实现我们要处理的逻辑函数
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message, exit to  press CTRL + C")
	<-forever
}

//==================================>

//订阅模式创建rabbitmq
func NewRabbitMqPubSub(exchangeName string) *RabbitMq {
	//创建mq实例
	rabbitmq := NewRabbitMq("", exchangeName, "")
	return rabbitmq
}

//订阅模式下生产者
func (r *RabbitMq) PublishPub(message string) {
	//尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//是在127.0.0.1：15672后台的exchange中看到
		"fanout",
		//是否持久化
		false,
		//
		false,
		//如果internal设置为true，表示exchage不可以被client用来推送消息，仅用来进行exchange和exchage之间的绑定
		false,
		false,
		nil)

	r.failOnErr(err, "failed to declare an exchange")

	//2. 发送消息
	r.channel.Publish(
		//交换机
		r.Exchange, //简单模式为空
		//路由key
		"",
		//mandatory如果为true会根据exchange类型和routekey规则，
		//如果无法找到符合条件的队列就把消息返还给发送者
		false,
		//immediate如果为true，当exchage发送消息到队列后，
		//如果发现队列上没有绑定消费者就把消息返还给发送者
		false,
		//发送的信息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

//订阅模式下消费者
func (r *RabbitMq) RecivedSub() {
	//尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//是在127.0.0.1：15672后台的exchange中看到
		"fanout",
		//是否持久化
		false,
		false,
		//如果internal设置为true，表示exchage不可以被client用来推送消息，仅用来进行exchange和exchage之间的绑定
		false,
		false,
		nil)

	r.failOnErr(err, "failed to declare an exchange")

	//2 尝试创建队列
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)

	r.failOnErr(err, "failed to declare an queue")

	//3 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//在pub和sub模式下这里的key要为空
		"",
		r.Exchange,
		false,
		nil)

	//消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		fmt.Printf("channel.Consume failed, error:%+v\n", err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//实现我们要处理的逻辑函数
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message, exit to  press CTRL + C")
	<-forever

}

//=====================================
//routing路由模式

//创建rabbitmq实例
func NewRabbitMqRouting(exchange, routingKey string) *RabbitMq {
	//创建rabbitmq实例
	rabbitmq := NewRabbitMq("", exchange, routingKey)
	return rabbitmq
}

//路由模式下发送消息
func (r *RabbitMq) PublishRouting(message string) {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		//是否持久化
		false,
		//
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")

	//2. 发送消息
	r.channel.Publish(
		//交换机
		r.Exchange, //简单模式为空
		//路由key
		r.Key,
		false,
		false,
		//发送的信息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}


//路由模式下接收消息
func(r *RabbitMq)RecieveRouting(){
	//尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		//是否持久化
		false,
		//
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")


	//2 尝试创建队列
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)

	r.failOnErr(err, "failed to declare an queue")

	//3 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil)

	//消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		fmt.Printf("channel.Consume failed, error:%+v\n", err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//实现我们要处理的逻辑函数
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message, exit to  press CTRL + C")
	<-forever


}



//topic 话题模式
//=======================================>


func NewRabbitMqTopic(exchangeName, routingKey string)*RabbitMq{
	return NewRabbitMq("", exchangeName, routingKey)
}


//话题模式下发送消息
func (r *RabbitMq)PublishTopic(message string){
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		//是否持久化
		false,
		//
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")

	//2. 发送消息
	r.channel.Publish(
		//交换机
		r.Exchange, //简单模式为空
		//key
		r.Key,
		false,
		false,
		//发送的信息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}


/*
要注意key规则
其中"*"用于匹配一个单词， "#"用于匹配多个单词(可以是零)
匹配imooc.* 表示可以匹配 imooc.hello ,但是imooc.hello.one需要使用imooc.#才可以匹配到
*/
//话题模式接收消息
func(r *RabbitMq)RecivedTopic(){
	//尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		//是否持久化
		false,
		//
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")


	//2 尝试创建队列
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)

	r.failOnErr(err, "failed to declare an queue")

	//3 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil)

	//消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		fmt.Printf("channel.Consume failed, error:%+v\n", err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//实现我们要处理的逻辑函数
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message, exit to  press CTRL + C")
	<-forever
}
