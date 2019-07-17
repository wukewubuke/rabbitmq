package main

import (
	"rabbitmq/RabbitMq"
	"strconv"
	"time"
)

func main(){
	rabbitmq := RabbitMq.NewRabbitMqPubSub("newProduct")
	for i:=0;i<10;i++ {
		rabbitmq.PublishPub("订阅发布模式生产消息 imooc" + strconv.Itoa(i))
		time.Sleep( 10 * time.Millisecond)
	}
}
