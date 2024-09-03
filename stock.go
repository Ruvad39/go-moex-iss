package iss

import (
	"fmt"
	"log/slog"
	"net/http"
)

// StockInfo параметры акций
// https://iss.moex.com/iss/engines/stock/markets/shares/securities/columns.json?iss.only=securities
type StockInfo struct {
	SecID               string  `json:"SECID"`               // Код инструмента
	BoardID             string  `json:"BOARDID"`             // Код режима
	ShortName           string  `json:"SHORTNAME"`           // Кратк. наим.
	PrevPrice           float64 `json:"PREVPRICE"`           // Последняя предыдущего дня""Цена последней сделки нормального периода предыдущего торгового дня
	LotSize             int     `json:"LOTSIZE"`             // Размер лота
	FaceValue           float64 `json:"FACEVALUE"`           // Номинал
	Status              string  `json:"STATUS"`              // Статус
	BoardName           string  `json:"BOARDNAME"`           // Режим торгов
	Decimals            int     `json:"DECIMALS"`            // Точность
	SecName             string  `json:"SECNAME"`             // Наименование
	Remarks             string  `json:"REMARKS"`             // Примечание
	MarketCode          string  `json:"MARKETCODE"`          // Рынок, Идентификатор рынка на котором торгуется финансовый инструмент
	InstrID             string  `json:"INSTRID"`             // Группа инструментов
	SectorID            string  `json:"SectorID"`            // Сектор (Устарело)"
	MinStep             float64 `json:"MINSTEP"`             // Мин. шаг цены
	PrevWaPrice         float64 `json:"PREVWAPRICE"`         // Значение оценки (WAPRICE) предыдущего торгового дня
	FaceUnit            string  `json:"FACEUNIT"`            // Код валюты, в которой выражен номинал ценной бумаги
	PrevDate            string  `json:"PREVDATE"`            // Дата предыдущего торгового дня
	IssueSize           int     `json:"ISSUESIZE"`           // Объем выпуска
	ISIN                string  `json:"ISIN"`                // Международный идентификационный код ценной бумаги
	LatName             string  `json:"LATNAME"`             // Наименование финансового инструмента на английском языке
	RegNumber           string  `json:"REGNUMBER"`           // Регистрационный номер
	PrevLegalClosePrice float64 `json:"PREVLEGALCLOSEPRICE"` // Официальная цена закрытия предыдущего дня, рассчитываемая по методике ФСФР
	CurrencyID          string  `json:"CURRENCYID"`          // Валюта расчетов
	SecType             string  `json:"SECTYPE"`             // Тип ценной бумаги
	LisTLevel           int     `json:"LISTLEVEL"`           // Уровень листинга
	SettLeDate          string  `json:"SETTLEDATE"`          // Дата расчетов сделки
}

