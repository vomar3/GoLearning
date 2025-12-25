package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"unicode"
)

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type ProductUpdate struct {
	Name          string `json:"name"`
	QuantityDelta int    `json:"quantityDelta"`
}

var (
	mtx      sync.Mutex
	products = make(map[string]Product)
)

// Simple check for name (letters only)
func ValidateString(str string) bool {
	if str == "" {
		return false
	}

	for _, letter := range str {
		if !unicode.IsLetter(letter) {
			return false
		}
	}

	return true
}

// Validate all params
func (p Product) validate() bool {
	if !ValidateString(p.Name) {
		return false
	}

	epsilon64 := math.Nextafter(1.0, 2.0) - 1.0

	if p.Price <= epsilon64 || p.Quantity <= 0 {
		return false
	}

	return true
}

// Validate name only
func (p ProductUpdate) validate() bool {
	if !ValidateString(p.Name) {
		return false
	}

	return true
}

// Add product in products map
func HandleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // Need post method
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Invalid read JSON\n"))
		return
	}

	if !p.validate() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid product data\n"))
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if _, ok := products[p.Name]; ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Data with that name already exists\n"))
		return
	}

	products[p.Name] = p
}

// Write info about product
func HandleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Error query parameter\n"))
		return
	}

	mtx.Lock()
	value, ok := products[name]
	mtx.Unlock()

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Name %s isn't exist\n", name)
		_, _ = w.Write([]byte(msg))
		return
	}

	ProductJSON, err := json.Marshal(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error with convert to JSON\n"))
		return
	}

	_, _ = w.Write(ProductJSON)
}

// Updating the quantity of an item by name
func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var update ProductUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error with JSON\n"))
		return
	}

	if !update.validate() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid product data\n"))
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	p, ok := products[update.Name]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Data with name %s isn't exitst\n", update.Name)
		_, _ = w.Write([]byte(msg))
		return
	} else {
		p.Quantity -= update.QuantityDelta
		products[update.Name] = p
	}

	_, _ = w.Write([]byte("OK"))
}

// Write all products
func HandleAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	AllProducts, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error with convert to JSON\n"))
		return
	}

	_, _ = w.Write(AllProducts)
}

func main() {
	fmt.Println("The server is running")

	http.HandleFunc("/product/add", HandleAdd)
	http.HandleFunc("/product/info", HandleInfo)
	http.HandleFunc("/product/update", HandleUpdate)
	http.HandleFunc("/product/all", HandleAll)

	if err := http.ListenAndServe(":9091", nil); err != nil {
		fmt.Println("Error with HTTP server: ", err)
	}

	fmt.Println("Server closed")
}
