package main

import (
	"fmt"
	iss "github.com/Ruvad39/go-moex-iss"
	"log/slog"
)

func main() {
	// получить данные можно только с авторизацией
	user := ""
	pwd := ""
	err, client := iss.NewClient(iss.WithUser(user), iss.WithPwd(pwd))
	//err, client := iss.NewClient()
	if err != nil {
		slog.Error("main", "NewClient", err.Error())
	}

	iss.SetLogLevel(slog.LevelDebug)

	//
	service := client.NewOrderBookService("stock", "shares", "tqbr", "SBER")
	//
	//slog.Info("OrderBookService", "URL", service.URL())
	orderBook, err := service.Do()
	if err != nil {
		slog.Error("main", "ошибка OrderBookService", err.Error())
	}
	bid, _ := orderBook.BestBid()
	ask, _ := orderBook.BestAsk()
	bidVolume := orderBook.Bids.SumDepth()
	askVolume := orderBook.Asks.SumDepth()

	//slog.Info("main", "OrderBookService", orderBook)
	fmt.Printf("BestAsk %f BestBid %f, объем асков= %d объем бидов= %d \n", ask.Price, bid.Price, askVolume, bidVolume)
	fmt.Println(orderBook.String())

}
