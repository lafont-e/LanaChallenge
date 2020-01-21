/* Library to manage the products at the store  */

package products

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// productsFile name of the file with the products to manage
const productsFile = "./products.json"

// Product struct that defines a product
type Product struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

var errNotFound = errors.New("Product Not found")

// As there are no DB, we use a map to store the items available on the store
// There is a small interface to the items storage just in case a real DB should appear later
// Not using a sync.Map as this map is read only.
var products = make(map[string]*Product)

func init() {
	loadProducts()
}

// NewProduct returns a new product
func NewProduct(code, name string, price float32) *Product {
	return &Product{code, name, price}
}

// GetProduct Retrieves the product associated with a given code
func GetProduct(code string) (*Product, error) {
	if p, ok := products[code]; ok {
		return p, nil
	}
	return nil, errNotFound
}

// SearchISearchProducttem Search the storage area for the coded product
func SearchProduct(code string) (p *Product, err error) {
	p, ok := products[code]
	if !ok {
		log.Printf("Error, Search %s :%s\n", code, errNotFound)
		return nil, errNotFound
	}
	return
}

// GetPrice Return the price of a product
func (p *Product) GetPrice() float32 {
	return p.Price
}

// GetName Return the Name of a product
func (p *Product) GetName() string {
	return p.Name
}

// GetCode Return the Code of a product
func (p *Product) GetCode() string {
	return p.Code
}

// loadProducts function that reads products from a json file
func loadProducts() error {
	// read file
	data, err := ioutil.ReadFile(productsFile)
	if err != nil {
		log.Println("Error, Not Loading Products :", err)
		pwd, _ := os.Getwd()
		log.Println("Error cont:", pwd)
		return err
	}

	var ps []Product // as the Json file is stored as an array not a hash, to avoid Code duplication issues

	// unmarshall it
	err = json.Unmarshal(data, &ps)

	if err != nil {
		log.Println("Error, not able to read products file :", err)
		return err
	}

	for ix, p := range ps {
		products[p.Code] = &ps[ix]
	}

	return nil
}
