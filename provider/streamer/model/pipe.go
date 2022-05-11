package model

import (
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

const (
	VIDEO_PIPE_CHANNEL_LENGTH int = 100
	AUDIO_PIPE_CHANNEL_LENGTH int = 100
	INPUT_PIPE_CHANNEL_LENGTH int = 100
)

/*
	@model: WebRTCPipe
	@description:
		model for transmitting webrtc streams
*/
type WebRTCPipe struct {
	StreamInstance *StreamInstanceDaemonModel
	ConsumerID     string
	RTCPeerConn    *webrtc.PeerConnection
	VideoChan      chan *rtp.Packet
	AudioChan      chan *rtp.Packet
	InputChan      chan []byte
}
