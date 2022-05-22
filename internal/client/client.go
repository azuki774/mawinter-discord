package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

type RecordsDetails struct {
	Id         int64          `json:"id"`
	Date       string         `json:"date"`
	CategoryId int64          `json:"categoryID"`
	Name       string         `json:"name"`
	Price      int64          `json:"price"`
	Memo       sql.NullString `json:"memo"`
}

type ServerInfo struct {
	Addr string
	User string // Basic auth User
	Pass string // Basic auth Pass
}

type ClientRepository interface {
	PostMawinter(info *ServerInfo, categoryID int64, price int64) (*RecordsDetails, error)
}

type clientRepo struct {
}

func NewClientRepo() *clientRepo {
	return &clientRepo{}
}

func (c *clientRepo) PostMawinter(info *ServerInfo, categoryID int64, price int64) (*RecordsDetails, error) {
	var sendData RecordsDetails
	sendData = RecordsDetails{CategoryId: categoryID, Price: price}
	sendDataJson, err := json.Marshal(sendData)
	if err != nil {
		Logger.Errorw("failed to Marshal", "data", sendData, "error", err)
		return nil, err
	}

	Logger.Infow("server info", "addr", info.Addr, "user", info.User, "pass", info.Pass)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", info.Addr, bytes.NewReader(sendDataJson))
	if err != nil {
		Logger.Errorw("create new request error", "data", sendData, "error", err)
		return nil, err
	}
	req.SetBasicAuth(info.User, info.Pass)

	res, err := client.Do(req)
	defer func() {
		if res.Body == nil {
			return
		}
		err = res.Body.Close()
		if err != nil {
			Logger.Errorw("failed to close response body", "error", err)
		}
	}()
	if err != nil {
		Logger.Errorw("post data error", "data", sendData)
		return nil, err
	}

	if res.StatusCode != 200 {
		Logger.Errorw("received error response", "statusCode", res.StatusCode)
		return nil, fmt.Errorf("error response")
	}

	var resData RecordsDetails
	body, err := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &resData)
	if err != nil {
		Logger.Errorw("failed to Unmarshal", "data", res.Body, "error", err)
		return nil, err
	}

	return &resData, nil
}
