package handlers

import (
 "log"
 "strings"
 "bottg/internal/models"
 "bottg/internal/services"
 tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramHandler struct {
 bot     *tgbotapi.BotAPI
 service *services.BotService
}

func NewTelegramHandler(bot *tgbotapi.BotAPI, service *services.BotService) *TelegramHandler {
 return &TelegramHandler{bot: bot, service: service}
}

func (h *TelegramHandler) Start() {
 u := tgbotapi.NewUpdate(0)
 u.Timeout = 60

 updates := h.bot.GetUpdatesChan(u)

 log.Println("Bot started and listening...")

 for update := range updates {
  if update.Message == nil || update.Message.Text == "" {
   continue
  }

  chatID := update.Message.Chat.ID
  originalText := update.Message.Text
  userName := update.Message.From.UserName
  if userName == "" {
   userName = "User"
  }

  log.Printf("[%d] @%s: %q", chatID, userName, originalText)

  normalized := normalizeCommand(originalText)
  log.Printf("Normalized: %q", normalized)

  var replyText string
  var replyMarkup tgbotapi.ReplyKeyboardMarkup

  if strings.HasPrefix(normalized, "/start") {
   replyText = h.service.HandleStart(models.ChatID(chatID), userName)
   replyMarkup = h.getMainMenuKeyboard()
  } else if strings.HasPrefix(normalized, "/menu") {
   replyText = h.service.HandleMenu()
   replyMarkup = h.getMainMenuKeyboard()
  } else if strings.HasPrefix(normalized, "/help") {
   replyText = h.service.HandleHelp()
   replyMarkup = h.getMainMenuKeyboard()
  } else if strings.HasPrefix(normalized, "/info") {
   replyText = h.service.HandleInfo(userName, models.ChatID(chatID))
   replyMarkup = h.getMainMenuKeyboard()
  } else if strings.HasPrefix(normalized, "/stats") {
   replyText = h.service.HandleStats()
   replyMarkup = h.getMainMenuKeyboard()
  } else {
   replyText = h.service.HandleMessage(models.ChatID(chatID), originalText, userName)
   replyMarkup = h.getMainMenuKeyboard()
  }

  if replyText == "" {
   replyText = "I received your message but couldn't process it. Try /help"
   log.Printf("Empty reply for: %q", normalized)
  }

  msg := tgbotapi.NewMessage(chatID, replyText)
  msg.ReplyMarkup = replyMarkup

  if !strings.HasPrefix(normalized, "/") {
   msg.ReplyToMessageID = update.Message.MessageID
  }

  if _, err := h.bot.Send(msg); err != nil {
   log.Printf("Error sending to %d: %v", chatID, err)
  } else {
   log.Printf("Sent to %d: %q", chatID, replyText)
  }
 }
}

func normalizeCommand(text string) string {
 if idx := strings.Index(text, "@"); idx != -1 {
  text = text[:idx]
 }
 return strings.ToLower(strings.TrimSpace(text))
}

func (h *TelegramHandler) getMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
 row1 := tgbotapi.NewKeyboardButtonRow(
  tgbotapi.NewKeyboardButton("/help"),
  tgbotapi.NewKeyboardButton("/info"),
 )
 row2 := tgbotapi.NewKeyboardButtonRow(
  tgbotapi.NewKeyboardButton("/stats"),
  tgbotapi.NewKeyboardButton("Say Hello"),
 )

 keyboard := tgbotapi.NewReplyKeyboard(row1, row2)
 keyboard.ResizeKeyboard = true
 keyboard.OneTimeKeyboard = false
 return keyboard
}