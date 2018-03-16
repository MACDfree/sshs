package common

import (
	"log"
	"os"
	"os/user"
)

// ConfigPath 配置文件路径
var ConfigPath string

func init() {
	if ConfigPath == "" {
		ConfigPath = HomePath() + "/.sshs.yml"
	}
}

// CheckError 检测是否有异常，如有则直接停止应用
func CheckError(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

// HomePath 返回当前用户的home路径
func HomePath() string {
	u, err := user.Current()
	CheckError(err)
	return u.HomeDir
}

// CheckFileIsExist 判断文件是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
