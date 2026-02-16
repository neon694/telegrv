package models

type ChatID int64

type Message struct {
 Name string `json:"name"`
 Text string `json:"text"`
}