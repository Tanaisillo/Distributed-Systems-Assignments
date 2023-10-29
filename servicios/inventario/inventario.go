//inventario.go
//debe recibir mensajes de una cola, un conjunto de productos
//por cada producto resta la cantidad al stock en mongodb (coleccion productos)

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"inventario/configs"
	"inventario/models"
	msgq "inventario/msgqueues"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	//creamos la conexion a mongodb
	configs.EnvMongoURI()
	var connection_string = configs.GetConnection_String()
	err := configs.Connect(connection_string)
	FailOnError(err, "Error al conectar a mongodb")

	db := configs.GetDatabase()
	//collection := db.Collection("orders")
	collection2 := db.Collection("products")
	//fmt.Println(collection.Name())

	//creamos la conexion con rabbitmq
	//esta conexion abstrae sockets, negociacion de protocolos y autenticacion

	connection, channel := msgq.CrearConexionRabbitMQ()

	defer channel.Close()
	defer connection.Close()

	fmt.Println("Conectado a rabbitmq")

	// default exchange
	queue := msgq.DeclararCola(channel, "colaInventario") // match con el nombre de la cola

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
			var mensaje models.MensajeInventario
			err := json.Unmarshal(d.Body, &mensaje)
			if err != nil {
				log.Printf("Error al convertir el mensaje: %s", err)
				continue
			}

			// Goes through each product
			for _, productToUpdate := range mensaje.Products {

				fmt.Printf("Product to update: %s", productToUpdate.Title) //prints each product for debugging
				// Subtract the quantity from the stock in the MongoDB collection
				update := bson.D{{Key: "$inc", Value: bson.D{{Key: "quantity", Value: -productToUpdate.Quantity}}}}

				filter := bson.D{{Key: "title", Value: productToUpdate.Title}} // Assuming "title" is a unique identifier

				// Changes the amount of stock in the collection products
				_, err = collection2.UpdateOne(context.Background(), filter, update)
				if err != nil {
					log.Printf("Error updating document: %v", err)
				}
			}
		}
	}()

	<-forever
}
