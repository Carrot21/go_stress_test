package config

import (
	"github.com/cihub/seelog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var g_config *config

type config struct {
	HostPort    string `yaml:"hostport"`
	DialTimeout int    `yaml:"dialtimeout"`
	HeartBeat   int    `yaml:"heartbeat"`
}

func LoadConfig(filename string) (conf *config) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		seelog.Error(err)
		return nil
	}

	// Expand env vars
	data = []byte(os.ExpandEnv(string(data)))

	// Decoding config
	if err = yaml.UnmarshalStrict(data, &conf); err != nil {
		seelog.Error(err)
		return nil
	}

	g_config = conf

	seelog.Infof("LoadConfig: %v", *conf)
	return
}

func GetConfig() *config {
	if g_config == nil {
		seelog.Error("CONFIG FILE IS NULL!")
	}
	return g_config
}
