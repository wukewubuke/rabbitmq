package main

import (
	"rabbitmq/RabbitMq"
	"strconv"
	"time"
)

func main() {

	imoocOne := RabbitMq.NewRabbitMqRouting("exImooc", "imooc_one")
	imoocTwo := RabbitMq.NewRabbitMqRouting("exImooc", "imooc_two")
	for i := 0; i < 10; i++ {
		imoocOne.PublishRouting("imooc one " + strconv.Itoa(i))
		imoocTwo.PublishRouting("imooc two " + strconv.Itoa(i))
		time.Sleep(10 * time.Millisecond)
	}
}
