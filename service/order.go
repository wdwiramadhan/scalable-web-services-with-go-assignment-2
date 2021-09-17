package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/wdwiramadhan/scalable-web-services-with-go-assignment-2/helper"
)

type Item struct {
	Item_id int `json:"item_id"`
	Item_code string `json:"item_code"`
	Description string `json:"description"`
	Quantity int `json:"quantity"`
	Order_id int `json:"order_id"`
}

type Order struct {
	Order_id int `json:"order_id"`
	Customer_name string `json:"customer_name"`
	Ordered_at time.Time `json:"ordered_at"`
	Items Item `json:"items"`
}

type OrderService struct {
	DB *sql.DB
}

func (orderService *OrderService) GetOrders(res http.ResponseWriter, req *http.Request) {
	var response helper.Response
	var orders []Order

	rows, err := orderService.DB.Query("SELECT * FROM orders AS `order` INNER JOIN `items` AS item ON order.order_id = item.order_id")
	if err != nil {
		response.Status = false
		response.Code = 500
		json.NewEncoder(res).Encode(response)
		return
	}

	for rows.Next() {
		var order = Order{}
 		err := rows.Scan(&order.Order_id, &order.Customer_name, &order.Ordered_at, &order.Items.Item_id, &order.Items.Item_code, &order.Items.Description, &order.Items.Quantity, &order.Items.Order_id)
		 if err != nil {
			response.Status = false
			response.Code = 500
			response.Data = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		 }
		 orders = append(orders, order)
	}
	response.Status = true
	response.Code = 200
	response.Data = orders
	json.NewEncoder(res).Encode(response)
}

func (orderService *OrderService) StoreOrder(res http.ResponseWriter, req *http.Request){
	var response helper.Response
	var order Order
	json.NewDecoder(req.Body).Decode(&order)
	stmt, err := orderService.DB.Prepare("INSERT INTO orders (customer_name, ordered_at) VALUES(?,?)")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = stmt.Exec(order.Customer_name, order.Ordered_at)
	if err != nil {
		log.Fatal(err.Error())
	}
	response.Status = true
	response.Code = 201
	json.NewEncoder(res).Encode(response)
}

func (orderService *OrderService) GetOrder(res http.ResponseWriter, req *http.Request){
	var response helper.Response
	var order Order
	params := mux.Vars(req)
	rows,err := orderService.DB.Query("SELECT * FROM orders WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		err := rows.Scan(&order.Order_id, &order.Customer_name, &order.Ordered_at)
		if err != nil {
			response.Status = false
			response.Code = 500
			return
		 }
	}
	
	response.Status = true
	response.Code = 200
	response.Data = order
	json.NewEncoder(res).Encode(response)
}