// StockData рыночные данные по акциям
// https://iss.moex.com/iss/engines/stock/markets/shares/securities/columns.json?iss.only=marketdata
type StockData struct {
	SecID                          string  `json:"SECID"`                          // Код инструмента
	BoardID                        string  `json:"BOARDID"`                        // Код режима
	Bid                            float64 `json:"BID"`                            // Лучшая котировка на покупку
	BidDepth                       string  `json:"BIDDEPTH"`                       // Лотов на покупку по лучшей  = null
	Offer                          float64 `json:"OFFER"`                          // Лучшая котировка на продажу
	OfferDepth                     string  `json:"OFFERDEPTH"`                     // Лотов на продажу по лучшей  = null
	Spread                         float64 `json:"SPREAD"`                         // Разница между лучшей котировкой на продажу и покупку (спред), руб
	BidDeptht                      int     `json:"BIDDEPTHT"`                      // объем всех заявок на покупку в очереди Торговой Системы, выраженный в лотах
	OfferDeptht                    int     `json:"OFFERDEPTHT"`                    // Объем всех заявок на продажу в очереди Торговой Системы, выраженный в лотах
	Open                           float64 `json:"OPEN"`                           // Цена первой сделки
	Low                            float64 `json:"LOW"`                            // Минимальная цена сделки
	High                           float64 `json:"HIGH"`                           // Максимальная цена сделки
	Last                           float64 `json:"LAST"`                           // Цена последней сделки
	LastChange                     float64 `json:"LASTCHANGE"`                     // Изменение цены последней сделки к цене предыдущей сделки, рублей
	LastChangePrcnt                float64 `json:"LASTCHANGEPRCNT"`                // Изменение цены последней сделки к цене предыдущей сделки, %
	QTY                            int     `json:"QTY"`                            // Объем последней сделки, в лотах
	Value                          float64 `json:"VALUE"`                          // Объем последней сделки, в руб
	Value_USD                      float64 `json:"VALUE_USD"`                      // Объем последней сделки, дол. США
	WapRrice                       float64 `json:"WAPRICE"`                        // Средневзвешенная цен
	LastCNGToLastWaPrice           float64 `json:"LASTCNGTOLASTWAPRICE"`           // Изменение цены последней сделки к средневзвешенной цене, рублей
	WapToPrevWaPricePrcnt          float64 `json:"WAPTOPREVWAPRICEPRCNT"`          // Изменение средневзвешенной цены относительно средневзвешенной цены предыдущего торгового дня, %
	WapToPrevWaPrice               float64 `json:"WAPTOPREVWAPRICE"`               // Изменение средневзвешенной цены к средневзвешенной цене предыдущего торгового дня, рублей
	ClosePrice                     float64 `json:"CLOSEPRICE"`                     // Цена послеторгового периода
	MarketPriceToday               float64 `json:"MARKETPRICETODAY"`               // Рыночная цена по результатам торгов сегодняшнего дня, за одну ценную бумагу
	MarkeTPrice                    float64 `json:"MARKETPRICE"`                    // Рыночная цена предыдущего дня
	LastToPrevPrice                float64 `json:"LASTTOPREVPRICE"`                // Изменение цены последней сделки к последней цене предыдущего дня, %
	NumTrades                      int     `json:"NUMTRADES"`                      // Количество сделок за торговый день
	VolToDay                       int64   `json:"VOLTODAY"`                       // Объем совершенных сделок, выраженный в единицах ценных бумаг
	ValToDay                       int64   `json:"VALTODAY"`                       // Объем совершенных сделок, в валюте расчетов
	ValToDay_USD                   int64   `json:"VALTODAY_USD"`                   // Объем заключенных сделок, дол. США
	ETFSETTLEPRICE                 float64 `json:"ETFSETTLEPRICE"`                 // Расчетная стоимость акции\/пая иностранного биржевого инвестиционного фонда
	TradingStatus                  string  `json:"TRADINGSTATUS"`                  // Индикатор состояния торговой сессии по инструменту
	UpdateTime                     string  `json:"UPDATETIME"`                     // Время последнего обновления
	LastBid                        float64 `json:"LASTBID"`                        // Лучшая котировка на покупку на момент завершения нормального периода торгов
	LastOffer                      float64 `json:"LASTOFFER"`                      // Лучшая котировка на продажу на момент завершения нормального периода торгов
	LClosePrice                    float64 `json:"LCLOSEPRICE"`                    // Официальная цена закрытия, рассчитываемая по методике ФСФР как средневзвешенная цена сделок за последние 10 минут торговой сессии
	LCurrentPrice                  float64 `json:"LCURRENTPRICE"`                  // Официальная текущая цена, рассчитываемая как средневзвешенная цена сделок заключенных за последние 10 минут
	MarketPrice2                   float64 `json:"MARKETPRICE2"`                   // Рыночная цена 2, рассчитываемая в соответствии с методикой ФСФР
	NumBids                        int     `json:"NUMBIDS"`                        // Количество заявок на покупку в очереди Торговой системы  (undefined)
	NumOffers                      int     `json:"NUMOFFERS"`                      // Количество заявок на продажу в очереди Торговой системы  (undefined)
	Change                         float64 `json:"CHANGE"`                         // Изменение цены последней сделки по отношению к цене последней сделки предыдущего торгового дня
	Time                           string  `json:"TIME"`                           // Время заключения последней сделки
	HighBid                        float64 `json:"HIGHBID"`                        // Наибольшая цена спроса в течение торговой сессии
	LowOffer                       float64 `json:"LOWOFFER"`                       // Наименьшая цена предложения в течение торговой сессии
	PriceMinusPrevWapRice          float64 `json:"PRICEMINUSPREVWAPRICE"`          // Цена последней сделки к оценке предыдущего дня
	OpenPeriodPrice                float64 `json:"OPENPERIODPRICE"`                // Цена предторгового периода
	SeqNum                         int64   `json:"SEQNUM"`                         // номер обновления (служебное поле)
	SysTime                        string  `json:"SYSTIME"`                        // Время загрузки данных системой
	ClosingAuctionPrice            float64 `json:"CLOSINGAUCTIONPRICE"`            // Цена послеторгового аукциона
	ClosingAuctionVolume           float64 `json:"CLOSINGAUCTIONVOLUME"`           // Количество в сделках послеторгового аукциона
	ISSUECapitalization            float64 `json:"ISSUECAPITALIZATION"`            // Текущая капитализация акции
	ETFSETTLECurrency              string  `json:"ETFSETTLECURRENCY"`              // Валюта расчетной стоимости акции\/пая иностранного биржевого инвестиционного фонда
	ValToday_RUR                   int64   `json:"VALTODAY_RUR"`                   // Объем совершенных сделок, рублей
	TradingSession                 string  `json:"TRADINGSESSION"`                 // Торговая сессия
	TrendISSUECapitalization       float64 `json:"TRENDISSUECAPITALIZATION"`       // Изменение капитализации к капитализации предыдущего дня
	ISSUECapitalization_UpdateTime string  `json:"ISSUECAPITALIZATION_UPDATETIME"` // Время обновления капитализации
}

// GetStockInfo получить параметры инструментов по акциям
func (c *Client) GetStockInfo(symbols string) ([]StockInfo, error) {
	var err error
	const op = "GetStockInfo"

	url := NewIssRequest().Stock().Json().MetaData(false).OnlySecurities().Symbols(symbols).URL()
	r := &request{
		method:  http.MethodGet,
		fullURL: url,
	}

	var resp Response
	err = c.getJSON(r, &resp)
	if err != nil {
		slog.Error(op+".getJSON", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]StockInfo, 0, len(resp.Securities.Data))
	err = Unmarshal(resp.Securities.Columns, resp.Securities.Data, &result)
	//err = UnmarshalCSV(resp.Securities.Columns, resp.Securities.Data, &result, DefaultTagKey)

	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// GetStockData получить рыночные данные по акциям
func (c *Client) GetStockData(symbols string) ([]StockData, error) {
	var err error
	const op = "GetStockData"

	url := NewIssRequest().Stock().Json().MetaData(false).OnlyMarketData().Symbols(symbols).URL()
	r := &request{
		method:  http.MethodGet,
		fullURL: url,
	}

	var resp Response
	err = c.getJSON(r, &resp)
	if err != nil {
		slog.Error(op+".getJSON", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]StockData, 0, len(resp.MarketData.Data))
	err = Unmarshal(resp.MarketData.Columns, resp.MarketData.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
