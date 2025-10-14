package blueconfigclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Url   string
	Token string
}

type Response struct {
	Result any
	Error  string
}

func NewClient(url, token string) *Client {
	return &Client{
		Url:   url,
		Token: token,
	}
}

func (c *Client) makecall(endpoint, method string) (*Response, error) {

	if c.Url == "" {
		return nil, errors.New("url is empty")
	}

	var finalurl = c.Url + endpoint

	parsedUrl, err := url.Parse(finalurl)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		query := parsedUrl.Query()
		query.Add("token", c.Token)
		parsedUrl.RawQuery = query.Encode()
	}

	req := &http.Request{
		Method: method,
		URL:    parsedUrl,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("response body:", string(body))
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetNodes(p string) ([]string, error) {
	response, err := c.makecall(p, "GET")
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	if response.Result == nil {
		return []string{}, nil
	}

	nodes := response.Result.([]string)
	return nodes, nil
}

func (c *Client) GetProps(p string) ([]string, error) {
	p, _ = url.JoinPath(p, "/props")
	response, err := c.makecall(p, "GET")
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	if response.Result == nil {
		return []string{}, nil
	}

	props := response.Result.([]string)
	return props, nil
}

func (c *Client) GetValue(p string) (string, error) {
	p, _ = url.JoinPath(p, "/value")
	response, err := c.makecall(p, "GET")
	if err != nil {
		return "", err
	}

	if response.Error != "" {
		return "", errors.New(response.Error)
	}

	if response.Result == nil {
		return "", nil
	}

	value := response.Result.(string)
	return value, nil
}

func (c *Client) GetValues(p string) (map[string]string, error) {
	p, _ = url.JoinPath(p, "/values")
	response, err := c.makecall(p, "GET")
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	if response.Result == nil {
		return map[string]string{}, nil
	}

	var result = make(map[string]string)
	values := response.Result.(map[string]any)

	for key, value := range values {
		if value == nil {
			result[key] = ""
			continue
		}
		result[key] = value.(string)
	}

	return result, nil
}

func (c *Client) SimpleGet(p string) (*Response, error) {
	return c.makecall(p, "GET")
}

func (c *Client) GetValueInt(p string, or ...int) (int, error) {
	value, err := c.GetValue(p)
	if err != nil {
		return 0, err
	}

	if value == "" {
		if len(or) > 0 {
			return or[0], nil
		}
		return 0, errors.New("value not found")
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (c *Client) GetValueFloat(p string, or ...float64) (float64, error) {
	value, err := c.GetValue(p)
	if err != nil {
		return 0, err
	}

	if value == "" {
		if len(or) > 0 {
			return or[0], nil
		}
		return 0, errors.New("value not found")
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

func (c *Client) GetValueBool(p string, or ...bool) (bool, error) {
	value, err := c.GetValue(p)
	if err != nil {
		return false, err
	}

	if value == "" {
		if len(or) > 0 {
			return or[0], nil
		}
		return false, errors.New("value not found")
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}

	return b, nil
}

func (c *Client) GetValueTS(p string, or ...time.Time) (time.Time, error) {
	value, err := c.GetValue(p)
	if err != nil {
		return time.Time{}, err
	}

	if value == "" {
		if len(or) > 0 {
			return or[0], nil
		}
		return time.Time{}, errors.New("value not found")
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func (c *Client) GetValueDuration(p string, or ...time.Duration) (time.Duration, error) {
	value, err := c.GetValue(p)
	if err != nil {
		return 0, err
	}

	if value == "" {
		if len(or) > 0 {
			return or[0], nil
		}
		return 0, errors.New("value not found")
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}

	return d, nil
}

func (c *Client) GetValueUnixTS(p string, or ...int64) (int64, error) {
	value, err := c.GetValue(p)
	if err != nil {
		return 0, err
	}

	if value == "" {
		if len(or) > 0 {
			return or[0], nil
		}
		return 0, errors.New("value not found")
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (c *Client) GetValueStruct(p string, target any) error {
	value, err := c.GetValue(p)
	if err != nil {
		return err
	}

	if value == "" {
		return errors.New("value not found")
	}

	if !strings.HasPrefix(value, "{") && !strings.HasPrefix(value, "[") {
		return errors.New("value is not a struct")
	}

	return json.Unmarshal([]byte(value), target)
}

func (c *Client) ParseValues(p string, target any) error {
	values, err := c.GetValues(p)
	if err != nil {
		return err
	}

	//convert values into a struct
	jsoned, _ := json.Marshal(values)

	return json.Unmarshal(jsoned, target)
}
