package configs

import (
	//_ "github.com/joho/godotenv/autoload"
	"os"
)

func Get(keyName string, defaultValue ...string) string {
	configVars := map[string]string{
		"SHOW_DEBUG_LOG": "false",

		"HIGHLIGHT_START_TAG": "<font color='#1274E6'>",
		"HIGHLIGHT_END_TAG":   "</font>",

		//gRPC port
		"GRPC_PORT": "8080",
	}

	v, isExist := os.LookupEnv(keyName)
	if !isExist {
		if cv, ok := configVars[keyName]; ok {
			return cv
		}

		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}

	return v
}
