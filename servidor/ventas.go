package main

import (
	"context"
	"fmt"
	"log"
	"net"

	configs "servidor/configs"
	msgq "servidor/msgqueues"

	pb "servidor/proto"

	"google.golang.org/grpc"

	amqp "github.com/rabbitmq/amqp091-go"
)

type server struct {
	pb.UnimplementedBookstoreServiceServer
	rabbitMQChannel *amqp.Channel
}

// funcion auxiliar para revisar errores
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (s *server) ProcessOrder(ctx context.Context, req *pb.BookstoreRequest) (*pb.BookstoreResponse, error) {

	///////////// Aqui va la logica del servidor

	// Simple order logic: Calculate the total price of all products
	totalPrice := float32(0)
	for _, product := range req.Products {
		totalPrice += product.Price * float32(product.Quantity)
	}

	//TENEMOS TODOS LOS ELEMENTOS. AHORA QUEREMOS:
	//1. INSERTAR LA ORDEN EN MONGODB. ALMACENAMOS LA ORDEN ID PARA EL RETORNO

	//ACCEDEMOS A LA COLECCION DE ORDERS
	db := configs.GetDatabase()
	collection := db.Collection("orders")
	collection2 := db.Collection("products")

	//DESERIALIZAMOS LA ORDEN A JSON STRUCT
	order := pb.ConvertRequestToJSON(req)

	// INSERTAMOS LOS DATOS DE LA ORDEN A LA COLECCION
	// Y CONVERTIMOS EL ID A STRING
	result, insertErr := collection.InsertOne(context.TODO(), order)
	FailOnError(insertErr, "error al insertar en mongodb")
	orderID := configs.ExtractObjectIdAsString(result)

	// PREPARAMOS LOS MENSAJES A SERVICIOS ANTES DE ALTERAR LAS CANTIDADES
	// DE STOCK, ESTO ES UN WORKAROUND BASTANTE FEO

	InventoryData := msgq.AddInventoryData(order.Products, orderID)
	InventoryMessage, err := pb.ConvertirStructToJSON(InventoryData)
	FailOnError(err, "error al convertir customer a json")

	NotificationData := msgq.AddNotificationData(orderID, order)
	NotificationMessage, err := pb.ConvertirStructToJSON(NotificationData)
	FailOnError(err, "error al convertir customer a json")

	// INSERTAMOS UN NUMERO ALTO DE STOCK A LA COLECCION PRODUCTS PARA SIMULAR STOCK
	products := order.Products

	fmt.Println("Products:", products)

	for i := range products {
		products[i].Quantity = 9999

		//	 INSERTAMOS EL PRODUCTO EN LA COLECCION
		_, insertErr2 := collection2.InsertOne(context.TODO(), products[i])

		if insertErr2 != nil {
			log.Panicf("%s: %s", "error al insertar en mongodb", insertErr2)
		}
	}

	//2. ENVIAR DATOS CLIENTE A SERVICIO DE DESPACHO CON ORDER ID.

	// Creamos el mensaje con los datos de orden y el id de orden
	DeliveryData := msgq.AddOrderID(order.Customer, orderID)
	DeliveryMessage, err := pb.ConvertirStructToJSON(DeliveryData)
	FailOnError(err, "error al convertir customer a json")

	// Definimos la routing key para el mensaje y enviamos
	routing_key := "colaDespacho"
	msgq.EnviarMensaje(routing_key, s.rabbitMQChannel, DeliveryMessage)

	//3. ENVIAR DATOS PRODUCTOS A SERVICIO DE INVENTARIO

	// Definimos la routing key para el inventory
	routing_key2 := "colaInventario"
	// Enviamos el mensaje. Notar que InventoryMessage lo definimos mas arriba por utilidad
	msgq.EnviarMensaje(routing_key2, s.rabbitMQChannel, InventoryMessage)

	//4. ENVIAR DATOS ORDER ID Y DATOS DE ORDEN A SERVICIO DE NOTIFICACIONES

	// Definimos la routing key para la notificacion
	routing_key3 := "colaNotificacion"
	msgq.EnviarMensaje(routing_key3, s.rabbitMQChannel, NotificationMessage)

	return &pb.BookstoreResponse{Message: orderID}, nil
}

func main() {

	//creamos la conexion con rabbitmq
	//esta conexion abstrae sockets, negociacion de protocolos y autenticacion

	connection, channel := msgq.CrearConexionRabbitMQ()

	defer channel.Close()
	defer connection.Close()

	fmt.Println("conectado a rabbitmq")

	//creamos la conexion a mongodb
	configs.EnvMongoURI()
	var connection_string = configs.GetConnection_String()
	err := configs.Connect(connection_string)
	FailOnError(err, "Error al conectar a mongodb")

	// Create a TCP listener on port 50051 (Puedes cambiar el puerto si quieres)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server
	s := grpc.NewServer()

	// Register the server with the gRPC server
	pb.RegisterBookstoreServiceServer(
		s,
		&server{rabbitMQChannel: channel},
	)

	log.Println("Server is listening on port 50051") // Print server is running...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
