package alice

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

//Global function
func AskQuestion(question string) (string, error) {
	queryParts := strings.Split(question, ", ")
	if len(queryParts) != 2 {
		return "", errInvalidQuestionFormat
	} else if queryParts[0] == "Алиса" {
		return askAlice(queryParts[1])
	}
	return askFriend(queryParts[0], queryParts[1])
}

var (
	errUndefinedQuestion     = errors.New("неизвестный вопрос")
	errInvalidQuestionFormat = errors.New("неверный формат вопроса")
)

var (
	database = map[string]string{
		"Алексей": "Калининград",
		"Алина":   "Красноярск",
		"Артём":   "Владивосток",
		"Дима":    "Челябинск",
		"Егор":    "Пермь",
		"Коля":    "Красноярск",
		"Миша":    "Москва",
		"Петя":    "Михайловка",
		"Сергей":  "Омск",
		"Соня":    "Москва",
	}
	
	uTCOffset = map[string]int{
		"Владивосток":     10,
		"Волгоград":       3,
		"Воронеж":         3,
		"Екатеринбург":    5,
		"Казань":          3,
		"Калининград":     2,
		"Краснодар":       3,
		"Красноярск":      7,
		"Москва":          3,
		"Нижний Новгород": 3,
		"Новосибирск":     7,
		"Омск":            6,
		"Пермь":           5,
		"Ростов-на-Дону":  3,
		"Самара":          4,
		"Санкт-Петербург": 3,
		"Уфа":             5,
		"Челябинск":       5,
	}
)

func formatFriendsCount(friendsCount int) string {
	if friendsCount == 1 {
		return "1 друг"
	}
	if friendsCount >= 2 && friendsCount <= 4 {
		return fmt.Sprintf("%d друга", friendsCount)
	}
	return fmt.Sprintf("%d друзей", friendsCount)
}

func whatTime(city string) string {
	offset := uTCOffset[city]
	currentTime := time.Now().UTC()
	cityTime := currentTime.Add(time.Duration(offset) * time.Hour)
	return cityTime.Format("15:04")
}

func whatWeather(city string) (string, error) {
	baseURL := "https://bba4h86i8icpvi5ot97n.containers.yandexcloud.net/go/weather"
	params := url.Values{}
	params.Add("0", "")
	params.Add("M", "")
	params.Add(city, "")
	
	urlInstance, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("сетевая ошибка: %v", err)
	}
	urlInstance.RawQuery = params.Encode()
	
	response, err := http.Get(urlInstance.String())
	if err != nil {
		return "", fmt.Errorf("сетевая ошибка: %v", err)
	}
	
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка на сервере погоды: %v", err)
	}
	content := string(body)
	return content, nil
}

func askFriend(friendName, question string) (string, error) {
	friendCity, exists := database[friendName]
	if !exists {
		return "", fmt.Errorf("у тебя нет друга по имени %s", friendName)
	}
	if question == "ты где?" {
		return fmt.Sprintf("%s в городе: %s", friendName, friendCity), nil
	}
	if question == "который час?" {
		_, exists := uTCOffset[friendCity]
		if !exists {
			return "", fmt.Errorf("не могу определить время в городе: %s", friendCity)
		}
		cityTime := whatTime(friendCity)
		return fmt.Sprintf("Там сейчас %s", cityTime), nil
	}
	// спроси друга о погоде
	if question == "как погода?" {
		_, exists := uTCOffset[friendCity]
		if !exists {
			return "", fmt.Errorf("не могу определить погоду в городе: %s", friendCity)
		}
		return whatWeather(friendCity)
	}
	
	return "", errUndefinedQuestion
}

func askAlice(question string) (string, error) {
	switch question {
	case "сколько у меня друзей?":
		friendsCount := len(database)
		humanizedAnswer := formatFriendsCount(friendsCount)
		return fmt.Sprintf("У тебя %s", humanizedAnswer), nil
	case "кто все мои друзья?":
		friendsNames := make([]string, 0, len(database))
		for friend := range database {
			friendsNames = append(friendsNames, friend)
		}
		sort.Strings(friendsNames)
		friendsString := strings.Join(friendsNames, ", ")
		return fmt.Sprintf("Твои друзья: %s", friendsString), nil
	case "где все мои друзья?":
		friendsCities := make([]string, 0, len(database))
		for _, city := range database {
			friendsCities = append(friendsCities, city)
		}
		sort.Strings(friendsCities)
		citiesString := strings.Join(friendsCities, ", ")
		return fmt.Sprintf("Твои друзья в городах: %s", citiesString), nil
	default:
		return "", errUndefinedQuestion
	}
}
