package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var Logger *zap.Logger

type RecordsDetails struct {
	Id         int64  `json:"id"`
	Date       string `json:"date"`
	CategoryId int64  `json:"category_id"`
	Name       string `json:"category_name"`
	Price      int64  `json:"price"`
	Memo       string `json:"memo"`
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
	return &RecordsDetails{Id: 123, Date: "2000-01-23", CategoryId: categoryID, Name: "category_name", Price: price}, nil
}

func (c *UnimplementedClient) DeleteMawinter(info *ServerInfo, ID int64) error {
	return nil
}

func NewClientRepo() *clientRepo {
	return &clientRepo{}
}

func (c *clientRepo) PostMawinter(info *ServerInfo, categoryID int64, price int64) (*RecordsDetails, error) {
	type sendRecordStruct struct {
		CategoryId int64  `json:"category_id"`
		Price      int64  `json:"price"`
		Form       string `json:"from"`
	}
	sendData := sendRecordStruct{CategoryId: categoryID, Price: price, Form: "mawinter-discord"}
	sendDataJson, err := json.Marshal(sendData)
	if err != nil {
		Logger.Error("failed to Marshal", zap.Error(err))
		return nil, err
	}

	postaddr := info.Addr + "record/"
	Logger.Info("server info", zap.String("addr", postaddr), zap.String("user", info.User), zap.String("pass", info.Pass))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", postaddr, bytes.NewReader(sendDataJson))
	if err != nil {
		Logger.Error("create new request error", zap.Error(err))
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
			Logger.Error("failed to close response body", zap.Error(err))
		}
	}()
	if err != nil {
		Logger.Error("post data error", zap.Error(err))
		return nil, err
	}

	if res.StatusCode != 201 {
		Logger.Error("received error response", zap.Int("statusCode", res.StatusCode))
		return nil, fmt.Errorf("error response")
	}

	var resData RecordsDetails
	body, err := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &resData)
	if err != nil {
		Logger.Error("failed to Unmarshal", zap.Error(err))
		return nil, err
	}

	return &resData, nil
}

func (c *clientRepo) DeleteMawinter(info *ServerInfo, ID int64) error {
	// Not worked because mawinter-server not implemented
	deleteaddr := info.Addr + "record/" + strconv.FormatInt(ID, 10)
	Logger.Info("server info", zap.String("addr", deleteaddr), zap.String("user", info.User), zap.String("pass", info.Pass))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("DELETE", deleteaddr, nil)
	if err != nil {
		Logger.Error("create new request error", zap.Error(err))
		return err
	}
	req.SetBasicAuth(info.User, info.Pass)

	res, err := client.Do(req)
	if err != nil {
		Logger.Error("received error", zap.Error(err))
		return err
	}

	if res.StatusCode != 204 { // Except No Contents
		Logger.Error("received error response", zap.Int("statusCode", res.StatusCode))
		return fmt.Errorf("error response")
	}

	return nil
}
