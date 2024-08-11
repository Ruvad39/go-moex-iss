package main

import (
	"github.com/Ruvad39/go-moex-iss"
	"log/slog"
	"strconv"
)

func main() {

	// создание клиента
	_, client := iss.NewClient()
	//iss.SetLogLevel(slog.LevelDebug)

	// получить список акций
	Sec, err := client.GetStockInfo("")
	//Sec, err := client.GetStockInfo("SBER,MOEX")
	if err != nil {
		slog.Error("main", "ошибка GetStockInfo", err.Error())
	}

	slog.Info("GetStockInfo", slog.Int("всего len(Sec)", len(Sec)))
	for row, sec := range Sec {
		slog.Info(strconv.Itoa(row),
			"sec", sec,
		)
	}

	// текущие рыночные данные
	SecData, err := client.GetStockData("SBER,ROSN,MOEX")
	if err != nil {
		slog.Error("main", "ошибка GetStockData", err.Error())
	}

	slog.Info("GetStockData", slog.Int("всего len(Sec)", len(SecData)))
	for row, sec := range SecData {

		slog.Info(strconv.Itoa(row),
			"SecData", sec,
		)

	}
}
