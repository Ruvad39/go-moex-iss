package iss

import (
	"errors"
	"net/url"
	"path"
	"strconv"
)

const (
	FortsBoard   = "RFUD" // фьючерсы
	StockBoard   = "TQBR" // акция
	OptionsBoard = "ROPD" // опционы

	AlgoPackStock = "eq" // акции
	AlgoPackForts = "fo" // фьючерсы
	AlgoPackFx    = "fx" // валюта
)

var (
	// EOF обозначает конец выгрузки
	EOF = errors.New("end of data")
)

// IssRequest построитель запроса  iss-moex
type IssRequest struct {
	history    bool   // /history – данные итогов торгов.
	engines    string // trade_engine_name
	markets    string // market_name
	boards     string //
	securities bool   // securities
	symbol     string // один символ (после securities) для запроса "candle"
	target     string //
	format     string // формат данных json xml csv
	// параметры запроса
	symbols         string // список инструментов для строки запроса
	iss_only        string // iss.only=block1,block2 = ответ может содержать несколько блоков данных и этот
	metadata        bool   // iss.meta=on|off  = включать или нет метаинформацию
	jsonFull        bool   // iss.json=compact|extended сокращенный или  расширенный  формат json;
	dateFrom        string // дата from
	dateTo          string // дата till
	date            string // дата date
	interval        int    // Интервал свечек
	start           int    // start =
	q               string // Поиск инструмента по части Кода, Названию, ISIN, Идентификатору Эмитента, Номеру гос.регистрации.
	algoPack        string // тип данных алгопака
	algoPackMarkets string // рынок для  алгопака eq = акции fo = фьючерсы	fx = валюта
	latest          bool   // Super Candles флаг latest=1 возвращает последнюю пятиминутку за указанную дату
	algoPackStock   bool   //
	algoPackForts   bool   //
	algoPackFx      bool   //
}

func NewIssRequest() *IssRequest {
	// значение по умолчанию
	u := &IssRequest{
		format:   "json",
		metadata: true,
	}

	return u
}

// URL создадим строку url
func (u *IssRequest) URL() string {
	_url, err := url.Parse(DefaultApiURL)
	if err != nil {
		return ""
	}

	if u.algoPack != "" {
		_url.Path = path.Join(_url.Path, DefaultAlgoPack)
		// выберем рынок
		if u.algoPackMarkets != "" {
			_url.Path = path.Join(_url.Path, u.algoPackMarkets)
		}
		if u.algoPackStock {
			_url.Path = path.Join(_url.Path, AlgoPackStock)
		}
		if u.algoPackForts {
			_url.Path = path.Join(_url.Path, AlgoPackForts)
		}
		if u.algoPackFx {
			_url.Path = path.Join(_url.Path, AlgoPackFx)
		}
		// что выбираем
		_url.Path = path.Join(_url.Path, u.algoPack)
	}
	if u.history {
		_url.Path = path.Join(_url.Path, "history")
	}
	if u.engines != "" {
		_url.Path = path.Join(_url.Path, "engines", u.engines)
	}
	if u.markets != "" {
		_url.Path = path.Join(_url.Path, "markets", u.markets)
	}
	if u.boards != "" {
		_url.Path = path.Join(_url.Path, "boards", u.boards)
	}
	if u.securities {
		_url.Path = path.Join(_url.Path, "securities")
	}
	// если не пустой символ
	if u.symbol != "" {
		_url.Path = path.Join(_url.Path, u.symbol)
	}
	if u.target != "" {
		target := u.target
		if u.format != "" {
			target = target + "." + u.format
		}
		_url.Path = path.Join(_url.Path, target)
	}
	// все равно проставим формат
	if u.target == "" {
		_url.Path = path.Join(_url.Path, "."+u.format)
	}

	// создаем параметры
	q := _url.Query()

	if !u.metadata {
		q.Set("iss.meta", "off")
	}
	if u.iss_only != "" {
		q.Set("iss.only", u.iss_only)
	}
	if u.jsonFull {
		//iss.json=compact|extended
		q.Set("iss.json", "extended")
	}
	// нет проверки на формат даты
	if u.date != "" {
		q.Set("date", u.date)
	}
	// нет проверки на формат даты
	if u.dateFrom != "" {
		q.Set("from", u.dateFrom)
	}
	// нет проверки на формат даты
	if u.dateTo != "" {
		q.Set("till", u.dateTo)
	}
	if u.interval != 0 {
		q.Set("interval", strconv.Itoa(u.interval))
	}
	if u.start != 0 {
		q.Set("start", strconv.Itoa(u.start))
	}
	// если не пустой список инструментов
	if u.symbols != "" && u.target != "candles" {
		q.Set("securities", u.symbols)
	}
	if u.latest {
		q.Set("latest", "1")
	}
	// // если не пустой список колонок
	// //q.Set("securities.columns",  params.columns)
	// if params.columns != "" {
	// 	q.Set(params.iss_only+".columns", params.columns)
	// }

	// добавляем к URL параметры
	_url.RawQuery = q.Encode()

	return _url.String()
}

// History
func (u *IssRequest) History() *IssRequest {
	u.history = true
	return u
}

