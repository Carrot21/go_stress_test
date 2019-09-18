package util

import (
	"encoding/json"
	"encoding/xml"
	"github.com/cihub/seelog"
	"io/ioutil"
	"path"
)

type Result struct {
	Person []Person `xml:"row" json:"RECORDS"`
}

type Person struct {
	//UserId string `xml:"id"`//体现了数据源一致性的好处，代码根本不用修改
	UserId int64 `xml:"user_id" json:"user_id"`
	LoginToken string `xml:"last_login_token" json:"last_login_token"`
	DeviceToken string `xml:"device_token" json:"device_token"`
	VoipToken string `xml:"voip_token" json:"voip_token"`
	Push_token string `xml:"push_token" json:"push_token"`
	ChannelType uint32 `xml:"channel_type" json:"channel_type"`
	VersionCode string
	SessionId	[12]byte
}

func GetUserInfoFromFile(result *Result, filename string) error {
	conntent , err := ioutil.ReadFile(filename)
	if err != nil {
		seelog.Error(err)
	}
	extent := path.Ext(filename)
	seelog.Info("Extent file name : ", extent)
	if extent == ".xml" {
		err = xml.Unmarshal(conntent, result)
	} else if extent == ".json" {
		err = json.Unmarshal(conntent, result)
	} else {
		seelog.Info("Unknow file type")
	}

	if err != nil {
		seelog.Error(err)
	}

	return err
}