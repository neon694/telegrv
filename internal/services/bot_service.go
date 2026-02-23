package services

import (
 "fmt"
 "telegrv/internal/models"
 "telegrv/internal/interfaces"
)

type BotService struct {
 repo interfaces.Repository
}

func NewBotService(repo interfaces.Repository) *BotService {
 return &BotService{repo: repo}
}

func (s *BotService) HandleStart(chatID models.ChatID, userName string) string {
 s.repo.AddChat(chatID)
 if userName != "" {
  return fmt.Sprintf("Hello, %s! I am a bot with menu. Try /menu command!", userName)
 }
 return "Hello! I am a bot with menu. Try /menu command!"
}

func (s *BotService) HandleMenu() string {
 return "Main menu:\n\n/start - Restart conversation\n/help - Help\n/info - Your info\n/stats - Statistics\n\nChoose an option below or type a command!"
}

func (s *BotService) HandleHelp() string {
 return "Available commands:\n\n/start - Start conversation\n/menu - Open main menu\n/help - This help\n/info - Your info\n/stats - Bot statistics\n\nJust write anything - I will repeat it!"
}

func (s *BotService) HandleInfo(userName string, chatID models.ChatID) string {
 if userName != "" {
  return fmt.Sprintf("Your info:\nName: %s\nChat ID: %d", userName, chatID)
 }
 return fmt.Sprintf("Your info:\nChat ID: %d", chatID)
}

func (s *BotService) HandleStats() string {
 chats := s.repo.GetChats()
 messages := s.repo.GetMessages()
 return fmt.Sprintf("Bot statistics:\nTotal chats: %d\nTotal messages: %d", len(chats), len(messages))
}

func (s *BotService) HandleMessage(chatID models.ChatID, text, userName string) string {
 s.repo.SaveMessage(models.Message{
  Name: userName,
  Text: text,
 })
 return fmt.Sprintf("Confirmed, you wrote: %s", text)
}