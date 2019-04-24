package consul

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// TODO: HttpRequest
func NewHttpRequest(service_name string, method string, path string, body io.Reader) (*http.Request, error) {
	u, err := GenerateURL(service_name, path)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, u.String(), body)
}

func NewPostRequest(service_name string, path string, body io.Reader) (*http.Request, error) {
	return NewHttpRequest(service_name, "POST", path, body)
}

func NewGetRequest(service_name string, path string, body io.Reader) (*http.Request, error) {
	return NewHttpRequest(service_name, "GET", path, body)
}

func HttpGet(service_name string, path string) (resp *http.Response, err error) {
	if u, err := GenerateURL(service_name, path); err != nil {
		return nil, err
	} else {
		return http.Get(u.String())
	}
}

func HttpPost(service_name string, path string, contentType string, body io.Reader) (resp *http.Response, err error) {
	if u, err := GenerateURL(service_name, path); err != nil {
		return nil, err
	} else {
		return http.Post(u.String(), contentType, body)
	}
}

func HttpPostForm(service_name string, path string, data url.Values) (resp *http.Response, err error) {
	if u, err := GenerateURL(service_name, path); err != nil {
		return nil, err
	} else {
		return http.PostForm(u.String(), data)
	}
}

func GenerateURL(service_name string, path string) (*url.URL, error) {
	serverInfo := PoolingServiceInfo(service_name)
	if serverInfo == nil {
		return nil, errors.New("no valid server")
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u.Host = fmt.Sprintf("%s:%d", serverInfo.IP, serverInfo.Port)
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u, nil
}
