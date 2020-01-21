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
	// MUG &{MUG Lana Cofee Mug 7}
	// PEN &{PEN Lana Pen 5}
	// TSHIRT &{TSHIRT Lana T-Shirt 20}
}

func TestCode(t *testing.T) {
	a := []int{0, 1, 2, 3, 5, 6, 7}
	var b []int
	b = append(b, a[4:]...)
	fmt.Println(b)
	a = append(a[:4], 4)
	fmt.Println(a, b)
	a = append(a, b...)
	fmt.Println(a)
	for ix, v := range a {
		if ix%2 == 1 {
			fmt.Println(ix, v)
		}
	}
}
