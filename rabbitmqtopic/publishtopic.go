package main

import (
	"rabbitmq/RabbitMq"
	"strconv"
	"time"
)

func main() {

	imoocOne := RabbitMq.NewRabbitMqTopic("exImoocTopic", "imooc.topic.one")
	imoocTwo := RabbitMq.NewRabbitMqRouting("exImoocTopic", "imooc.topic.two")
	for i := 0; i < 10; i++ {
		imoocOne.PublishTopic("imooc topic one " + strconv.Itoa(i))
		imoocTwo.PublishTopic("imooc topic two " + strconv.Itoa(i))
		time.Sleep(10 * time.Millisecond)
	}
}
