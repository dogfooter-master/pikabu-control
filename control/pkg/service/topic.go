package service

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/mgo.v2/bson"
	"os"
)

type TopicObject struct {
	Id                  bson.ObjectId          `json:"_id,omitempty"`
	UserId              bson.ObjectId          `json:"user,omitempty"`
	ClientToken         string                 `json:"client_token,omitempty"`
	ClientType          string                 `json:"client_type,omitempty"`
	ChannelType         string                 `json:"channel_type,omitempty"`
	LiveId              string                 `json:"live_id,omitempty"`
	MateToken           string                 `json:"mate_token,omitempty"`
	OpponentClientToken string                 `json:"opponent_client_token,omitempty"`
	AccessToken         string                 `json:"access_token,omitempty"`
	Service             string                 `json:"service,omitempty"`
	Account             string                 `json:"account,omitempty"`
	RemoteAddr          string                 `json:"remote_addr,omitempty"`
	Sdp                 map[string]interface{} `json:"sdp,omitempty"`
	Candidate           map[string]interface{} `json:"candidate,omitempty"`
}

var TopicRequest chan TopicObject
var RedisPool *redis.Pool

func init() {
	switch GetPlatform() {
	case "linux":
		//err := initializeSubscribe()
		//if err != nil {
		//	panic(err)
		//}
		RedisPool = newPool()
		fallthrough
	case "windows":
		topic := TopicObject{}
		topic.Subscribe()
	}
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConfig.Hosts)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

