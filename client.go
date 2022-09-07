package panasonicfridge

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const UI_VERSION = 4.0

type Client struct {
	Phone        string
	Password     string
	conn         *http.Client
	messageId    int
	mac          string
	userId       string
	realFamilyId int
	familyId     int
}

type RequestParams struct {
	Id        int                    `json:"id"`
	UIVersion float32                `json:"uiVersion"`
	Params    map[string]interface{} `json:"params"`
}

type ResponseData struct {
	Id      int                    `json:"id"`
	Results map[string]interface{} `json:"results"`
	Error   map[string]interface{} `json:"error"`
}

func NewClient(phone string, password string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		Phone:    phone,
		Password: password,
		conn: &http.Client{Jar: jar, Timeout: 10 * time.Second, Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
		messageId: 1,
		mac:       getMacAddress(),
	}
}

func (c *Client) GetToken() (string, error) {
	params := &RequestParams{
		Id:        c.getNextMessageId(),
		UIVersion: UI_VERSION,
		Params: map[string]interface{}{
			"usrId": c.Phone,
		},
	}
	response, err := c.sendRequest(TOKEN_URL, params)
	if err != nil {
		return "", err
	}
	return response.Results["token"].(string), nil
}

func (c *Client) Login() error {
	token, err := c.GetToken()
	if nil != err {
		return err
	}
	params := &RequestParams{
		Id:        c.getNextMessageId(),
		UIVersion: UI_VERSION,
		Params: map[string]interface{}{
			"telId":          c.mac,
			"usrId":          c.Phone,
			"checkFailCount": 1,
			"pwd":            encodePassword(c.Password, c.Phone, token),
		},
	}
	response, err := c.sendRequest(LOGIN_URL, params)
	if err != nil {
		return err
	}
	c.userId = response.Results["usrId"].(string)
	c.realFamilyId = int(response.Results["realFamilyId"].(float64))
	c.familyId = int(response.Results["familyId"].(float64))
	return nil
}

func (c *Client) GetDevices() (*ResponseData, error) {
	params := &RequestParams{
		Id:        c.getNextMessageId(),
		UIVersion: UI_VERSION,
		Params: map[string]interface{}{
			"familyId":     fmt.Sprintf("%d", c.familyId),
			"usrId":        c.userId,
			"realFamilyId": fmt.Sprintf("%d", c.realFamilyId),
		},
	}

	return c.sendRequest(DEVICES_URL, params)
}

func (c *Client) sendRequest(url string, params interface{}) (*ResponseData, error) {
	payload, err := json.Marshal(params)
	if nil != err {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if nil != err {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Host", HOST)
	request.Header.Add("User-Agent", "SmartApp")
	response, err := c.conn.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result := &ResponseData{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetDeviceStatus(deviceId string) (*ResponseData, error) {
	sToken := getSToken(deviceId)
	params := map[string]interface{}{
		"id":       c.getNextMessageId(),
		"deviceId": deviceId,
		"token":    sToken,
		"usrId":    c.userId,
	}

	response, err := c.sendRequest(GET_STATUS_URL, params)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) SetDeviceStatus(deviceId string, data interface{}) (*ResponseData, error) {
	sToken := getSToken(deviceId)
	params := map[string]interface{}{
		"id":       c.getNextMessageId(),
		"deviceId": deviceId,
		"token":    sToken,
		"usrId":    c.userId,
		"params":   data,
	}
	return c.sendRequest(SET_STATUS_URL, params)
}

func (c *Client) getNextMessageId() int {
	id := c.messageId
	c.messageId += 1
	return id
}
