package msgq

import (
	"context"
	"log"
	"servidor/models"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// CrearConexionRabbitMQ() (*amqp.Connection, *amqp.Channel)
// No recibe argumentos, retorna objetos de RabbitMQ de conexion y canal.
func CrearConexionRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Error al conectar a rabbitmq")
	canal, err := conn.Channel()
	FailOnError(err, "Error al abrir canal")

	return conn, canal
}

// DeclararCola(ch *amqp.Channel) (q amqp.Queue)
// Recibe un canal de RabbitMQ y devuelve una Queue de RabbitMQ.
func DeclararCola(ch *amqp.Channel) (q amqp.Queue) {
	q, err := ch.QueueDeclare(
		"ColaSalida", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Error al declarar queue")
	return q
}

// EnviarMensaje(routing_key string, ch *amqp.Channel, body []byte)
// Recibe un string con el nombre de la cola de destino, un canal de RabbitMQ y un arreglo de bytes
// Envía el arreglo de bytes por la cola. No tiene retorno.
func EnviarMensaje(routing_key string, ch *amqp.Channel, body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := ch.PublishWithContext(ctx,
		"",          //default exchange
		routing_key, //define implicitamente la cola
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	FailOnError(err, "Error al enviar mensaje")
}

// AddOrderID(customer models.Customer, orderID string) models.MensajeDespacho
// Combina datos de cliente y ordenID en un mensaje para el servicio de despacho
func AddOrderID(customer models.Customer, orderID string) models.MensajeDespacho {
	var mensaje models.MensajeDespacho
	mensaje.OrderID = orderID
	mensaje.Customer = customer
	return mensaje
}
	
// AddInventoryData(products []models.Product, orderID string) models.MensajeInventario
// Combina datos de productos y ordenID en un mensaje para el servicio de inventario
func AddInventoryData(products []models.Product, orderID string) models.MensajeInventario {
	var mensaje models.MensajeInventario
	mensaje.Products = products
	return mensaje
}

// AddNotificationData(orderID string, order models.Bookstore_Order) models.MensajeNotificacion
// Combina datos de la orden y ordenID en un mensaje para el servicio de notificacion
func AddNotificationData(orderID string, order models.Bookstore_Order) models.MensajeNotificacion {
	var mensaje models.MensajeNotificacion
	mensaje.OrderID = orderID
	mensaje.Order = order
	return mensaje
}
