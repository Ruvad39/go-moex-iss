package iss

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// OptionInfo параметры инструментов по опционам
// https://iss.moex.com/iss/engines/futures/markets/options/columns.json?iss.only=securities
// 29 полей
//  `db:"id" json:"id" gorm:"primaryKey,autoIncrement"`

type OptionInfo struct {
	SecID                 string  `json:"SECID" csv:"SECID"`                                 // Код инструмента
	BoardID               string  `json:"BOARDID" csv:"BOARDID"`                             // Код режима
	ShortName             string  `json:"SHORTNAME" csv:"SHORTNAME"`                         // Серия срочного инструмента"
	SecName               string  `json:"SECNAME" csv:"SECNAME"`                             // Наименование срочного инструмента
	OptionType            string  `json:"OPTIONTYPE" csv:"OPTIONTYPE"`                       // Вид опциона
	Strike                float64 `json:"STRIKE" csv:"STRIKE"`                               // Цена страйка
	CentralStrike         float64 `json:"CENTRALSTRIKE" csv:"CENTRALSTRIKE"`                 // Центральный страйк
	PrevSettlePrice       float64 `json:"PREVSETTLEPRICE" csv:"PREVSETTLEPRICE"`             // Расчетная цена предыдущего дня, рублей
	Decimals              int     `json:"DECIMALS" csv:"DECIMALS"`                           // Точность
	MinStep               float64 `json:"MINSTEP" csv:""MINSTEP"`                            // Мин. шаг цены
	LastTradeDate         string  `json:"LASTTRADEDATE" csv:"LASTTRADEDATE"`                 // Последний торговый день
	LastDelDate           string  `json:"LASTDELDATE" csv:"LASTDELDATE"`                     // День исполнения
	PrevPrice             float64 `json:"PREVPRICE" csv:"PREVPRICE"`                         // Цена последней сделки предыдущего торгового дня
	StepPrice             float64 `json:"STEPPRICE" csv:"STEPPRICE"`                         // Стоимость шага цены
	LatName               string  `json:"LATNAME" csv:"LATNAME"`                             // Наименование финансового инструмента на английском языке
	ImNp                  float64 `json:"IMNP" csv:"IMNP"`                                   // ГО по непокрытой позиции на первом уровне лимита концентрации
	ImP                   float64 `json:"IMP" csv:"IMP"`                                     // ГО по синтетической позиции на первом уровне лимита концентрации
	ImBuy                 float64 `json:"IMBUY" csv:"IMBUY"`                                 // ГО под покупку опциона на первом уровне лимита концентрации
	ImTime                string  `json:"IMTIME" csv:"IMTIME"`                               // Данные по ГО на (datetime)
	BuySellFee            float64 `json:"BUYSELLFEE" csv:"BUYSELLFEE"`                       // Сбор за регистрацию сделки
	ScalperFee            float64 `json:"SCALPERFEE" csv:"SCALPERFEE"`                       // Сбор за скальперскую сделку
	NegotiatedFee         float64 `json:"NEGOTIATEDFEE" csv:"NEGOTIATEDFEE"`                 // Сбор за адресную сделку
	ExerciseFee           float64 `json:"EXERCISEFEE" csv:"EXERCISEFEE"`                     // Клиринговая комиссия за исполнение контракта
	AssetCode             string  `json:"ASSETCODE" csv:ASSETCODE"`                          // Код базового актива
	UnderlyingAsset       string  `json:"UNDERLYINGASSET" csv:"UNDERLYINGASSET"`             // Базовый актив
	UnderlyingType        string  `json:"UNDERLYINGTYPE" csv:"UNDERLYINGTYPE"`               // Тип базового актива (F - фьючерс, S - акции)
	UnderlyingSettlePrice float64 `json:"UNDERLYINGSETTLEPRICE" csv:"UNDERLYINGSETTLEPRICE"` // Котировка базового актива
}

// OptionData рыночные данные по опционам
// https://iss.moex.com/iss/engines/futures/markets/options/columns.json?iss.only=marketdata
// TODO добавить csv теги
type OptionData struct {
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
	SettleToPrevSettle    float64 `json:"SETTLETOPREVSETTLE"`    // Изменение текущей расчетной цены
	NumTrades             int64   `json:"NUMTRADES"`             // Количество совершенных сделок, штук
	VolToDay              int64   `json:"VOLTODAY"`              // Объем совершенных сделок, контрактов
	ValToDay              float64 `json:"VALTODAY"`              // Объем совершенных сделок, рублей
	ValToDay_USD          float64 `json:"VALTODAY_USD"`          // Объем совершенных сделок, дол. США
	UpdateTime            string  `json:"UPDATETIME"`            // Время последнего обновления
	LastChangePrcnt       float64 `json:"LASTCHANGEPRCNT"`       // Изменение цены последней сделки к предыдущей, %"
	BidDepth              int     `json:"BIDDEPTH"`              // Объем заявок на покупку по лучшей котировке, выраженный в лотах null
	BidDepthT             int     `json:"BIDDEPTHT"`             // Суммарный объем заявок на покупку null
	NumBids               int     `json:"NUMBIDS"`               // Количество заявок на покупку null
	OfferDepth            int     `json:"OFFERDEPTH"`            // Объем заявки на продажу по лучшей котировке null
	OfferDepthT           int     `json:"OFFERDEPTHT"`           // Суммарный объем заявок на продажу null
	NumOffers             int     `json:"NUMOFFERS"`             // Количество заявок на продажу null
	Time                  string  `json:"TIME"`                  // Время заключения последней сделки
	SettleToPrevSettlePrc float64 `json:"SETTLETOPREVSETTLEPRC"` // Изменение текущей расчетной цены относительно
	SeqNum                int64   `json:"SEQNUM"`                // Номер обновления (служебное поле)
	SysTime               string  `json:"SYSTIME"`               // Время загрузки данных системой
	OiChange              int64   `json:"OICHANGE"`              // Изменение открытых позиций к предыдущему закрытию, контр.
	OpenPosition          int64   `json:"OPENPOSITION"`          // Открытые позиции, контрактов
}

