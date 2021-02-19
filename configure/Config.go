package configure

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

type serverConf struct {
	SvrIp string `yaml:"svr_ip"`
	SvrPort int32 `yaml:"svr_port"`
}

type redisConf struct {
	Ip string `yaml:"ip"`
	Port int `yaml:"port"`
	Passwd string `yaml:"passwd"`
}

type mysqlConf struct {
	Ip string `yaml:"ip"`
	Port int `yaml:"port"`
	UserName string `yaml:"username"`
	Passwd string `yaml:"password"`
	Proto string `yaml:"proto"`
}

type adminConfigure struct {
	Server serverConf `yaml:"server"`
	Redis redisConf `yaml:"redis"`
	MySQL mysqlConf `yaml:"mysql"`
}

var configInstance adminConfigure
var once sync.Once

func ParseConfigure(path string) {
	once.Do(func (){
		yamlFile, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Read file err: %v", err.Error())
		}
		err = yaml.Unmarshal(yamlFile, &configInstance)

		if err != nil {
			log.Fatalf("Parse config err: %v", err.Error())
		}
	})
}

func GetConfigInstance() *adminConfigure {
	return &configInstance
}
