package Init

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	MongoDB struct {
		Connection string `json:"connection"`
		Database   string `json:"database"`
		Collection string `json:"collection"`
	} `json:"mongodb"`
	RabbitMQ struct {
		Connection string `json:"connection"`
	} `json:"rabbitmq"`
	JWT struct {
		Secret string `json:"secret"`
	} `json:"jwt"`
	App struct {
		Url string `json:"url"`
	} `json:"app"`
}

var SetConfig Config

func LoadConfig(env string) error {
	filename := "conf_dev.json"
	if env == "production" {
		filename = "conf_prod.json"
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &SetConfig)

	return nil
}
