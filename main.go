package main

import (
	"flag"
	"go_stress_test/config"
	"go_stress_test/logic"
	"log"
)

var (
	confFile = flag.String("confFile", "config/go_stress_test.yml", "Configuration file")
	csvFile  = flag.String("csvFile", "", "请输入csv格式的文件")
)

func main() {
	flag.Parse()

	if *csvFile == "" {
		log.Fatalln("参数不正确！请输入要解析的csv格式的文件，如-csvFile=xxx.csv")
		return
	}

	config.LoadConfig(*confFile)

	csvSlice := logic.ParseCSVFile(*csvFile)

	logic.SimulateLogin(csvSlice)

	//发心跳包的
	logic.SimulateHeartBeat(csvSlice)

	//select {
	//case <-time.After(10 * time.Second):
	//	log.Println("send over")
	//}
}