// Stock проставим параметры для акций
func (u *IssRequest) Stock() *IssRequest {
	u.engines = "stock"
	u.markets = "shares"
	u.boards = StockBoard
	u.target = "securities"
	return u
}

// Forts проставим параметры для фьючерсов
func (u *IssRequest) Forts() *IssRequest {
	u.engines = "futures"
	u.markets = "forts"
	u.target = "securities"
	u.boards = ""
	return u
}

// Options проставим параметры для опционов
func (u *IssRequest) Options() *IssRequest {
	u.engines = "futures"
	u.markets = "options"
	u.target = "securities"
	u.boards = ""
	return u
}

//https://iss.moex.com/iss/engines/stock/markets/bonds/boards/tqob/securities.json

// Bondsпроставим параметры для Bonds
func (u *IssRequest) Bonds() *IssRequest {
	u.engines = "stock"
	u.markets = "bonds"
	u.target = "securities"
	u.boards = ""
	return u
}

// Candle проставим параметры для запроса свечи
func (u *IssRequest) Candle() *IssRequest {
	u.target = "candles"
	u.securities = true
	return u
}

// Engines /engines/(trade_engine_name)
func (u *IssRequest) Engines(param string) *IssRequest {
	u.engines = param
	return u
}

// Markets /markets/(market_name)
func (u *IssRequest) Markets(param string) *IssRequest {
	u.markets = param
	return u
}

// Boards /boards/(boardid)
func (u *IssRequest) Boards(param string) *IssRequest {
	u.boards = param
	return u
}

// Target
func (u *IssRequest) Target(param string) *IssRequest {
	u.target = param
	//если target == candles всегда вставляем securities
	if param == "candles" {
		u.securities = true
	}
	return u
}

// WithSecurities добавлять securities в строку
func (u *IssRequest) WithSecurities(param bool) *IssRequest {
	u.securities = param
	return u
}

// Json проставим формат данных: json xml
func (u *IssRequest) Json() *IssRequest {
	u.format = "json"
	return u
}

//// Xml проставим формат данных: xml
//func (u *IssRequest) Xml() *IssRequest {
//	u.format = "xml"
//	return u
//}
//
//// формат данных: csv
//func (u *IssRequest) Csv() *IssRequest {
//	u.format = "csv"
//	return u
//}
//
//// формат данных: html
//func (u *IssRequest) Html() *IssRequest {
//	u.format = "html"
//	return u
//}

// MetaData включать или нет метаинформацию перечень,
// тип данных и размер полей (столбцов или атрибутов xml)
// iss.meta=on|off
func (u *IssRequest) MetaData(param bool) *IssRequest {
	u.metadata = param
	return u
}

// JsonFull расширенный формат json
// iss.json=extended
func (u *IssRequest) JsonFull() *IssRequest {
	u.jsonFull = true
	return u
}

// Only iss.only=
func (u *IssRequest) Only(param string) *IssRequest {
	u.iss_only = param
	return u
}

// OnlySecurities iss.only=securities
func (u *IssRequest) OnlySecurities() *IssRequest {
	u.iss_only = "securities"
	return u
}

// OnlyMarketData iss.only=marketdata
func (u *IssRequest) OnlyMarketData() *IssRequest {
	u.iss_only = "marketdata"
	return u
}

// MarketData iss.only=marketdata
func (u *IssRequest) MarketData() *IssRequest {
	u.iss_only = "marketdata"
	return u
}

func (u *IssRequest) Start(param int) *IssRequest {
	u.start = param
	return u
}

// From нет проверки на формат даты
func (u *IssRequest) From(param string) *IssRequest {
	u.dateFrom = param
	return u
}

// To нет проверки на формат даты
func (u *IssRequest) To(param string) *IssRequest {
	u.dateTo = param
	return u
}

// Date нет проверки на формат даты
func (u *IssRequest) Date(param string) *IssRequest {
	u.date = param
	return u
}

// Symbols список символов в строке запроса
func (u *IssRequest) Symbols(param string) *IssRequest {
	u.symbols = param
	return u
}

// Symbol один символ. после securities
func (u *IssRequest) Symbol(param string) *IssRequest {
	u.symbol = param
	return u
}

// Interval
func (u *IssRequest) Interval(param int) *IssRequest {
	u.interval = param
	return u
}

// что выбираем из алгопака
func (u *IssRequest) AlgoPack(param string) *IssRequest {
	u.algoPack = param
	return u
}

// какой рынок в алгопаке
func (u *IssRequest) AlgoPackMarkets(param string) *IssRequest {
	u.algoPackMarkets = param
	return u
}

// флаг latest=1 возвращает последнюю пятиминутку за указанный период
func (u *IssRequest) Latest(param bool) *IssRequest {
	u.latest = param
	return u
}

func (u *IssRequest) AlgoPackStock(param bool) *IssRequest {
	u.algoPackStock = param
	return u
}

func (u *IssRequest) AlgoPackForts(param bool) *IssRequest {
	u.algoPackForts = param
	return u
}

func (u *IssRequest) AlgoPackFx(param bool) *IssRequest {
	u.algoPackFx = param
	return u
}
