package iss

import "time"

var TzMsk = initMoscow()

func initMoscow() *time.Location {
	var loc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", int(3*time.Hour/time.Second))
	}
	return loc
}

// Response структура ответа от iss
type Response struct {
	Candles struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"candles"`
	MarketData struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"marketdata"`
	Securities struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"securities"`
	OrderBook struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"orderbook"`
	History struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"history"`
}
