package main

import "rabbitmq/RabbitMq"

func main(){
	rabbitmq := RabbitMq.NewRabbitMqSimple("imoocSimple")

	rabbitmq.ConsumeSimple()
}