type OptionHistory struct {
	TradeDate         string  `json:"TRADEDATE"`         // 	Дата за которую предоставляются данные
	SecID             string  `json:"SECID"`             // Уникальный краткий код инструмента
	BoardID           string  `json:"BOARDID"`           // Идентификатор режима торгов
	Open              float64 `json:"OPEN"`              // Цена открытия
	Low               float64 `json:"LOW"`               // Минимальная цена сделки
	High              float64 `json:"HIGH"`              // Максимальная цена сделки
	Close             float64 `json:"CLOSE"`             // Цена последней сделки
	Value             float64 `json:"VALUE"`             // Оборот в рублях
	Volume            int64   `json:"VOLUME"`            // Оборот в контрактах
	OpenPositionValue float64 `json:"OPENPOSITIONVALUE"` // Объем открытых позиций, руб.
	OpenPosition      int64   `json:"OPENPOSITION"`      // Открытые позиции в контрактах
	SettlePrice       float64 `json:"SETTLEPRICE"`       // Расчетная цена текущего дня
	WфPrice           float64 `json:"WAPRICE"`           // Средневзвешенная цена (рыночные сделки)
	SettlePriceDay    float64 `json:"SETTLEPRICEDAY"`    // Теоретическая цена в дневном клиринге (пункты)
	Change            float64 `json:"CHANGE"`            // Изменение цены последней сделки по отношению к цене последней сделки предыдущего торгового, %
	QTY               int64   `json:"QTY"`               // Объем последней сделки, контрактов
	NumTrades         int64   `json:"NUMTRADES"`         // 	Количество сделок
}

// GetOptionInfo получить параметры инструментов по опционам
func (c *Client) GetOptionInfo(symbols string) ([]OptionInfo, error) {
	var err error
	const op = "GetOptionInfo"
	url := NewIssRequest().Options().Json().MetaData(false).OnlySecurities().Symbols(symbols).URL()
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

	result := make([]OptionInfo, 0, len(resp.Securities.Data))
	err = Unmarshal(resp.Securities.Columns, resp.Securities.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// GetOptionData получить рыночные данные по опционам
func (c *Client) GetOptionData(symbols string) ([]OptionData, error) {
	var err error
	const op = "GetOptionMarketData"

	url := NewIssRequest().Options().Json().MetaData(false).OnlyMarketData().Symbols(symbols).URL()
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

	result := make([]OptionData, 0, len(resp.MarketData.Data))
	err = Unmarshal(resp.MarketData.Columns, resp.MarketData.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// GetOptionHistory получить исторические данные по одному символу
func (c *Client) GetOptionHistory(symbols string, from, to string) ([]OptionHistory, error) {
	return c.NewOptionHistoryService(symbols, from, to, "").Do()
}

// GetOptionHistoryAllDate получить исторические данные по всем символам за заданную дату
func (c *Client) GetOptionHistoryAllDate(date string) ([]OptionHistory, error) {
	return c.NewOptionHistoryService("", "", "", date).Do()
}

// OptionHistoryService сервис для получения исторических данных
type OptionHistoryService struct {
	client     *Client
	issRequest *IssRequest
}

// параметры должны быть
// или symbol + from + to = данные по одному символу
// или  date = данные по всем символам за определенную дату
func (c *Client) NewOptionHistoryService(symbol string, from, to string, date string) *OptionHistoryService {
	iss := NewIssRequest().History().Options().WithSecurities(true).Json().MetaData(false).Target(symbol).From(from).To(to).Date(date)

	return &OptionHistoryService{
		client:     c,
		issRequest: iss,
	}
}

// Do выполняет выгрузку History
func (s *OptionHistoryService) Do() ([]OptionHistory, error) {
	const op = "OptionHistoryService.Do"

	result := make([]OptionHistory, 0)
	count := 1
	for {
		// "fetch candles: item(s) processed"
		s.client.log.Debug(op, "запрос истории: номер запроса", count)

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

func (s *OptionHistoryService) Next() ([]OptionHistory, error) {
	var err error
	const op = "OptionHistoryService.Next"

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

	result := make([]OptionHistory, 0, len(resp.History.Data))
	err = Unmarshal(resp.History.Columns, resp.History.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// TODO поменять на константу (не равную 0. а к примеру 10)
	if len(result) == 0 {
		return nil, EOF
	}
	s.client.log.Debug(op, "len(result)", len(result))

	// увеличим параметр start на кол-во полученных данных
	s.issRequest.start += len(result)
	s.client.log.Debug(op, "s.issRequest.start", s.issRequest.start)

	return result, nil
}
