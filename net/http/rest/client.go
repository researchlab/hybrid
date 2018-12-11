package rest

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/go-resty/resty"
)

// Page
type Page struct {
	Data      json.RawMessage `json:data`
	Page      int             `json:page`
	PageCount int             `json:pageCount`
	PageSize  int             `json:pageSize`
}

//GetData
func (p *Page) GetData(data interface{}) error {
	return json.Unmarshal(p.Data, data)
}

// R return rest client for rest controller
func R() Client {
	return &client{headers: map[string]string{}}
}

// Client is client api for rest controller
type Client interface {
	SetHeader(header, value string) Client
	Create(rootURL, class string, obj interface{}) (string, error)
	Get(rootURL, class string, id uint, assoc []string) (string, error)
	Query(rootURL, class string, where string, values []string, orderBy []string, page int, pageSize int) (*Page, error)
	InvokeService(rootURL, class, method string, args ...interface{}) (string, error)
}

type client struct {
	headers map[string]string
}

func (p *client) SetHeader(header, value string) Client {
	p.headers[header] = value
	return p
}

// Create post obj to api server.
func (p *client) Create(rootURL, class string, obj interface{}) (string, error) {
	url := fmt.Sprintf("%s/objs/%s", rootURL, class)
	body, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	res, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(body)).
		Post(url)
	if err != nil {
		return "", err
	}
	resBody := string(res.Body())
	if res.StatusCode() != 200 {
		return "", fmt.Errorf("command is failed. %s", resBody)
	}

	return resBody, nil
}

func (p *client) Get(rootURL, class string, id uint, assoc []string) (string, error) {

	var url string
	if len(assoc) > 0 {
		ass := strings.Join(assoc, ",")
		url = fmt.Sprintf("%s/objs/%s/%d?associations=%s", rootURL, class, id, ass)
	} else {
		url = fmt.Sprintf("%s/objs/%s/%d", rootURL, class, id)
	}

	res, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return "", err
	}
	resBody := string(res.Body())
	if res.StatusCode() != 200 {
		return "", fmt.Errorf("get is failed. %s", resBody)
	}

	return resBody, nil
}

func (p *client) Query(rootURL, class string, where string, values []string, orderBy []string, page int, pageSize int) (*Page, error) {
	if page <= 0 {
		page = 0
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	url := fmt.Sprintf("%s/objs/%s?page=%d&pageSize=%d", rootURL, class, page, pageSize)
	if where != "" && values != nil && len(values) > 0 {
		url = fmt.Sprintf("%s&where=%s&values=%s", url, where, strings.Join(values, ","))
	}
	if orderBy != nil && len(orderBy) > 0 {
		url = fmt.Sprintf("%s&ordery=%s", url, strings.Join(orderBy, ","))
	}
	res, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}
	resBody := string(res.Body())
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("query is failed. %s", resBody)
	}
	var data Page
	err = json.Unmarshal([]byte(resBody), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (p *client) InvokeService(rootURL, class, method string, args ...interface{}) (string, error) {

	b, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	body := string(b)
	url := fmt.Sprintf("%s/objs/%s/%s", rootURL, class, method)
	req := resty.R().
		SetHeader("Content-Type", "application/json")
	for k, v := range p.headers {
		req.SetHeader(k, v)
	}
	res, err := req.
		SetBody(body).
		Post(url)
	if err != nil {
		return "", err
	}
	resBody := string(res.Body())
	if res.StatusCode() != 200 {
		return "", fmt.Errorf("%s is failed. %s", method, resBody)
	}
	return resBody, nil
}

func (p *client) Watch(rootURL, class string, id uint, assoc []string, waitKey int64, waitSecond int) chan *WatchEvent {
	var url string
	if len(assoc) > 0 {
		ass := strings.Join(assoc, ",")
		url = fmt.Sprintf("%s/objs/%s/%d?associations=%s&waitKey=%d&waitSecond=%d", rootURL, class, id, ass, waitKey, waitSecond)
	} else {
		url = fmt.Sprintf("%s/objs/%s/%d&waitKey=%d&waitSecond=%d", rootURL, class, id, waitKey, waitSecond)
	}
	wec := make(chan *WatchEvent)
	go func(wec chan *WatchEvent) {
		res, err := resty.R().
			SetHeader("Content-Type", "application/json").
			Get(url)
		if err != nil {
			wec <- &WatchEvent{Error: err.Error()}
			return
		}
		resBody := string(res.Body())
		if res.StatusCode() != 200 {
			wec <- &WatchEvent{Error: fmt.Sprintf("wait is failed. %s", resBody)}
			return
		}
		watchEvent := &WatchEvent{}
		err = json.Unmarshal([]byte(resBody), watchEvent)
		if err != nil {
			wec <- &WatchEvent{Error: err.Error()}
			return
		}
		wec <- watchEvent
	}(wec)
	return wec
}
