package engine

import (
	"context"
	"log"
	"time"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
)

func RunTradeQueries(Trade Trade, queries *db.Queries) {
	log.Println("Starting trade process at", time.Now().Format("2006-01-02 15:04:05"))
	tradeErr := queries.CreateTrade(context.Background(), db.CreateTradeParams{
		Buyer:    Trade.BuyerID,
		Seller:   Trade.SellerID,
		Quantity: Trade.Amount,
		Stock:    Trade.Stock,
		Price:    Trade.Price,
	})
	if tradeErr != nil {
		log.Println(tradeErr.Error())
		panic(tradeErr)
	}
	queries.BeginTransaction(context.Background())
	sellerBalanceErr := queries.UpdateBalance(context.Background(), db.UpdateBalanceParams{
		ID:      Trade.SellerID,
		Balance: float32(Trade.Amount) * Trade.Price,
	})
	if sellerBalanceErr != nil {
		log.Println(sellerBalanceErr.Error())
		panic(sellerBalanceErr)
	}
	buyerPortfolioErr := queries.AddOrUpdateStockToPortfolio(context.Background(), db.AddOrUpdateStockToPortfolioParams{
		User:     Trade.BuyerID,
		Stock:    Trade.Stock,
		Quantity: Trade.Amount,
	})
	if buyerPortfolioErr != nil {
		log.Println(buyerPortfolioErr.Error())
		panic(buyerPortfolioErr)
	}
	queries.EndTransaction(context.Background())
	stock, _ := queries.GetStockById(context.Background(), Trade.Stock)
	var stockTrend string
	var stockPercentageChange float32
	if stock.Price < Trade.Price {
		stockTrend = "up"
		stockPercentageChange = (Trade.Price - stock.Price) / stock.Price * 100
		stockPercentageChange = float32(int(stockPercentageChange*100)) / 100
	} else if stock.Price == Trade.Price {
		stockTrend = "unchanged"
		stockPercentageChange = 0
	} else {
		stockTrend = "down"
		stockPercentageChange = (stock.Price - Trade.Price) / stock.Price * 100
		stockPercentageChange = float32(int(stockPercentageChange*100)) / 100
	}
	queries.BeginTransaction(context.Background())
	updateStockPriceErr := queries.UpdateStockPrice(context.Background(), db.UpdateStockPriceParams{
		ID:               Trade.Stock,
		Price:            Trade.Price,
		Trend:            stockTrend,
		PercentageChange: stockPercentageChange,
	})
	if updateStockPriceErr != nil {
		log.Println(updateStockPriceErr.Error())
		panic(updateStockPriceErr)
	}
	queries.EndTransaction(context.Background())
	priceHistoryErr := queries.CreatePriceHistory(context.Background(), db.CreatePriceHistoryParams{
		Stock: Trade.Stock,
		Price: Trade.Price,
	})
	if priceHistoryErr != nil {
		log.Println(priceHistoryErr.Error())
		panic(priceHistoryErr)
	}
	updateBuyerOrderErr := queries.UpdatePendingOrder(context.Background(), db.UpdatePendingOrderParams{
		ID:                Trade.BuyerOrderID,
		FulfilledQuantity: Trade.Amount,
	})
	if updateBuyerOrderErr != nil {
		log.Println(updateBuyerOrderErr.Error())
		panic(updateBuyerOrderErr)
	}
	updateSellerOrderErr := queries.UpdatePendingOrder(context.Background(), db.UpdatePendingOrderParams{
		ID:                Trade.SellerOrderID,
		FulfilledQuantity: Trade.Amount,
	})
	if updateSellerOrderErr != nil {
		log.Println(updateSellerOrderErr.Error())
		panic(updateSellerOrderErr)
	}
	log.Println("Trade process completed at", time.Now().Format("2006-01-02 15:04:05"))
}
