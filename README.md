# Библиотека, которая позволяет работать с функционалом [IIS Московской Биржи](https://www.moex.com/a2193)



## Установка


```bash
go get github.com/Ruvad39/go-moex-iss
```

## api реализован на текущий момент:

```go
// GetStockInfo получить параметры инструментов по акциям
GetStockInfo(symbols string) ([]StockInfo, error)
// GetStockMarketData получить рыночные данные по фондовому рынка
GetStockData(symbols string) ([]StockData, error)
// GetStockCandles получить историю свечей по акциям
GetStockCandles(symbols string, interval int, from, to string) (Candles, error)

// GetBondsInfo получить параметры инструментов 
GetBondsInfo(board string) ([]BondInfo, error)
// TODO получить рыночные данные

// GetFortsInfo получить параметры инструментов по фьючерсам
GetFortsInfo(symbols string) ([]FortsInfo, error)
// GetFortsMarketData получить рыночные данные по фьючерсам
GetFortsData(symbols string) ([]FortsData, error)
// GetFortsCandles получить историю свечей по фьючерсам
GetFortsCandles(symbols string, interval int, from, to string) (Candles, error)

// GetOptionInfo получить параметры инструментов по опционам
GetOptionInfo(symbols string) ([]OptionInfo, error)
// GetOptionData получить рыночные данные по опционам
GetOptionData(symbols string) ([]OptionData, error)



```

## Примеры


### Пример получения свечей

```go
// создание клиента

    client := iss.NewClient()

    // получить свечи через создание сервиса
    service := client.NewCandlesService("stock", "shares", "tqbr", "SBER",
        int(iss.Interval_D1),
        "2023-05-01",
        "2025-01-01",
    )
    // запрос данных
    candles, err := service.Do()
    if err != nil {
        slog.Error("main", "ошибка CandlesService", err.Error())
    }

    slog.Info("CandlesService",
        "Symbol", candles.Symbol,
        "Interval", candles.Interval,
        "всего свечей", candles.Len(),
        "mindate", candles.First().Begin,
        "maxdate", candles.Last(0).Begin,
    )
    // цикл по списку свечей
    for row, candle := range candles.Data {
        slog.Info(strconv.Itoa(row),
            "candleData", candle,
        )

    }
```

### Другие примеры смотрите [тут](/example)

---
