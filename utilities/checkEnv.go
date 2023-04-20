package utilities

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	TelegramToken string
	ChatGptToken  string
	AccessToken   string
)

func CheckEnvFile() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	LookupEnv(&TelegramToken, "TELEGRAM_TOKEN")
	LookupEnv(&ChatGptToken, "CHATGPT_TOKEN")
	LookupEnv(&AccessToken, "ACCESS_TOKEN")
}

func LookupEnv(out *string, key string, defaultVal ...string) {
	val, exist := os.LookupEnv(key)
	if !exist {
		if len(defaultVal) > 0 {
			*out = defaultVal[0]
		} else {
			panic(key + " not found in .env file")
		}
	}
	*out = val
}
