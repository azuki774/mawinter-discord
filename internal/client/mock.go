package client

type mockClientRepo struct {
}

func NewMockClientRepo() *mockClientRepo {
	return &mockClientRepo{}
}

func (c *mockClientRepo) PostMawinter(info ServerInfo, categoryID int64, price int64) (*RecordsDetails, error) {
	return &RecordsDetails{Id: 123, Date: "2000-01-23", Name: "testcategory", CategoryId: categoryID, Price: price}, nil
}
