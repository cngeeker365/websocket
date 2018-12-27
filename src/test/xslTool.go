package main

import (
	"encoding/json"
	"fmt"
	"github.com/Luxurioust/excelize"
)

type Point struct {
	DeviceName 	string
	PointName 	string
	Frequency 	string
	Mode 		string
	Name 		string
	RedisKEY 	string
	Expire 		string
}

func main(){
	xlsx, _ := excelize.OpenFile(`C:\Users\taobaibai\Desktop\前后台数据模板.xlsx`)
	index := xlsx.GetSheetIndex("chongtian")
	rows := xlsx.GetRows("chongtian")

	fmt.Println(index)

	data := make([]*Point,15)

	for rn, row:= range rows {
		var point Point
		for cn, col:=range row{
			if rn >1 && cn >1 {
				fmt.Printf("%v-%v-%v,	", rn,cn,col)
				switch cn {
				case 2: point.DeviceName = col
				case 3: point.PointName = col
				case 5:
					point.Name = col
					point.RedisKEY=col
				case 8: point.Frequency = col
				case 9: point.Mode = col
				case 10: point.Expire = col
				}
			}
		}
		if point.Name != ""{
			data = append(data, &point)
		}
		fmt.Println()
	}

	tmp,_ := json.Marshal(data)
	fmt.Println(string(tmp))
}