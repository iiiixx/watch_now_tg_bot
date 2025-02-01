package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"read_adviser_tg_bot/lib/e"
	"strconv"
)

type Client struct {
	host     string      //хост api сервеса тг
	basePath string      //префиск, с которого начинаются все запросы //tg-bot.com/bot<token>, токен используется для аутентификации запросов.
	client   http.Client //экземпляр http.Client, используемый для отправки HTTP-запросов.
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	//Функция-конструктор, создающая новый экземпляр Client. Она принимает:
	//host: Хост сервиса (например, "api.telegram.org").
	//token: Токен бота, используемый для авторизации.

	//Возвращает новый объект Client с заполненными полями:
	//host устанавливается из аргумента.
	//basePath формируется через функцию newBasePath.
	//client инициализируется стандартным HTTP-клиентом.

	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}

}

func newBasePath(token string) string {
	//Формирует базовый путь, начинающийся с "bot<token>". Пример:
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	//offset int: Смещение для получения обновлений (например, начиная с какого ID или времени обрабатывать обновления).
	//limit int: Максимальное количество обновлений, которые нужно получить за один запрос.

	//1. Формирование параметров запроса:
	//Создаётся объект url.Values (представляет параметры строки запроса).
	q := url.Values{}
	//С помощью Add добавляются два параметра:
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	//2. Выполнение HTTP-запроса:
	//Метод вызывает doRequest, передавая: getUpdatesMethod — имя метода API, которое, определено строковая константа
	//q — параметры строки запроса.
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	//3. Парсинг JSON-ответа:
	//Создаётся переменная res типа UpdatesResponse. Это структура, которая, вероятно, содержит данные ответа API в структурированном виде.
	var res UpdatesResponse
	//json.Unmarshal преобразует JSON-данные из data в структуру res.
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	/*log.Printf("Received %d updates from Telegram", len(res.Result))

	// Теперь добавим логи для каждого обновления
	for _, update := range res.Result {
		log.Printf("Message text: %s", update.Message.Text)
	}*/

	//Из структуры UpdatesResponse возвращается поле Result, которое, является срезом объектов типа Update.
	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	//chatID int: Уникальный идентификатор чата, в который нужно отправить сообщение.
	//text string: Текст сообщения, которое нужно отправить.

	//1. Формирование параметров запроса:
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	//2. Метод doRequest отправляет HTTP-запрос с параметрами,
	//используя sendMessageMethod — константу, содержащую название API-метода
	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	//Здесь используется отложенное выполнение (defer) для обработки ошибок.
	//Если метод завершится с ошибкой, она будет обёрнута дополнительным сообщением "can't do request" через функцию WrapIfErr.
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	//1. Создаётся объект URL
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	//2. Создание HTTP-запроса:
	//Метод создаёт GET-запрос к указанному URL
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	//3. Добавление параметров запроса:
	//Здесь параметры query (объект url.Values) кодируются и добавляются к URL как строка запроса.
	req.URL.RawQuery = query.Encode()

	//4. Выполнение запроса:
	//Запрос отправляется через c.client
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	//5. Чтение тела ответа:
	//Тело ответа считывается в переменную body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
