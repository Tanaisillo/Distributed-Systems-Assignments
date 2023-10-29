package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// connection string to the MongoDB database
var connection_string string

// GetConnection_String returns the connection string to the MongoDB database from env file

func GetConnection_String() string {
	return connection_string
}

// SetConnection_String sets the connection string to the variable
func SetConnection_String(connection string) {
	connection_string = connection
}

// EnvMongoURI loads the connection string from the env file
func EnvMongoURI() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	SetConnection_String(os.Getenv("MONGO_URI"))
}
