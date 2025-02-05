package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	setWebhook()
	getMe()
	deleteWebhook()
	getWebhookInfo()
	setWebhookWithUpdates()
}

const (
	token     = "7925094881:AAFEsIkrN_iVNkS9hgZ4geX9_9YKxnEXfl8"
	sheetId   = "13n4_A9X0iD6Vuyfc4CpDXe_KFvyg-XV5E5IQiGLSEVE"
	webAppUrl = "https://script.google.com/macros/s/AKfycbx8JxNbTgyZ8_DplCdaGmzlooZZYkiPlipW5hj4E-CNQyDqyfjCFV1sWgUOKfnZS9ES/exec"
)

func setWebhook() {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", token, webAppUrl)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error setting webhook:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func getMe() {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting bot info:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func deleteWebhook() {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/deleteWebhook", token)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error deleting webhook:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func getWebhookInfo() {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getWebhookInfo", token)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting webhook info:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func setWebhookWithUpdates() {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)

	updates := []string{
		"message",
		"edited_message",
		"channel_post",
		"edited_channel_post",
		"message_reaction",
		"message_reaction_count",
		"callback_query",
		"poll",
		"poll_answer",
		"my_chat_member",
		"chat_member",
		"chat_join_request",
	}

	payload := map[string]interface{}{
		"url":             webAppUrl,
		"allowed_updates": updates,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error setting webhook with updates:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
