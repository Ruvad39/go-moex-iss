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
		//slog.Error(err.Error())
	}
}

func main() {
	// создание клиента
	// получить некоторые данные можно только с авторизацией
	user, _ := os.LookupEnv("MOEX_USER")
	pwd, _ := os.LookupEnv("MOEX_PWD")
	err, client := iss.NewClient(iss.WithUser(user), iss.WithPwd(pwd))
	if err != nil {
		slog.Error("main", "NewClient", err.Error())
	}
	iss.SetLogLevel(slog.LevelDebug)

	// если не указать дату = то данные за последний день
	//oi, err := client.GetFutOIAll("", 1)
	// latest =1 возвращает последнюю пятиминутку за указанную дату
	//oi, err := client.GetFutOIAll("2024-08-09", 1)
	//if err != nil {
	//	slog.Error("main", "ошибка GetFutOIAll", err.Error())
	//}
	//slog.Info("GetFutOIAll", slog.Int("всего len(Sec)", len(oi)))
	//for row, sec := range oi {
	//	slog.Info(strconv.Itoa(row),
	//		"oi", sec,
	//	)
	//}

	// по заданному тикеру
	oi, err := client.GetFutOI("Si", "2024-08-01", "2024-09-01", 0)
	if err != nil {
		slog.Error("main", "ошибка GetFutOI", err.Error())
	}
	slog.Info("GetFutOIAll", slog.Int("всего len(Sec)", len(oi)))
	for row, sec := range oi {
		slog.Info(strconv.Itoa(row),
			"oi", sec,
		)
	}
}
