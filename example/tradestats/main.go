package main

import (
	iss "github.com/Ruvad39/go-moex-iss"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
		slog.Error(err.Error())
	}
}

func main() {
	// получить данные можно только с авторизацией
	user, _ := os.LookupEnv("MOEX_USER")
	pwd, _ := os.LookupEnv("MOEX_PWD")

	client, err := iss.NewClient(iss.WithUser(user), iss.WithPwd(pwd))
	if err != nil {
		slog.Error("main", "NewClient", err.Error())
	}
	//iss.SetLogLevel(slog.LevelDebug)

	stats, err := client.GetStockTradeStats("SBER", "2024-09-01", "2024-09-03", false)
	//stats, err := client.GetStockTradeStatsAll("2024-09-03", false)
	//stats, err := client.GetFortsTradeStats("SiU4", "2024-09-01", "2024-09-03", false)
	//stats, err := client.GetFortsTradeStatsAll("2024-09-03", false)
	//stats, err := client.GetFxTradeStatsAll("2024-09-03", false)

	if err != nil {
		slog.Error("main", "ошибка GetTradeStats", err.Error())
	}

	// цикл по списку свечей
	for row, candle := range stats {
		//fmt.Println(row, candle)
		slog.Info(strconv.Itoa(row),
			"stats", candle,
		)
	}
}
