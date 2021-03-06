package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	// port := os.Getenv("PORT")

	// if port == "" {
	// 	log.Fatal("$PORT is not set")
	// }

	// dns := os.Getenv("DATABASE_URL")

	// if dns == "" {
	// 	log.Fatal("$DATABASE_URL is not set")
	// }

	db, err = gorm.Open("mysql", "vsb6ig2gv6zc066m:szectndz2gcu0r57@tcp(eyw6324oty5fsovx.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/f1uje8ak4fcxhh3l")

	if err != nil {
		log.Println("failed", err)
	} else {
		log.Println("Success")
	}

	db.AutoMigrate(&Product{})

	handleRequest()
}

func handleRequest() {
	// log.Print("server start at http://127.0.0.1:7777")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Success create data"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}
	db.Find(&products)
	res := Result{Code: 200, Data: products, Message: "Success get data"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	var product Product
	db.First(&product, productID)
	res := Result{Code: 200, Data: product, Message: "Success get data"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	payloads, _ := ioutil.ReadAll(r.Body)

	var productUpdate Product
	json.Unmarshal(payloads, &productUpdate)

	var product Product
	db.Find(&product, productID)
	db.Model(&product).Updates(productUpdate)

	res := Result{Code: 200, Data: product, Message: "Success update data"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)
	db.Delete(&product)

	res := Result{Code: 200, Message: "Success delete data"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
