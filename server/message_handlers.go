package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strconv"
)

// MessagesHandler is a route for /unread_messages that returns all messages sent to the logged in user in JSON format
func MessagesHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	currentUser, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var recentMessages []Message
	db.Where("receiver_id = ?", currentUser.ID).Order("created_at desc").Limit(10).Find(&recentMessages)
	if len(recentMessages) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
		return
	}

	messagesJSON, err := json.Marshal(recentMessages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(messagesJSON)
}

// PastMessagesHandler is a route for /past_messages/{id} that returns all messages sent by either the logged in user or the user with the id in JSON format
func PastMessagesHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	receiverID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	currentUser, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var results []Message
	messagesQuery := "(receiver_id = ? and sender_id = ?) or (receiver_id = ? and sender_id = ?)"
	chatMessages := db.Where(messagesQuery, currentUser.ID, receiverID, receiverID, currentUser.ID)
	if res := chatMessages.Limit(20).Order("created_at desc").Find(&results); res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}

	resultsJSON, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resultsJSON)
}

// ChatHandler is a route for /message/{id} that displays the chat messaging page between the logged in user and the user with the id
func ChatHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	currentUser, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	receiverID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if currentUser.ID == receiverID {
		http.Error(w, "You cannot message yourself", http.StatusUnauthorized)
		return
	}

	messagesQuery := "sender_id = ? and receiver_id = ? and read = ?"
	unreadMessages := db.Model(&Message{}).Where(messagesQuery, receiverID, currentUser.ID, false)
	if res := unreadMessages.UpdateColumn("read", true); res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusUnauthorized)
		return
	}

	var user User
	if res := db.Find(&user, receiverID); res.Error != nil {
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles("templates/boilerplate/nothing_boilerplate.html",
		"templates/navbar.html", "templates/chat.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t.Execute(w, MessageTemplateType{
		UserTemplateType: UserTemplateType{
			CurrentUser:    currentUser,
			HasCurrentUser: true,
		},
		UserID: receiverID,
	})
}
