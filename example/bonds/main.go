package main

import (
	"github.com/Ruvad39/go-moex-iss"
	"log/slog"
)

func main() {

	// создание клиента
	client, err := iss.NewClient()
	if err != nil {
		slog.Error("main", "NewClient", err.Error())
	}
	iss.SetLogLevel(slog.LevelDebug)

	Sec, err := client.GetBondsInfo("tqob") // ОФЗ
	//Sec, err := client.GetBondsInfo("TQCB") // Корпоративные
	// выборка по boardGroup = все облигации
	//Sec, err := client.GetBondsInfo("", 58)

	if err != nil {
		slog.Error("main", "ошибка GetBondsInfo", err.Error())
	}

	for n, sec := range Sec {
		slog.Info("bonds",
			"row", n,
			"SecID", sec.SecID,
			"SecName", sec.SecName,
			"board", sec.BoardID,
			"CouponDate", sec.NextCoupon,
			"MatDate", sec.MatDate,
			"CouponValue", sec.CouponValue,
			"PrevPrice", sec.PrevPrice,
		)
	}
	slog.Info("GetBondsInfo", slog.Int("всего len(Sec)", len(Sec)))

	// текущие рыночные данные по облигации
	//SecData, err := client.GetBondsData("tqob", "SU26218RMFS6")
	//SecData, err := client.GetBondsData("tqob", "")
	//SecData, err := client.GetBondsData("", "", 58)
	//if err != nil {
	//	slog.Error("main", "ошибка GetBondsData", err.Error())
	//}
	//
	//slog.Info("GetBondsData", slog.Int("всего len(Sec)", len(SecData)))
	//for row, sec := range SecData {
	//	slog.Info(strconv.Itoa(row),
	//		"SecData", sec,
	//	)
	//}

	// Исторические данные
	SecHistory, err := client.GetBondHistory("SU26207RMFS9", "2025-05-19", "2025-05-31")
	//SecHistory, err := client.GetBondsHistoryDate("2025-02-21")
	if err != nil {
		slog.Error("main", "ошибка GetBondHistory", err.Error())
	}

	slog.Info("GetBondHistory", slog.Int("всего len(Sec)", len(SecHistory)))
	slog.Info("GetBondHistory", slog.Any("SecHistory", SecHistory))
	// for row, sec := range SecHistory {
	// 	slog.Info(strconv.Itoa(row),
	// 		"SecHistory", sec,
	// 	)
	// }

	// свечи
	// candles, err :=
	// service := client.NewCandlesService("stock", "shares", "tqbr", "SBER",
	// 	int(iss.Interval_D1),
	// 	"2025-05-01",
	// 	"2025-06-01",
	//).Do()
	// candles, err := service.Do()

}
