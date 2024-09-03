package iss

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

/*
https://moexalgo.github.io/des/supercandles/

// все за заданный день
https://iss.moex.com/iss/datashop/algopack/eq/tradestats.json?date=2024-04-17
// заданная акция за период
https://iss.moex.com/iss/datashop/algopack/eq/tradestats/SBER.json?from=2024-08-30&till=2024-09-03&latest=1


datashop/algopack/eq/tradestats
/datashop/algopack/eq/obstats

eq = акции
fo = фьючерсы
fx = валюта

*/

// TradeStats
type TradeStats struct {
	TradeDate  string  `csv:"tradedate" json:"tradedate"`   // дата сделки
	TradeTime  string  `csv:"tradetime" json:"tradetime"`   // время сделки
	SecID      string  `csv:"secid" json:"secid"`           // код инструмента
	AssetCode  string  `csv:"asset_code" json:"asset_code"` // Код базового актива (для фьючерсов)
	Open       float64 `csv:"pr_open" json:"pr_open"`       // цена открытия
	High       float64 `csv:"pr_high" json:"pr_high"`       // максимальная цена за период
	Low        float64 `csv:"pr_low" json:"pr_low"`         // минимальная цена за период
	Close      float64 `csv:"pr_close" json:"pr_close"`     // последняя цена за период
	Std        float64 `csv:"pr_std" json:"pr_std"`         // стандартное отклонение цены
	Volume     int64   `csv:"vol" json:"vol"`               // объем в лотах
	Value      float64 `csv:"val" json:"val"`               // объем в рублях
	Trades     int64   `csv:"trades" json:"trades"`         // количество сделок
	Vwap       float64 `csv:"pr_vwap" json:pr_vwap"`        // взвешенная средняя цена
	Change     float64 `csv:"pr_change" json:"pr_change"`   // изменение цены за период, %
	TradesBuy  int64   `csv:"trades_b" json:"trades_b"`     // кол-во сделок на покупку
	TradesSell int64   `csv:"trades_s" json:"trades_s"`     // кол-во сделок на продажу
	ValueBuy   float64 `csv:"val_b" json:"val_b"`           // объем покупок в рублях
	ValueSell  float64 `csv:"val_s" json:"val_s"`           // объем продаж в рублях
	VolumeBuy  int64   `csv:"vol_b" json:"vol_b"`           // объем покупок в лотах
	VolumeSell int64   `csv:"vol_s" json:"vol_s"`           // объем продаж в лотах
	Disb       float64 `csv:"disb" json:"disb"`             // соотношение объема покупок и продаж
	VwapBuy    float64 `csv:"pr_vwap_b" json:"pr_vwap_b"`   // средневзвешенная цена покупки
	VwapSell   float64 `csv:"pr_vwap_s" json:"pr_vwap_s"`   // средневзвешенная цена продажи
	OiOpen     int64   `csv:"oi_open" json:"oi_open"`       // ОИ на открытии (для фьючерсов)
	OiHigh     int64   `csv:"oi_high" json:"oi_high"`       // максимальный ОИ (для фьючерсов)
	OiLow      int64   `csv:"oi_low" json:"oi_low"`         // минимальный ОИ (для фьючерсов)
	OiClose    int64   `csv:"oi_close" json:"oi_close"`     // ОИ на закрытии (для фьючерсов)
	SYSTIME    string  `csv:"SYSTIME" json:"SYSTIME"`       // время системы
}

// TradeStatsService сервис для получения супер свечей (TradeStats)
type TradeStatsService struct {
	client     *Client
	issRequest *IssRequest
}

// NewTradeStatsService создание сервиса
// или все за заданную дату symbol == "" + указана date
// или по одному символу за период symbol != "" + указаны from, to (если не указаны = то за текущий день)
func (c *Client) NewTradeStatsService(markets, symbol string, from, to string, date string, latest bool) *TradeStatsService {
	// eq = акции
	iss := NewIssRequest().
		AlgoPackMarkets(markets).
		AlgoPack("tradestats").
		Target(symbol).
		From(from).
		To(to).
		Date(date).
		Latest(latest).
		Json().MetaData(false)

	return &TradeStatsService{
		client:     c,
		issRequest: iss,
	}
}

// Next загружает следующую страницу данных
// Если данных больше нет, то возвращается ошибка EOF
// TODO что возвращать данные или ссылку?
func (s *TradeStatsService) Next() ([]TradeStats, error) {
	var err error
	const op = "TradeStatsService.Next"

	r := &request{
		method:  http.MethodGet,
		fullURL: s.issRequest.URL(),
	}

	var resp Response
	err = s.client.getJSON(r, &resp)
	if err != nil {
		slog.Error(op+".getJSON", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]TradeStats, 0, len(resp.Data.Data))
	err = Unmarshal(resp.Data.Columns, resp.Data.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// TODO поменять на константу (не равную 0. а к примеру 10)
	if len(result) == 0 {
		//if len(result) < 10 {
		return result, EOF
	}
	//s.client.log.Debug(op,
	//	"len(result)", len(result),
	//	"mindate", result[0].Begin,
	//	"maxdate", result[len(result)-1].Begin,
	//)

	// увеличим параметр start на кол-во полученных данных
	s.issRequest.start += len(result)
	s.client.log.Debug(op, "s.issRequest.start", s.issRequest.start)

	return result, nil
}

// Do выполняет выгрузку свечей
func (s *TradeStatsService) Do() ([]TradeStats, error) {
	const op = "TradeStatsService.Do"

	result := make([]TradeStats, 0)
	count := 1
	for {
		// "fetch candles: item(s) processed"
		s.client.log.Debug(op, "запрос свечей: номер запроса", count)

		t_result, err := s.Next()
		if err != nil {
			if errors.Is(err, EOF) {
				break
			}
			return result, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, t_result...)
		count++

	}

	return result, nil
}

// GetStockTradeStats получим данные TradeStats по заданной акции
func (c *Client) GetStockTradeStats(symbol string, from, to string, latest bool) ([]TradeStats, error) {
	service := c.NewTradeStatsService(AlgoPackStock, symbol, from, to, "", latest)
	return service.Do()
}

// GetStockTradeStatsAll получим данные TradeStats по всем акция за заданный день
func (c *Client) GetStockTradeStatsAll(date string, latest bool) ([]TradeStats, error) {
	service := c.NewTradeStatsService(AlgoPackStock, "", "", "", date, latest)
	return service.Do()
}

// GetFortsTradeStats получим данные TradeStats по заданному фьючерсу
func (c *Client) GetFortsTradeStats(symbol string, from, to string, latest bool) ([]TradeStats, error) {
	service := c.NewTradeStatsService(AlgoPackForts, symbol, from, to, "", latest)
	return service.Do()
}

// GetStockTradeStatsAll получим данные TradeStats по всем фьючерсам за заданный день
func (c *Client) GetFortsTradeStatsAll(date string, latest bool) ([]TradeStats, error) {
	service := c.NewTradeStatsService(AlgoPackForts, "", "", "", date, latest)
	return service.Do()
}

// GetFxTradeStats получим данные TradeStats по заданной валюте
func (c *Client) GetFxTradeStats(symbol string, from, to string, latest bool) ([]TradeStats, error) {
	service := c.NewTradeStatsService(AlgoPackFx, symbol, from, to, "", latest)
	return service.Do()
}

// GetFxTradeStatsAll получим данные TradeStats по всем валютам за заданный день
func (c *Client) GetFxTradeStatsAll(date string, latest bool) ([]TradeStats, error) {
	service := c.NewTradeStatsService(AlgoPackFx, "", "", "", date, latest)
	return service.Do()
}