//func initializeSubscribe() (err error) {
//	c := mgoSession.DB(mgoConfig.Database).C(mgoConfig.TopicCollection)
//
//	err = c.Create(&mgo.CollectionInfo{
//		Capped:   true,
//		MaxBytes: 4096,
//	})
//	if err == nil {
//		topic := TopicObject{
//			Service: "Init",
//		}
//		if err = c.Insert(&topic); err != nil {
//			return
//		}
//	} else {
//		err = nil
//	}
//
//	return
//}
func (t *TopicObject) Subscribe() {
	switch GetPlatform() {
	case "linux":
		go func() {
			c := RedisPool.Get()
			defer c.Close()
			var channels []string
			channels = append(channels, redisConfig.WebSocketChannel)
			psc := redis.PubSubConn{Conn: c}
			_ = psc.Subscribe(redis.Args{}.AddFlat(channels)...)
			for {
				switch v := psc.Receive().(type) {
				case redis.Message:
					fmt.Fprintf(os.Stderr, "%s: message: %s\n", v.Channel, v.Data)
					result := TopicObject{}
					result.Decode(v.Data)
					result.Process()
				case redis.Subscription:
					fmt.Fprintf(os.Stderr, "%s: %s %d\n", v.Channel, v.Kind, v.Count)
				case error:
					fmt.Fprintf(os.Stderr, "err - %v\n", v)
				}
			}
		}()

		//c := mgoSession.DB(mgoConfig.Database).C(mgoConfig.TopicCollection)
		//
		//timeout := -1 * time.Second
		//iter := c.Find(nil).Sort("$natural").Limit(10).Tail(timeout)
		//go func() {
		//	for {
		//		result := TopicObject{}
		//		for iter.Next(&result) {
		//			result.Process()
		//		}
		//		if iter.Err() != nil {
		//			iter.Close()
		//			return
		//		}
		//		if iter.Timeout() {
		//			continue
		//		}
		//	}
		//	iter.Close()
		//}()
	case "windows":
		TopicRequest = make(chan TopicObject, 10)
		go func() {
			for {
				result := <-TopicRequest
				result.Process()
			}
		}()
	}
}
func (t *TopicObject) Encode() (rs string) {
	b, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return ""
	}
	rs = string(b)

	return
}
func (t *TopicObject) Decode(b []byte) {
	_ = json.Unmarshal(b, t)
}
func (t *TopicObject) Publish() (err error) {
	t.Id = bson.NewObjectId()
	switch GetPlatform() {
	case "linux":
		c := RedisPool.Get()
		defer c.Close()
		_, _ = c.Do("PUBLISH", redisConfig.WebSocketChannel, t.Encode())
		//c := mgoSession.DB(mgoConfig.Database).C(mgoConfig.TopicCollection)
		//if err = c.Insert(t); err != nil {
		//	return
		//}
	case "windows":
		TopicRequest <- *t
	}
	return
}
func (t *TopicObject) Process() {
	switch t.Service {
	case "Register":
		// 클라이언트가 접속했음
		// WebSocketHub.CloseExceptMyAddr(t.ClientToken, t.AccessToken)
		rs := WebSocketMessage{
			Data: Payload{
				Category:    "ws",
				Service:     "RegisterComplete",
				Account:     t.Account,
				ClientToken: t.ClientToken,
			},
		}
		WebSocketHub.SendToClient(t.ClientToken, rs)
		rs = WebSocketMessage{
			Data: Payload{
				Category:            "ws",
				Service:             "ReadyToLive",
				Account:             t.Account,
				OpponentClientToken: t.ClientToken,
			},
		}
		WebSocketHub.BroadcastToPikabu(t.UserId, rs)
	case "Unregister":
	case "SignInMate":
		rs := WebSocketMessage{
			Data: Payload{
				Category:    "ws",
				Service:     "SignInMate",
				AccessToken: t.AccessToken,
			},
		}
		fmt.Fprintf(os.Stderr, "DEBUG: SignInMate %v\n", t.ClientToken)
		WebSocketHub.SendToClient(t.ClientToken, rs)
	//case "RegisterAgent":
	//	// 더메이트가 접속했다.
	//	// 동일한 계정으로 접속된 모든 웹에 알림
	//	rs := WebSocketMessage{
	//		Data: Payload{
	//			Category:            "ws",
	//			Service:             "ReadyToLive",
	//			Account:             t.Account,
	//			OpponentClientToken: t.ClientToken,
	//		},
	//	}
	//	WebSocketHub.BroadcastToPikabu(t.UserId, rs)
	case "GetMate":
		// 현재 접속한 더메이트 정보를 요청한다.
	case "StartToLive":
		// 웹에서 모바일로 Offer 해달라고 콜함.
		fmt.Fprintf(os.Stderr, "DEBUG: StartToLive %v %v -> %v\n", t.Account, t.ClientToken, t.OpponentClientToken)
		rs := WebSocketMessage{
			Data: Payload{
				Category:            "ws",
				Service:             "RequestOffer",
				OpponentClientToken: t.ClientToken,
				LiveId:              t.LiveId,
			},
		}
		WebSocketHub.SetLiveId(t.OpponentClientToken, t.LiveId)
		WebSocketHub.SendToOpponent(t.OpponentClientToken, rs)
	case "Offer":
		fmt.Fprintf(os.Stderr, "DEBUG: Offer %v %v -> %v\n", t.Account, t.ClientToken, t.OpponentClientToken)
		rs := WebSocketMessage{
			Data: Payload{
				Category:    "ws",
				Service:     "Offer",
				Sdp:         t.Sdp,
				ChannelType: t.ChannelType,
				ClientToken: t.ClientToken,
				LiveId:      t.LiveId,
			},
		}
		WebSocketHub.SendToOpponent(t.OpponentClientToken, rs)
	case "Answer":
		fmt.Fprintf(os.Stderr, "DEBUG: Answer %v\n", t.Account)
		rs := WebSocketMessage{
			Data: Payload{
				Category:            "ws",
				Service:             "Answer",
				Candidate:           t.Candidate,
				Sdp:                 t.Sdp,
				ClientToken:         t.ClientToken,
				ChannelType:         t.ChannelType,
				OpponentClientToken: t.OpponentClientToken,
			},
		}
		WebSocketHub.SendToOpponent(t.OpponentClientToken, rs)
	case "Candidate":
		fmt.Fprintf(os.Stderr, "DEBUG: Candidate %v\n", t.Account)
		rs := WebSocketMessage{
			Data: Payload{
				Category:            "ws",
				Service:             "Candidate",
				Candidate:           t.Candidate,
				ClientToken:         t.ClientToken,
				ChannelType:         t.ChannelType,
				OpponentClientToken: t.OpponentClientToken,
			},
		}
		WebSocketHub.SendToOpponent(t.OpponentClientToken, rs)
	case "OnReceiveLiveImage":
		fmt.Fprintf(os.Stderr, "DEBUG: OnReceiveLiveImage %v, %v, %v\n", t.Account, t.LiveId, t.OpponentClientToken)
		rs := WebSocketMessage{
			Data: Payload{
				Category: "ws",
				Service:  "OnReceiveLiveImage",
				LiveId:   t.LiveId,
			},
		}
		WebSocketHub.SendToOpponent(t.OpponentClientToken, rs)
	case "UnableToLive":
		fmt.Fprintf(os.Stderr, "DEBUG: UnableToLive %v, %v\n", t.Account, t.LiveId)
		rs := WebSocketMessage{
			Data: Payload{
				Category: "ws",
				Service:  "UnableToLive",
				LiveId:   t.LiveId,
			},
		}
		WebSocketHub.SendToOpponentByLiveId(t.LiveId, rs)
	}
}
