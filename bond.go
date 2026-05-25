package iss

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// TODO рыыночные данные
// https://iss.moex.com/iss/engines/stock/markets/bonds/columns.json?iss.only=marketdata

// BondInfo
// https://iss.moex.com/iss/engines/stock/markets/bonds/columns.json?iss.only=securities
type BondInfo struct {
	SecID                 string  `json:"SECID"`                 // Код инструмента
	BoardID               string  `json:"BOARDID"`               // Код режима
	ShortName             string  `json:"SHORTNAME"`             // Кратк. наим.
	SecName               string  `json:"SECNAME"`               // Наименование финансового инструмента
	PrevWaPrice           float64 `json:"PREVWAPPRICE"`          // Средневзвешенная цена предыдущего дня, % к номиналу
	YielDatPrevWaPrice    float64 `json:"YIELDATPREVWAPRICE"`    // Доходность по оценке пред. дня
	CouponValue           float64 `json:"COUPONVALUE"`           // Сумма купона, в валюте номинала
	NextCoupon            string  `json:"NEXTCOUPON"`            // Дата окончания купона
	MatDate               string  `json:"MATDATE"`               // Дата погашения
	Accruedint            float64 `json:"ACCRUEDINT"`            // НКД на дату расчетов, в валюте расчетов
	PrevPrice             float64 `json:"PREVPRICE"`             // Цена последней сделки пред. дня, % к номиналу
	LotSize               float64 `json:"LOTSIZE"`               // Размер лота
	FaceValue             float64 `json:"FACEVALUE"`             // Непогашенный долг
	Status                string  `json:"STATUS"`                // Статус
	Decimals              int     `json:"DECIMALS"`              // Точность, знаков после запятой
	CouponPeriod          int     `json:"COUPONPERIOD"`          // Длительность купона
	IssueSize             int64   `json:"ISSUESIZE"`             // Объем выпуска, штук
	PrevLegalClosePrice   float64 `json:"PREVLEGALCLOSEPRICE"`   // Официальная цена закрытия предыдущего дня
	PrevDate              string  `json:"PREVDATE"`              // Дата предыдущего торгового дня
	Remarks               string  `json:"REMARKS"`               // Примечание
	MarketCode            string  `json:"MARKETCODE"`            // Рынок
	InstrID               string  `json:"INSTRID"`               // Группа инструментов
	MinStep               float64 `json:"MINSTEP"`               // Мин. шаг цены
	FaceUnit              string  `json:"FACEUNIT"`              // Валюта номинала
	BuyBackPrice          float64 `json:"BUYBACKPRICE"`          // Цена оферты
	BuyBackDate           string  `json:"BUYBACKDATE"`           // Дата, к которой рассчитывается доходность (если данное поле не заполнено, то \"Доходность посл.сделки\" рассчитывается к Дате погашения
	ISIN                  string  `json:"ISIN"`                  // ISIN
	LatName               string  `json:"LATNAME"`               // Англ. наименование
	RegNumber             string  `json:"REGNUMBER"`             // Регистрационный номер
	CurrencyID            string  `json:"CURRENCYID"`            // Сопр. валюта инструмента
	IssueSizePlaced       float64 `json:"ISSUESIZEPLACED"`       // Количество ценных бумаг в обращении
	LisTLevel             int     `json:"LISTLEVEL"`             // Уровень листинга
	SecType               string  `json:"SECTYPE"`               // Тип ценной бумаги
	CouponPercent         float64 `json:"COUPONPERCENT"`         // Ставка купона, %"
	OfferDate             string  `json:"OFFERDATE"`             // Дата Оферты
	SettleDate            string  `json:"SETTLEDATE"`            // Дата расчетов сделки
	LotValue              float64 `json:"LOTVALUE"`              // Номинальная стоимость лота, в валюте номинала
	FaceValueOnSettleDate float64 `json:"FACEVALUEONSETTLEDATE"` // Номинальная стоимость на дату расчетов (значение, известное на дату заключения сделки)
}

