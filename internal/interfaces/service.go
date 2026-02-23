package interfaces

import "telegrv/internal/models"

type BotService interface {
    HandleStart(chatID models.ChatID, userName string) string
    HandleMenu() string
    HandleHelp() string
    HandleInfo(userName string, chatID models.ChatID) string
    HandleStats() string
    HandleMessage(chatID models.ChatID, text, userName string) string
}