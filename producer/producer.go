package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Conecta ao servidor RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Cria um canal de comunicação
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declara uma fila
	q, err := ch.QueueDeclare(
		"task_queue", // nome da fila
		true,         // durável
		false,        // deletar quando não usada
		false,        // exclusiva
		false,        // sem esperar
		nil,          // argumentos adicionais
	)
	failOnError(err, "Failed to declare a queue")

	// Envia várias mensagens
	for i := 1; i <= 5; i++ {
		body := "Mensagem " + time.Now().Format("15:04:05")
		err = ch.Publish(
			"",     // exchange
			q.Name, // chave de roteamento
			false,  // mensagem obrigatória
			false,  // imediata
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // faz a mensagem persistir no RabbitMQ
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Enviado %s", body)
		time.Sleep(1 * time.Second) // Aguarda um segundo entre as mensagens
	}
}
