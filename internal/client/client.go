package client

import "database/sql"

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
	PostMawinter(info ServerInfo, categoryID int64, price int64) (*RecordsDetails, error)
}
