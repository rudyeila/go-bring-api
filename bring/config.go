package bring

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BaseURL        string        `envconfig:"BRING_BASE_URL" default:"https://api.getbring.com/rest/v2"`
	ClientID       string        `envconfig:"BRING_CLIENT_ID" default:"webApp"`
	Country        string        `envconfig:"BRING_COUNTRY" default:"DE"`
	ApiKey         string        `envconfig:"BRING_API_KEY" default:"cof4Nc6D8saplXjE3h3HXqHH8m7VU2i1Gs0g85Sp"`
	User           string        `envconfig:"BRING_USER"`
	Password       string        `envconfig:"BRING_PASSWORD"`
	DefaultTimeout time.Duration `default:"10s"`
}

func NewConfig(filenames ...string) (Config, error) {
	conf := Config{}
	_ = godotenv.Load(filenames...)

	err := envconfig.Process("", &conf)
	if err != nil {
		return conf, fmt.Errorf("process envconfig :%w", err)
	}

	return conf, nil
}
