package user

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// var (
// 	userInfoAPI      = map[string]string{"metod": "GET", "url": "user/lass"}
// )

// var Appid = "10000"
// var App_Token = "123456789"
// var userBaseInfoAPI = []string{"post", ""}
var (
	userLassAPI      = []string{"GET", "http://localhost:8800/user/lass"}
	ErrorHttpRequest = errors.New("http request error")

	// authCodeAPI = []string{}
)

type PublicResponse struct {
	Backcode int    `json:"backcode"`
	Backmsg  string `json:"backmsg"`
	Backdata User   `json:"backdata"`
}
type User struct {
	UserId string `json:"userId"`
	Name   string `json:"uname"`
	Icon   string `json:"head"`
}

// func init() {
// }

// func GetUserInfo(token string) (error, interface{}) {
// 	hclient := &http.Client{}
// 	req, err := http.NewRequest(userInfoAPI["method"], userInfoAPI["url"], nil)
// 	req.Header.Set("Token", token)
// 	if err != nil {
// 		return err, nil
// 	}
// 	resp, err := hclient.Do(req)
// 	if err != nil {
// 		return err, nil
// 	}
// 	defer resp.Body.Close()
// 	buffer, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err, nil
// 	}
// 	var backData muserResponse
// 	if err := json.Unmarshal(buffer, &backData); err != nil {
// 		return err, nil
// 	}
// 	// backData.()
// 	return nil, backData.Backdata
// }

// func userBaseInfo() {
// 	`{""}`
// 	do(userBaseInfoAPI[0], userBaseInfoAPI[1], strings.NewReader())
// }

// func doLogin()  {
// 	do('')
// }

// func getAuthCode() (string, error) {
// 	result, err := do(authCodeAPI[0], authCodeAPI[1], nil, nil)
// 	if err != nil {
// 		return "", ErrorHttpRequest
// 	}
// 	if result.backcode != 1001 {
// 		return "", errors.New(result.backmsg)
// 	}
// 	return result.backdata.(string), nil
// }

type LoginRequest struct {
	Token string `json:"token"`
}

var (
	ErrorParamSyntax = errors.New("params syntax error")
)

func GetUserLass(token string) (u *User, err error) {

	header := map[string][]string{}
	header["Token"] = []string{token}
	header["Content-Type"] = []string{"application/json"}
	result, err := do(userLassAPI[0], userLassAPI[1], nil, header)
	if err != nil {
		err = ErrorHttpRequest
		return
	}
	if result.Backcode != 1001 {
		err = errors.New(result.Backmsg)
		return
	}
	u = &result.Backdata
	return
}

func do(method string, url string, body io.Reader, header http.Header) (result *PublicResponse, err error) {
	request, err := http.NewRequest(method, url, body)
	request.Header = header
	if err != nil {
		return
	}
	response, err := client().Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return
	}
	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	log.Print(string(buffer))
	if err = json.Unmarshal(buffer, &result); err != nil {
		return
	}
	return result, nil
}

func client() *http.Client {
	return &http.Client{}
}

// func sign() string {
// 	h := md5.New()
// 	h.Write([]byte(Appid + App_Token))
// 	return fmt.Sprint(h)
// }
