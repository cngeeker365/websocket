package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

var Conf *Config

type Config struct {
	Redis 			Redis	`yaml: "redis"`
	Pages  			[]Pages	`yaml: "pages"`
}

type Redis struct {
	Pool 			Pool	`yaml: "pool"`
}

type Pool struct {
	Host 			string	`yaml:"host"`
	MaxIdle			int		`yaml:"maxIdle"`
	MaxActive		int		`yaml:"maxActive"`
	IdleTimeout		int		`yaml:"idleTimeout"`
	Wait			int		`yaml:"wait"`
	MaxConnLifeTime	int		`yaml:"maxConnLifeTime"`
}

type Pages struct {
	Page 			Page	`yaml: "page"`
}

type Page struct {
	Name			string	`yaml: "name"`
	Heartbeat		int		`yaml: "heartbeat"`
	Key				string	`yaml: "key"`
}

func (c *Config) getConf(path string) *Config {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.UnmarshalStrict(yamlFile, c)
	return c
}

func (c *Config) RedisPool() *Pool{
	return &c.Redis.Pool
}

func (c *Config) Page(name string) *Page{
	for i:=0; i< len(c.Pages);i++{
		if strings.EqualFold(c.Pages[i].Page.Name, name){
			return &c.Pages[i].Page
		}
	}
	return nil
}

func init() {
	var c Config
	Conf = c.getConf(`conf.yaml`)
}

func main()  {
	conf := Conf
	data, _ := json.Marshal(conf)
	fmt.Println(string(data))

	var a = conf.Redis.Pool.IdleTimeout
	fmt.Printf("%T -- %v",a,a)
}
