package main

import "rabbitmq/RabbitMq"

func main(){

	rabbitmq := RabbitMq.NewRabbitMqRouting("exImoocTopic", "#")
	rabbitmq.RecivedTopic()
}

