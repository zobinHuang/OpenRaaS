package model

/*
	@model: WSPacket
	@description:
		represent a websocket packet
*/
type WSPacket struct {
	PacketType string `json:"packet_type"`
	PacketID   string `json:"packet_id"`
	Data       string `json:"data"`
}

/*
	@var: EmptyPacket
	@description:
		represent a pure empty packet, use as comparison
		object or return value
*/
var EmptyPacket = WSPacket{}
