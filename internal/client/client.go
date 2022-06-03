package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

type RecordsDetails struct {
	Id         int64  `json:"id"`
	Date       string `json:"date"`
	CategoryId int64  `json:"categoryID"`
	// Name       string         `json:"name"`
	Price int64          `json:"price"`
	Memo  sql.NullString `json:"memo"`
}

type ServerInfo struct {
	Addr string
	User string // Basic auth User
	Pass string // Basic auth Pass
}

type ClientRepository interface {
	PostMawinter(info *ServerInfo, categoryID int64, price int64) (*RecordsDetails, error)
	DeleteMawinter(info *ServerInfo, ID int64) error
	mustEmbedUnimplementedClient()
}

type clientRepo struct {
	UnimplementedClient
}

type UnimplementedClient struct {
}

func (*UnimplementedClient) mustEmbedUnimplementedClient() {}

type UnsafeElectConsumeService interface {
	mustEmbedUnimplementedClient()
}

func (c *UnimplementedClient) PostMawinter(info *ServerInfo, categoryID int64, price int64) (*RecordsDetails, error) {
	return &RecordsDetails{Id: 123, Date: "2000-01-23", CategoryId: categoryID, Price: price}, nil
}

func (c *UnimplementedClient) DeleteMawinter(info *ServerInfo, ID int64) error {
	return nil
}

func NewClientRepo() *clientRepo {
	return &clientRepo{}
}

func (c *clientRepo) PostMawinter(info *ServerInfo, categoryID int64, price int64) (*RecordsDetails, error) {
	var sendData RecordsDetails = RecordsDetails{CategoryId: categoryID, Price: price}
	sendDataJson, err := json.Marshal(sendData)
	if err != nil {
		Logger.Errorw("failed to Marshal", "data", sendData, "error", err)
		return nil, err
	}

	postaddr := info.Addr + "record/"
	Logger.Infow("server info", "addr", postaddr, "user", info.User, "pass", info.Pass)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", postaddr, bytes.NewReader(sendDataJson))
	if err != nil {
		Logger.Errorw("create new request error", "data", sendData, "error", err)
		return nil, err
	}
	req.SetBasicAuth(info.User, info.Pass)

	res, err := client.Do(req)
	defer func() {
		if res == nil {
			return
		}
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

func (c *clientRepo) DeleteMawinter(info *ServerInfo, ID int64) error {
	deleteaddr := info.Addr + "record/" + strconv.FormatInt(ID, 10)
	Logger.Infow("server info", "addr", deleteaddr, "user", info.User, "pass", info.Pass)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("DELETE", deleteaddr, nil)
	if err != nil {
		Logger.Errorw("create new request error", "error", err)
		return err
	}
	req.SetBasicAuth(info.User, info.Pass)

	res, err := client.Do(req)
	if err != nil {
		Logger.Errorw("received error", "error", err)
		return err
	}

	if res.StatusCode != 204 { // Except No Contents
		Logger.Errorw("received error response", "statusCode", res.StatusCode)
		return fmt.Errorf("error response")
	}

	return nil
}
