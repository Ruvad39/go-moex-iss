package iss

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Candle структура свечи
type Candle struct {
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Value  float64 `json:"value"`
	Volume float64 `json:"volume"`
	Begin  string  `json:"begin"`
	End    string  `json:"end"`
}

var layout = "2006-01-02 15:04:05"

// Time время начала свечи
func (k Candle) Time() time.Time {
	var t time.Time
	t, err := time.Parse(layout, k.Begin)
	if err != nil {
		slog.Error("Time", "err", err.Error())
	}
	return t
}

// Candles слайс свечей
type Candles struct {
	Symbol   string
	Interval string
	Data     []Candle
}

// Len
func (k Candles) Len() int {
	return len(k.Data)
}

// First первая свеча
func (k Candles) First() Candle {
	return k.Data[0]
}

// Last вернем заданную по номеру свечу. 0 = последняя
func (k Candles) Last(position int) Candle {
	return k.Data[len(k.Data)-1-position]
}

// Index вернем заданную по номеру свечу. 0 = последняя
func (k Candles) Index(position int) Candle {
	return k.Last(position)
}

// Доступные интервалы свечей.
const (
	Interval_M1  = 1
	Interval_M10 = 10
	Interval_H1  = 60
	Interval_D1  = 24
	Interval_W1  = 7
	Interval_MN1 = 31
	Interval_Q1  = 4
)

// ParseInterval
func ParseInterval(input string) (int, error) {
	switch input {
	case "M1":
		return Interval_M1, nil
	case "M10":
		return Interval_M10, nil
	case "H1":
		return Interval_H1, nil
	case "D1":
		return Interval_D1, nil
	case "W1":
		return Interval_W1, nil
	case "MN1":
		return Interval_MN1, nil
	case "Q1":
		return Interval_Q1, nil
	default:
		return -1, fmt.Errorf("не поддерживаемый формат периода свечи %s", input)
	}
}

// IntervalToString конвертация Interval свечей  в строку
func IntervalToString(i int) string {
	switch i {
	case Interval_M1:
		return "M1"
	case Interval_M10:
		return "M10"
	case Interval_H1:
		return "H1"
	case Interval_D1:
		return "D1"
	case Interval_W1:
		return "W1"
	case Interval_MN1:
		return "MN1"
	case Interval_Q1:
		return "Q1"

	}
	return "неизвестно"
}

// CandlesService сервис для получения исторических свечей
type CandlesService struct {
	client     *Client
	issRequest *IssRequest
}

// NewCandlesService создание сервиса
func (c *Client) NewCandlesService(engines, markets, board, symbol string, interval int, from, to string) *CandlesService {
	iss := NewIssRequest().Candle().
		Engines(engines).
		Markets(markets).
		Boards(board).
		Symbol(symbol).
		Interval(interval).
		From(from).
		To(to).
		Json().MetaData(false)

	return &CandlesService{
		client:     c,
		issRequest: iss,
	}
}

// Do выполняет выгрузку свечей
func (s *CandlesService) Do() (Candles, error) {
	const op = "CandlesService.Do"

	candles := Candles{
		Symbol:   s.issRequest.symbol,
		Interval: IntervalToString(s.issRequest.interval),
	}
	candlesData := make([]Candle, 0)
	count := 1
	for {
		// "fetch candles: item(s) processed"
		s.client.log.Debug(op, "запрос свечей: номер запроса", count)

		t_candles, err := s.Next()
		if err != nil {
			if errors.Is(err, EOF) {
				break
			}
			return candles, fmt.Errorf("%s: %w", op, err)
		}
		candlesData = append(candlesData, t_candles...)
		count++

	}
	candles.Data = candlesData
	return candles, nil
}

// Next загружает следующую страницу данных
// Если данных больше нет, то возвращается ошибка EOF
// TODO что возвращать данные или ссылку?
func (s *CandlesService) Next() ([]Candle, error) {
	var err error
	const op = "CandlesService.Next"

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

	result := make([]Candle, 0, len(resp.OrderBook.Data))
	err = Unmarshal(resp.Candles.Columns, resp.Candles.Data, &result)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(result) == 0 {
		return nil, EOF
	}
	s.client.log.Debug(op, "len(result)", len(result), "mindate", result[0].Begin,
		"maxdate", result[len(result)-1].Begin)

	s.issRequest.start += len(result)
	s.client.log.Debug(op, "s.issRequest.start", s.issRequest.start)

	return result, nil
}

// GetStockCandles получить историю свечей по акциям
// Board TQBR
func (c *Client) GetStockCandles(symbols string, interval int, from, to string) (Candles, error) {
	return c.NewCandlesService("stock", "shares", StockBoard, symbols, int(interval), from, to).Do()

}

// GetFortsCandles получить историю свечей по акциям
// Board RFUD
func (c *Client) GetFortsCandles(symbols string, interval int, from, to string) (Candles, error) {
	return c.NewCandlesService("futures", "forts", FortsBoard, symbols, int(interval), from, to).Do()

}
