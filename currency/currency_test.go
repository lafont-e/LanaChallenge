package currency

import (
	"fmt"
)

func ExamplePrint() {
	var a Currency = 12305
	fmt.Printf("%6.2f", a)
	// output:
	// 123.05
}
