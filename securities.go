/*
весь список доступных инструментов
https://iss.moex.com/iss/securities.json?&securities.columns=secid,shortname,name,type,group,primary_boardid

поиск инструмента
https://iss.moex.com/iss/securities.json?q=sbrf-3.25&iss.json=extended&securities.columns=secid,shortname,name,type,group,primary_boardid

--------------------------------------------------------
данные по акции
https://iss.moex.com/iss/securities/SBER.json?iss.json=extended

https://iss.moex.com/iss/securities/SBER.json?iss.meta=off&iss.only=description&description.columns=name,title,value

	["SECID", "Код ценной бумаги", "SBER"],
	["NAME", "Полное наименование", "Сбербанк России ПАО ао"],
	["SHORTNAME", "Краткое наименование", "Сбербанк"],
	["ISIN", "ISIN код", "RU0009029540"],
	["REGNUMBER", "Номер государственной регистрации", "10301481B"],
	["ISSUESIZE", "Объем выпуска", "21586948000"],
	["FACEVALUE", "Номинальная стоимость", "3"],
	["FACEUNIT", "Валюта номинала", "SUR"],
	["ISSUEDATE", "Дата начала торгов", "2007-07-20"],
	["LATNAME", "Английское наименование", "Sberbank"],
	["LISTLEVEL", "Уровень листинга", "1"],
	["ISQUALIFIEDINVESTORS", "Бумаги для квалифицированных инвесторов", "0"],
	["MORNINGSESSION", "Допуск к утренней дополнительной торговой сессии", "1"],
	["EVENINGSESSION", "Допуск к вечерней дополнительной торговой сессии", "1"],
	["TYPENAME", "Вид\/категория ценной бумаги", "Акция обыкновенная"],
	["GROUP", "Код типа инструмента", "stock_shares"],
	["TYPE", "Тип бумаги", "common_share"],
	["GROUPNAME", "Типа инструмента", "Акции"],
	["EMITTER_ID", "Код эмитента", "484"]

-------------------------------------------------------
данные по фьючу
https://iss.moex.com/iss/securities/SRH5.json?iss.json=extended
выбор по названию а не по коду
https://iss.moex.com/iss/securities/Si-6.26.json?iss.json=extended&shortname=1

https://iss.moex.com/iss/securities/SRM6.json?iss.meta=off&iss.only=description&description.columns=name,title,value

	["SECID", "Краткий код", "SRH5"],
	["NAME", "Наименование серии инструмента", "Фьючерсный контракт SBRF-3.25"],
	["SHORTNAME", "Краткое наименование контракта", "SBRF-3.25"],
	["DELIVERYTYPE", "Исполнение", "Поставка ценных бумаг путем заключения сделки в Секции фондового рынка в порядке, предусмотренном Правилами проведения торгов на фондовом рынке ПАО Московская Биржа (до 19.12.2016 – ЗАО «ФБ ММВБ»), по цене, равной результату деления РЦ Контракта, определенной по итогам вечернего Расчетного периода последнего дня заключения Контракта, на лот Контракта."],
	["FRSTTRADE", "Начало обращения", "2024-03-06"],
	["LSTTRADE", "Последний день обращения", "2025-03-20"],
	["LSTDELDATE", "Дата исполнения", "2025-03-21"],
	["ASSETCODE", "Код базового актива", "SBRF"],
	["EXECTYPE", "Тип контракта", "Поставочный"],
	["LOTSIZE", "Лот", "100"],
	["CONTRACTNAME", "Наименование контракта", "Фьючерсный контракт на обыкновенные акции ПАО Сбербанк"],
	["GROUPTYPE", "Группа контрактов", "Акции"],
	["UNIT", "Котировка", "в рублях за лот"],
	["EXPIRATIONCLRSESS", "Клиринг исполнения", "Основной"],
	["TYPENAME", "Вид контракта", "Фьючерс"],
	["GROUP", "Код типа инструмента", "futures_forts"],
	["TYPE", "Тип бумаги", "futures"],
	["GROUPNAME", "Типа инструмента", "Фьючерсы"],
	["EMITTER_ID", "Код эмитента", "484"]

-------------------------------------------------------------
данные по опциону
https://iss.moex.com/iss/securities/SR29000BA5D.json?iss.json=extended
https://iss.moex.com/iss/securities/SR29000BA5D.json?iss.meta=off&iss.only=description&description.columns=name,title,value

	["SECID", "Краткий код", "SR29000BA5D"],
	["NAME", "Наименование серии инструмента", "Нед. марж. амер. Call 29000 с исп. 22 янв. на фьюч. контр. SBRF-3.25"],
	["SHORTNAME", "Краткое наименование контракта", "SBRF-3.25M220125CA29000"],
	["LSTTRADE", "Последний день обращения", "2025-01-22"],
	["SETTLETYPE", "Тип опциона", "Поставочный"],
	["LSTDELDATE", "Дата исполнения", "2025-01-22"],
	["ASSETCODE", "Код базового актива", "SBRF"],
	["FRSTTRADE", "Начало обращения", "2024-12-27"],
	["OPTIONTYPE", "Вид опциона пут или колл", "C"],
	["STRIKE", "Цена Страйк", "29000"],
	["DELIVERYTYPE", "Исполнение", "Исполнение в течение срока действия опциона по заявлению держателя. Автоматическое исполнение опциона в деньгах в последний день срока действия опциона.  Подробнее об <a href=\"http:\/\/fs.moex.com\/files\/14812\">автоэкспирации опционов<\/a> (раздел 1).  <br>  Для того, чтобы определить, находится ли опцион в деньгах, цена базового актива принимается равной расчетной цене контракта. Расчетная цена определяется в соответствии со Спецификацией и приложением 1 к Правилам торгов.  <br>  При исполнении опциона заключается фьючерс, являющийся базовым активом опциона, по цене, равной цене исполнения опциона. Держатель опциона становится покупателем фьючерса, подписчик опциона становится продавцом фьючерса."],
	["EXECTYPE", "Способ исполнения опциона", "Американский"],
	["MARGINSTYLE", "Способ маржирования опциона", "Маржируемый"],
	["LOTSIZE", "Лот", "1"],
	["CONTRACTNAME", "Наименование контракта", "Маржируемый опцион колл на фьючерсный контракт на обыкновенные акции ПАО Сбербанк"],
	["GROUPTYPE", "Группа контрактов", "Акции"],
	["UNIT", "Котировка", "в рублях за лот"],
	["UNDERLYINGASSET", "Базовый актив", "SRH5"],
	["EXPIRATIONCLRSESS", "Клиринг исполнения", "Основной"],
	["SERIES_NAME", "Опционная серия", "SBRF-3.25M220125XA"],
	["TYPENAME", "Вид контракта", "Опцион"],
	["GROUP", "Код типа инструмента", "futures_options"],
	["TYPE", "Тип бумаги", "option"],
	["GROUPNAME", "Типа инструмента", "Опционы"]

-------------------------------------------------------------
данные по облигации
https://iss.moex.com/iss/securities/RU000A1090K0.json?iss.json=extended

https://iss.moex.com/iss/securities/RU000A1090K0.json?iss.meta=off&iss.only=description&description.columns=name,title,value

	["SECID", "Код ценной бумаги", "RU000A1090K0"],
	["NAME", "Полное наименование", "Магнит ПАО БО-004Р-06"],
	["SHORTNAME", "Краткое наименование", "Магнит4P06"],
	["REGNUMBER", "Номер государственной регистрации", "4B02-06-60525-P-004P"],
	["ISIN", "ISIN код", "RU000A1090K0"],
	["ISSUEDATE", "Дата начала торгов", "2024-07-22"],
	["MATDATE", "Дата погашения", "2026-07-12"],
	["BUYBACKDATE", "Дата к которой рассчитывается доходность", "2025-03-19"],
	["INITIALFACEVALUE", "Первоначальная номинальная стоимость", "1000"],
	["FACEUNIT", "Валюта номинала", "SUR"],
	["LATNAME", "Английское наименование", "Magnit BO-004P-06"],
	["STARTDATEMOEX", "Дата начала торгов на Московской Бирже", "2024-07-22"],
	["PROGRAMREGISTRYNUMBER", "Государственный регистрационный номер программы облигации", "4-60525-P-004P-02E"],
	["LISTLEVEL", "Уровень листинга", "2"],
	["DAYSTOREDEMPTION", "Дней до погашения", "503"],
	["ISSUESIZE", "Объем выпуска", "25000000"],
	["FACEVALUE", "Номинальная стоимость", "1000"],
	["ISQUALIFIEDINVESTORS", "Бумаги для квалифицированных инвесторов", "0"],
	["COUPONFREQUENCY", "Периодичность выплаты купона в год", "12"],
	["COUPONDATE", "Дата выплаты купона", "2025-03-19"],
	["COUPONPERCENT", "Ставка купона, %", "22"],
	["COUPONVALUE", "Сумма купона, в валюте номинала", "18.08"],
	["TYPENAME", "Вид\/категория ценной бумаги", "Биржевая облигация"],
	["GROUP", "Код типа инструмента", "stock_bonds"],
	["TYPE", "Тип бумаги", "exchange_bond"],
	["GROUPNAME", "Типа инструмента", "Облигации"],
	["EMITTER_ID", "Код эмитента", "886"]

создать общий тип инструмента
*/
package iss

