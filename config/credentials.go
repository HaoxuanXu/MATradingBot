package config

import (
	"strings"

	"github.com/HaoxuanXu/MATradingBot/config/credentials/live"
	"github.com/HaoxuanXu/MATradingBot/config/credentials/paper"
)

type Credentials struct {
	API_KEY    string `json:"api_key"`
	API_SECRET string `json:"api_secret"`
	BASE_URL   string `json:"base_url"`
}

func GetCredentials(accountType, serverType string) Credentials {
	var credentials Credentials
	if strings.ToLower(accountType) == "live" {
		if strings.ToLower(serverType) == "staging" {
			credentials = Credentials{
				API_KEY:    live.API_KEY_STAGING,
				API_SECRET: live.API_SECRET_STAGING,
				BASE_URL:   live.BASE_URL,
			}
		} else if strings.ToLower(serverType) == "production" {
			credentials = Credentials{
				API_KEY:    live.API_KEY_PROD,
				API_SECRET: live.API_SECRET_PROD,
				BASE_URL:   live.BASE_URL,
			}
		}
	} else if strings.ToLower(accountType) == "paper" {
		credentials = Credentials{
			API_KEY:    paper.API_KEY,
			API_SECRET: paper.API_SECRET,
			BASE_URL:   paper.BASE_URL,
		}
	}
	return credentials
}
