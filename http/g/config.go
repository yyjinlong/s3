package g

import (
	"io/ioutil"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

type ServerInfo struct {
	Address string `yaml:"address"`
	Logfile string `yaml:"logfile"`
}

type PostgresInfo struct {
	Master string `yaml:"master"`
	Slave1 string `yaml:"slave1"`
	Slave2 string `yaml:"slave2"`
}

type Settings struct {
	Server   ServerInfo   `yaml:"server"`
	Postgres PostgresInfo `yaml:"postgres"`
}

var (
	lock    = new(sync.RWMutex)
	setting Settings
)

func ParseConfig(file string) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	lock.Lock()
	defer lock.Unlock()

	if err := yaml.Unmarshal(buf, &setting); err != nil {
		panic(err)
	}
}

func Config() Settings {
	lock.RLock()
	defer lock.RUnlock()
	return setting
}
