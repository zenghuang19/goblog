package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/view"
	"net/http"
)

// AuthController 处理静态页面
type AuthController struct {
}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "auth.register")
}

// DoRegister 注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}
	// 2.表单规则
	errs := requests.ValidateRegistrationForm(_user)
	if len(errs) > 0 {
		//验证未通过
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		//  验证通过 —— 入库，并跳转到首页
		_user.Create()
		if _user.ID > 0 {
			// 登录用户并跳转到首页
			flash.Success("恭喜注册成功")
			auth.Login(_user)
			http.Redirect(w,r,"/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建用户失效，请联系管理员")
		}
	}
}

// Login 登录表单
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "auth.login")
}

func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request)  {
	// 1.初始化表单
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")


	//2.尝试登录
	if err := auth.Attempt(email,password);err == nil {
		//登录成功
		flash.Success("欢迎回来！")
		http.Redirect(w,r,"/", http.StatusFound)
	}else {
		// 失败
		view.RenderSimple(w,view.D{
			"Error": err.Error(),
			"Email": email,
			"Password": password,
		},"auth.login")
	}
}

func (*AuthController) Logout(w http.ResponseWriter,r *http.Request)  {
	auth.Logout()
	flash.Success("您已退出成功")
	http.Redirect(w,r,"/",http.StatusFound)
}
