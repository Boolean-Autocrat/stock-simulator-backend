package engine

type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

// Add buy order to order book
func (book *OrderBook) addBuyOrder(order Order) {
	n := len(book.BuyOrders)
	var i int
	for i = n - 1; i >= 0; i-- {
		buyOrder := book.BuyOrders[i]
		if buyOrder.Price < order.Price {
			break
		}
	}
	if i == n-1 {
		book.BuyOrders = append(book.BuyOrders, order)
	} else {
		if i >= 0 {
			copy(book.BuyOrders[i+1:], book.BuyOrders[i:])
			book.BuyOrders[i] = order
		} else {
			book.BuyOrders = append([]Order{order}, book.BuyOrders...)
		}
	}
}

// Add sell order to order book
func (book *OrderBook) addSellOrder(order Order) {
	n := len(book.SellOrders)
	var i int
	for i = n - 1; i >= 0; i-- {
		sellOrder := book.SellOrders[i]
		if sellOrder.Price > order.Price {
			break
		}
	}
	if i == n-1 {
		book.SellOrders = append(book.SellOrders, order)
	} else {
		if i >= 0 {
			copy(book.SellOrders[i+1:], book.SellOrders[i:])
			book.SellOrders[i] = order
		} else {
			book.SellOrders = append([]Order{order}, book.SellOrders...)
		}
	}
}

// Remove buy order at given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove sell order at given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
