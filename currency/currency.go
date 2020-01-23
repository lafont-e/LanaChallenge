// currency is an integer that represents a 2 decimal Float value

package currency

import (
	"fmt"
	"math/big"
)

type Currency int64

// Format formats a currency as a float
func (c Currency) Format(f fmt.State, r rune) {
	b := big.NewFloat(float64(c) / 100)
	b.Format(f, r)

}

// NewCurrency codifies a float as a currency
func NewCurrency(f float64) Currency {
	return Currency(f * 100)
}