// marketdata
// https://iss.moex.com/iss/engines/stock/markets/bonds/boardgroups/58/securities/columns.json?iss.meta=off&iss.only=marketdata

type BondData struct {
	SecID           string  `json:"SECID"`           // Код инструмента
	BoardID         string  `json:"BOARDID"`         // Код режима
	Bid             float64 `json:"BID"`             // Лучшая котировка на покупку
	BidDepth        string  `json:"BIDDEPTH"`        // Лотов на покупку по лучшей  = null
	Offer           float64 `json:"OFFER"`           // Лучшая котировка на продажу
	OfferDepth      string  `json:"OFFERDEPTH"`      // Лотов на продажу по лучшей  = null
	Spread          float64 `json:"SPREAD"`          // Разница между лучшей котировкой на продажу и покупку (спред), руб
	BidDeptht       int     `json:"BIDDEPTHT"`       // объем всех заявок на покупку в очереди Торговой Системы, выраженный в лотах
	OfferDeptht     int     `json:"OFFERDEPTHT"`     // Объем всех заявок на продажу в очереди Торговой Системы, выраженный в лотах
	Open            float64 `json:"OPEN"`            // Цена первой сделки
	Low             float64 `json:"LOW"`             // Минимальная цена сделки
	High            float64 `json:"HIGH"`            // Максимальная цена сделки
	Last            float64 `json:"LAST"`            // Цена последней сделки
	LastChange      float64 `json:"LASTCHANGE"`      // Изменение цены последней сделки к цене предыдущей сделки, рублей
	LastChangePrcnt float64 `json:"LASTCHANGEPRCNT"` // Изменение цены последней сделки к цене предыдущей сделки, %
	QTY             int     `json:"QTY"`             // Объем последней сделки, в лотах
	Value           float64 `json:"VALUE"`           // Объем последней сделки, в руб
	Value_USD       float64 `json:"VALUE_USD"`       // Объем последней сделки, дол. США
	WapRrice        float64 `json:"WAPRICE"`         // Средневзвешенная цен
	NumTrades       int     `json:"NUMTRADES"`       // Количество сделок за торговый день
	VolToDay        int64   `json:"VOLTODAY"`        // Объем совершенных сделок, выраженный в единицах ценных бумаг
	ValToDay        int64   `json:"VALTODAY"`        // Объем совершенных сделок, в валюте расчетов
	CBRCLOSE        float64 `json:"CBRCLOSE"`        // Вмененная ключевая ставка Банка России (CBR)"
	Time            string  `json:"TIME"`            // Время заключения последней сделки
	SeqNum          int64   `json:"SEQNUM"`          // номер обновления (служебное поле)
	SysTime         string  `json:"SYSTIME"`         // Время загрузки данных системой

}

// history
//https://iss.moex.com/iss/history/engines/stock/markets/bonds/boardgroups/58/securities/SU26218RMFS6.json?iss.meta=off&from=2025-02-20
// "BOARDID", "TRADEDATE", "SHORTNAME", "SECID", "NUMTRADES", "VALUE", "LOW", "HIGH", "CLOSE", "LEGALCLOSEPRICE", "ACCINT", "WAPRICE", "YIELDCLOSE", "OPEN", "VOLUME", "MARKETPRICE2", "MARKETPRICE3", "ADMITTEDQUOTE", "MP2VALTRD", "MARKETPRICE3TRADESVALUE", "ADMITTEDVALUE", "MATDATE", "DURATION", "YIELDATWAP", "IRICPICLOSE", "BEICLOSE", "COUPONPERCENT", "COUPONVALUE", "BUYBACKDATE", "LASTTRADEDATE", "FACEVALUE", "CURRENCYID", "CBRCLOSE", "YIELDTOOFFER", "YIELDLASTCOUPON", "OFFERDATE", "FACEUNIT", "TRADINGSESSION", "CALLOPTIONDATE", "CALLOPTIONYIELD", "CALLOPTIONDURATION", "PUTOPTIONDATE", "DATEYIELDFROMISSUER"

