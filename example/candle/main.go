package main

import (
	"fmt"
	iss "github.com/Ruvad39/go-moex-iss"
	"log/slog"
	"os"
	"strconv"
)

func main() {

	// создание клиента
	err, client := iss.NewClient()
	if err != nil {
		slog.Error("main", "NewClient", err.Error())
	}
	//iss.SetLogLevel(slog.LevelDebug)

	// пример через создание сервиса
	// service := client.NewCandlesService("stock", "shares", "tqbr", "SBER",
	//	int(iss.Interval_D1),
	//	"2000-01-01",
	//	"2025-01-01",
	//).Do()
	// candles, err := service.Do()

	// пример: получим интервал через парсинг
	// M1 M10 H1 D1 W1 MN1 Q1
	interval, err := iss.ParseInterval("D1")
	if err != nil {
		slog.Error("main", "ошибка ParseInterval", err.Error())
		return
	}

	// по акциям
	candles, err := client.GetStockCandles("SBER", interval, "2024-05-01", "2025-01-01")
	// по фючерсам
	//candles, err := client.GetFortsCandles("SiU4", interval, "2024-08-09", "2025-01-01")
	if err != nil {
		slog.Error("main", "ошибка GetCandles", err.Error())
		return
	}

	//
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

	// запишем в файл
	//err = fileWrite_2(&candles)
	//err = fileWrite(&candles)
	//if err != nil {
	//	slog.Error("main", "ошибка fileWrite", err.Error())
	//}

}

// time.Now().Format("2006-01-02")
func fileWrite(candles *iss.Candles) error {
	//fileName := candles.Symbol +_"" + candles.Interval + "-" + candles.First().Begin + "_" + candles.Last(0).Begin + ".csv"
	fileName := fmt.Sprint(candles.Symbol,
		"_",
		candles.Interval,
		"_",
		candles.First().Time().Format("20060102"),
		"-",
		candles.Last(0).Time().Format("20060102"),
		".csv")

	file, err := os.Create(fileName)
	if err != nil {
		slog.Error("fileWrite.os.Create", "err", err.Error())
		return err
	}
	_, err = file.WriteString("<BEGIN>,<OPEN>,<HIGH>,<LOW>,<CLOSE>,<VOLUME>" + "\n")
	if err != nil {
		slog.Error("fileWrite.file.WriteString", "err", err.Error())
		return err
	}

	delimiter := ","
	// цикл по списку свечей
	for _, candle := range candles.Data {

		str := fmt.Sprint(
			candle.Begin, delimiter,
			strconv.FormatFloat(candle.Open, 'f', -1, 64), delimiter,
			strconv.FormatFloat(candle.High, 'f', -1, 64), delimiter,
			strconv.FormatFloat(candle.Low, 'f', -1, 64), delimiter,
			strconv.FormatFloat(candle.Low, 'f', -1, 64), delimiter,
			strconv.FormatFloat(candle.Close, 'f', -1, 64), delimiter,
			strconv.FormatInt(int64(candle.Volume), 10),
		)
		file.WriteString(str + "\n")

	}

	file.Close()

	slog.Info("fileWrite", "файл успешно создан:", fileName)
	return nil

}

func fileWrite_2(candles *iss.Candles) error {
	//fileName := candles.Symbol +_"" + candles.Interval + "-" + candles.First().Begin + "_" + candles.Last(0).Begin + ".csv"
	fileName := fmt.Sprint(candles.Symbol,
		"_",
		candles.Interval,
		"_",
		candles.First().Time().Format("20060102"),
		"-",
		candles.Last(0).Time().Format("20060102"),
		".csv")

	file, err := os.Create(fileName)
	if err != nil {
		slog.Error("fileWrite.os.Create", "err", err.Error())
		return err
	}

	// без заголовка "time": 0, "open": 1, "close": 2, "low": 3, "high": 4, "volume": 5
	//_, err = file.WriteString("<BEGIN>,<OPEN>,<HIGH>,<LOW>,<CLOSE>,<VOLUME>" + "\n")
	//if err != nil {
	//	slog.Error("fileWrite.file.WriteString", "err", err.Error())
	//	return err
	//}

	delimiter := ","
	// цикл по списку свечей
	for _, candle := range candles.Data {

		// "time": 0, "open": 1, "close": 2, "low": 3, "high": 4, "volume": 5
		str := fmt.Sprint(
			candle.Time().Unix(), delimiter,
			strconv.FormatFloat(candle.Open, 'f', -1, 64), delimiter,
			strconv.FormatFloat(candle.Close, 'f', -1, 64), delimiter,
			strconv.FormatFloat(candle.High, 'f', -1, 64), delimiter,
			strconv.FormatFloat(candle.Low, 'f', -1, 64), delimiter,
			strconv.FormatFloat(candle.Low, 'f', -1, 64), delimiter,
			strconv.FormatInt(int64(candle.Volume), 10),
		)
		file.WriteString(str + "\n")

	}

	file.Close()

	slog.Info("fileWrite", "файл успешно создан:", fileName)
	return nil

}
