// Deals with the checkout process

package tickets

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/elafont/LanaChallenge/products"
)

// line represents a product sold on a ticket
type line struct {
	quantity int
	code     string
	name     string
	price    float32
}

// Ticket is a set of products sold expressed as lines
type Ticket struct {
	sync.Mutex
	lines      []*line               // Lines in the ticket
	discounts  []*line               // Discounts in the ticket
	ts         time.Time             // Creation date for the ticket
	total      float32               // Total for the ticket
	discTotal  float32               // Total of discounts
	promotions map[string]*promotion // Promotions to apply to this ticket
	up2date    bool                  // is the ticket up 2 date ?
}

// NewTicket returns an empty ticket
func NewTicket(promo map[string]*promotion) *Ticket {
	if promo == nil {
		promo = NoPromotions
	}
	return &Ticket{ts: time.Now(), promotions: promo, lines: make([]*line, 0, 0), discounts: make([]*line, 0, 0)}
}

// Add products to a ticket
func (t *Ticket) Add(q int, product *products.Product) {
	if q != 0 { // does not add a line if quantity is zero
		t.Lock()
		t.addLocked(q, product)
		t.Unlock()
	}
}

// addLocked products to a ticket when ticket is Locked
func (t *Ticket) addLocked(q int, product *products.Product) {
	if q != 0 { // does not add a line if quantity is zero
		t.lines = append(t.lines, &line{quantity: q, code: product.Code, name: product.Name, price: product.Price})
		t.total += float32(q) * product.Price
		t.up2date = false
	}
}

// AddDiscountLocked discounts products in a ticket when ticket is Locked
func (t *Ticket) addDiscountLocked(q int, product *products.Product) {
	if q != 0 { // does not add a line if quantity is zero
		t.discounts = append(t.discounts, &line{quantity: q, code: product.Code, name: product.Name, price: product.Price})
		t.discTotal += float32(q) * product.Price
	}
}

// clearDiscountsLocked removes all discounts of a ticket
func (t *Ticket) clearDiscountsLocked() {
	t.discounts = t.discounts[:0] // Resets the slice len
	t.discTotal = 0               // zeroes the total discount counter
}

// Total calc the discounts and total price of the ticket
func (t *Ticket) Total() float32 {
	if !t.up2date {
		t.Lock()
		t.clearDiscountsLocked()

		for promoCode, promo := range t.promotions {
			var linesCode []*line
			for ix, l := range t.lines {
				if l.code == promoCode {
					linesCode = append(linesCode, t.lines[ix])
				}
			}

			// This function will evaluate the lines with the same CODE and will add discounts to the ticket if needed
			promo.function(t, linesCode)
		}

		t.Unlock()
		t.up2date = true
	}
	return (t.total + t.discTotal)

}

// Status shows the status of the ticket
func (t *Ticket) Status() string {
	return fmt.Sprintf("ticket created on %s, amount %f", t.ts, t.Total())
}

// String, returns the whole ticket in printed form
func (t *Ticket) String() string {
	var tks strings.Builder
	var tkTotal = t.Total()

	header := "LANA STORE              date:%19s\n\n"
	subhead := "Qty Article                          Price   Total "
	separator := "--- ------------------------------ ------- -------"
	itemLine := "%3d %30s  %6.2f  %6.2f\n"
	footer := "                                    Total:  %6.2f\n"

	tks.WriteString(fmt.Sprintln())
	tks.WriteString(fmt.Sprintf(header, t.ts.Format("02, Jan/2006 15:04")))
	tks.WriteString(fmt.Sprintln(subhead))
	tks.WriteString(fmt.Sprintln(separator))

	t.Lock()
	for _, line := range t.lines {
		quantity, price := line.quantity, line.price
		tks.WriteString(fmt.Sprintf(itemLine, quantity, line.name, price, float32(quantity)*price))
	}
	if len(t.discounts) > 0 {
		tks.WriteString(fmt.Sprintln(separator))
		for _, line := range t.discounts {
			quantity, price := line.quantity, line.price
			tks.WriteString(fmt.Sprintf(itemLine, quantity, line.name, price, float32(quantity)*price))
		}
	}
	t.Unlock()

	tks.WriteString(fmt.Sprintln(separator))
	tks.WriteString(fmt.Sprintf(footer, tkTotal))
	tks.WriteString(fmt.Sprintln())
	return tks.String()
}
