package logic

import (
	"encoding/csv"
	"fmt"
	"go_stress_test/entity"
	"log"
	"os"
	"strconv"
	"time"
)

func HandleReponseResults(csvSlice [][]string, ch chan *entity.ResponseResults, isGenFile bool, onLineTime int) {
	var (
		processingTime uint64 // 处理总时间
		maxTime        uint64 // 最大时长
		minTime        uint64 // 最小时长
		successNum     uint64 // 成功处理数
		failureNum     uint64 // 处理失败数
	)

	// 定时输出一次计算结果
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		select {
		case <-ticker.C:
			go calculateData(uint64(len(csvSlice)), processingTime, maxTime, minTime, successNum, failureNum, isGenFile, onLineTime)
			ticker.Stop()
		}
	}()

	header()

	close(ch)

	for data := range ch {
		processingTime = processingTime + data.Time

		if maxTime <= data.Time {
			maxTime = data.Time
		}

		if minTime == 0 {
			minTime = data.Time
		} else if minTime > data.Time {
			minTime = data.Time
		}

		// 是否请求成功
		if data.IsSucceed == true {
			successNum = successNum + 1
		} else {
			failureNum = failureNum + 1
		}
	}
}

// 打印表头信息
func header() {
	// 打印的时长都为毫秒
	fmt.Println("───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────")
	result := fmt.Sprintf(" 并发数│ 成功数│ 失败数│   QPS  │最长耗时│最短耗时│平均耗时│  总耗时")
	fmt.Println(result)
	fmt.Println("───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────")

	return
}

// 计算数据
func calculateData(concurrent, processingTime, maxTime, minTime, successNum, failureNum uint64, isGenFile bool, onLineTime int) {
	if processingTime == 0 {
		processingTime = 1
	}

	var (
		qps          float64
		averageTime  float64
		maxTimeFloat float64
		minTimeFloat float64
	)

	// 平均 每个协程成功数*总协程数据/总耗时 (每秒)
	if processingTime != 0 {
		qps = float64(successNum*1e9*concurrent) / float64(processingTime)
	}

	// 平均时长 总耗时/总请求数 纳秒=>毫秒
	if successNum != 0 && concurrent != 0 {
		averageTime = float64(processingTime) / float64(1e6*concurrent)
	}

	// 纳秒=>毫秒
	maxTimeFloat = float64(maxTime) / 1e6
	minTimeFloat = float64(minTime) / 1e6

	//总耗时(纳秒)
	processingTimeNa := float64(processingTime) / 1e6

	// 打印的时长都为毫秒
	table(successNum, failureNum, qps, averageTime, maxTimeFloat, minTimeFloat, processingTimeNa, concurrent)

	if isGenFile {
		generateCSVFile(successNum, failureNum, qps, averageTime, maxTimeFloat, minTimeFloat, processingTimeNa, concurrent)
	}

	fmt.Printf("\n  		模拟用户在线时长为%d分钟，请耐心等待...\n", onLineTime)
}

// 打印表格
func table(successNum, failureNum uint64, qps, averageTime, maxTimeFloat, minTimeFloat, processTimeFloat float64, concurrentNum uint64) {
	// 打印的时长都为毫秒
	result := fmt.Sprintf("%7d│%7d│%7d│%8.2f│%8.2f│%8.2f│%8.2f│%8.2f", concurrentNum, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime, processTimeFloat)
	fmt.Println(result)

	return
}

func generateCSVFile(successNum, failureNum uint64, qps, averageTime, maxTimeFloat, minTimeFloat, processingTimeNa float64, concurrentNum uint64) {
	format := "2006-01-02.15.04"

	timeStr := fmt.Sprintf("%s", time.Now().Format(format))

	File := "stressTest." + timeStr + ".csv"

	newFile, err := os.Create(File)
	if err != nil {
		log.Fatalln("ERR:", err)
	}

	defer func() {
		newFile.Close()
	}()

	newFile.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(newFile)
	header := []string{"并发数", "成功数", "失败数", "QPS", "最长耗时", "最短耗时", "平均耗时", "总耗时"}
	data := [][]string{
		header,
	}

	successNumStr := intToStr(successNum)
	failureNumStr := intToStr(failureNum)
	concurrentNumStr := intToStr(concurrentNum)

	qpsStr := floatToStr(qps)
	averageTimeStr := floatToStr(averageTime)
	maxTimeFloatStr := floatToStr(maxTimeFloat)
	minTimeFloatStr := floatToStr(minTimeFloat)
	processingTimeStr := floatToStr(processingTimeNa)

	content := []string{concurrentNumStr, successNumStr, failureNumStr, qpsStr, maxTimeFloatStr, minTimeFloatStr, averageTimeStr, processingTimeStr}

	data = append(data, content)

	w.WriteAll(data)
	w.Flush()
}

func intToStr(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func floatToStr(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64)
}
