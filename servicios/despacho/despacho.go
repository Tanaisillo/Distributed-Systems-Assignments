//despacho.go
//debe recibir mensajes de una cola
//e insertar datos de despacho en la orden correspondiente en mongodb

package main

import (
	"context"
	"despacho/configs"
	"despacho/models"
	msgq "despacho/msgqueues"
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"
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
	collection := db.Collection("orders")
	ctx := context.TODO()
	fmt.Println(collection.Name())

	//creamos la conexion con rabbitmq
	//esta conexion abstrae sockets, negociacion de protocolos y autenticacion

	connection, channel := msgq.CrearConexionRabbitMQ()

	defer channel.Close()
	defer connection.Close()

	fmt.Println("conectado a rabbitmq")

	queue := msgq.DeclararCola(channel, "colaDespacho")

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
			var mensaje models.MensajeDespacho
			err := json.Unmarshal(d.Body, &mensaje)
			FailOnError(err, "Error al convertir mensaje")

			// Extract relevant information
			orderID := mensaje.OrderID
			customer := mensaje.Customer

			deliveryInfo := models.Delivery{}

			// Fill in the information
			deliveryInfo.ShippingAddress.Name = customer.Name
			deliveryInfo.ShippingAddress.Lastname = customer.Lastname
			deliveryInfo.ShippingAddress.Address1 = customer.Location.Address1
			deliveryInfo.ShippingAddress.Address2 = customer.Location.Address2
			deliveryInfo.ShippingAddress.City = customer.Location.City
			deliveryInfo.ShippingAddress.State = customer.Location.State
			deliveryInfo.ShippingAddress.PostalCode = customer.Location.PostalCode
			deliveryInfo.ShippingAddress.Country = customer.Location.Country
			deliveryInfo.ShippingAddress.Phone = customer.Phone
			deliveryInfo.ShippingMethod = "USPS" //estas van harcodeadas como en el ejemplo
			deliveryInfo.TrackingNumber = "12345678901234567890"

			fmt.Println(deliveryInfo)
			// Convert the deliveryInfo to a BSON document
			// deliveryBSON, err := bson.Marshal(deliveryInfo)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// Update the existing document with the new deliveries information
			objID := configs.ConvertStringToObjectId(orderID)
			filter := bson.M{"_id": objID}
			update := bson.M{"$set": bson.M{"deliveries": deliveryInfo}}
			_, err = collection.UpdateOne(ctx, filter, update)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("DeliveryInfo appended to the orders collection")
		}
	}()

	<-forever
}
