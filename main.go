package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Order struct {
	OrderId      uint      `json:"orderId" gorm:"primary_key"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `json:"orderedAt"`
	Items        []Item    `json:"items" gorm:"foreignkey:OrderID"`
}

type Item struct {
	ItemId      uint   `json:"lineItemId" gorm:"primary_key"`
	IitemCode   string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderId     uint   `json:"-"`
}

var db *gorm.DB

func DBInit() {
	var err error
	db, err = gorm.Open("mysql", "root:administrator@tcp(localhost:3306)/orders_by?parseTime=True")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect to database")
	}
	db.AutoMigrate(&Order{}, &Item{})
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orders []Order
	db.Preload("Items").Find(&orders)
	json.NewEncoder(w).Encode(orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputOrderID := params["orderId"]
	var order Order
	db.Preload("Items").First(&order, inputOrderID)
	json.NewEncoder(w).Encode(order)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	db.Create(&order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var orderUpdate Order
	json.NewDecoder(r.Body).Decode(&orderUpdate)
	db.Save(&orderUpdate)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderUpdate)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inputOrderID := params["orderId"]
	id64, _ := strconv.ParseUint(inputOrderID, 10, 64)
	deleteID := uint(id64)
	db.Where("order_id = ?", deleteID).Delete(&Item{})
	db.Where("order_id = ?", deleteID).Delete(&Order{})
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	DBInit()
	router := mux.NewRouter()
	router.HandleFunc("/orders", CreateOrder).Methods("POST")
	router.HandleFunc("/orders", GetOrders).Methods("GET")
	router.HandleFunc("/orders/{orderId}", GetOrder).Methods("GET")
	router.HandleFunc("/orders/{orderId}", UpdateOrder).Methods("PUT")
	router.HandleFunc("/orders/{orderId}", DeleteOrder).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}
