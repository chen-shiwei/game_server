package user

import (
	"config"
	"fmt"
	"net/http"
	"strings"
)

func FetchUser(token string) interface{} {
	userServer := config.GetString("userServer")
	if !strings.HasPrefix(userServer, "http://") {
		userServer = fmt.Sprintf("http://%s", userServer)
	}
	if !strings.HasSuffix(userServer, "/") {
		userServer = fmt.Sprintf("%s/", userServer)
	}
	client, _ := http.NewRequest("POST", userServer, nil)
	return client
}
