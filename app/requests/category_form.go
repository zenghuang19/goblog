package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/category"
)

func ValidateCategoryForm(data category.Category) map[string][]string {
	// 定制规则
	rules := govalidator.MapData{
		"name": []string{"required", ",min_cn:2", "max_cn:8", "not_exists:categories,name"},
	}

	//错误消息
	messages := govalidator.MapData{
		"name" : []string{
			"required:分类名称未必填",
			"min_cn:名称长度至少2个字",
			"max_cn:名称长度不能超过8个字符",
		},
	}

	// 配置初始化
	opts := govalidator.Options{
		Data: &data,
		Rules: rules,
		TagIdentifier: "valid",
		Messages: messages,
	}

	//开始验证
	return govalidator.New(opts).ValidateStruct()
}
