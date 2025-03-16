package main

import (
	"flag"
	"log"
	"os"
	"time"

	"TelegramBot/internal/bot"
	"TelegramBot/internal/storage"
	"TelegramBot/internal/telegram"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке .env файла: %v", err)
	}

	connString := os.Getenv("DB_CONN_STRING")
	if connString == "" {
		log.Fatal("DB_CONN_STRING не указана в .env файле")
	}

	storage, err := storage.NewStorage(connString)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	token := mustToken()

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