type BondHistory struct {
	SecID         string  `json:"SECID"`         // Код инструмента
	BoardID       string  `json:"BOARDID"`       // Код режима
	ShortName     string  `json:"SHORTNAME"`     // Кратк. наим.
	TRADEDATE     string  `json:"TRADEDATE"`     //
	NumTrades     int     `json:"NUMTRADES"`     // Количество сделок за торговый день
	Value         float64 `json:"VALUE"`         // Объем п в руб
	Open          float64 `json:"OPEN"`          // Цена первой сделки
	Low           float64 `json:"LOW"`           // Минимальная цена сделки
	High          float64 `json:"HIGH"`          // Максимальная цена сделки
	Close         float64 `json:"CLOSE"`         // Цена последней сделки
	MatDate       string  `json:"MATDATE"`       // Дата погашения
	COUPONPERCENT float64 `json:"COUPONPERCENT"` // "Ставка купона, %
	COUPONVALUE   float64 `json:"COUPONVALUE"`   // Сумма купона, в валюте номинала
	FACEVALUE     float64 `json:"FACEVALUE"`     //
	DURATION      int     `json:"DURATION"`      // Дюрация, дней
	YIELDCLOSE    float64 `json:"YIELDCLOSE"`    // Доходность по последней сделке
	LASTTRADEDATE string  `json:"LASTTRADEDATE"` // Дата последней сделки
}

// GetBondsInfo
func (c *Client) GetBondsInfo(board string, boardGroup ...int) ([]BondInfo, error) {
	var err error
	const op = "GetBondsInfo"
	issReq := NewIssRequest().Bonds().Boards(board).Json().MetaData(false).OnlySecurities()

	if len(boardGroup) > 0 {
		issReq.BoardGroups(boardGroup[0])
	}
	url := issReq.URL()
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

	result := make([]BondInfo, 0, len(resp.Securities.Data))
	err = Unmarshal(resp.Securities.Columns, resp.Securities.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}

// GetBondsData получить рыночные данные по облигации
func (c *Client) GetBondsData(board string, symbols string, boardGroup ...int) ([]BondData, error) {
	var err error
	const op = "GetBondsData"
	issReq := NewIssRequest().Bonds().Boards(board).Json().MetaData(false).OnlyMarketData().Symbols(symbols)
	if len(boardGroup) > 0 {
		issReq.BoardGroups(boardGroup[0])
	}
	url := issReq.URL()
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

	result := make([]BondData, 0, len(resp.MarketData.Data))
	err = Unmarshal(resp.MarketData.Columns, resp.MarketData.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// GetBondHistory получить исторические данные по одному символу
func (c *Client) GetBondHistory(symbols string, from, to string) ([]BondHistory, error) {
	return c.NewBondHistoryService(symbols, from, to, "").Do()
}

// GetBondsHistoryDate получить исторические данные по всем облигациям за заданную дату
func (c *Client) GetBondsHistoryDate(date string) ([]BondHistory, error) {
	return c.NewBondHistoryService("", "", "", date).Do()
}

// BondHistoryService сервис для получения исторических данных
type BondHistoryService struct {
	client     *Client
	issRequest *IssRequest
}

// параметры должны быть
// или symbol + from + to = данные по одному символу
// или  date = данные по всем символам за определенную дату
// todo добавить параметр boardgroups необязательный
func (c *Client) NewBondHistoryService(symbol string, from, to string, date string) *BondHistoryService {
	// BoardGroups(58)  = Т+: Основной режим - безадрес.
	iss := NewIssRequest().History().Bonds().BoardGroups(58).WithSecurities(true).Json().MetaData(false).Target(symbol).From(from).To(to).Date(date)
	return &BondHistoryService{
		client:     c,
		issRequest: iss,
	}
}

func (s *BondHistoryService) Next() ([]BondHistory, error) {
	var err error
	const op = "BondHistoryService.Next"

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

	result := make([]BondHistory, 0, len(resp.History.Data))
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

// Do выполняет выгрузку History
func (s *BondHistoryService) Do() ([]BondHistory, error) {
	const op = "BondHistoryService.Do"

	result := make([]BondHistory, 0)
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
