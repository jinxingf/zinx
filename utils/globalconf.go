package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
Global Conf for all modules
*/

type GlobalConfig struct {
	// About for Server
	TcpServer ziface.IServer
	Host      string
	TCPPort   int
	Name      string

	// About for zinx
	Version        string // zinx server version
	MaxConn        int    // TCPServer max connection
	MaxPackageSize uint32 // package max size
}

var GlobalConf *GlobalConfig

// Reload load config from file
func (g *GlobalConfig) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalConf)
	if err != nil {
		panic(err)
	}

}

func init() {
	// default config value
	GlobalConf := &GlobalConfig{
		Name:           "ZinxServerAPP",
		Version:        "v0.4",
		TCPPort:        8999,
		Host:           "127.0.0.1",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// try to load config from config file
	GlobalConf.Reload()
}
