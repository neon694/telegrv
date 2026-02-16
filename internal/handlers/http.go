package handlers

import (
 "encoding/json"
 "log"
 "net/http"
 "strconv"

 "bottg/internal/repositories"
 tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HTTPHandler struct {
 bot  *tgbotapi.BotAPI
 repo *repositories.MemoryRepo
}

func NewHTTPHandler(bot *tgbotapi.BotAPI, repo *repositories.MemoryRepo) *HTTPHandler {
 return &HTTPHandler{bot: bot, repo: repo}
}

func (h *HTTPHandler) Start(addr string) {
 http.HandleFunc("/send", h.send)
 http.HandleFunc("/broadcast", h.broadcast)
 http.HandleFunc("/chats", h.getChats)

 log.Printf("HTTP API запущен на http://localhost%s", addr)
 http.ListenAndServe(addr, nil)
}

func (h *HTTPHandler) send(w http.ResponseWriter, r *http.Request) {
 chatIDStr := r.URL.Query().Get("chat_id")
 text := r.URL.Query().Get("text")

 if chatIDStr == "" || text == "" {
  http.Error(w, `{"error": "chat_id and text are required"}`, http.StatusBadRequest)
  return
 }

 chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
 if err != nil {
  http.Error(w, `{"error": "invalid chat_id"}`, http.StatusBadRequest)
  return
 }

 _, err = h.bot.Send(tgbotapi.NewMessage(chatID, text))
 if err != nil {
  http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
  return
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(map[string]string{
  "status":  "ok",
  "chat_id": chatIDStr,
  "text":    text,
 })
}

func (h *HTTPHandler) broadcast(w http.ResponseWriter, r *http.Request) {
 text := r.URL.Query().Get("text")
 if text == "" {
  http.Error(w, `{"error": "text is required"}`, http.StatusBadRequest)
  return
 }

 chats := h.repo.GetChats()
 sentCount := 0

 for _, chatID := range chats {
  if _, err := h.bot.Send(tgbotapi.NewMessage(int64(chatID), text)); err == nil {
   sentCount++
  }
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(map[string]interface{}{
  "status":  "ok",
  "sent_to": sentCount,
  "total":   len(chats),
  "text":    text,
 })
}

func (h *HTTPHandler) getChats(w http.ResponseWriter, r *http.Request) {
 chats := h.repo.GetChats()

 chatIDs := make([]int64, len(chats))
 for i, id := range chats {
  chatIDs[i] = int64(id)
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(map[string]interface{}{
  "chats": chatIDs,
  "count": len(chats),
 })
}