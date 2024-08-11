package iss

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// только по авторизации
// https://iss.moex.com/iss/engines/stock/markets/shares/boards/tqbr/securities/sber/orderbook.json

// OrderBookData структура стакана которая приходил с сервера
type OrderBookData struct {
	SecID      string  `json:"SECID"`
	BoardID    string  `json:"BOARDID"`
	BuySell    string  `json:"BUYSELL"`
	Price      float64 `json:"PRICE"`
	Quantity   int32   `json:"QUANTITY"`
	SEQNUM     int64   `json:"SEQNUM"`
	UpdateTime string  `json:"UPDATETIME"`
	Decimals   int32   `json:"DECIMALS"`
}

type PriceVolume struct {
	Price  float64 // цена
	Volume int32   // объем
}

func (p PriceVolume) String() string {
	return fmt.Sprintf("PriceVolume{ Price: %v, Volume: %v }", p.Price, p.Volume)
}

// PriceVolumeSlice Биды  Аски
type PriceVolumeSlice []PriceVolume

func (slice PriceVolumeSlice) Len() int { return len(slice) }

// First вернем первый элемент
func (slice PriceVolumeSlice) First() (PriceVolume, bool) {
	if len(slice) > 0 {
		return slice[0], true
	}
	return PriceVolume{}, false
}

// Last вернем заданный элемент
func (slice PriceVolumeSlice) Last(position int) (PriceVolume, bool) {
	if len(slice) > position {
		return slice[position], true
	}
	return PriceVolume{}, false
}

// SumDepth вернем объем стакана
func (slice PriceVolumeSlice) SumDepth() int32 {
	var total int32
	for _, pv := range slice {
		total = total + pv.Volume
	}

	return total
}

// OrderBook биржевой стакан
type OrderBook struct {
	Bids       PriceVolumeSlice `json:"bids"` // Биды
	Asks       PriceVolumeSlice `json:"asks"` // Аски
	SEQNUM     int64            `json:"SEQNUM"`
	UpdateTime string           `json:"UPDATETIME"`
	Decimals   int32            `json:"DECIMALS"`
}

// LastTime время обновление стакана
func (b OrderBook) LastTime() time.Time {
	// дата = 8 первых цифр из SEQNUM 202407111 ;  время = UpdateTime "10:20:38"
	dt := strconv.FormatInt(b.SEQNUM, 10)[:8] + " " + b.UpdateTime
	t, _ := time.ParseInLocation("20060102 15:04:05", dt, TzMsk)

	//t.In(TzMsk)
	return t
}

func (b OrderBook) BestBid() (PriceVolume, bool) {
	if len(b.Bids) == 0 {
		return PriceVolume{}, false
	}
	return b.Bids[0], true
}

func (b OrderBook) BestAsk() (PriceVolume, bool) {
	if len(b.Asks) == 0 {
		return PriceVolume{}, false
	}
	return b.Asks[0], true
}

func (b OrderBook) String() string {
	sb := strings.Builder{}
	sb.WriteString("BOOK ")
	//sb.WriteString(b.Symbol)
	sb.WriteString("\n")
	sb.WriteString(b.LastTime().Format("2006-01-02T15:04:05-0700"))
	//sb.WriteString(b.LastTime().String())
	sb.WriteString("\n")

	if len(b.Asks) > 0 {
		sb.WriteString("ASKS:\n")
		for i := len(b.Asks) - 1; i >= 0; i-- {
			sb.WriteString("- ASK: ")
			sb.WriteString(b.Asks[i].String())
			sb.WriteString("\n")
		}
	}

	if len(b.Bids) > 0 {
		sb.WriteString("BIDS:\n")
		for _, bid := range b.Bids {
			sb.WriteString("- BID: ")
			sb.WriteString(bid.String())
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// NewOrderBook соберем стакан из входящего OrderBookData
func NewOrderBook(in []OrderBookData) OrderBook {
	book := OrderBook{}
	if len(in) == 0 {
		return book
	}
	book.SEQNUM = in[0].SEQNUM
	book.UpdateTime = in[0].UpdateTime
	book.Decimals = in[0].Decimals

	// пройдем в цикле и заполним биды и аскм
	// Cтакан по акциям доступен с глубиной в 10 уровней
	// Cтакан по фьючерсам доступен с глубиной в 20 уровней
	// поставим максимальный размер = 20
	bids := make(PriceVolumeSlice, 0, 20)
	asks := make(PriceVolumeSlice, 0, 20)

	for _, data := range in {
		pv := PriceVolume{data.Price, data.Quantity}
		if data.BuySell == "B" {
			bids = append(bids, pv)
		} else {
			asks = append(asks, pv)
		}
	}
	// TODO сделать замеры производительности разных сортировок
	// asks  прямая сортировка
	sort.SliceStable(asks, func(a, b int) bool {
		return asks[a].Price < asks[b].Price
	})
	// bids обратная сортировка
	sort.SliceStable(bids, func(a, b int) bool {
		return bids[b].Price < bids[a].Price
	})

	book.Bids = bids
	book.Asks = asks
	return book

}

// OrderBookService сервис для получения стакана
type OrderBookService struct {
	client     *Client
	issRequest *IssRequest
}

func (c *Client) NewOrderBookService(engines, markets, board, symbol string) *OrderBookService {
	iss := NewIssRequest().Candle().
		Engines(engines).
		Markets(markets).
		Boards(board).
		Symbol(symbol).
		Target("orderbook").
		Json().MetaData(false)

	return &OrderBookService{
		client:     c,
		issRequest: iss,
	}
}

// test
func (s OrderBookService) URL() string {
	return s.issRequest.URL()
}

// Do выполняет загрузку стакана
// запрос должен выполнятся только с авторизацией
func (s OrderBookService) Do() (OrderBook, error) {
	var err error
	const op = "OrderBookService.Do"

	result := OrderBook{}
	r := &request{
		method:            http.MethodGet,
		fullURL:           s.issRequest.URL(),
		authorizationOnly: true,
	}

	var resp Response
	err = s.client.getJSON(r, &resp)
	if err != nil {
		slog.Error(op+".getJSON", "err", err.Error())
		return result, fmt.Errorf("%s: %w", op, err)
	}

	data := make([]OrderBookData, 0, len(resp.OrderBook.Data))
	err = Unmarshal(resp.OrderBook.Columns, resp.OrderBook.Data, &data)
	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return NewOrderBook(data), nil
}
