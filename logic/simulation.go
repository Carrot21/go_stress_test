package logic

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"go_stress_test/config"
	"go_stress_test/entity"
	msgcmdproto "go_stress_test/proto"
	"go_stress_test/util"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func ConnTCPserver() net.Conn {
	conn, err := net.DialTimeout("tcp", config.GetConfig().HostPort,
		time.Duration(config.GetConfig().DialTimeout)*time.Second)
	if err != nil {
		log.Fatalln("Fatal error:", err.Error())
		os.Exit(1)
	}

	log.Println(conn.RemoteAddr().String(), "connection succcess!")

	return conn
}

func SimulateLogin(csvSlice [][]string) {
	var wg sync.WaitGroup

	for i := 0; i < len(csvSlice); i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			conn := ConnTCPserver()

			msgHead := &entity.Header{
				NPID:      binary.LittleEndian.Uint16([]byte{'U', 'S'}),
				NVersion:  2,
				SessionId: [12]byte{},
				BEncrypt:  0,
				NCmdId:    0xa001,
				NBodySize: 0,
			}

			msgBody := msgcmdproto.CMLoginV1{
				SUserId:      csvSlice[i][0],  // string
				SLoginToken:  csvSlice[i][1],  // string
				SDeviceToken: "tokeninfotest", // string
				NPushType:    254,             // uint32   `protobuf:"varint,
				SPushToken:   "",              //result.Person[i].Push_token,// string
				SVersionCode: "2.3.0",         // string
			}

			log.Println(msgBody)
			// 对数据进行序列化
			msgBodyProto, err := proto.Marshal(&msgBody)
			if err != nil {
				log.Fatalln("Mashal data error:", err)
			}

			msgHead.NBodySize = uint16(len(msgBodyProto))

			sendData := util.StructToByte(msgHead)
			sendData = append(sendData, msgBodyProto...)

			//conn := conns[i]
			conn.Write(sendData)

			//read ack data
			recvData := make([]byte, 1024)
			reqLen, err := conn.Read(recvData)
			if err != nil {
				log.Println("Error to read message", err.Error())
			}

			log.Printf("Recv data from %s, data len = %d", conn.RemoteAddr(), reqLen)

			loginAck := msgcmdproto.CMLoginV1Ack{}

			proto.Unmarshal(recvData[20:], &loginAck)

			if loginAck.NErr != msgcmdproto.ErrCode_NON_ERR {
				log.Printf("user %s login error , errorcode = %d\n", loginAck.GetSUserId(), loginAck.GetNErr())
			}

			//copy(result.Person[i].SessionId[:], recvData[3:15])
			log.Printf("user %s login at %d , status = %d\n",
				loginAck.GetSUserId(), loginAck.GetNLastLoginTime(), loginAck.GetNErr())
			//log.Println("sessionid: ", result.Person[i].SessionId)
			//conns = append(conns, conn)
		}(i)
	}

	wg.Wait()
}

//发心跳包的
func SimulateHeartBeat(csvSlice [][]string) {
	var g errgroup.Group

	for i := 0; i < len(csvSlice); i++ {
		g.Go(func() error {
			conn := ConnTCPserver()

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
			//conn := conns[i]
			wLen, err := conn.Write(sendData)
			if err != nil {
				log.Println("Write Data Error: ", error(err))
				return err
			}

			log.Printf("Write data to %s, len = %d\n", conn.RemoteAddr(), wLen)

			time.Sleep(time.Duration(config.GetConfig().HeartBeat) * time.Second)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatalln("Err:", err)
	}
}
