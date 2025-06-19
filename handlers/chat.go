package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"morethancoder/t3-clone/db"
	"morethancoder/t3-clone/services"
	"morethancoder/t3-clone/utils"
	"morethancoder/t3-clone/views/components"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	datastar "github.com/starfederation/datastar/sdk/go"
	"github.com/yuin/goldmark"
)

type Data struct {
	Chat []struct {
		Content []struct {
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"content"`
		Role string `json:"role"`
	} `json:"chat"`
	Model string `json:"model"`
}

func randomString() string {
	return uuid.NewString()
}

func GETNewChat(w http.ResponseWriter, r *http.Request) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		utils.Log.Error("No jwt cookie found")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	authRefreshResponse, err := db.Db.AuthRefresh(jwtCookie.Value)
	if err != nil {
		utils.Log.Error(err.Error())
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Error(w, "Failed to auth refresh", http.StatusUnauthorized)
		return
	}
	services.UserSSEHub.ExcuteScript(authRefreshResponse.Record.ID, "window.location.reload()")
	return
}

func POSTChat(w http.ResponseWriter, r *http.Request) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		utils.Log.Error("No jwt cookie found")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	authRefreshResponse, err := db.Db.AuthRefresh(jwtCookie.Value)
	if err != nil {
		utils.Log.Error(err.Error())
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Error(w, "Failed to auth refresh", http.StatusUnauthorized)
		return
	}

	userData := &Data{}
	err = datastar.ReadSignals(r, userData)
	if err != nil {
		utils.Log.Error(err.Error())
		return
	}

	var chatBubbles []templ.Component
	var msgs []services.Message

	for _, chat := range userData.Chat {
		chatBubbles = append(chatBubbles, components.ChatBubble(components.ChatBubbleData{
			ID:        randomString(),
			Company:   strings.Split(userData.Model, "/")[0],
			Role:      chat.Role,
			Type:      chat.Content[0].Type,
			Text:      chat.Content[0].Text,
			Timestamp: time.Now().Format("3:04 PM"),
			Model:     userData.Model,
		}))
		msgs = append(msgs, services.Message{
			Role:    chat.Role,
			Content: []services.Content{services.Content{Type: chat.Content[0].Type, Text: chat.Content[0].Text}},
		})
	}

	chatBubbles = append(chatBubbles, components.ChatBubble(components.ChatBubbleData{
		ID:      randomString(),
		Model:   userData.Model,
		Role:    "loading",
		Company: strings.Split(userData.Model, "/")[0],
	}))

	//rendering the prompt
	services.UserSSEHub.BroadcastFragments(authRefreshResponse.Record.ID,
		components.Chat(chatBubbles),
	)

	chatBubbles = chatBubbles[:len(chatBubbles)-1]

	req := services.OpenRouterRequest{
		Model:    userData.Model,
		Messages: msgs,
	}

	res, err := services.OpenRouter.Request(req)
	if err != nil {
		utils.Log.Error(err.Error())
		chatBubbles = append(chatBubbles, components.ChatBubble(components.ChatBubbleData{
			ID:        randomString(),
			Role:      "assistant",
			Type:      "text",
			Text:      err.Error(),
			Company:   strings.Split(userData.Model, "/")[0],
			Model:     userData.Model,
			Timestamp: time.Now().Format("3:04 PM"),
		}))

		services.UserSSEHub.BroadcastFragments(authRefreshResponse.Record.ID,
			components.Chat(chatBubbles),
		)

		chatBubbles = chatBubbles[:len(chatBubbles)-1]
		return
	}

	for _, choice := range res.Choices {
		var buf bytes.Buffer
		err = goldmark.Convert([]byte(choice.Message.Content.(string)), &buf)
		if err != nil {
			utils.Log.Error(err.Error())
			return
		}
		chatBubbles = append(chatBubbles, components.ChatBubble(components.ChatBubbleData{
			ID:        randomString(),
			Role:      choice.Message.Role,
			Type:      "text",
			Text:      buf.String(),
			Company:   res.Provider,
			Timestamp: time.Now().Format("3:04 PM"),
			Model:     userData.Model,
		}))
		//update userdata chat
		userData.Chat = append(userData.Chat, struct {
			Content []struct {
				Text string "json:\"text\""
				Type string "json:\"type\""
			} "json:\"content\""
			Role string "json:\"role\""
		}{
			Role: "Assistance",
			Content: []struct {
				Text string "json:\"text\""
				Type string "json:\"type\""
			}{
				{
					Text: buf.String(),
					Type: choice.Message.Role,
				},
			},
		})
	}

	updatedChat, err := json.Marshal(userData.Chat)
	if err != nil {
		utils.Log.Error(err.Error())
		return
	}
	services.UserSSEHub.BroadcastSignals(authRefreshResponse.Record.ID, []byte(fmt.Sprintf(`{ chat : %s }`, string(updatedChat))))

	services.UserSSEHub.BroadcastFragments(authRefreshResponse.Record.ID,
		components.Chat(chatBubbles),
	)
	return
}
