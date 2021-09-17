package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/wdwiramadhan/scalable-web-services-with-go-assignment-2/config"
	"github.com/wdwiramadhan/scalable-web-services-with-go-assignment-2/helper"
	"github.com/wdwiramadhan/scalable-web-services-with-go-assignment-2/service"
)


func main(){
	var err error
	config.DB, err = sql.Open("mysql", config.DBUrl(config.BuildDBConfig()))
	if err != nil {
		fmt.Println("error", err)
	}
	defer config.DB.Close()

	orderService := service.OrderService{ DB: config.DB}
	router := mux.NewRouter()
	router.HandleFunc("/", func (res http.ResponseWriter, req *http.Request){
		var response helper.Response
		response.Status = true
		response.Code = 200
		json.NewEncoder(res).Encode(response)
	}).Methods("get")
	router.HandleFunc("/orders", orderService.GetOrders).Methods("get")
	router.HandleFunc("/orders", orderService.StoreOrder).Methods("post")
	router.HandleFunc("/orders/{id}", orderService.GetOrder).Methods("get")

	fmt.Println("server running on port 8080")
	http.ListenAndServe(":8080", router)
}