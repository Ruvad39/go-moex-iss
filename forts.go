package iss

import (
	"fmt"
	"log/slog"
	"net/http"
)

// FortsInfo параметры инструментов по рынку фортс
// https://iss.moex.com/iss/engines/futures/markets/forts/columns.json?iss.only=securities
// TODO добавить csv таг
type FortsInfo struct {
	SecID            string  `json:"SECID"`            // Код инструмента
	BoardID          string  `json:"BOARDID"`          // Код режима
	ShortName        string  `json:"SHORTNAME"`        // Кратк. наим.
	SecName          string  `json:"SECNAME"`          // Наименование срочного инструмента
	PrevSetTlePrice  float64 `json:"PREVSETTLEPRICE"`  // Расчетная цена предыдущего дня, рублей
	Decimals         int     `json:"DECIMALS"`         // Точность
	MinStep          float64 `json:"MINSTEP"`          // Мин. шаг цены
	LastTradeDate    string  `json:"LASTTRADEDATE"`    // Последний торговый день
	LastDelDate      string  `json:"LASTDELDATE"`      // День исполнения
	SecType          string  `json:"SECTYPE"`          // Тип инструмента
	LatName          string  `json:"LATNAME"`          // Наименование финансового инструмента на английском языке
	AssetCode        string  `json:"ASSETCODE"`        // Код базового актива
	PrevOpenPosition int     `json:"PREVOPENPOSITION"` // Открытые позиции предыдущего дня, контр.
	LotVolume        int     `json:"LOTVOLUME"`        // К-во единиц базового актива в инструменте
	InitialMargin    float64 `json:"INITIALMARGIN"`    // Гарантийное обеспечение на первом уровне лимита концентрации
	HighLimit        float64 `json:"HIGHLIMIT"`        // Верхний лимит
	LowLimit         float64 `json:"LOWLIMIT"`         // Нижний лимит
	StepPrice        float64 `json:"STEPPRICE"`        // Стоимость шага цены
	LastSettlePrice  float64 `json:"LASTSETTLEPRICE"`  // Расчетная цена последнего клиринга
	PrevPrice        float64 `json:"PREVPRICE"`        // Цена последней сделки предыдущего торгового дня
	IMTime           string  `json:"IMTIME"`           // Данные по ГО на
	BuySellFee       float64 `json:"BUYSELLFEE"`       // Сбор за регистрацию сделки*, руб.
	ScalPerFee       float64 `json:"SCALPERFEE"`       // Сбор за скальперскую сделку*, руб.
	NegotiatedFee    float64 `json:"NEGOTIATEDFEE"`    // Сбор за адресную сделку*, руб.
	ExerciseFee      float64 `json:"EXERCISEFEE"`      // Клиринговая комиссия за исполнение контракта*, руб.

}

