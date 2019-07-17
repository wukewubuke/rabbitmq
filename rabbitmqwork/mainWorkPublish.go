package main

import (
	"rabbitmq/RabbitMq"
	"strconv"
	"time"
)

func main(){
	rabbitmq := RabbitMq.NewRabbitMqSimple("imoocSimple")
	for i:=0;i<10;i++ {
		rabbitmq.PublishSimple("one hello imooc" + strconv.Itoa(i))
		time.Sleep( 100 * time.Millisecond)
	}
}

