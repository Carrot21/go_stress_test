package entity

type MsgHeader struct {
	Flags uint16		//协议标识:0x535a5951 "SZYQ"  判断到消息类型
	Length uint16 		//消息体长度
	AppId uint8			//app_id  用于用户的app标识
	Version uint8		//版本号  消息版本(协议表示的子版本号)
	ProType uint8		//消息体类型  0:pb  1: json
	MsgType uint8		//0：单聊 1：群聊  2：超级大群 ...
	DstUserId uint32	//目标id
	SrcUserId uint32	//源id
	MsgSeq uint32		//消息唯一标识
	Command uint16		//命令字  0xFFFF
	ReasonCode uint16	//错误码 --- 24B
}

type Header struct {
	NPID uint16 // pdu identification , default as 'US'
	NVersion uint8 // pdu version
	SessionId [12]byte //UidCode_t fixed length 12 byte size;
	BEncrypt uint8// one byte for identifing whether or not encrypt content.
	NCmdId uint16// command id of protobuf protocol.
	NBodySize uint16 // size
}

type UserInfo struct {
	UserId string
	LoginToken string
	DeviceToken string
	VoipToken string
	PushToken string
	ChannelType uint32
	VersionCode string
	SessionId	[12]byte
}
//
//type Person struct {
//	//UserId string `xml:"user_id"`
//	UserId string `xml:"id"`
//	LoginToken string `xml:"last_login_token"`
//	DeviceToken string `xml:"device_token"`
//	VoipToken string `xml:"voip_token"`
//	Push_token string `xml:"push_token"`
//	ChannelType uint32 `xml:"channel_type"`
//	VersionCode string
//	SessionId	[12]byte
//}