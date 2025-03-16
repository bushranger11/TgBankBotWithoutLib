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
		balance, err := b.storage.GetBalance(userID)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Ошибка при получении баланса")
			return
		}
		b.telegramAPI.SendMessage(chatID, fmt.Sprintf("Ваш баланс: %.2f", balance))
	case strings.HasPrefix(text, "/deposit "):
		amountStr := text[len("/deposit "):]
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Неверный формат суммы")
			return
		}
		err = b.storage.Deposit(userID, amount)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Ошибка при пополнении баланса")
			return
		}
		balance, err := b.storage.GetBalance(userID)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Ошибка при получении баланса")
			return
		}
		b.telegramAPI.SendMessage(chatID, fmt.Sprintf("Счет пополнен на %.2f. Новый баланс: %.2f", amount, balance))
	case strings.HasPrefix(text, "/withdraw "):
		amountStr := text[len("/withdraw "):]
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Неверный формат суммы")
			return
		}
		err = b.storage.Withdraw(userID, amount)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Недостаточно средств на счете")
			return
		}
		balance, err := b.storage.GetBalance(userID)
		if err != nil {
			b.telegramAPI.SendMessage(chatID, "Ошибка при получении баланса")
			return
		}
		b.telegramAPI.SendMessage(chatID, fmt.Sprintf("Снято %.2f. Новый баланс: %.2f", amount, balance))
	default:
		b.telegramAPI.SendMessage(chatID, "Неизвестная команда")
	}
}
