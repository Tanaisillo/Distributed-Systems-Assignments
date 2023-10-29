package configs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// Connect establishes a connection to the MongoDB database

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Connect(uri string) error {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	FailOnError(err, "Error al conectar a mongodb")

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	FailOnError(err, "Error al hacer ping a mongodb")

	dbName := ExtractDatabaseName(uri)

	//set the db variable to the extracted database
	db = client.Database(dbName)

	return nil
}

// CreateCollection creates a new collection in the database
// This is not used

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
