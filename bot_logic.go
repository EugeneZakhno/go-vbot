package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	token      = "7925094881:AAFEsIkrN_iVNkS9hgZ4geX9_9YKxnEXfl8"
	sheetID    = "13n4_A9X0iD6Vuyfc4CpDXe_KFvyg-XV5E5IQiGLSEVE"
	webAppURL  = "https://script.google.com/macros/s/AKfycbx8JxNbTgyZ8_DplCdaGmzlooZZYkiPlipW5hj4E-CNQyDqyfjCFV1sWgUOKfnZS9ES/exec"
	portalIDRx = `^[A-Z0-9]{32}$`
	phoneRx    = `^[0-9]{12}$`
	//passwordRx = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*()_+={}\[\]:;"'<>,.?\/\\|-]{8,20}$`
	balanceCmd = "balance"

	startMsg = "<strong>Авторизация:</strong> \n 🪙<strong>Шаг 1:</strong> Пришлите мне <strong>portalId</strong>, \n который получили от тех.поддержки. \n (Пример: RT3C4E58636929057532709E1B39OPR1)"
	loginMsg = "📱<strong>Шаг 2:</strong> Введите логин: (Пример: 754861154414)"
	passMsg  = "🗝️<strong>Шаг 3:</strong> Введите свой пароль в виде Spoiler:\n (Пример:<tg-spoiler> gRteS1Rb </tg-spoiler>) \n *Ваш пароль маскируется в виде ********,\n и не будет известен никому!"
	timeMsg  = "🕒Время сессии 15 минут (по умолчанию), \n по истечении времени сессии кнопка работать не будет.🕒"
	yesMsg   = "🎉Позравляем! Вы авторизованы."
	wrongMsg = "Что-то пошло не так... Давайте начнем сначала! \n Пробуем /start"

	urlAddressHostKZ = "https://openapi.mypay.kz/api/v4/"
)

var (
	portalIDRegexp = regexp.MustCompile(portalIDRx)
	phoneRegexp    = regexp.MustCompile(phoneRx)
	//passwordRegexp = regexp.MustCompile(passwordRx)
)

// These structs are minimal and just represent what's needed for this example.  You'll likely
// want to expand them to include other fields from the Telegram API.

type Update struct {
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	Date int    `json:"date"`
	From User   `json:"from"`
}

type CallbackQuery struct {
	From User   `json:"from"`
	Data string `json:"data"`
}

type Chat struct {
	ID int64 `json:"id"`
}
type User struct {
	ID int `json:"id"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// sendTelegramMessage sends a message to the Telegram API.
func sendTelegramMessage(chatID int64, text string, replyMarkup *InlineKeyboardMarkup) error {
	botURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	reqBody := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}
	if replyMarkup != nil {
		reqBody["reply_markup"] = replyMarkup
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(botURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close() // Close response once it's handled

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// This func replaces writeValueEditSheetDebug(e) in Google Apps Script
func logRequest(r *http.Request) {

	rBody, err := r.GetBody() // Get Body io.ReadCloser and reassign to r.Body
	if err != nil {
		log.Println("Error getting request body:", err)
		return
	}
	defer rBody.Close()

	bodyBytes, _ := io.ReadAll(rBody)
	log.Println("Request Body:", string(bodyBytes))

}

// doPost handles incoming webhook requests.
func doPost(w http.ResponseWriter, r *http.Request) {

	logRequest(r)

	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var chatID int64
	var text string
	var date int
	var userID int64 // User's ID

	if update.Message != nil {
		chatID = update.Message.Chat.ID
		text = update.Message.Text
		date = update.Message.Date
		userID = int64(update.Message.From.ID) //Store UserID
		fmt.Println(userID)
	} else if update.CallbackQuery != nil {
		chatID = int64(update.CallbackQuery.From.ID)
		if update.CallbackQuery.Data == balanceCmd {
			// dat := getVirtualAccountlist() // Implement getVirtualAccountlist
			dat := "Balance information would go here" // Placeholder
			sendTelegramMessage(chatID, dat, nil)      // Send balance info
		}
		return // Exit early for callback queries after handling balance
	}

	currentDate := time.Now().Format(time.RFC3339)

	// Handle different commands/inputs
	var err error
	switch {
	case text == "/start":
		err = sendTelegramMessage(chatID, startMsg, nil)
		writeValueEditSheetMessages(chatID, "", "", "", currentDate, date, text) // Moved function call to avoid unnecessary logging

	case portalIDRegexp.MatchString(text):
		err = sendTelegramMessage(chatID, loginMsg, nil)
		writeValueEditSheetMessages(chatID, "", "", "", currentDate, date, text)
		writeValueEditSheetAuthorization(chatID, text) //Added an extra param

	case phoneRegexp.MatchString(text):
		err = sendTelegramMessage(chatID, passMsg, nil)
		writeValueEditSheetMessages(chatID, "", "", "", currentDate, date, text)
		writeValueEditSheetAuthorization(chatID, text)

		//case passwordRegexp.MatchString(text):

		inlineKeyboard := InlineKeyboardMarkup{
			InlineKeyboard: [][]InlineKeyboardButton{
				{{Text: "Узнать баланс", CallbackData: balanceCmd}},
			},
		}

		err = sendTelegramMessage(chatID, timeMsg, &inlineKeyboard) //Added inline keyboard
		maskedPass := maskPassword(text)

		writeValueEditSheetMessages(chatID, "", "", "", currentDate, date, maskedPass)

		//getDataFromTable(chatID, text) // Implement getDataFromTable
		//send(chatID)                 // Implement send

	default:
		err = sendTelegramMessage(chatID, wrongMsg, nil)

	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// Implement more robust error handling, such as logging
	}
}

func maskPassword(password string) string {
	return strings.Repeat("*", len(password))
}

// Placeholder functions - You'll need to implement these to interact with Google Sheets/other services
func writeValueEditSheetMessages(chatID int64, firstname, lastname, username, currentDate string, date int, text string) {
	// Implement actual sheet writing logic here
	log.Println("writeValueEditSheetMessages called:", chatID, firstname, lastname, username, currentDate, date, text)
}

func writeValueEditSheetAuthorization(chatID int64, text string) {
	// Implement actual sheet writing logic here
	log.Println("writeValueEditSheetAuthorization called:", chatID, text)
}