import (
	"fmt"
	"log/slog"
	"net/http"
)

/*
type Securities struct {
	SecID                string  `json:"SECID"`                // Код инструмента
	BoardID              string  `json:"BOARDID"`              // Код режима
	ShortName            string  `json:"SHORTNAME"`            // Краткое наименование
	Name                 string  `json:"NAME"`                 // Полное название
	TYPENAME             string  `json:"TYPENAME"`             // Вид / категория ценной бумаги
	GROUP                string  `json:"GROUP"`                // Код типа инструмента
	GROUPNAME            string  `json:"GROUPNAME"`            // Типа инструмента
	FACEVALUE            float64 `json:"FACEVALUE"`            // Номинальная стоимость
	FACEUNIT             string  `json:"FACEUNIT"`             // Валюта номинала
	ISQUALIFIEDINVESTORS bool    `json:"ISQUALIFIEDINVESTORS"` // Бумаги для квалифицированных инвесторов
	LISTLEVEL            int     `json:"LISTLEVEL"`            // Уровень листинга
	ISSUEDate            string  `json:"ISSUEDATE"`            // Дата начала торгов
	MORNINGSESSION       bool    `json:"MORNINGSESSION"`       // Допуск к утренней дополнительной торговой сессии
	EVENINGSESSION       bool    `json:"EVENINGSESSION"`       // Допуск к вечерней дополнительной торговой сессии
	MatDate              string  `json:"MATDATE"`              // "Дата погашения
	DAYSTOREDEMPTION     int     `json:"DAYSTOREDEMPTION"`     // Дней до погашения (bond)
	COUPONFREQUENCY      int     `json:"COUPONFREQUENCY"`      // Периодичность выплаты купона в год (bond)
	COUPONDATE           string  `json:"COUPONDATE"`           // Дата выплаты купона (bond)
	COUPONPERCENT        float64 `json:"COUPONPERCENT"`        // Ставка купона, % (bond)
	COUPONVALUE          float64 `json:"COUPONVALUE"`          // Сумма купона, в валюте номинала (bond)
	FRSTTRADE            string  `json:"FRSTTRADE"`            // Начало обращения (futures) (option)
	LSTTRADE             string  `json:"LSTTRADE"`             // Последний день обращения (futures) (option)
	LSTDELDATE           string  `json:"LSTDELDATE"`           // Дата исполнения (futures)
	ASSETCODE            string  `json:"ASSETCODE"`            // Код базового актива (futures) (option)
	UNDERLYINGASSET      string  `json:"UNDERLYINGASSET"`      // Базовый актив  (option)
	LOTSIZE              int     `json:"LOTSIZE"`              // Лот
	GROUPTYPE            string  `json:"GROUPTYPE"`            // Группа контрактов (futures)
	SETTLETYPE           string  `json:"SETTLETYPE"`           // Тип опциона (option)
	OPTIONTYPE           string  `json:"OPTIONTYPE"`           // Вид опциона пут или колл (option)
	STRIKE               float64 `json:"STRIKE"`               // Цена Страйк (option)
	EXECTYPE             string  `json:"EXECTYPE"`             // Способ исполнения опциона (option)
	MARGINSTYLE          string  `json:"MARGINSTYLE"`          // Способ маржирования опциона (option)
	SERIES_NAME          string  `json:"SERIES_NAME"`          // Опционная серия

}
*/

// GetTickersAll список всех инструментов
func (c *Client) GetTickersAll() ([]Ticker, error) {

	var err error
	const op = "GetTickersAl"
	url := NewIssRequest().Target("securities").Json().MetaData(false).OnlySecurities().URL()
	slog.Debug("GetTickersAll", "url", url)
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

	result := make([]Ticker, 0, len(resp.Securities.Data))
	err = Unmarshal(resp.Securities.Columns, resp.Securities.Data, &result)
	//err = UnmarshalCSV(resp.Securities.Columns, resp.Securities.Data, &result, DefaultTagKey)

	if err != nil {
		slog.Error(op+".Unmarshal", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil

}
