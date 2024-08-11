package main

import (
	iss "github.com/Ruvad39/go-moex-iss"
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

	// получить список опционов
	Sec, err := client.GetOptionInfo("Si90000BI4")
	//Sec, err := client.GetOptionInfo("")
	if err != nil {
		slog.Error("main", "ошибка GetOptionInfo", err.Error())
	}

	slog.Info("GetOptionInfo", slog.Int("всего len(Sec)", len(Sec)))
	for row, sec := range Sec {

		slog.Info(strconv.Itoa(row),
			"sec", sec,
		)

	}

	// получить по опционам рыночные данные
	SecData, _ := client.GetOptionData("Si90000BI4")
	//SecData, _ := client.GetOptionMarketData("")
	slog.Info("GetOptionMarketData", slog.Int("всего len(Sec)", len(SecData)))
	for row, sec := range SecData {
		slog.Info(strconv.Itoa(row),
			"sec", sec,
		)
	}

}
