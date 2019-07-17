package main

import "rabbitmq/RabbitMq"

func main(){
	rabbitmq := RabbitMq.NewRabbitMqPubSub("newProduct")
	rabbitmq.RecivedSub()
}
