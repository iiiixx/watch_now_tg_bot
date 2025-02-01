package main

import (
	"flag"
	"log"

	tgClient "read_adviser_tg_bot/clients/telegram"
	"read_adviser_tg_bot/events/telegram"
	"read_adviser_tg_bot/storage/files"

	eventconsumer "read_adviser_tg_bot/consumer/event_consumer"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	//fetcher - получает, processor - обрабатывает

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)

	}
}

func mustToken() string {
	//must пишется ,так как мы не будем обрабатывать ошибку, токен обязателен

	token := flag.String(
		"token-bot-token",                  //имя флага, который ожидается в командной строке.
		"",                                 //значение по умолчанию (пустая строка, если флаг не указан).
		"token for access to telegram bot", //usage
	)

	flag.Parse() //обрабатывает переданные аргументы командной строки и присваивает их значения соответствующим переменным.

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
