package main

import (
	"fmt"
	"rabbitmq/RabbitMq"
)

func main(){
	rabbitmq := RabbitMq.NewRabbitMqSimple("imoocSimple")

	rabbitmq.PublishSimple("hello imooc")
	fmt.Printf("发送成功!\n")
}
