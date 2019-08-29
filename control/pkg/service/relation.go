package service

import "gopkg.in/mgo.v2/bson"

type RelationObject struct {
	HospitalId bson.ObjectId `bson:"hospital_id,omitempty"`
	PcUserId   bson.ObjectId `bson:"pc_user_id,omitempty"`
}
