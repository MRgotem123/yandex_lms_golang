package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
)

const ()

func main() {
	bot, err := tgbotapi.NewBotAPI(TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Photo != nil {
			// Получаем наибольшее по размеру фото
			photo := update.Message.Photo[len(update.Message.Photo)-1]

			// Получаем файл Telegram
			file, err := bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
			if err != nil {
				log.Println(err)
				continue
			}

			// Скачиваем фото
			url := file.Link(TelegramToken)
			resp, err := http.Get(url)
			if err != nil {
				log.Println("Ошибка скачивания фото:", err)
				continue
			}
			defer resp.Body.Close()

			imageData, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Ошибка чтения фото:", err)
				continue
			}

			// Отправляем в Clarifai
			foodNames, err := recognizeFoodClarifai(imageData)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось распознать еду 😢"))
				log.Println("ERROR:", err)
				continue
			}

			msg := "Похоже, это:\n"
			finalCal := 0
			for _, food := range foodNames {
				cal := estimateCalories(food)
				finalCal += cal
				msg += fmt.Sprintf("• *%s* — ~%d ккал\n", food, cal)
			}
			msg += fmt.Sprintf("Всего: %d ккал\n", finalCal)

			message := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
			message.ParseMode = "Markdown"
			bot.Send(message)
		}
	}
}

func recognizeFoodClarifai(imageData []byte) ([]string, error) {
	client := resty.New()
	apiURL := fmt.Sprintf("https://api.clarifai.com/v2/workflows/food-item-recognition-workflow-96325q/results")

	requestBody := map[string]interface{}{
		"inputs": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"image": map[string]interface{}{
						"base64": encodeToBase64(imageData),
					},
				},
			},
		},
	}

	res, err := client.R().
		SetHeader("Authorization", "Key "+ClarifaiAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(apiURL)

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		log.Println("Clarifai API error:", res.Status(), string(res.Body()))
		return nil, fmt.Errorf("clarifai error: %s", res.Status())
	}

	var result struct {
		Results []struct {
			Outputs []struct {
				Data struct {
					Concepts []struct {
						Name  string  `json:"name"`
						Value float64 `json:"value"`
					} `json:"concepts"`
				} `json:"data"`
			} `json:"outputs"`
		} `json:"results"`
	}

	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		log.Println("Unmarshal error:", string(res.Body()))
		return nil, err
	}

	if len(result.Results) == 0 || len(result.Results[0].Outputs) == 0 {
		return nil, fmt.Errorf("пустой результат")
	}

	concepts := result.Results[0].Outputs[0].Data.Concepts

	var top3 []string
	for i := 0; i < len(concepts) && i < 3; i++ {
		top3 = append(top3, concepts[i].Name)
	}

	if len(top3) == 0 {
		return nil, fmt.Errorf("ничего не найдено")
	}

	return top3, nil
}

func estimateCalories(food string) int {
	// Примерно оцениваем калории
	foodCalories := map[string]int{
		"pizza":     266,
		"burger":    295,
		"apple":     52,
		"banana":    89,
		"salad":     100,
		"sushi":     200,
		"ice cream": 207,
	}

	if val, ok := foodCalories[food]; ok {
		return val
	}

	return 150 // среднее значение по умолчанию
}

// утилита: конвертируем []byte в base64 string
func encodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func portionKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("100 г", "100"),
			tgbotapi.NewInlineKeyboardButtonData("200 г", "200"),
			tgbotapi.NewInlineKeyboardButtonData("300 г", "300"),
			tgbotapi.NewInlineKeyboardButtonData("400 г", "400"),
			tgbotapi.NewInlineKeyboardButtonData("500 г", "500"),
		),
	)
}
