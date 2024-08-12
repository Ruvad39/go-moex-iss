/*
функционал по аналогии с https://github.com/moexalgo/moexalgo

Methods
info :   Возвращает информацию об инструменте.
marketdata :   Возвращает рыночную информацию и статистику об инструменте.
candles :  Возвращает итератор свечей инструмента по заданным параметрам.
orderbook :  Возвращает текущий стакан лучших цен.

# выбираем акции Сбера
sber = Ticker('SBER')
# получим дневные свечи с 2020 года
sber.candles(start='2020-01-01', end='2023-11-01').head()

engines/stock
*/
package iss

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

var ErrTickerNotFound = errors.New("Ticker not found")
var ErrTickerSymbol = errors.New("код тикера должен быть больше 3-х символов")

//type TickerOption = func(*Ticker)
//
//func WithBoard(board string) TickerOption {
//	return func(t *Ticker) {
//		t.issRequest.boards = board
//	}
//}

type Ticker struct {
	symbol      string
	client      *Client
	issRequest  *IssRequest
	SecID       string
	ShortName   string  // Краткое наименование ценной бумаги (Серия срочного инструмента)
	SecName     string  // Наименование финансового инструмента (срочного инструмента)
	Decimals    int     // Точность, знаков после запятой
	MinStep     float64 // Шаг цены
	SecType     string  // Тип инструмента
	AssetCode   string  // Базовый актив
	LastDelDate string  // День исполнения
}

type TickerInfo struct {
	SecID         string  `json:"SECID"`         // Код инструмента
	BoardID       string  `json:"BOARDID"`       // Код режима
	ShortName     string  `json:"SHORTNAME"`     // наименование ценной бумаги (Серия срочного инструмента)
	SecName       string  `json:"SECNAME"`       // Наименование финансового инструмента (срочного инструмента)
	Decimals      int     `json:"DECIMALS"`      // Точность
	MinStep       float64 `json:"MINSTEP"`       // Мин. шаг цены
	LotVolume     int     `json:"LOTVOLUME"`     // К-во единиц базового актива в инструменте
	LastTradeDate string  `json:"LASTTRADEDATE"` // Последний торговый день
	LastDelDate   string  `json:"LASTDELDATE"`   // День исполнения
	SecType       string  `json:"SECTYPE"`       // Тип инструмента
	AssetCode     string  `json:"ASSETCODE"`     // Код базового актива
	StepPrice     float64 `json:"STEPPRICE"`     // Стоимость шага цены
	PrevPrice     float64 `json:"PREVPRICE"`     // Цена последней сделки предыдущего торгового дня
	FaceValue     float64 `json:"FACEVALUE"`     // Номинал (акции)
	LisTLevel     int     `json:"LISTLEVEL"`     // Уровень листинга (акции)

}

type TickerData struct {
	SecID           string  `json:"SECID"`           // Код инструмента
	BoardID         string  `json:"BOARDID"`         // Код режима
	Bid             float64 `json:"BID"`             // Лучшая котировка на покупку
	BidDepth        float64 `json:"BIDDEPTH"`        // Лотов на покупку по лучшей = null
	Offer           float64 `json:"OFFER"`           // Лучшая котировка на продажу
	OfferDepth      float64 `json:"OFFERDEPTH"`      // Лотов на продажу по лучшей = null
	Spread          float64 `json:"SPREAD"`          // Разница между лучшей котировкой на продажу и покупку (спред), руб
	BidDeptht       int     `json:"biddeptht"`       // Oбъем всех заявок на покупку в очереди Торговой Системы, выраженный в лотах
	OfferDeptht     int     `json:"OFFERDEPTHT"`     // Объем всех заявок на продажу в очереди Торговой Системы, выраженный в лотах
	Open            float64 `json:"OPEN"`            // Цена первой сделки
	Low             float64 `json:"LOW"`             // Минимальная цена сделки
	High            float64 `json:"HIGH"`            // Максимальная цена сделки
	Last            float64 `json:"LAST"`            // Цена последней сделки
	LastChange      float64 `json:"LASTCHANGE"`      // Изменение цены последней сделки к цене предыдущей сделки, рублей
	LastChangePrcnt float64 `json:"LASTCHANGEPRCNT"` // Изменение цены последней сделки к цене предыдущей сделки, %
	QTY             int     `json:"QTY"`             // Объем последней сделки, в лотах
	Value           float64 `json:"VALUE"`           // Объем последней сделки, в руб
	WapRrice        float64 `json:"WAPRICE"`         // Средневзвешенная цен
	NumTrades       int     `json:"NUMTRADES"`       // Количество сделок за торговый день
	VolToDay        int64   `json:"VOLTODAY"`        // Объем совершенных сделок, выраженный в единицах ценных бумаг
	ValToDay        int64   `json:"VALTODAY"`        // Объем совершенных сделок, в валюте расчетов
	OpenPosition    float64 `json:"OPENPOSITION"`    // Открытые позиции, контрактов
	OiChange        int64   `json:"OICHANGE"`        // Изменение открытых позиций к предыдущему закрытию, контр.
	UpdateTime      string  `json:"UPDATETIME"`      // Время последнего обновления
	SysTime         string  `json:"SYSTIME"`         // Время загрузки данных системой
	TradingSession  string  `json:"TRADINGSESSION"`  // Торговая сессия
}

