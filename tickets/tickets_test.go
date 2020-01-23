package tickets

import (
	"fmt"
	"testing"

	"github.com/lafont-e/LanaChallenge/currency"
	"github.com/lafont-e/LanaChallenge/products"
)

var PEN, TSHIRT, MUG *products.Product

func loadProducts() (err error) {
	PEN, err = products.GetProduct("PEN")
	if err == nil {
		TSHIRT, err = products.GetProduct("TSHIRT")
	}
	if err == nil {
		MUG, err = products.GetProduct("MUG")
	}
	return
}

func TestTicket(t *testing.T) {
	if err := loadProducts(); err != nil {
		t.Error("can not load products", err)
	}
	type cline struct {
		quantity int
		prod     *products.Product
	}
	var buyList = []struct {
		list     string
		total    currency.Currency
		articles []*cline
	}{
		{"NoDiscount", currency.NewCurrency(32.50), []*cline{&cline{1, PEN}, &cline{1, TSHIRT}, &cline{1, MUG}}},
		{"2x1.A", currency.NewCurrency(25.00), []*cline{&cline{1, PEN}, &cline{1, TSHIRT}, &cline{1, PEN}}},
		{"2x1.B", currency.NewCurrency(25.00), []*cline{&cline{2, PEN}, &cline{1, TSHIRT}}},
		{"2x1.C", currency.NewCurrency(30.00), []*cline{&cline{2, PEN}, &cline{1, TSHIRT}, &cline{1, PEN}}},
		{"3plus.A", currency.NewCurrency(52.50), []*cline{&cline{1, TSHIRT}, &cline{1, TSHIRT}, &cline{1, MUG}, &cline{1, TSHIRT}}},
		{"3plus.B", currency.NewCurrency(52.50), []*cline{&cline{3, TSHIRT}, &cline{1, MUG}}},
		{"3plus.C", currency.NewCurrency(67.50), []*cline{&cline{3, TSHIRT}, &cline{1, MUG}, &cline{1, TSHIRT}}},
		{"3plus.D", currency.NewCurrency(67.50), []*cline{&cline{4, TSHIRT}, &cline{1, MUG}}},
		{"AllDisc.A", currency.NewCurrency(62.50), []*cline{&cline{1, PEN}, &cline{1, TSHIRT}, &cline{1, PEN}, &cline{1, PEN}, &cline{1, MUG}, &cline{1, TSHIRT}, &cline{1, TSHIRT}}},
		{"AllDisc.B", currency.NewCurrency(62.50), []*cline{&cline{3, PEN}, &cline{3, TSHIRT}, &cline{1, MUG}}},
		{"AllDisc.C", currency.NewCurrency(77.50), []*cline{&cline{3, PEN}, &cline{4, TSHIRT}, &cline{1, MUG}}},
	}

	for _, bl := range buyList {
		t.Run(bl.list, func(t *testing.T) {
			ticket := NewTicket(Promotions)
			for _, codedLine := range bl.articles {
				ticket.Add(codedLine.quantity, codedLine.prod)
			}

			if ticket.Total() != bl.total {
				fmt.Println(bl.list, ticket.String())
				t.Error("Bad Total")
			}

		})
	}
}

func TestDiscount2x1(t *testing.T) {
	if err := loadProducts(); err != nil {
		t.Error("can not load products", err)
	}

	ticket := NewTicket(Promotions)
	price := PEN.GetPrice()
	ticket.Add(2, PEN)
	if ticket.Total() != price {
		fmt.Println(ticket.String())
		t.Error("Discount not applied")
	}

	ticket.Add(1, PEN)
	if ticket.Total() != (2 * price) {
		fmt.Println(ticket.String())
		t.Error("Total badly calculated")
	}

	ticket.Add(1, PEN)
	if ticket.Total() != (2 * price) {
		fmt.Println(ticket.String())
		t.Error("Total badly calculated")
	}

}

func TestDiscount3plus(t *testing.T) {
	if err := loadProducts(); err != nil {
		t.Error("can not load products", err)
	}

	ticket := NewTicket(Promotions)
	price := TSHIRT.GetPrice()
	ticket.Add(2, TSHIRT)
	if ticket.Total() != (2 * price) {
		fmt.Println(ticket.String())
		t.Error("Total badly calculated")
	}

	ticket.Add(1, TSHIRT)
	if ticket.Total() != 3*(price-price/4) {
		fmt.Println(ticket.String())
		t.Error("Discount not applied")
	}

	ticket.Add(2, TSHIRT)
	if ticket.Total() != 5*(price-price/4) {
		fmt.Println(ticket.String())
		t.Error("Discount not applied")
	}
}
