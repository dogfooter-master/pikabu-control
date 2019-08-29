package service

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
)

var mgoSession *mgo.Session

func init() {
	err := initializeMongo()
	if err != nil {
		panic(err)
	}
}

func initializeMongo() (err error) {
	info := &mgo.DialInfo{
		Addrs:    []string{mgoConfig.Hosts},
		Timeout:  60 * time.Second,
		Database: mgoConfig.Database,
		Username: mgoConfig.Username,
		Password: mgoConfig.Password,
	}

	mgoSession, err = mgo.DialWithInfo(info)
	if err != nil {
		err = fmt.Errorf("fail to DialWithInfo(%#v) error - %v", info, err)
		return
	}

	return
}
