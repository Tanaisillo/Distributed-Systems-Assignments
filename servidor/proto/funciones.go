package proto

import (
	"encoding/json"
	"servidor/models"

	protojson "google.golang.org/protobuf/encoding/protojson"
)

// Convertir la payload de gRPC a un modelo json
func ConvertRequestToJSON(req *BookstoreRequest) models.Bookstore_Order {
	var order models.Bookstore_Order

	jsonBytes, _ := protojson.Marshal(req)
	_ = json.Unmarshal(jsonBytes, &order)
	return order
}

func ConvertirStructToJSON(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	return jsonData, err
}