// GetTicker поиск тикера
// func (c *Client) NewTicker(symbol string, opts ...TickerOption) (*Ticker, error) {
func (c *Client) GetTicker(symbol string) (*Ticker, error) {
	// сразу отфильтруем короткий код символа
	if len(symbol) < 3 {
		return nil, ErrTickerSymbol
	}
	iss := NewIssRequest().Json().MetaData(false)
	t := &Ticker{
		symbol:     symbol,
		client:     c,
		issRequest: iss,
	}

	//for _, opt := range opts {
	//	opt(t)
	//}

	//  раздельный поиск по акциям и фьючам
	// поиск среди акций
	exists, err := t.getStock()
	if err != nil {
		return nil, err
	}
	if exists {
		return t, nil
	}
	// поиск среди фьючерсов
	exists, err = t.getForts()
	if err != nil {
		return nil, err
	}
	if exists {
		return t, nil
	}

	// если дошли до сюда = значит НЕ нашли такой тикер
	return t, ErrTickerNotFound
}

// getStock поиск тикера среди акций
func (t *Ticker) getStock() (bool, error) {

	sec, err := t.client.GetStockInfo(t.symbol)
	if err != nil {
		return false, err
	}
	if len(sec) == 1 {
		t.issRequest.securities = true
		t.issRequest.symbol = t.symbol
		t.issRequest.engines = "stock"
		t.issRequest.markets = "shares"
		t.issRequest.boards = StockBoard

		t.SecID = sec[0].SecID
		t.ShortName = sec[0].ShortName
		t.SecName = sec[0].SecName
		t.MinStep = sec[0].MinStep
		t.SecType = sec[0].SecType
		t.Decimals = sec[0].Decimals
		return true, nil
	}
	return false, nil
}

// getForts поиск тикера среди фьючерсов
func (t *Ticker) getForts() (bool, error) {
	// TODO получился не очень красивый код = подумать как сделать рефакторинг
	// поиск по коду (SiU4)
	sec, err := t.client.GetFortsInfo(t.symbol)
	if err != nil {
		return false, err
	}
	if len(sec) == 1 {
		t.issRequest.securities = true
		t.issRequest.symbol = t.symbol
		t.issRequest.engines = "futures"
		t.issRequest.markets = "forts"
		t.issRequest.boards = FortsBoard

		t.SecID = sec[0].SecID
		t.ShortName = sec[0].ShortName
		t.SecName = sec[0].SecName
		t.MinStep = sec[0].MinStep
		t.SecType = sec[0].SecType
		t.Decimals = sec[0].Decimals
		t.AssetCode = sec[0].AssetCode
		t.LastDelDate = sec[0].LastDelDate

		return true, nil
	}
	// поиск по названию "Si-9.24"
	// поиск перебором по списку и поиск по ShortName
	sec, err = t.client.GetFortsInfo("")
	if err != nil {
		return false, err
	}

	for _, _sec := range sec {
		// если нашли = выйдем
		if _sec.ShortName == t.symbol {
			t.issRequest.securities = true
			t.issRequest.symbol = _sec.SecID
			t.issRequest.engines = "futures"
			t.issRequest.markets = "forts"
			t.issRequest.boards = FortsBoard

			t.SecID = _sec.SecID
			t.ShortName = _sec.ShortName
			t.SecName = _sec.SecName
			t.MinStep = _sec.MinStep
			t.SecType = _sec.SecType
			t.Decimals = _sec.Decimals
			t.AssetCode = _sec.AssetCode
			t.LastDelDate = _sec.LastDelDate
			return true, nil
		}

	}

	return false, nil
}

// Info Информация по тикеру
func (t *Ticker) Info() (TickerInfo, error) {
	var err error
	const op = "Ticker.Info"
	result := TickerInfo{}

	url := t.issRequest.OnlySecurities().URL()
	r := &request{
		method:  http.MethodGet,
		fullURL: url,
	}

	var resp Response
	err = t.client.getJSON(r, &resp)
	if err != nil {
		slog.Error(op+".getJSON", "err", err.Error())
		return result, fmt.Errorf("%s: %w", op, err)
	}

	list := make([]TickerInfo, 0, len(resp.Securities.Data))
	err = Unmarshal(resp.Securities.Columns, resp.Securities.Data, &list)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return result, fmt.Errorf("%s: %w", op, err)
	}
	result = list[0]
	return result, nil
}

// Data текущие рыночные данные
func (t *Ticker) Data() (TickerData, error) {
	var err error
	const op = "Ticker.Data"
	result := TickerData{}

	url := t.issRequest.OnlyMarketData().URL()
	r := &request{
		method:  http.MethodGet,
		fullURL: url,
	}

	var resp Response
	err = t.client.getJSON(r, &resp)
	if err != nil {
		slog.Error(op+".getJSON", "err", err.Error())
		return result, fmt.Errorf("%s: %w", op, err)
	}

	list := make([]TickerData, 0, len(resp.MarketData.Data))
	err = Unmarshal(resp.MarketData.Columns, resp.MarketData.Data, &list)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return result, fmt.Errorf("%s: %w", op, err)
	}
	result = list[0]
	return result, nil
}

// Candles получим исторические свечи
func (t *Ticker) Candles(interval int, from, to string) (Candles, error) {
	//return t.client.NewCandlesService("stock", "shares", StockBoard, symbols, int(interval), from, to).Do()
	return t.client.NewCandlesService(t.issRequest.engines, t.issRequest.markets, t.issRequest.boards, t.issRequest.symbol, int(interval), from, to).Do()

}

// OrderBook Возвращает текущий стакан лучших цен
// нужна авторизация
func (t *Ticker) OrderBook() (OrderBook, error) {
	return t.client.NewOrderBookService(t.issRequest.engines, t.issRequest.markets, t.issRequest.boards, t.issRequest.symbol).Do()
}
