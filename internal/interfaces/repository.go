package interfaces

import "telegrv/internal/models"

type Repository interface {
    AddChat(id models.ChatID)
    GetChats() []models.ChatID
    SaveMessage(msg models.Message)
    GetMessages() []models.Message
}