// Discounts is the package that calcs if any discount is applicable on a ticket

package tickets

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/lafont-e/LanaChallenge/products"
)

// ppromotionsFile name of the file with the promotions to manage
const promotionsFile = "./promotions.json"
const discountCode = "-@"

// Discounts, rule to apply for discounts, every rule will be applied to every line
// of the ticket, if no items are declared, the rule will be applied to the whole ticket
type discfunc func(*Ticket, []*line) int
type promotion struct {
	code     string
	function discfunc
}

// Promotions acceptable on a ticket
var Promotions = make(map[string]*promotion)

// NoPromotions empty Promotions table to be used when no promotions are offered
var NoPromotions = make(map[string]*promotion)

// List of available discount functions by name
// THIS ARRAY MUST BE DECLARED BEFORE USE AS ITS THE ONE THAT links the json file with the promotions
var funcCodes = map[string]discfunc{
	"2*1": discount2x1,
	"3+":  discount3plus,
}

func init() { // Load the promotions offered
	loadPromotions()
}

// discount2x2 rule to discount half of the products (ticket must be locked)
func discount2x1(t *Ticket, lines []*line) int {
	var q int

	if len(lines) > 0 { // lines could not be empty if any discount should be applied
		var price = lines[0].price
		var code = lines[0].code

		for _, l := range lines {
			q += l.quantity
		}

		// Add a line with the discounts in the ticket
		q /= 2
		t.addDiscountLocked(q, products.NewProduct(discountCode, code+" 2x1 Discount", -price))
	}

	return q // Number of discounts applied
}

// discount3plus rule to discount 25% when 3 or more products are bought (ticket must be locked)
func discount3plus(t *Ticket, lines []*line) int {
	var q int

	if len(lines) > 0 { // lines could not be empty if any discount should be applied
		var price = lines[0].price
		var code = lines[0].code

		for _, l := range lines {
			q += l.quantity
		}

		if q > 2 { // if more than 2 products are bought, then apply the discount
			// Add a line with the discounts in the ticket
			t.addDiscountLocked(q, products.NewProduct(discountCode, code+" 25% on +3 Discount", -price/4)) // 25% discount
		}
	}
	return q
}

// loadPromotions function that reads products from a json file
func loadPromotions() error {
	// read file
	data, err := ioutil.ReadFile(promotionsFile)
	if err != nil {
		log.Println("Error, Not Loading Promotions :", err)
		return err
	}

	type jspromotion struct{ Code, Function string }

	var pm []jspromotion

	// unmarshall it
	err = json.Unmarshal(data, &pm)

	if err != nil {
		log.Println("Error, not able to read promotions file :", err)
		return err
	}

	for _, p := range pm {
		if f, ok := funcCodes[p.Function]; ok {
			Promotions[p.Code] = &promotion{code: p.Code, function: f}
		}
	}

	return nil
}
