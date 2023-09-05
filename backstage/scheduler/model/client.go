package model

import (
	"encoding/json"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

/*
@enum
@description: client type
*/
const (
	CLIENT_TYPE_UNKNOWN    string = "unknown"
	CLIENT_TYPE_PROVIDER   string = "provider"
	CLIENT_TYPE_CONSUMER   string = "consumer"
	CLIENT_TYPE_DEPOSITARY string = "depository"
	CLIENT_TYPE_FILESTORE  string = "filestore"
)

/*
@enum
@description: consumer type
*/
const (
	CONSUMER_TYPE_STREAM   string = "stream"
	CONSUMER_TYPE_TERMINAL string = "terminal"
)

/*
@enum
@description: provider type
*/
const (
	PROVIDER_TYPE_STREAM   string = "stream"
	PROVIDER_TYPE_TERMINAL string = "terminal"
)

/*
@model: Consumer
@description: consumer client
*/
type Consumer struct {
	ConsumerCore
	Client
}

/*
@model: ConsumerCore
@description: metadata for consumer client
*/
type ConsumerCore struct {
	ConsumerType string
}

/*
@model: Client
@description: client core for websocket communication
*/
type Client struct {
	// client meta data
	ClientID string

	// websocket
	WebsocketConnection *websocket.Conn
	SendLock            sync.Mutex

	SendCallbackList     map[string]func(req WSPacket)
	SendCallbackListLock sync.Mutex

	RecvCallbackList map[string]func(req WSPacket)

	Done chan struct{}
}

/*
@func: Send
@description:

	[1] send websocket packet
	[2] register send callback (optional)
*/
func (c *Client) Send(request WSPacket, callback func(response WSPacket)) {
	// generate packet id
	request.PacketID = uuid.Must(uuid.NewV4()).String()

	// transfer WSPacket to json string
	requestPacketString, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"ClientID":    c.ClientID,
			"Packet Data": request.Data,
			"error":       err,
		}).Warn("Failed to marshal websocket packet into json string during requesting")
		return
	}

	// register send callback
	if callback != nil {
		// wrap callback with packetID
		wrapperCallback := func(resp WSPacket) {
			defer func() {
				if err := recover(); err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Warn("Recovered from error in client callback")
				}
			}()

			resp.PacketID = request.PacketID
			callback(resp)
		}
		c.SendCallbackListLock.Lock()
		c.SendCallbackList[request.PacketID] = wrapperCallback
		c.SendCallbackListLock.Unlock()
	}

	// send request data
	c.SendLock.Lock()
	c.WebsocketConnection.SetWriteDeadline(time.Now().Add(20 * time.Second))
	c.WebsocketConnection.WriteMessage(websocket.TextMessage, requestPacketString)
	c.SendLock.Unlock()
}

/*
@func: Receive
@description:

	register receive callback based on packet type
*/
func (c *Client) Receive(packetType string, callback func(request WSPacket) (response WSPacket)) {
	c.RecvCallbackList[packetType] = func(request WSPacket) {
		// invoke receive callback
		resp := callback(request)

		// skip response if it is EmptyPacket
		if resp == EmptyPacket {
			return
		}

		// add meta data
		resp.PacketID = request.PacketID

		// transfer WSPacket to json string
		respPacketString, err := json.Marshal(resp)
		if err != nil {
			log.WithFields(log.Fields{
				"ClientID":    c.ClientID,
				"Packet Data": resp.Data,
				"error":       err,
			}).Warn("Failed to marshal websocket packet into json string during responding")
		}

		// send response data
		c.SendLock.Lock()
		c.WebsocketConnection.SetWriteDeadline(time.Now().Add(20 * time.Second))
		c.WebsocketConnection.WriteMessage(websocket.TextMessage, respPacketString)
		c.SendLock.Unlock()
	}
}

/*
@func: Listen
@description:

	listen loop for the client
*/
func (c *Client) Listen() {
	for {
		// config recv deadline
		c.WebsocketConnection.SetReadDeadline(time.Now().Add(20 * time.Second))

		// recv data from ws connection
		_, rawMsg, err := c.WebsocketConnection.ReadMessage()
		if err != nil {
			log.WithFields(log.Fields{
				"ClientID": c.ClientID,
				"error":    err,
			}).Warn("Failed to receive data from websocket connection, close websocket connection")
			close(c.Done)
			break
		}

		// parse websocket packet
		wspacket := WSPacket{}
		err = json.Unmarshal(rawMsg, &wspacket)
		if err != nil {
			log.WithFields(log.Fields{
				"ClientID":    c.ClientID,
				"Raw Message": rawMsg,
				"error":       err,
			}).Warn("Failed to parse websocket")
			continue
		}

		// check send callback based on websocket packet id
		c.SendCallbackListLock.Lock()
		callback, ok := c.SendCallbackList[wspacket.PacketID]
		c.SendCallbackListLock.Unlock()
		if ok {
			go callback(wspacket)
			c.SendCallbackListLock.Lock()
			delete(c.SendCallbackList, wspacket.PacketID)
			c.SendCallbackListLock.Unlock()
			continue
		}

		// check send callback based on websocket packet type
		if callback, ok := c.RecvCallbackList[wspacket.PacketType]; ok {
			go callback(wspacket)
		}
	}
}

/*
@func: Close
@description:

	close websocket connection
*/
func (c *Client) Close() {
	if c == nil || c.WebsocketConnection == nil {
		return
	}
	c.WebsocketConnection.Close()
}
