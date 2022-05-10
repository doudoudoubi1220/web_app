package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"web_app/models"
)

//把每一步数据库操作封装成函数
//待logic层根据业务需求调用

const secret = "xsq"

// CheckUserExist 检查指定用户名的用户是否存在数据库中是否重复
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	fmt.Println(count)
	if count > 0 {
		return errors.New("用户已存在")
	}
	fmt.Println("rrr")
	return

}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	fmt.Println("acs")
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行SQL语句入库
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassWord := user.Password
	//执行SQL语句入库
	sqlStr := ` select user_id, username,password from user where username = ?`

	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return errors.New("用户不存在")
	}
	if err != nil {
		return err
	}
	//判断密码是否正确
	password := encryptPassword(oPassWord)
	if password != user.Password {
		return errors.New("密码错误")
	}
	return
}

func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
