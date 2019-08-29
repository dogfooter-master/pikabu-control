package service

import (
	"errors"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type SecretTokenObject struct {
	Token      string        `bson:"token,omitempty"`
	Count      int32         `bson:"count,omitempty"`
	CreateTime time.Time     `bson:"create_time,omitempty"`
}

func (s *SecretTokenObject) GenerateCertificationNumber() {
	s.Token = strconv.Itoa(random(100000, 999999))
	s.CreateTime = time.Now().UTC()
}
func (s *SecretTokenObject) Refresh() {
	s.CreateTime = time.Now().UTC()
}
func (s *SecretTokenObject) IsExpiredPc() bool {
	if time.Since(s.CreateTime).Hours() > 4 {
		return true
	}
	return false
}
func (s *SecretTokenObject) IsExpired() bool {
	if time.Since(s.CreateTime).Hours() > 24 {
		return true
	}
	return false
}
func (s *SecretTokenObject) IsExpiredLong() bool {
	//TODO: 유효기간을 얼마로 할 지
	if time.Since(s.CreateTime).Hours() > 24*365 {
		return true
	}
	return false
}
func (s *SecretTokenObject) GenerateAccessToken() (err error) {
	s.Token = strings.Replace(uuid.Must(uuid.NewUUID()).String(), "-", "", -1)
	s.CreateTime = time.Now().UTC()
	return
}
func (s *SecretTokenObject) Authenticate() (user UserObject, err error) {
	//defer TimeTrack(time.Now(), GetFunctionName())
	if len(s.Token) == 0 {
		err = errors.New("'access_token' is mandatory")
		return
	}
	if user, err = s.Read(); err != nil {
		err = errors.New("unauthorized")
		return
	} else {
		if user.SecretToken.IsExpiredLong() {
			err = errors.New("expired")
			return
		}
	}
	return
}
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
func (s *SecretTokenObject) Read() (obj UserObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	readBson := bson.M{}
	if len(s.Token) > 0 {
		readBson["secret_token.token"] = s.Token
	}
	err = collection.Find(readBson).One(&obj)
	if err != nil {
		return
	}

	return
}

