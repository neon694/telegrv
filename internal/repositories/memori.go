package repositories

import "telegrv/internal/models"

type MemoryRepo struct {
 chats    []models.ChatID
 messages []models.Message
}

func (r *MemoryRepo) AddChat(id models.ChatID) {
 for _, chat := range r.chats {
  if chat == id {
   return
  }
 }
 r.chats = append(r.chats, id)
}

func (r *MemoryRepo) GetChats() []models.ChatID {
 result := make([]models.ChatID, len(r.chats))
 copy(result, r.chats)
 return result
}

func (r *MemoryRepo) SaveMessage(msg models.Message) {
 r.messages = append(r.messages, msg)
}

func (r *MemoryRepo) GetMessages() []models.Message {
 result := make([]models.Message, len(r.messages))
 copy(result, r.messages)
 return result
}