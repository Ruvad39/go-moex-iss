package main

import (
	"log/slog"
	"strconv"

	iss "github.com/Ruvad39/go-moex-iss"
)

func main() {
	// создание клиента
	err, client := iss.NewClient()
	if err != nil {
		slog.Error("main", "NewClient", err.Error())
	}
	//iss.SetLogLevel(slog.LevelDebug)

	// получить список фьчерсов
	Sec, err := client.GetFortsInfo("SiU4,MXU4")
	//Sec, err := client.GetFortsInfo("")
	if err != nil {
		slog.Error("main", "ошибка GetFortsInfo", err.Error())
	}

	slog.Info("GetFortsInfo", slog.Int("всего len(Sec)", len(Sec)))
	for row, sec := range Sec {

		slog.Info(strconv.Itoa(row),
			"sec", sec,
		)

	}

	// получить по фючерсам рыночные данные
	SecData, _ := client.GetFortsData("CRU4,MXU4")

	slog.Info("GetFortsData", slog.Int("всего len(Sec)", len(SecData)))
	for row, sec := range SecData {
		slog.Info(strconv.Itoa(row),
			"sec", sec,
		)
	}

}
