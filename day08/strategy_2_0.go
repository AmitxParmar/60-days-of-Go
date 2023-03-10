package main

import "fmt"

// struct is the way go group some state
// and encapsulation is done by first letter
// upper case letter indicates public, lower case letters indicates private
// this encapsulation is valid on the whole package(further discussion)
// we talk more about POO in another day.

// Customer is a single struct with name and fidelty points
type Customer struct {
	name     string
	fidelity int
}

// LIneItem represents a item in a cart
type LineItem struct {
	product  string
	quantity int
	price    float64
}

// Methods with receiver(item in function below) are binded with struct

// Total returns quantity of items multiplied by price
func (item LineItem) Total() float64 {
	return float64(item.quantity) * item.price
}

// String is a better representation of an item
func (item LineItem) String() string {
	return fmt.Sprintf("<LineItem product:%s quantity:%d price:%.2f>", item.product, item.quantity, item.price)
}

// Order is the relationship of a customer, the cart and possible promo
type Order struct {
	ctm   Customer
	cart  []LineItem
	promo func(Order) float64 // promo becomes a function
}

// Total is the sum of items purchased
func (order Order) Total() float64 {
	total := 0.0
	for _, item := range order.cart {
		total += item.Total()
	}
	return total
}

// Due calculate order value considering discount
func (order Order) Due() float64 {
	discount := 0.0
	if order.promo != nil {
		discount = order.promo(order)
	}
	return order.Total() - discount
}

// String returns the order representation when is printed
func (order Order) String() string {
	return fmt.Sprintf("<Order total: %.2f due: %.2f>", order.Total(), order.Due())
}

// FidelityPromo receives an order and return a discount
func FidelityPromo(o Order) float64 {
	if o.ctm.fidelity >= 1000 {
		return o.Total() * 0.05
	}
	return 0.0
}

// BulkItemPromo receives an order return a discount
func BulkItemPromo(o Order) float64 {
	discount := 0.0
	for _, item := range o.cart {
		if item.quantity >= 20 {
			discount += item.Total() * .1
		}
	}
	return discount
}

// LargeOrderPromo receives an order and return a discount
func LargeOrderPromo(o Order) float64 {
	set := map[string]bool{}
	for _, item := range o.cart {
		set[item.product] = true
	}
	if len(set) >= 10 {
		return o.Total() * 0.07
	}
	return 0.0
}

func main() {
	// This example is the same present in excellent python book, Fluent Python, written by Luciano Ramalho
	// This example uses function as parameter
	joe := Customer{"John Doe", 0}
	ann := Customer{"Ann Smith", 1100}
	cart := []LineItem{
		LineItem{"banana", 4, 0.50},
		LineItem{"apple", 10, 1.50},
		LineItem{"banana", 5, 5.00},
	}
	// Joe don't have fidelity points, he don't win a discount
	fmt.Printf("\n% have %d fidelity points\n", joe.name, joe.fidelity)
	fmt.Println(Order{joe, cart, FidelityPromo})
	// Ann have fidelity points, this guarantees a discount.
	fmt.Printf("\n% have %d fidelity points\n", joe.name, joe.fidelity)
	fmt.Println(Order{ann, cart, FidelityPromo})

	// 30 bananas??
	bananaCart := []LineItem{
		LineItem{"banana", 30, .5},
		LineItem{"apple", 10, 1.5},
	}
	// Ok, many items guarantees discount on BulkItemPromo
	fmt.Printf("\n%s buy many items of the same product %s\n", joe.name,
		bananaCart)
	fmt.Println(Order{joe, bananaCart, BulkItemPromo})
	// 10 random items
	largeOrder := []LineItem{}
	for i := 0; i < 10; i++ {
		largeOrder = append(largeOrder, LineItem{string(65 + i), 1, 1.0})
	}
	// only to check LargeOrderPromo
	fmt.Printf("\n%s represents an order with many distinct items %s", joe.name, largeOrder)
	fmt.Println(Order{joe, largeOrder, LargeOrderPromo})
	// only 3 distinct items, no discount here!
	fmt.Printf("\nonly 3 distinct items, no discount here!")
	fmt.Println(Order{joe, cart, LargeOrderPromo})
}
