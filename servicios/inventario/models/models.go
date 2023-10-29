package models

type Bookstore_Order struct {
	Products []Product `json:"products"`
	Customer Customer  `json:"customer"`
}

type Product struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Genre       string  `json:"genre"`
	Pages       int     `json:"pages"`
	Publication string  `json:"publication"`
	Quantity    int     `json:"quantity"`
	Price       float32 `json:"price"`
}

type Customer struct {
	Name     string   `json:""`
	Lastname string   `json:""`
	Email    string   `json:""`
	Location Location `json:""`
	Phone    string   `json:""`
}

type Location struct {
	Address1   string `json:""`
	Address2   string `json:""`
	City       string `json:""`
	State      string `json:""`
	PostalCode string `json:""`
	Country    string `json:""`
}

type MensajeInventario struct {
	Products []Product `json:"products"`
}
