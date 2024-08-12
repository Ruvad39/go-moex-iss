# Библиотека, которая позволяет работать с функционалом [IIS Московской Биржи](https://www.moex.com/a2193)
без авторизации

## Установка

```bash
go get github.com/Ruvad39/go-moex-iss
```

## api который реализован на текущий момент:

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
//GetBondsData(board string)

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

// GetTicker поиск тикера по коду
GetTicker(symbol string) (*Ticker, error)
// Info Информация по тикеру
Ticker.Info() (TickerInfo, error)
// Data текущая рыночная информация по тикеру
Ticker.Data() (TickerData, error)
// Candles исторические свечи по тикеру 
Ticker.Candles(interval int, from, to string) (Candles, error) 
// OrderBook получить стакан. 
// Нужна аторизация
Ticker.OrderBook() (OrderBook, error)

// algopack 

// GetFutOIAll Открытые позиции физ. и юр. лиц по всем инструментам
// date = за указанную дату; latest =1 возвращает последнюю пятиминутку за указанную дату
GetFutOIAll(date string, latest int) ([]FutOI, error)
// GetFutOI данные по заданному тикеру
// ticker = Краткий код базового актива (Si, RI, GD, ...); from = Дата начала периода; to = Дата окончания периода; флаг latest=1 возвращает последнюю пятиминутку за указанный период
GetFutOI(ticker string, from, to string, latest int) ([]FutOI, error)

// TODO другие данные algopack https://moexalgo.github.io
```

## Примеры

### создание клиента
```go
// без авторизации, задержка данных по времени 15 минут.
// и не доступны некоторые сервисы algopack
user := ""
pwd := ""
err, client := iss.NewClient(iss.WithUser(user), iss.WithPwd(pwd))
//err, client := iss.NewClient()
if err != nil {
    slog.Error("main", "NewClient", err.Error())
}
```
### Данные по акциям

```go
err, client := iss.NewClient()
if err != nil {
    slog.Error("main", "NewClient", err.Error())
}
// получить список акций
Sec, err := client.GetStockInfo("")
// можно указать список через запятую
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
// получить рыночные данные по акции
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
// исторические свечи по акции
candles, err := client.GetStockCandles("SBER", iss.Interval_D1, "2024-05-01", "2025-01-01")
if err != nil {
    slog.Error("main", "ошибка GetCandles", err.Error())
	return
}

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
	
```
### Данные по тикеру
```go
// создание (поиск) тикера
//ticker, err := client.GetTicker("SBER") 
ticker, err := client.GetTicker("RTS-9.24")
if err != nil {
	slog.Error("main", "ошибка NewTicker", err.Error())
}

// параметры инструмента
info, err := ticker.Info()
if err != nil {
    slog.Error("main", "ошибка ticker.Info", err.Error())
}
slog.Info("ticker.Info", slog.Any("t_info", info))

// текущие рыночные данные
data, err := ticker.Data()
if err != nil {
    slog.Error("main", "ошибка ticker.data", err.Error())
}
slog.Info("ticker.Info", slog.Any("t_data", data))

// свечи
candles, err := ticker.Candles(iss.Interval_D1, "2024-07-01", "2025-01-01")
if err != nil {
    slog.Error("main", "ошибка Candles", err.Error())
    return
}
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
// стакан. Нужна авторизация
orderBook, err := ticker.OrderBook()
if err != nil {
    slog.Error("main", "ticker.OrderBook", err.Error())
return
}
bid, _ := orderBook.BestBid()
ask, _ := orderBook.BestAsk()
bidVolume := orderBook.Bids.SumDepth()
askVolume := orderBook.Asks.SumDepth()

slog.Info("orderBook", "BestAsk", ask.Price, "BestBid", bid.Price, "объем асков", askVolume, "объем бидов", bidVolume)
fmt.Println(orderBook.String())


```


### Другие примеры смотрите [тут](/example)


