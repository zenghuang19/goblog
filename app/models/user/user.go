package user

import (
	"goblog/app/models"
	"goblog/pkg/model"
	"goblog/pkg/password"
	"goblog/pkg/route"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email    string `gorm:"type:varchar(255);unique;" valid:"email"`
	Password string `gorm:"type:varchar(255)" valid:"password"`

	// gorm:"-" —— 设置 GORM 在读写时略过此字段，仅用于表单验证
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

// ComparePassword 对比密码是否匹配
func (user *User) ComparePassword(_password string)bool  {
	return password.CheckHash(_password,user.Password)
}

// Link 方法生成用户链接
func (user *User)Link()string  {
	return route.Name2URL("users.show", "id",user.GetStringID())
}

// All 获取所有用户数据
func All()([]User,error)  {
	var users []User
	if err := model.DB.Find(&users).Error;err != nil {
		return users,err
	}

	return users,nil
}