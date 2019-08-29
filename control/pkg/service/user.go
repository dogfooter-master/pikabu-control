package service

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserObject struct {
	Id           bson.ObjectId      `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Nickname     string             `bson:"nickname,omitempty"`
	Login        LoginObject        `bson:"login,omitempty"`
	Relation     RelationObject     `bson:"relation,omitempty"`
	SecretToken  SecretTokenObject  `bson:"secret_token,omitempty"`
	Client       ClientObject       `bson:"client,omitempty"`
	Avatar       FileObject         `bson:"avatar,omitempty"`
	Status       string             `bson:"status,omitempty"`
	Time         TimeLogObject      `bson:"time,omitempty"`
	CustomConfig CustomConfigObject `bson:"custom_config,omitempty"`
}

func (d *UserObject) Create() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	index := mgo.Index{
		Key:    []string{"login.account"},
		Unique: true,
	}
	if err = collection.EnsureIndex(index); err != nil {
		err = fmt.Errorf("collection.EnsureIndex: %v", err)
		return
	}

	d.Time.Initialize()

	err = collection.Insert(d)
	if err != nil {
		err = fmt.Errorf("insert: %v", err)
		return
	}

	return
}

func (d *UserObject) Read() (obj UserObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	readBson := bson.M{}
	if len(d.Id) > 0 {
		readBson["_id"] = d.Id
	}
	if len(d.Login.Account) > 0 {
		readBson["login.account"] = d.Login.Account
	}

	err = collection.Find(readBson).One(&obj)
	if err != nil {
		return
	}
	return
}
func (d *UserObject) ReadAll(skip int, limit int) (objList []UserObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	readBson := bson.M{}
	if len(d.Login.Account) > 0 {
		readBson["login.account"] = d.Login.Account
	}

	err = collection.Find(readBson).Skip(skip).Limit(limit).All(&objList)
	if err != nil {
		return
	}
	return
}
func (d *UserObject) Update() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	//object := UserObject{
	//	Id: d.Id,
	//}
	//err = collection.Find(bson.M{"_id": object.Id}).One(&object)
	//if err != nil {
	//	return
	//}

	updateBson := bson.M{}

	if len(d.Login.Password) > 0 {
		updateBson["login.password"] = d.Login.Password
	}
	if len(d.Name) > 0 {
		updateBson["name"] = d.Name
	}
	if len(d.Nickname) > 0 {
		updateBson["nickname"] = d.Nickname
	}
	if len(d.Avatar.Path) > 0 {
		updateBson["avatar"] = d.Avatar
	}
	if len(d.SecretToken.Token) > 0 {
		updateBson["secret_token"] = d.SecretToken
	}
	if len(d.Status) > 0 {
		updateBson["status"] = d.Status
	}
	if d.Time.LoginTime.IsZero() == false {
		updateBson["time.login_time"] = d.Time.LoginTime
	}
	if len(d.Relation.HospitalId) > 0 {
		updateBson["relation.hospital_id"] = d.Relation.HospitalId
	}
	if len(d.Relation.PcUserId) > 0 {
		updateBson["relation.pc_user_id"] = d.Relation.PcUserId
	}

	d.Time.Update()
	updateBson["time.update_time"] = d.Time.UpdateTime

	err = collection.Update(bson.M{"_id": d.Id}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}

	return
}
func (d *UserObject) UpdatePcUser() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	updateBson := bson.M{}
	updateBson["relation.pc_user_id"] = d.Relation.PcUserId

	d.Time.Update()
	updateBson["time.update_time"] = d.Time.UpdateTime

	err = collection.Update(bson.M{"_id": d.Id}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}

	return
}
func (d *UserObject) Delete() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	err = collection.Remove(bson.M{"_id": d.Id})
	if err != nil {
		return
	}

	return
}
func (d *UserObject) Validate() (err error) {
	if len(d.Name) < 1 {
		err = errors.New("'name' is mandatory")
		return
	}
	if len(d.Nickname) == 0 {
		d.Nickname = d.Name
	}

	return
}