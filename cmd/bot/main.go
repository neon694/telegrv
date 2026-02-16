package main

import (
 "log"
 "os"
 "bottg/internal/repositories"
 "bottg/internal/services"
  "bottg/internal/handlers"

 tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
 "github.com/joho/godotenv"
)

func main() {
 err := godotenv.Load()
 if err != nil {
  log.Println("file.env not found")
 }

 token := os.Getenv("TELEGRAM_TOKEN")
 if token == "" {
  log.Fatal("Error: The TELEGRAM_TOKEN variable is not set. Create a .env file with the following contents: TELEGRAM_TOKEN=your_token")
 }

 bot, err := tgbotapi.NewBotAPI(token)
 if err != nil {
  log.Fatal("Error creating bot:", err)
 }
 log.Printf("Bot @%s launched", bot.Self.UserName)

 repo := &repositories.MemoryRepo{}

 service := services.NewBotService(repo)

 telegramHandler := handlers.NewTelegramHandler(bot, service)
 httpHandler := handlers.NewHTTPHandler(bot, repo)

 go telegramHandler.Start()
 httpHandler.Start(":8080")
}