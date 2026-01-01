package rest

import (
	"fmt"
	"net/http"
	"strings"
)

// GET http://localhost:3000/carts/1/price
// GET http://localhost:3000/carts/1
// DELETE http://localhost:3000/carts/1/items/1
// POST http://localhost:3000/carts/1/items -d '{
//  "product": "Shoes",
//  "price": 2500.50
//}'
//POST http://localhost:3000/carts -d '{}'

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/carts", Handler)
	return router
}

func Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	switch r.Method {
	case http.MethodGet:
		fmt.Println(params)
		switch len(params) {
		case 1:
			// Делаем запрос в бд
		case 2:
			// Тоже делаем запрос в бд
		}
	case http.MethodPost:
		fmt.Println("POST")
	case http.MethodDelete:
		fmt.Println("DELETE")
	}
}
