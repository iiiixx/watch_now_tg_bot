package telegram

import (
	"io"
	"json"
	"net/http"
	"net/url"
	"path"
	"read_adviser_tg_bot/lib/e"
	"strconv"
)

type Client struct {
	host     string //хост api сервеса тг
	basePath string //префиск, с которого начинаются все запросы //tg-bot.com/bot<token>
	client   http.Client
}

const (
	getUpdatesMethod = "getUpdates"
)

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}

}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("offset", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	//Здесь используется отложенное выполнение (defer) для обработки ошибок.
	//Если метод завершится с ошибкой, она будет обёрнута дополнительным сообщением "can't do request" через функцию WrapIfErr.
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	//Создаётся объект URL
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	//Метод создаёт GET-запрос к указанному URL
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	//Здесь параметры query (объект url.Values) кодируются и добавляются к URL как строка запроса.
	req.URL.RawQuery = query.Encode()

	//Запрос отправляется через c.client
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	//Тело ответа считывается в переменную body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) SendMessage() {

}
