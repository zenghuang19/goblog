package password

import (
	"fmt"
	"goblog/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// Hash 使用bcrypt 对密码进行加密
func Hash(password string) string  {
	bytes,err := bcrypt.GenerateFromPassword([]byte(password),14)
	logger.LogError(err)
	fmt.Println(string(bytes))
	return string(bytes)
}

// CheckHash 对比明文密码和数据库的哈希值
func CheckHash(password, hash string)bool  {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	logger.LogError(err)
	return err == nil
}

//IsHashed 判断字符串是否是哈希过的数据
func IsHashed(str string)bool  {
	// bcrypt 加密后的长度等于60
	return len(str) == 60
}