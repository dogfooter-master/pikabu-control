package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strings"

	// "bytes"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 8192
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 클라이언트 정보
type Client struct {
	ClientToken string
	LiveId      string
	ClientType  string
	Conn        *websocket.Conn
	User        UserObject
	Send        chan []byte
}

type Hub struct {
	Clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type WebSocketMessage struct {
	Data Payload `json:"data,omitempty"`
	Err  string  `json:"err,omitempty"`
}

func (w *WebSocketMessage) Encode() (b []byte, err error) {
	b, err = json.MarshalIndent(w, "", " ")
	if err != nil {
		return
	}

	return
}

var WebSocketHub *Hub

func init() {
	WebSocketHub = newHub()
	go WebSocketHub.run()
}

func (c *Client) readPump() {
	defer func() {
		WebSocketHub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message WebSocketMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		fmt.Fprintf(os.Stderr, "RECEIVE: %v\n", message.Data.Service)
		switch message.Data.Service {
		case "RegisterMate":
			c.ClientType = message.Data.ClientType
		case "StartToLive":
			WebSocketHub.GenerateLiveId(c.ClientToken)
		case "UpdateAccessToken":
			c.User.SecretToken.Token = message.Data.AccessToken
			continue
			//c.LiveId = strings.Replace(uuid.Must(uuid.NewUUID()).String(), "-", "", -1)
		}

		fmt.Fprintf(os.Stderr, "Token: %v\n", message.Data.AccessToken)
		accessToken := SecretTokenObject{
			Token: message.Data.AccessToken,
		}
		var user UserObject
		if user, err = accessToken.Authenticate(); err != nil {
			fmt.Fprintf(os.Stderr, "unauthorized: %v\n", err)
			res := WebSocketMessage{
				Data: Payload{
					Service:     message.Data.Service,
					ClientToken: c.ClientToken,
				},
			}
			message, err := res.Encode()
			if err != nil {
				message = []byte(`{ "err" : "` + err.Error() + `" }`)
			}
			c.Send <- message
			continue
		}
		c.User = user
		topic := TopicObject{
			Service:             message.Data.Service,
			UserId:              c.User.Id,
			AccessToken:         c.User.SecretToken.Token,
			Account:             c.User.Login.Account,
			RemoteAddr:          c.Conn.RemoteAddr().String(),
			ClientToken:         c.ClientToken,
			LiveId:              c.LiveId,
			MateToken:           message.Data.MateToken,
			Sdp:                 message.Data.Sdp,
			Candidate:           message.Data.Candidate,
			OpponentClientToken: message.Data.OpponentClientToken,
		}
		if err = topic.Publish(); err != nil {
			fmt.Fprintf(os.Stderr, "fail to topic.Publish: %v\n", err)
			continue
		}
		//message.Data.Debug(">> " + GetFunctionName() + " WebSocket Incoming: ")
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			//fmt.Fprintf(os.Stderr, "Write:\n %#v\n", message)
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
}

func (h *Hub) GetCount(userId string) (count int) {
	count = 0
	for _, c := range h.Clients {
		if c.User.Id.Hex() == userId {
			count += 1
		}
	}

	return count
}
func (h *Hub) CompareRemoteAddr(src, tgt string) bool {
	srcTokens := strings.Split(src, ":")
	tgtTokens := strings.Split(tgt, ":")
	srcAddr := ""
	for i, e := range srcTokens {
		if i == len(srcTokens)-1 {
			break
		}
		srcAddr += e
	}
	tgtAddr := ""
	for i, e := range tgtTokens {
		if i == len(tgtTokens)-1 {
			break
		}
		tgtAddr += e
	}
	if strings.Compare(srcAddr, tgtAddr) == 0 {
		return true
	} else {
		return false
	}
}
func (h *Hub) CloseExceptMyAddr(clientToken, accessToken string) {
	if _, ok := h.Clients[clientToken]; ok {
		c := h.Clients[clientToken]
		for k, v := range h.Clients {
			fmt.Fprintf(os.Stderr, "DEBUG: %v %v %v %v\n", k, v.User.Login.Account, v.Conn.RemoteAddr(), c.User.SecretToken.Token)
			if clientToken == k {
				continue
			}
			if h.CompareRemoteAddr(v.Conn.RemoteAddr().String(), c.Conn.RemoteAddr().String()) {
				// IP가 같다면 패스
				continue
			}
			if v.User.SecretToken.Token == accessToken {
				WebSocketHub.unregister <- v
			}
		}
	}
}
func (h *Hub) SendToClient(clientToken string, res WebSocketMessage) {
	if _, ok := h.Clients[clientToken]; ok {
		client := h.Clients[clientToken]
		message, err := res.Encode()
		if err != nil {
			message = []byte(`{ "err" : "` + err.Error() + `" }`)
		}
		client.Send <- message
	}
}
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			fmt.Fprintf(os.Stderr, "DEBUG: WebSocket Opened(%v)\n", client.Conn.RemoteAddr())
			h.Clients[client.ClientToken] = client
			//go func() {
			//	for {
			//		if _, ok := h.Clients[client.ClientToken]; ok {
			//			rs := WebSocketMessage{
			//				Data: Payload{
			//					Category:  "ws",
			//					Service:   "ReadyToLive",
			//					Account:   h.Clients[client.ClientToken].User.Login.Account,
			//					MateToken: h.Clients[client.ClientToken].ClientToken,
			//				},
			//			}
			//			WebSocketHub.BroadcastToPikabu(h.Clients[client.ClientToken].User.Id, rs)
			//			time.Sleep(10 * time.Second)
			//		} else {
			//			break
			//		}
			//	}
			//}()
		case client := <-h.unregister:
			if _, ok := h.Clients[client.ClientToken]; ok {
				//userid := client.User.UserId
				//remember := client.remember
				//sessionUuid := client.sessionUuid
				fmt.Fprintf(os.Stderr, "DEBUG: WebSocket Closed(%v:%v)\n", client.Conn.RemoteAddr(), client.ClientType)
				c := client
				delete(h.Clients, client.ClientToken)
				close(client.Send)

				if c.ClientType == "mate" {
					// 상대방에게 접속 종료되었음을 전송한다.
					topic := TopicObject{
						Service:             "UnableToLive",
						ClientToken:         c.ClientToken,
						LiveId:              c.LiveId,
					}
					if err := topic.Publish(); err != nil {
						fmt.Fprintf(os.Stderr, "fail to topic.Publish: %v\n", err)
					}
					// 접속 종료한 것이 mate 라면, 해당 유저의 mate가 남아있는 지 체크하고 없으면 Live 못한다고 보낸다.
					//if WebSocketHub.HasMate(c.User.Id) == false {
					//	rs := WebSocketMessage{
					//		Data: Payload{
					//			Category: "ws",
					//			Service:  "UnableToLive",
					//		},
					//	}
					//	WebSocketHub.BroadcastToPikabu(c.User.Id, rs)
					//}
				}

				//logger.Debugf("GetCount: %v", h.GetCount(userid))
				//logger.Infof("remeber: %v", remember)
				//logger.Debugf("sessionUuid: %v", sessionUuid)

				// if h.GetCount(userid) == 0 && remember == false {
				// 	clearSession(sessionUuid)
				// }
			}
		case message := <-h.broadcast:
			for token, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					//logger.Info("Closed")
					close(client.Send)
					delete(h.Clients, token)
				}
			}
		}
	}
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}
func (h *Hub) BroadcastToPikabu(userId bson.ObjectId, res WebSocketMessage) {
	for _, v := range h.Clients {
		fmt.Fprintf(os.Stderr, "BroadcastToPikabu - DEBUG: UserId(%v) - UserId(%v)\n", v.User.Id.Hex(), userId.Hex())
		// 나 자신은 제외
		if v.ClientType == "mate" {
			continue
		}
		if v.User.Id == userId {
			message, err := res.Encode()
			if err != nil {
				message = []byte(`{ "err" : "` + err.Error() + `" }`)
			}
			v.Send <- message
		}
	}
}
func (h *Hub) HasMate(userId bson.ObjectId) bool {
	for _, v := range h.Clients {
		if v.ClientType == "mate" && v.User.Id == userId {
			return true
		}
	}
	return false
}
func (h *Hub) SendToOpponent(clientToken string, msg WebSocketMessage) {
	if _, ok := h.Clients[clientToken]; ok {
		c := h.Clients[clientToken]
		fmt.Fprintf(os.Stderr, "DEBUG: %v to %v(Live:%v)\n", msg.Data.Service, c.User.Login.Account, c.LiveId)
		message, err := msg.Encode()
		if err != nil {
			message = []byte(`{ "err" : "` + err.Error() + `" }`)
		}
		c.Send <- message
	}
}
func (h *Hub) SendToOpponentByLiveId(liveId string, msg WebSocketMessage) {
	for _, c := range h.Clients {
		if c.LiveId == liveId {
			fmt.Fprintf(os.Stderr, "DEBUG: %v to %v(Live:%v)\n", msg.Data.Service, c.LiveId)
			message, err := msg.Encode()
			if err != nil {
				message = []byte(`{ "err" : "` + err.Error() + `" }`)
			}
			c.Send <- message
		}
	}
}
func (h *Hub) GenerateLiveId(clientToken string) {
	if c, ok := h.Clients[clientToken]; ok {
		c.LiveId = bson.NewObjectId().Hex()
	}
}
func (h *Hub) SetLiveId(clientToken, liveToken string) {
	if c, ok := h.Clients[clientToken]; ok {
		c.LiveId = liveToken
	}
}
func ServeWebSocket(w http.ResponseWriter, r *http.Request) {

	//_, err := session(w, r)
	//if err != nil {
	//	//logger.Error("websocket session check fail")
	//	return
	//}

	//cookie, _ := r.Cookie("userid")
	//var userid string
	//if cookie != nil {
	//	userid = cookie.Value
	//} else {
	//	userid = ""
	//}
	//
	//cookie, _ = r.Cookie("sessionUuid")
	//var session_uuid string
	//if cookie != nil {
	//	session_uuid = cookie.Value
	//} else {
	//	session_uuid = ""
	//}
	//
	//cookie, _ = r.Cookie("remember")
	//remember := false
	//if cookie != nil {
	//	if cookie.Value == "true" {
	//		remember = true
	//	} else {
	//		remember = false
	//	}
	//}

	// 인증 성공했으므로 반드시 있다.
	//User, err := data.UserByUserId(userid)
	//if err != nil {
	//	//logger.Errorf("Not found: %v", userid)
	//	w.WriteHeader(http.StatusNotAcceptable)
	//	return
	//}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	//var respBody Payload
	//if err := json.NewDecoder(r.Body).Decode(&respBody); err != nil {
	//	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	//	return
	//}
	//respBody.Debug("< " + GetFunctionName() + " WebSocket Connection Incoming")

	client := &Client{
		ClientToken: strings.Replace(uuid.Must(uuid.NewUUID()).String(), "-", "", -1),
		Conn:        conn,
		User:        UserObject{},
		Send:        make(chan []byte, 4096),
	}
	WebSocketHub.register <- client

	go client.writePump()
	go client.readPump()
}
