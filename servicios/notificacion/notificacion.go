//notificacion.go
//debe recibir mensajes de una cola, los datos de la orden con ID
//y envia por la api

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notificacion/models"
	msgq "notificacion/msgqueues"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	//creamos la conexion con rabbitmq
	//esta conexion abstrae sockets, negociacion de protocolos y autenticacion

	connection, channel := msgq.CrearConexionRabbitMQ()

	defer channel.Close()
	defer connection.Close()

	
	fmt.Println("Conectado a rabbitmq")
	
	// default exchange
	queue := msgq.DeclararCola(channel, "colaNotificacion") // match con el nombre

	msgs, err := channel.Consume(
		queue.Name, //nombre de la cola
		"",         //consumer
		true,       //auto-ack
		false,      //exclusive
		false,      //no-local
		false,      //no-wait
		nil,        //args
	)
	FailOnError(err, "Error al consumir mensajes")

	//forever es un canal que bloquea la ejecucion de main
	//y permite que el programa se mantenga a la escucha de mensajes

	var forever chan struct{}

	//iniciamos una rutina para consumir los mensajes

	go func() {
		for d := range msgs {
			//los mensajes llegan como un arreglo de bytes
			fmt.Printf("Received a message: %s", d.Body)

			//los convertiremos a la estructura correspondiente
			var mensaje models.MensajeNotificacion
			err := json.Unmarshal(d.Body, &mensaje)
			FailOnError(err, "Error al convertir mensaje")

			// Defines the payload for the mail
			payload := models.Payload{
				OrderID:  mensaje.OrderID,
				GroupID:  "T9q!6B3j#5", // contraseña de la primera maquina
				Products: mensaje.Order.Products,
				Customer: mensaje.Order.Customer,
			}
			// Convert payload to JSON
			jsonData, err := json.Marshal(payload)
			FailOnError(err, "Error al convertir payload a JSON")
			fmt.Println("Payload:", string(jsonData))

			// Make the HTTP request
			url := "https://sjwc0tz9e4.execute-api.us-east-2.amazonaws.com/Prod" // url
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			FailOnError(err, "Error creando la solicitud HTTP")

			// Set request headers
			req.Header.Set("Content-Type", "application/json")

			// Make the request
			client := &http.Client{}
			resp, err := client.Do(req)
			FailOnError(err, "Error haciendo la solicitud HTTP")

			// Print the response status and body
			fmt.Println("Status Code:", resp.Status)
			// Read and print the response body
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(resp.Body)
			fmt.Println("Response Body:", buf.String())
			resp.Body.Close()
		}
	}()

	<-forever
}
