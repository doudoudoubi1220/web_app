package mysql

import (
	"testing"
	"time"
	"web_app/models"
	"web_app/settings"
)

func init() {
	dbCfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "",
		DB:           "bluebell",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdConns:   10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          123,
		AuthorID:    123,
		CommunityID: 1,
		Status:      1,
		Title:       "test",
		Content:     "just a test",
		CreateTime:  time.Time{},
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatal("CreatePost insert record into mysql failed,err: ", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
