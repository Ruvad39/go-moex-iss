package main

import (
	"fmt"
	"log/slog"
	"strconv"

	iss "github.com/Ruvad39/go-moex-iss"
)

func main() {
	// создание клиента
	// получить данные по стакану можно только с авторизацией
	user := ""
	pwd := ""
	client, err := iss.NewClient(iss.WithUser(user), iss.WithPwd(pwd))
	if err != nil {
		slog.Error("main", "NewClient", err.Error())
	}
	//iss.SetLogLevel(slog.LevelDebug)

	// TODO
	// Sec, err := client.GetTickersAll()
	// slog.Info("GetTickersAll", slog.Int("всего len(Sec)", len(Sec)))
	// for row, sec := range Sec {
	// 	slog.Info(strconv.Itoa(row),
	// 		"sec", sec,
	// 	)
	// }
	//https://iss.moex.com/iss/securities.json

	// создание (поиск) тикера
	//ticker, err := client.GetTicker("SBER")
	ticker, err := client.GetTicker("RTS-6.26")
	if err != nil {
		slog.Error("main", "ошибка NewTicker", err.Error())
	}
	slog.Info("NewTicker", slog.Any("Ticker", ticker))

	// параметры инструмента
	info, err := ticker.Info()
	if err != nil {
		slog.Error("main", "ошибка ticker.Info", err.Error())
	}
	slog.Info("ticker.Info", slog.Any("t_info", info))

	// текущие рыночные данные
	data, err := ticker.Data()
	if err != nil {
		slog.Error("main", "ошибка ticker.data", err.Error())
	}
	slog.Info("ticker.Info", slog.Any("t_data", data))

	// свечи
	// candles(ticker)
	// стакан
	orderBook(ticker)
}

// получим список свечей
func candles(ticker *iss.Ticker) {
	// свечи
	candles, err := ticker.Candles(iss.Interval_M1, "2026-01-01", "2026-12-01")
	if err != nil {
		slog.Error("main", "ошибка Candles", err.Error())
		return
	}
	slog.Info("Candles",
		"всего len(candles)", candles.Len(),
		"mindate", candles.First().Begin,
		"maxdate", candles.Last(0).Begin,
		"Symbol", candles.Symbol,
		"Interval", candles.Interval,
	)

	// цикл по списку свечей
	for row, candle := range candles.Data {
		slog.Info(strconv.Itoa(row),
			"time", candle.Time(),
			"candleData", candle,
		)
	}
}

// получим стакан
// Нужна авторизация
func orderBook(ticker *iss.Ticker) {

	orderBook, err := ticker.OrderBook()
	if err != nil {
		slog.Error("main", "ticker.OrderBook", err.Error())
		return
	}
	bid, _ := orderBook.BestBid()
	ask, _ := orderBook.BestAsk()
	bidVolume := orderBook.Bids.SumDepth()
	askVolume := orderBook.Asks.SumDepth()

	slog.Info("orderBook", "BestAsk", ask.Price, "BestBid", bid.Price, "объем асков", askVolume, "объем бидов", bidVolume)
	fmt.Println(orderBook.String())
}
