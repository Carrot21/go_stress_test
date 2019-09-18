package logic

import (
	"encoding/binary"
	"github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	"go_stress_test/config"
	"go_stress_test/entity"
	msgcmdproto "go_stress_test/proto"
	"go_stress_test/util"
	"net"
	"os"
	"sync"
	"time"
)

func ConnTCPserver() net.Conn {
	conn, err := net.DialTimeout("tcp", config.GetConfig().HostPort,
		time.Duration(config.GetConfig().DialTimeout)*time.Second)
	if err != nil {
		seelog.Error("Fatal error:", err.Error())
		os.Exit(1)
	}

	return conn
}

func SimulateLogin(csvSlice [][]string, ch chan<- *entity.ResponseResults) {
	var wg sync.WaitGroup

	for i := 0; i < len(csvSlice); i++ {
		wg.Add(1)

		go func(i int, ch chan<- *entity.ResponseResults) {
			conn := ConnTCPserver()

			defer wg.Done()
			defer conn.Close()

			msgHead := &entity.Header{
				NPID:      binary.LittleEndian.Uint16([]byte{'U', 'S'}),
				NVersion:  2,
				SessionId: [12]byte{},
				BEncrypt:  0,
				NCmdId:    0xa001,
				NBodySize: 0,
			}

			msgBody := msgcmdproto.CMLoginV1{
				SUserId:      csvSlice[i][0],
				SLoginToken:  csvSlice[i][1],
				SDeviceToken: "tokeninfotest",
				NPushType:    254,
				SPushToken:   "",
				SVersionCode: "2.3.0",
			}

			seelog.Info(msgBody)
			// 对数据进行序列化
			msgBodyProto, err := proto.Marshal(&msgBody)
			if err != nil {
				seelog.Error("Mashal data error:", err)
			}

			msgHead.NBodySize = uint16(len(msgBodyProto))

			sendData := util.StructToByte(msgHead)
			sendData = append(sendData, msgBodyProto...)

			var (
				startTime = time.Now()
				isSucceed = true
			)

			if _, err := conn.Write(sendData); err != nil {
				isSucceed = false

				seelog.Error("Write CM Server Failed:", err)
			}

			//read ack data
			recvData := make([]byte, 1024)
			reqLen, err := conn.Read(recvData)

			spentTime := uint64(DiffNano(startTime))

			if err != nil {
				isSucceed = false

				seelog.Info("Error to read message", err.Error())
			}

			seelog.Infof("Recv data from %s, data len = %d", conn.RemoteAddr(), reqLen)

			loginAck := msgcmdproto.CMLoginV1Ack{}

			proto.Unmarshal(recvData[20:], &loginAck)

			if loginAck.NErr != msgcmdproto.ErrCode_NON_ERR {
				seelog.Infof("user %s login error , errorcode = %d\n", loginAck.GetSUserId(), loginAck.GetNErr())
			}

			seelog.Infof("user %s login at %d , status = %d\n",
				loginAck.GetSUserId(), loginAck.GetNLastLoginTime(), loginAck.GetNErr())

			responseResults := &entity.ResponseResults{
				Time:      spentTime,
				IsSucceed: isSucceed,
			}

			ch <- responseResults
		}(i, ch)
	}

	wg.Wait()
}

// 时间差，纳秒
func DiffNano(startTime time.Time) (diff int64) {

	startTimeStamp := startTime.UnixNano()
	endTimeStamp := time.Now().UnixNano()

	diff = endTimeStamp - startTimeStamp

	return
}

//发心跳包的
func SimulateHeartBeat(csvSlice [][]string, onLineTime int) {
	var wg sync.WaitGroup

	for i := 0; i < len(csvSlice); i++ {
		wg.Add(1)
		go func() {
			conn := ConnTCPserver()

			defer wg.Done()
			defer conn.Close()

			ticker := time.NewTicker(time.Duration(onLineTime) * time.Minute)

			for {
				select {
				case <-ticker.C:
					ticker.Stop()
					return
				default:
					sendHeartBeat(conn)
				}
			}
		}()
	}

	wg.Wait()
}

//var n uint32

func sendHeartBeat(conn net.Conn) {
	data := []byte{'U', 'S'}

	msgHead := &entity.Header{
		NPID:      binary.LittleEndian.Uint16(data),
		NVersion:  2,
		SessionId: [12]byte{},
		BEncrypt:  0,
		NCmdId:    0xa001,
		NBodySize: 0,
	}

	// 对数据进行序列化
	sendData := util.StructToByte(msgHead)

	wLen, err := conn.Write(sendData)
	if err != nil {
		seelog.Info("Write Data Error: ", error(err))
	}

	seelog.Infof("Write data to %s, len = %d\n", conn.RemoteAddr(), wLen)

	time.Sleep(time.Duration(config.GetConfig().HeartBeat) * time.Second)

	//atomic.AddUint32(&n, 1)
	//println(n)
}
