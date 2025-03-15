package bot

import (
	"fmt"
	"strconv"
	"strings"

	"TelegramBot/internal/storage"
	"TelegramBot/internal/telegram"
)

type Bot struct {
	telegramAPI *telegram.API
	storage     *storage.Storage
}

func NewBot(telegramAPI *telegram.API, storage *storage.Storage) *Bot {
	return &Bot{
		telegramAPI: telegramAPI,
		storage:     storage,
	}
}

func (b *Bot) HandleUpdate(update telegram.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	userID := update.Message.From.ID

	switch {
	case text == "/start":
		b.telegramAPI.SendMessage(chatID, "Добро пожаловать в банк! Используйте команды:\n/balance - проверить баланс\n/deposit <сумма> - пополнить счет\n/withdraw <сумма> - снять деньги")
	case text == "/balance":
		balance := b.storage.GetBalance(userID)
		b.telegramAPI.SendMessage(chatID, fmt.Sprintf("Ваш баланс: %.2f", balance))
	case strings.HasPrefix(text, "/deposit "):
		amountStr := text[len("/deposit "):]
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Неверный формат суммы")
			return
		}
		b.storage.Deposit(userID, amount)
		b.telegramAPI.SendMessage(chatID, fmt.Sprintf("Счет пополнен на %.2f. Новый баланс: %.2f", amount, b.storage.GetBalance(userID)))
	case strings.HasPrefix(text, "/withdraw "):
		amountStr := text[len("/withdraw "):]
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Неверный формат суммы")
			return
		}
		if !b.storage.Withdraw(userID, amount) {
			b.telegramAPI.SendMessage(chatID, "Недостаточно средств на счете")
			return
		}
		b.telegramAPI.SendMessage(chatID, fmt.Sprintf("Снято %.2f. Новый баланс: %.2f", amount, b.storage.GetBalance(userID)))
	default:
		b.telegramAPI.SendMessage(chatID, "Неизвестная команда")
	}
}
