package service

import (
	"database/sql"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

var mgoSession *mgo.Session
var mySqlDB *sql.DB
func init() {
	err := initializeMongo()
	if err != nil {
		panic(err)
	}

	// sql.DB 객체 생성
	dataSourceName := mySqlConfig.Username +
		":" + mySqlConfig.Password +
		"@tcp("+mySqlConfig.Hosts+")/"+
		mySqlConfig.Database+")"
	mySqlDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	var name string
	err = mySqlDB.QueryRow("SELECT mb_name FROM g5_member WHERE mb_id = 'admin'").Scan(&name)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "mb_name=%v\n", name)
}

func initializeMongo() (err error) {
	//info := &mgo.DialInfo{
	//	Addrs:    []string{mgoConfig.Hosts},
	//	Timeout:  60 * time.Second,
	//	Database: mgoConfig.Database,
	//	Username: mgoConfig.Username,
	//	Password: mgoConfig.Password,
	//
	//}

	url := "mongodb://" + mgoConfig.Username + ":" + mgoConfig.Password + "@" + mgoConfig.Hosts + "/" + mgoConfig.Database + "?authSource=admin"
	mgoSession, err = mgo.Dial(url)
	if err != nil {
		err = fmt.Errorf("fail to Dial(%#v) error - %v", url, err)
		return
	}

	return
}
