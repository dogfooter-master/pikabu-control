package service

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type ClientObject struct {
	ClientAccessCode map[string]SecretTokenObject `bson:"client_access_code,omitempty"`
}

func (c *ClientObject) GenerateAccessCode(userId bson.ObjectId, clientToken string) (accessCode string, err error) {
	var newOne = make(map[string]SecretTokenObject)
	for k, v := range c.ClientAccessCode {
		if v.IsExpiredPc() == false {
			newOne[k] = v
		}
	}
	if code, ok := c.ClientAccessCode[clientToken]; ok {
		accessCode = code.Token
	} else {
		code.GenerateCertificationNumber()
		accessCode = code.Token
		newOne[clientToken] = code
	}
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	updateBson := bson.M{}
	updateBson["client.client_access_code"] = newOne
	err = collection.Update(bson.M{"_id": userId}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}
	return
}
func (c *ClientObject) ConfirmAccessCode(userId bson.ObjectId, accessCode string) (clientToken string, err error) {
	var newOne = make(map[string]SecretTokenObject)
	for k, v := range c.ClientAccessCode {
		if v.Token == accessCode {
			if v.IsExpiredPc() {
				err = errors.New("expired")
				return
			}
			clientToken = k
		} else {
			if v.IsExpiredPc() == false {
				newOne[k] = v
			}
		}
	}
	if len(clientToken) == 0 {
		err = errors.New("not found")
		return
	}
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	updateBson := bson.M{}
	updateBson["client.client_access_code"] = newOne
	err = collection.Update(bson.M{"_id": userId}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}

	return
}
