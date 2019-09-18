package main

import (
	"flag"
	"github.com/cihub/seelog"
	"go_stress_test/config"
	"go_stress_test/entity"
	"go_stress_test/logic"
	"log"
	"time"
)

var (
	confFile = flag.String("confFile", "config/go_stress_test.yml", "Configuration file")
	csvFile  = flag.String("csvFile", "", "请输入csv格式的文件，如./go_stress_test -csvFile=xxx.csv")
	IsGenerateFile = flag.Bool("genFile",false,"若要生成压力测试文件，则在命令行上添加-genFile=true，默认不生成")
)

func main() {
	InitLog()

	ch := make(chan *entity.ResponseResults, 10000)

	flag.Parse()

	if *csvFile == "" {
		log.Fatalln("参数不正确！请输入要解析的csv格式的文件，如./go_stress_test -csvFile=xxx.csv")
		seelog.Error("参数不正确！请输入要解析的csv格式的文件，如./go_stress_test -csvFile=xxx.csv")
		return
	}

	config.LoadConfig(*confFile)

	csvSlice := logic.ParseCSVFile(*csvFile)

	logic.SimulateLogin(csvSlice, ch)

	logic.HandleReponseResults(csvSlice,ch,*IsGenerateFile)

	//发心跳包的
	logic.SimulateHeartBeat(csvSlice)

	select {
	case <-time.After(10 * time.Second):
		log.Println("所有用户都已退出！")
	}
}

func InitLog() {
	defer seelog.Flush()

	//加载配置文件
	logger, err := seelog.LoggerFromConfigAsFile("config/log_config.xml")

	if err != nil {
		panic("parse log_config.xml error")
	}

	//替换记录器
	seelog.ReplaceLogger(logger)
}
