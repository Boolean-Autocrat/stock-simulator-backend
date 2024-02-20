package engine

// Process and return trades generated to market
func (book *OrderBook) Process(order Order) []Trade {
	if order.Side == 1 {
		return book.processLimitBuy(order)
	}
	return book.processLimitSell(order)
}

// Process a limit buy order
func (book *OrderBook) processLimitBuy(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	n := len(book.SellOrders)
	// find min. one matching order
	if n != 0 && book.SellOrders[n-1].Price <= order.Price {
		// loop through all matching orders
		for i := n - 1; i >= 0; i-- {
			sellOrder := book.SellOrders[i]
			if sellOrder.Price > order.Price {
				break
			}
			// fill entire order
			if sellOrder.Amount >= order.Amount {
				trades = append(trades, Trade{
					BuyerOrderID:  order.OrderID,
					BuyerID:       order.UserID,
					SellerOrderID: sellOrder.OrderID,
					SellerID:      sellOrder.UserID,
					Amount:        order.Amount,
					Price:         sellOrder.Price,
					Stock:         order.Stock,
				})
				sellOrder.Amount -= order.Amount
				if sellOrder.Amount == 0 {
					book.removeSellOrder(i)
				}
				return trades
			}
			// fill partial order
			if sellOrder.Amount < order.Amount {
				trades = append(trades, Trade{
					BuyerOrderID:  order.OrderID,
					BuyerID:       order.UserID,
					SellerOrderID: sellOrder.OrderID,
					SellerID:      sellOrder.UserID,
					Amount:        sellOrder.Amount,
					Price:         sellOrder.Price,
					Stock:         order.Stock,
				})
				order.Amount -= sellOrder.Amount
				book.removeSellOrder(i)
				continue
			}
		}
	}
	// add remaining order
	book.addBuyOrder(order)
	return trades
}

// Process a limit sell order
func (book *OrderBook) processLimitSell(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	n := len(book.BuyOrders)
	// find min. one matching order
	if n != 0 && book.BuyOrders[n-1].Price >= order.Price {
		// loop through all matching orders
		for i := n - 1; i >= 0; i-- {
			buyOrder := book.BuyOrders[i]
			if buyOrder.Price < order.Price {
				break
			}
			// fill entire order
			if buyOrder.Amount >= order.Amount {
				trades = append(trades, Trade{
					BuyerID:  buyOrder.UserID,
					SellerID: order.UserID,
					Amount:   order.Amount,
					Price:    buyOrder.Price,
					Stock:    order.Stock,
				})
				buyOrder.Amount -= order.Amount
				if buyOrder.Amount == 0 {
					book.removeBuyOrder(i)
				}
				return trades
			}
			// fill partial order
			if buyOrder.Amount < order.Amount {
				trades = append(trades, Trade{
					BuyerID:  buyOrder.UserID,
					SellerID: order.UserID,
					Amount:   buyOrder.Amount,
					Price:    buyOrder.Price,
					Stock:    order.Stock,
				})
				order.Amount -= buyOrder.Amount
				book.removeBuyOrder(i)
				continue
			}
		}
	}
	// add remaining order
	book.addSellOrder(order)
	return trades
}
