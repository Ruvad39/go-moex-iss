package iss

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

const (
	libraryName     = "MOEX ISS"
	libraryVersion  = "v0.1.0"
	DefaultApiURL   = "https://iss.moex.com/iss/"
	DefaultAuthURL  = "https://passport.moex.com/authenticate"
	DefaultAlgoPack = "/datashop/algopack"
)

const (
	autCookiesName = "MicexPassportCert"
	autHeaderName  = "X-MicexPassport-Marker"
)

var defaultHTTPTimeout = time.Second * 10

var logLevel = &slog.LevelVar{} // INFO

// SetLogLevel проставим уровень логирования
func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}

// HTTPClient interface Do
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient        HTTPClient
	log               *slog.Logger
	userName          string
	password          string
	micexPassportCert string
}

func NewClient(opts ...ClientOption) (*Client, error) {
	var err error

	client := &Client{
		httpClient: DefaultHTTPClient(),
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})).With(slog.String("package", "moex-iss")),
	}
	//  входящие параметры
	for _, opt := range opts {
		opt(client)
	}
	// если не пустое имя пользователя и пароль = проведем авторизацию
	if client.userName != "" && client.password != "" {
		err = client.Connect()
	}
	return client, err
}

// callAPI запрос к http серверу
func (c *Client) callAPI(r *request) (data []byte, err error) {
	err = c.parseRequest(r)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	ctx := context.Background()
	req = req.WithContext(ctx)
	req.Header = r.header

	//c.log.Debug("callAPI", slog.Any("request", req))

	//req.SetBasicAuth(c.userName, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// если запрос должен делаться только с авторизацией
	if r.authorizationOnly {
		// найдем нужный заголовок
		val := resp.Header.Get(autHeaderName)
		if val != "granted" {
			// вернем ошибку 403
			return nil, fmt.Errorf("Ошибка HTTP 403 - Доступ запрещен. У вас нет прав на просмотр этого каталога или страницы с использованием предоставленных вами учетных данны")
		}

	}
	//c.log.Debug("callAPI", "resp", resp)
	defer func() { _ = resp.Body.Close() }()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.log.Debug("callAPI", "status code", resp.StatusCode, "body", string(data))
	//c.log.Debug("callAPI", "resp.Header", resp.Header)
	if resp.StatusCode >= http.StatusBadRequest {
		c.log.Error("callAPI", "Error response body", resp.StatusCode, slog.Any("data", data))
	}
	return data, err
}

func (c *Client) parseRequest(r *request) (err error) {
	err = r.validate()
	if err != nil {
		return err
	}

	// если мы сделали клиент с cookiejar = всегда должны что то добавить, иначе ошибка?
	r.setHeader("Cookie", c.micexPassportCert)

	//queryString := r.query.Encode()
	if r.baseURL == "" {
		r.baseURL = DefaultApiURL
	}
	if r.fullURL == "" {
		fullURL, err := url.Parse(r.baseURL)
		if err != nil {
			return err
		}
		fullURL.RawQuery = r.query.Encode()
		r.fullURL = fullURL.String()
	}
	c.log.Debug("parseRequest", slog.Any("fullURL", r.fullURL))
	return nil
}

// getJSON выполним запрос и распарсим ответ в JSON
func (c *Client) getJSON(r *request, v interface{}) error {
	var err error
	const op = "getJSON"

	body, err := c.callAPI(r)
	if err != nil {
		slog.Error("getJSON.callAPI", "err", err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}
	if err = json.Unmarshal(body, &v); err != nil {
		slog.Error("getJSON.json.Unmarshal", "err", err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}
	return err
}

/*
Для аутентификации пользователей используется basic-аутентификация.
и передаются серверу в заголовке запроса на https://passport.moex.com/authenticate

При успешной аутентификации сервер возвращает cookie с именем MicexPassportCert.
Далее этот токен должен передаваться при последующих запросах
*/

// Connect Подключение (авторизация) к информационно-статистическому серверу Московской Биржи (ИСС/ISS)
func (c *Client) Connect() error {
	var err error
	method := http.MethodGet
	fullURL := DefaultAuthURL
	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		c.log.Error("Connect.http.NewRequest", "err", err.Error())
		return err
	}
	req.SetBasicAuth(c.userName, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("Connect.httpClient.Do", "err", err.Error())
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == autCookiesName {
			c.micexPassportCert = cookie.Value
			//slog.Debug("resp.Cookies", "Cookie: ", autCookiesName+"="+cookie.Value)
		}
	}

	if resp.StatusCode >= http.StatusBadRequest {
		c.log.Error("Connect", "Error response body", resp.StatusCode, slog.Any("data", data))
	}
	return nil
}

// DefaultHTTPClient создадим и вернем  стандартный  http.Client
func DefaultHTTPClient() HTTPClient {
	var dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	var defaultTransport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   100,
		ExpectContinueTimeout: 0,
		ForceAttemptHTTP2:     true,
		TLSClientConfig:       &tls.Config{},
	}
	var jar, _ = cookiejar.New(nil)
	var defaultHttpClient = &http.Client{
		Timeout:   defaultHTTPTimeout,
		Transport: defaultTransport,
		Jar:       jar,
	}
	return defaultHttpClient
}

// ClientOption установка параметров клиента
type ClientOption func(c *Client)

// WithLogger Logger по умолчанию
func WithLogger(logger *slog.Logger) ClientOption {
	return func(client *Client) {
		client.log = logger
	}
}

// WithUser установим пользователя
func WithUser(user string) ClientOption {
	return func(client *Client) {
		client.userName = user
	}
}

// WithPwd установим пароль пользователя
func WithPwd(pwd string) ClientOption {
	return func(client *Client) {
		client.password = pwd
	}
}

// WithTimeout установить значение Timeout для http.Client
// по умолчанию стоит 10s
func WithTimeout(valie time.Duration) ClientOption {
	return func(client *Client) {
		defaultHTTPTimeout = valie
	}
}

// WithHttpClient установить httpClient
func WithHttpClient(client HTTPClient) ClientOption {
	return func(c *Client) { c.httpClient = client }
}
