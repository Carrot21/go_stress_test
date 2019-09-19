package main

import (
	"flag"
	"fmt"
	"github.com/cihub/seelog"
	"go_stress_test/config"
	"go_stress_test/entity"
	"go_stress_test/logic"
	"log"
)

var (
	confFile       = flag.String("confFile", "config/go_stress_test.yml", "配置文件，不要添加在命令行上，使用默认即可")
	csvFile        = flag.String("csvFile", "", "请输入csv格式的文件，如./go_stress_test -csvFile=xxx.csv")
	onLineTime     = flag.Int("time", 1, "模拟用户在线时长，单位为分钟，默认1分钟")
	IsGenerateFile = flag.Bool("genFile", false, "若要生成压力测试文件，则在命令行上添加-genFile=true，默认不生成")
)

func main() {
	InitLog()

	ch := make(chan *entity.ResponseResults, 100000)
	connChan := make(chan *entity.UserConnInfo, 100000)

	flag.Parse()

	if *csvFile == "" {
		log.Fatalln("参数不正确！请输入要解析的csv格式的文件，如./go_stress_test -csvFile=xxx.csv")
		seelog.Error("参数不正确！请输入要解析的csv格式的文件，如./go_stress_test -csvFile=xxx.csv")
		return
	}

	config.LoadConfig(*confFile)

	csvSlice := logic.ParseCSVFile(*csvFile)

	println("  		开始压测，请耐心等待见证奇迹的时刻！")

	logic.SimulateLogin(csvSlice, ch, connChan)

	logic.HandleReponseResults(csvSlice, ch, *IsGenerateFile, *onLineTime)

	//发心跳包的
	logic.SimulateHeartBeat(*onLineTime, connChan)

	fmt.Println("  		所有用户都已经成功退出！")
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
