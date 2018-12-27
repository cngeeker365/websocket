package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"

	//"time"
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
	Conf = c.getConf(`D:\goProjects\wserverYGZ\src\config\test.yaml`)
}

func main()  {
	//user := &ucenter.User{ID:123,Name:"lalala",Password:"123456"}
	//s,_ := json.Marshal(user)
	//fmt.Println(string(s))
	//
	//fmt.Println(time.Now().Unix())

	//t:=time.Now()
	//	//fmt.Println(t)
	//	//s:=fmt.Sprintf("%d-%d-%d %d:%d:00", t.Year(),t.Month(),t.Day(),t.Hour(),t.Minute()+1)
	//	//fmt.Println(s)
	//	//fmt.Println(t.Round(time.Duration(time.Second)))
	//	//fmt.Println(  (t.Unix()-int64(60*t.Second())))

	//var c Config
	//conf := c.getConf(`D:\goProjects\wserverYGZ\src\config\conf.yaml`)
	conf := Conf
	data, _ := json.Marshal(conf)
	fmt.Println(string(data))

	//fmt.Println(conf.Pages[1].Page.Name)
	//fmt.Println(len(conf.Pages))
	var a = conf.Redis.Pool.IdleTimeout
	fmt.Printf("%T -- %v",a,a)

}
