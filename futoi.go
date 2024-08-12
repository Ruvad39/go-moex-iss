/*
Открытые позиции по фьючерсным контрактам в разрезе физ. и юр. лиц
https://moexalgo.github.io/des/futoi/

AUTHORIZED	Real-time данные (обновление каждые 5 мин)
NOT AUTHORIZED	Данные задержаны на 15 дней

по всем инстументам
/iss/analyticalproducts/futoi/securities.json

# получить последнюю пятиминутку
https://iss.moex.com/iss/analyticalproducts/futoi/securities.json?date=2024-04-17&latest=1

по заданному инструменту
/iss/analyticalproducts/futoi/securities/{ticker}.json

# пример
https://iss.moex.com/iss/analyticalproducts/futoi/securities/si.json?date=2024-04-17

# получить последнюю пятиминутку
https://iss.moex.com/iss/analyticalproducts/futoi/securities/si.json?date=2024-04-17&latest=1

https://iss.moex.com/iss/analyticalproducts/futoi/securities/si.json?from=2024-08-12&till=2024-08-12&latest=1


TODO выдает данные по 1000 записей. сделать итератор запросов
*/

package iss

import (
	"fmt"
	"log/slog"
	"net/http"
)

// FutOI Открытые позиции по фьючерсным контрактам в разрезе физ. и юр. лиц
type FutOI struct {
	SessID      int32  `json:"sess_id"`       // номер торговой сессии
	SeqNum      int32  `json:"seqnum"`        // Номер пакета данных. Техническое поле
	TradeDate   string `json:"tradedate"`     // Дата date:10
	TradeTime   string `json:"tradetime"`     // Время последней сделки, которая была учтена при расчете time:10
	Ticker      string `json:"ticker"`        // Двухсимвольный код контракта
	ClGroup     string `json:"clgroup"`       // группа клиентов: fiz – физические лица, yur – юридические лица
	Pos         int64  `json:"pos"`           // Величина открытых позиций
	PosLong     int64  `json:"pos_long"`      // Величина длинных открытых позиций
	PosShort    int64  `json:"pos_short"`     // Величина коротких открытых позиций
	PosLongNum  int64  `json:"pos_long_num"`  // Количество лиц, имеющих длинную открытую позицию
	PosShortNum int64  `json:"pos_short_num"` // Количество лиц, имеющих короткую открытую позицию
	SysTime     string `json:"systime"`       // Время публикации данных datetime:19
}

// GetFutOIAll Открытые позиции физ. и юр. лиц по всем инструментам
// date = за дату ; latest =1 возвращает последнюю пятиминутку за указанную дату
func (c *Client) GetFutOIAll(date string, latest int) ([]FutOI, error) {
	var err error
	const op = "GetFutOIAll"
	url := "https://iss.moex.com/iss/analyticalproducts/futoi/securities.json"
	r := &request{
		method:  http.MethodGet,
		baseURL: url,
	}
	if date != "" {
		r.setParam("date", date)
	}
	if latest == 1 {
		r.setParam("latest", latest)
	}

	type requestData struct {
		Sec struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"futoi"`
	}
	var resp requestData
	err = c.getJSON(r, &resp)
	if err != nil {
		slog.Error(op+".getJSON", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]FutOI, 0, len(resp.Sec.Data))
	err = Unmarshal(resp.Sec.Columns, resp.Sec.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// GetFutOI по заданному тикеру
func (c *Client) GetFutOI(ticker string, from, to string, latest int) ([]FutOI, error) {
	var err error
	const op = "GetFutOI"
	url := "https://iss.moex.com/iss/analyticalproducts/futoi/securities/" + ticker + ".json"
	r := &request{
		method:  http.MethodGet,
		baseURL: url,
	}
	if from != "" {
		r.setParam("from", from)
	}
	if to != "" {
		r.setParam("till", to)
	}
	if latest == 1 {
		r.setParam("latest", latest)
	}

	type requestData struct {
		Sec struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"futoi"`
	}
	var resp requestData
	err = c.getJSON(r, &resp)
	if err != nil {
		slog.Error(op+".getJSON", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]FutOI, 0, len(resp.Sec.Data))
	err = Unmarshal(resp.Sec.Columns, resp.Sec.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
