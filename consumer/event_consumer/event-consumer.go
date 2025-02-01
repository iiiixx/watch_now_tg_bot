package eventconsumer

import (
	"log"
	"read_adviser_tg_bot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	//log.Println("starting consumer...") // Логируем начало работы
	for {
		//log.Println("fetching events...") // Логируем начало запроса на получение событий
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		//log.Printf("fetched %d events", len(gotEvents)) // Логируем количество полученных событий

		if len(gotEvents) == 0 {

			//log.Println("no events received, sleeping for 1 second...") // Логируем, если нет новых событий

			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}

	}
}

/*
Проблемы:

1. Потеря событий: ретраи, возвращение в хранилище, фоллбэк, подтверждение для фетчера
2. Обработка всей пачки: останавливаться полсе первой ошибки, счетчик ошибок
3. Параллельная обработка
*/
func (c *Consumer) handleEvents(events []events.Event) error {
	//log.Printf("handling %d events", len(events)) // Логируем количество событий для обработки

	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
