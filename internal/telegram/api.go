package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type API struct {
	baseURL string
	token   string
}

func NewAPI(token string) *API {
	return &API{
		baseURL: "https://api.telegram.org/bot",
		token:   token,
	}
}

func (a *API) GetUpdates(offset int) ([]Update, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s/getUpdates?offset=%d", a.baseURL, a.token, offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (a *API) SendMessage(chatID int64, text string) error {
	data := url.Values{}
	data.Set("chat_id", strconv.FormatInt(chatID, 10))
	data.Set("text", text)

	resp, err := http.PostForm(fmt.Sprintf("%s%s/sendMessage", a.baseURL, a.token), data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
