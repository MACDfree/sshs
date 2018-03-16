package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/MACDfree/sshs/common"
	"gopkg.in/yaml.v2"
)

// Session 表示ssh连接信息
type Session struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

// readConfigData 解析yaml文件内容返回session集合
func readConfigData(data []byte) map[string]Session {
	sessions := make(map[string]Session)
	err := yaml.Unmarshal(data, &sessions)
	common.CheckError(err)
	return sessions
}

func writeConfigStr(sessions map[string]Session) string {
	d, err := yaml.Marshal(sessions)
	common.CheckError(err)
	return string(d)
}

// ReadConfig 读取配置文件并解析成map
func ReadConfig() (map[string]Session, bool) {
	if !common.CheckFileIsExist(common.ConfigPath) {
		return make(map[string]Session), false
	}
	content, err := ioutil.ReadFile(common.ConfigPath)
	common.CheckError(err)
	return readConfigData(content), true
}

// WriteConfig 将sessions写会配置文件中
func WriteConfig(sessions map[string]Session) {
	str := writeConfigStr(sessions)
	file, err := os.OpenFile(common.ConfigPath, os.O_RDWR|os.O_CREATE, 0600)
	common.CheckError(err)
	defer file.Close()
	_, err = io.WriteString(file, str)
}
