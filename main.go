package main

import (
	"context"
	"flag"
	"log"

	tgClient "tg_bot/clients/telegram"
	"tg_bot/events/telegram"
	"tg_bot/storage/sqlite"

	eventconsumer "tg_bot/consumer/event_consumer"
)

//Для остановки службы: - sudo launchctl unload /Library/LaunchDaemons/com.example.bot.plist

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	//storagePath = "files_storage"
	batchSize = 100
)

// run with:
// -tg-bot-token ''

func main() {
	// s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		//files.New(storagePath),
		s,
	)

	log.Print("service started")

	//fetcher - получает, processor - обрабатывает

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	log.Println("starting event consumer")
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)

	}
}

func mustToken() string {
	//must пишется ,так как мы не будем обрабатывать ошибку, токен обязателен

	token := flag.String(
		"tg-bot-token",                     //имя флага, который ожидается в командной строке.
		"",                                 //значение по умолчанию (пустая строка, если флаг не указан).
		"token for access to telegram bot", //usage
	)

	flag.Parse() //обрабатывает переданные аргументы командной строки и присваивает их значения соответствующим переменным.

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
