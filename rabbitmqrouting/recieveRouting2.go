package main

import (
	"rabbitmq/RabbitMq"
)

func main(){
	rabbitmq := RabbitMq.NewRabbitMqRouting("exImooc", "imooc_two")
	rabbitmq.RecieveRouting()
}
