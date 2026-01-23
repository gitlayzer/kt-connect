package options

import (
	"strings"
)

var customizeKubeConfig = ""
var customizeKtConfig = ""

func GetCustomizeKubeConfig() (string, bool) {
	if len(customizeKubeConfig) > 50 {
		return customizeKubeConfig, true
	}
	return "", false
}

func GetCustomizeKtConfig() (string, bool) {
	if strings.Contains(customizeKtConfig, ":") {
		return customizeKtConfig, true
	}
	return "", false
}
