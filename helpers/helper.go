package helpers

import (
	"backend-technoscape/initializers"
	"log"
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type Config struct {
	API      string
	Username string
	Password string
}

func GetEnv() (Config, error) {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Println("================= Could not load environment variables ", err, " ================¸")
		return Config{}, err
	}

	return Config{
		API:      config.API,
		Username: config.Username,
		Password: config.Password,
	}, nil
}
