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
	ItemId int `json:"itemId"`
	ItemCode string `json:"itemCode"`
	Description string `json:"description"`
	Quantity int `json:"quantity"`
	OrderId int `json:"orderId"`
}

type Order struct {
	OrderId int `json:"orderId"`
	CustomerName string `json:"customerName"`
	OrderedAt time.Time `json:"orderedAt"`
	Item Item `json:"item"`
}

type OrderItems struct {
	OrderId int `json:"orderId"`
	CustomerName string `json:"customerName"`
	OrderedAt time.Time `json:"orderedAt"`
	Items []Item `json:"items"`
}

type OrderService struct {
	DB *sql.DB
}

func (orderService *OrderService) GetOrders(res http.ResponseWriter, req *http.Request) {
	var response helper.Response
	var orders []Order
	var ordersItems []OrderItems
	var orderItems OrderItems

	rows, err := orderService.DB.Query("SELECT * FROM orders AS `order` JOIN `items` AS item ON order.order_id = item.order_id")
	if err != nil {
		response.Status = false
		response.Code = 500
		response.Data = err.Error()
		json.NewEncoder(res).Encode(response)
		return
	}

	for rows.Next() {
		var order = Order{}
 		err := rows.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.Item.ItemId, &order.Item.ItemCode, &order.Item.Description, &order.Item.Quantity, &order.Item.OrderId)
		 if err != nil {
			response.Status = false
			response.Code = 500
			response.Data = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		 }
		 orders = append(orders, order)
	}
	for _, order := range orders {
		if orderItems.OrderId != order.OrderId{
			orderItems = OrderItems{}
			orderItems = OrderItems{order.OrderId, order.CustomerName, order.OrderedAt, append(orderItems.Items, order.Item)}
			ordersItems = append(ordersItems, orderItems)
		}else{	
			ordersItems[len(ordersItems)-1] = OrderItems{order.OrderId, order.CustomerName, order.OrderedAt, append(orderItems.Items, order.Item)}
		}
		
	}
	response.Status = true
	response.Code = 200
	response.Data = ordersItems
	json.NewEncoder(res).Encode(response)
}

func (orderService *OrderService) StoreOrder(res http.ResponseWriter, req *http.Request){
	var response helper.Response
	var order OrderItems
	json.NewDecoder(req.Body).Decode(&order)

	stmt, err := orderService.DB.Prepare("INSERT INTO orders (customer_name, ordered_at) VALUES(?,?) ON DUPLICATE KEY UPDATE order_id=LAST_INSERT_ID(order_id)")
	if err != nil {
		log.Fatal(err.Error())
	}
	result, err := stmt.Exec(order.CustomerName, order.OrderedAt)
	if err != nil {
		log.Fatal(err.Error())
	}

	orderId, err := result.LastInsertId()
	for _, value := range order.Items {
		stmt, err := orderService.DB.Prepare("INSERT INTO items (item_code, description, quantity, order_id) VALUES(?,?,?,?)")
		if err != nil {
			log.Fatal(err.Error())
		}
		_,err = stmt.Exec(value.ItemCode, value.Description, value.Quantity, orderId)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	response.Status = true
	response.Code = 201
	json.NewEncoder(res).Encode(response)
}

func (orderService *OrderService) GetOrder(res http.ResponseWriter, req *http.Request){
	var response helper.Response
	var orders []Order
	var orderItems OrderItems
	params := mux.Vars(req)
	rows,err := orderService.DB.Query("SELECT * FROM orders AS `order` INNER JOIN items AS `item` ON order.order_id = item.order_id WHERE order.order_id = ?", params["orderId"])
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.Item.ItemId, &order.Item.ItemCode, &order.Item.Description, &order.Item.Quantity, &order.Item.OrderId)
		if err != nil {
			response.Status = false
			response.Code = 500
			return
		 }
		 orders = append(orders, order)
	}

	for _, order := range orders {
		orderItems = OrderItems{order.OrderId, order.CustomerName, order.OrderedAt, append(orderItems.Items, order.Item)}
	}
	
	response.Status = true
	response.Code = 200
	response.Data = orderItems
	json.NewEncoder(res).Encode(response)
}

func (orderService *OrderService) UpdateOrder(res http.ResponseWriter, req *http.Request){
	var response helper.Response
	var orderItems OrderItems

	json.NewDecoder(req.Body).Decode(&orderItems)
	stmt, err := orderService.DB.Prepare("UPDATE orders SET customer_name = ?, ordered_at = ? WHERE order_id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(orderItems.CustomerName, orderItems.OrderedAt, orderItems.OrderId)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, value := range orderItems.Items {
		stmt, err := orderService.DB.Prepare("UPDATE items SET item_code = ?, description = ?, quantity = ? WHERE item_id = ?")
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = stmt.Exec(value.ItemCode, value.Description, value.Quantity, value.ItemId)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	response.Status = true
	response.Code = 201
	json.NewEncoder(res).Encode(response)
}

func (orderService *OrderService) DeleteOrder(res http.ResponseWriter, req *http.Request){
	var response helper.Response
	params := mux.Vars(req)

	_, err := orderService.DB.Exec("DELETE FROM orders WHERE order_id = ?", params["orderId"])	
	if err != nil {
		response.Status = false
		response.Code = 500
		json.NewEncoder(res).Encode(response)
		return
	}
	response.Status = true
	response.Code = 201
	json.NewEncoder(res).Encode(response)
}