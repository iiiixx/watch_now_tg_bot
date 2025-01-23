package main

import (
	"flag"
	"log"
	"read_adviser_tg_bot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {

	tgClient = telegram.New(mustToken())

	// fearcher = featcher.New()

	// processor = processor.New()

	//consumer.Start(fetcher - получает, processor - обрабатывает)
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
