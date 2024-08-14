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

	// по одному символу
	//history, err := client.GetOptionHistory("GD2300BU4", "2024-01-01", "")
	//if err != nil {
	//	slog.Error("main", "ошибкаGetOptionHistory", err.Error())
	//	return
	//}

	history, err := client.GetOptionHistoryAllDate("2024-08-01")
	if err != nil {
		slog.Error("main", "ошибка GetOptionHistory", err.Error())
		return
	}

	slog.Info("History", "всего len(history)", len(history))
	for row, h := range history {
		if h.OpenPosition != 0 {
			slog.Info(strconv.Itoa(row),
				"Data", h,
			)
		}
	}

}
