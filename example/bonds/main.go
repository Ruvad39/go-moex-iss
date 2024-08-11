package main

import (
	"github.com/Ruvad39/go-moex-iss"
	"log/slog"
	"strconv"
)

func main() {

	// создание клиента
	err, client := iss.NewClient()
	if err != nil {
		slog.Error("main", "NewClient", err.Error())
	}
	iss.SetLogLevel(slog.LevelDebug)

	// получить список бондов
	// ОФЗ: board = tqob
	// Корпоративные: board = TQIR
	Sec, err := client.GetBondsInfo("tqob")
	//Sec, err := client.GetBondsInfo("TQIR")

	if err != nil {
		slog.Error("main", "ошибка GetBondsInfo", err.Error())
	}

	slog.Info("GetBondsInfo", slog.Int("всего len(Sec)", len(Sec)))
	for row, sec := range Sec {

		slog.Info(strconv.Itoa(row),
			"sec", sec,
		)
		//slog.Info("sec", slog.Any("sec", sec))
	}

	//// текущиме рыночные ланыне
	//SecData, err := client.GetStockMarketData("SBER,GAZP")
	//if err != nil {
	//	slog.Error("main", "ошибка GetStockMarketData", err.Error())
	//}
	//
	//slog.Info("GetStockMarketData", slog.Int("всего len(Sec)", len(SecData)))
	//for row, sec := range SecData {
	//
	//	slog.Info(strconv.Itoa(row),
	//		"SecData", sec,
	//	)
	//
	//}
}
