package logic

import (
	"fmt"
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

//存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	//1.判断用户存不存在
	fmt.Println("ccc")
	if err := mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询错误
		return err
	}
	fmt.Println("aaaa")
	//2.生成uid
	userID := snowflake.GenID()
	fmt.Println("雪花")

	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	fmt.Println("构造实例")
	//3.保存到数据库
	return mysql.InsertUser(user)

}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