// FortsData рыночные данные по инструментам фортс
// https://iss.moex.com/iss/engines/futures/markets/forts/columns.json?iss.only=marketdata
// TODO добавить csv таг
type FortsData struct {
	SecID                 string  `json:"SECID"`                 // Код инструмента
	BoardID               string  `json:"BOARDID"`               // Код режима
	Bid                   float64 `json:"BID"`                   // Лучшая котировка на покупку
	Offer                 float64 `json:"OFFER"`                 // Лучшая котировка на продажу
	Spread                float64 `json:"SPREAD"`                // Разница между лучшей котировкой на продажу и покупку (спред), руб
	Open                  float64 `json:"OPEN"`                  // Цена первой сделки
	Low                   float64 `json:"LOW"`                   // Минимальная цена сделки
	High                  float64 `json:"HIGH"`                  // Максимальная цена сделки
	Last                  float64 `json:"LAST"`                  // Цена последней сделки
	Quantity              int     `json:"QUANTITY"`              // Объем последней сделки, контрактов
	LastChange            float64 `json:"LASTCHANGE"`            // Изменение цены последней сделки к предыдущей цене
	SettlePrice           float64 `json:"SETTLEPRICE"`           // Текущая расчетная цена
	SettleTopRevSettle    float64 `json:"SETTLETOPREVSETTLE"`    // Изменение текущей расчетной цены
	NumTrades             int     `json:"NUMTRADES"`             // Количество совершенных сделок, штук
	VolToDay              int64   `json:"VOLTODAY"`              // Объем совершенных сделок, контрактов
	ValToDay              int64   `json:"VALTODAY"`              // Объем совершенных сделок, рублей
	ValToDay_USD          int64   `json:"VALTODAY_USD"`          // Объем совершенных сделок, дол. США
	UpdateTime            string  `json:"UPDATETIME"`            // Время последнего обновления
	LastChangePrcnt       float64 `json:"LASTCHANGEPRCNT"`       // Изменение цены последней сделки к предыдущей, %"
	BidDepth              int     `json:"BIDDEPTH"`              // Объем заявок на покупку по лучшей котировке, выраженный в лотах null
	BidDepthT             int     `json:"BIDDEPTHT"`             // Суммарный объем заявок на покупку null
	NumBids               int     `json:"NUMBIDS"`               // Количество заявок на покупку null
	OfferDepth            int     `json:"OFFERDEPTH"`            // Объем заявки на продажу по лучшей котировке null
	OfferDepthT           int     `json:"OFFERDEPTHT"`           // Суммарный объем заявок на продажу null
	NumOffers             int     `json:"NUMOFFERS"`             // Количество заявок на продажу null
	Time                  string  `json:"TIME"`                  // Время заключения последней сделки
	SETTLETOPREVSETTLEPRC float64 `json:"SETTLETOPREVSETTLEPRC"` // Изменение текущей расчетной цены относительно расчетной цены предыдущего торгового дня, %
	SEQNUM                int64   `json:"SEQNUM"`                // Номер обновления (служебное поле)
	SysTime               string  `json:"SYSTIME"`               // Время загрузки данных системой
	TradeDate             string  `json:"TRADEDATE"`             // Дата последней сделки
	LastToPrevPrice       float64 `json:"LASTTOPREVPRICE"`       // Изменение цены последней сделки к последней цене предыдущего дня, %
	OpenPosition          float64 `json:"OPENPOSITION"`          // Открытые позиции, контрактов
	OiChange              int64   `json:"OICHANGE"`              // Изменение открытых позиций к предыдущему закрытию, контр.
	OpenPeriodPrice       float64 `json:"OPENPERIODPRICE"`       // Цена аукциона открытия
	SwapRate              float64 `json:"SWAPRATE"`              // Фандинг в рублях (величина SwapRate, согласно спецификации контракта)
}

// GetFortsInfo получить параметры инструментов по фьючерсам
func (c *Client) GetFortsInfo(symbols string) ([]FortsInfo, error) {
	var err error
	const op = "GetFortsInfo"

	url := NewIssRequest().Forts().Json().MetaData(false).OnlySecurities().Symbols(symbols).URL()
	r := &request{
		method:  http.MethodGet,
		fullURL: url,
	}
	var resp Response
	err = c.getJSON(r, &resp)
	if err != nil {
		slog.Error("GetFortsInfo.getJSON", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]FortsInfo, 0, len(resp.Securities.Data))
	err = Unmarshal(resp.Securities.Columns, resp.Securities.Data, &result)
	if err != nil {
		slog.Error("GetFortsInfo.Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// GetFortsData получить рыночные данные по фьючерсам
func (c *Client) GetFortsData(symbols string) ([]FortsData, error) {
	var err error
	const op = "GetFortsData"

	url := NewIssRequest().Forts().Json().MetaData(false).OnlyMarketData().Symbols(symbols).URL()
	r := &request{
		method:  http.MethodGet,
		fullURL: url,
	}
	var resp Response
	err = c.getJSON(r, &resp)
	if err != nil {
		slog.Error("GetFortsData.getJSON", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]FortsData, 0, len(resp.MarketData.Data))
	err = Unmarshal(resp.MarketData.Columns, resp.MarketData.Data, &result)
	if err != nil {
		slog.Error("GetFortsData.Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
