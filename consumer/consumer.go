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

	// Declara a mesma fila usada pelo produtor
	q, err := ch.QueueDeclare(
		"task_queue", // nome da fila
		true,         // durável
		false,        // deletar quando não usada
		false,        // exclusiva
		false,        // sem esperar
		nil,          // argumentos adicionais
	)
	failOnError(err, "Failed to declare a queue")

	// Configura o consumidor
	msgs, err := ch.Consume(
		q.Name, // nome da fila
		"",     // nome do consumidor
		true,   // auto-ack (confirmação automática)
		false,  // exclusivo
		false,  // sem espera
		false,  // no local
		nil,    // argumentos adicionais
	)
	failOnError(err, "Failed to register a consumer")

	// Processa as mensagens recebidas
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] Recebido %s", d.Body)
			time.Sleep(2 * time.Second) // Simula tempo de processamento
			log.Printf(" [x] Processado %s", d.Body)
		}
	}()

	log.Printf(" [*] Aguardando mensagens. Para sair pressione CTRL+C")
	<-forever
}
