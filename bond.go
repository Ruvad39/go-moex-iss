package iss

import (
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
	YielDatPrevWaPrice    float64 `json:"YIELDATPREVWAPPRICE"`   // Доходность по оценке пред. дня
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

func (c *Client) GetBondsInfo(board string) ([]BondInfo, error) {
	var err error
	const op = "GetBondsInfo"

	url := NewIssRequest().Bonds().Boards(board).Json().MetaData(false).OnlySecurities().URL()
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
