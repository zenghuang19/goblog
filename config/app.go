package config

import "goblog/pkg/config"

func init()  {
	config.Add("app", config.StrMap{
		//应用名称
		"name" : config.Env("APP_NAME","Goblog"),

		//当前环境
		"env" : config.Env("APP_ENV", "production"),

		//调试模式
		"debug":config.Env("APP_DEBUG", false),

		//应用服务端端口
		"port":config.Env("APP_PORT", "3000"),

		//加密key
		"key":config.Env("APP_KEY","33446a9dcf9ea060a0a6532b166da32f304af0de"),

		//域名
		"url":config.Env("APP_URL","127.0.0.1:"),
	})
}
