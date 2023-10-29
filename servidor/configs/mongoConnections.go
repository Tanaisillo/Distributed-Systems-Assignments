package configs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// Connect establishes a connection to the MongoDB database
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Connect(uri string) error {

	// Connection to mongo
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	failOnError(err, "Error al conectar a mongodb")

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	failOnError(err, "Error al hacer ping a mongodb")

	// Gets the db name
	dbName := ExtractDatabaseName(uri)

	// Set the db variable to the extracted database
	db = client.Database(dbName)

	return error(nil)
}

// CreateCollection creates a new collection in the database
// This is not used (collections are created when declared)
func Create_Collection(name string) {
	createCollectionOptions := options.CreateCollection()
	if err := db.CreateCollection(context.TODO(), name, createCollectionOptions); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created collection:", name)
}

// GetDatabase returns the MongoDB database instance
func GetDatabase() *mongo.Database {
	return db
}

// ExtractDatabaseName extracts the database name from the connection string
func ExtractDatabaseName(uri string) string {
	//split the connection string by /
	split := strings.Split(uri, "/")
	//get the last element of the slice
	databaseName := split[len(split)-1]
	return databaseName
}

// Converts a primitive object from mongodb to a string type
func ExtractObjectIdAsString(result *mongo.InsertOneResult) string {
	orderID := result.InsertedID.(primitive.ObjectID).Hex()
	return orderID
}

// Converts the string type into an object
func ConvertStringToObjectId(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Panicf("%s: %s", "error al convertir string a objectid", err)
	}
	return objID
}
