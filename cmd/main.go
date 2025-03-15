package main

import (
	"flag"
	"log"
	"time"

	"TelegramBot/internal/bot"
	"TelegramBot/internal/storage"
	"TelegramBot/internal/telegram"
)

func main() {
	token := mustToken()

	storage := storage.NewStorage()

	telegramAPI := telegram.NewAPI(token)

	bot := bot.NewBot(telegramAPI, storage)

	offset := 0
	for {
		updates, err := telegramAPI.GetUpdates(offset)
		if err != nil {
			log.Println("Ошибка при получении обновлений:", err)
			continue
		}

		for _, update := range updates {
			bot.HandleUpdate(update)
			offset = update.UpdateID + 1
		}

		time.Sleep(1 * time.Second)
	}
}

func mustToken() string {
	token := flag.String("token", "", "Токен вашего Telegram бота (обязательный)")
	flag.Parse()

	if *token == "" {
		flag.Usage() // Показывает справку по использованию флагов
		log.Fatal("Токен бота обязателен. Используйте флаг -token=YOUR_BOT_TOKEN")
	}

	return *token
}
