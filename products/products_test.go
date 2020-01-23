package products

import (
	"fmt"
	"testing"
)

func TestLoadProducts(t *testing.T) {
	// the products map should be loaded at testing
	if products == nil {
		t.Fail()
	}
}

func ExampleCheckProducts() {
	for code := range products {
		p, _ := SearchProduct(code)
		fmt.Println(code, p)
	}
	// Unordered output:
	// MUG &{MUG Lana Cofee Mug 7.50}
	// PEN &{PEN Lana Pen 5.00}
	// TSHIRT &{TSHIRT Lana T-Shirt 20.00}
}
