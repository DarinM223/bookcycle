package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strconv"
)

// UnreadMessagesHandler Route: /unread_messages
func UnreadMessagesHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	currentUser, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var recentMessages []Message
	db.Where("receiver_id = ? and read = ?", currentUser.ID, 0).
		Order("created_at desc").Limit(10).Find(&recentMessages)
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

// PastMessagesHandler Route: /past_messages/{id}
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
	res := db.Where("(receiver_id = ? and sender_id = ?) or (receiver_id = ? and sender_id = ?)",
		currentUser.ID, receiverID, receiverID, currentUser.ID).Limit(20).Order("created_at desc").Find(&results)
	if res.Error != nil {
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

// ChatHandler Route: /message/{id}
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

	var result *gorm.DB
	result = db.Model(&Message{}).Where("sender_id = ? and receiver_id = ? and read = ?", receiverID, currentUser.ID, false).UpdateColumn("read", true)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}

	var user User
	result = db.Find(&user, receiverID)
	if result.Error != nil {
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
