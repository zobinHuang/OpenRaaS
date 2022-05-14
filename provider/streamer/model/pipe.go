package model

import (
	"fmt"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"github.com/zobinHuang/BrosCloud/provider/streamer/utils"

	log "github.com/sirupsen/logrus"
)

const (
	// length of pipe buffer
	VIDEO_PIPE_CHANNEL_LENGTH int = 100
	AUDIO_PIPE_CHANNEL_LENGTH int = 100
	INPUT_PIPE_CHANNEL_LENGTH int = 100

	// video encode
	VCODEC_H264 string = "h264"
	VCODEC_VPX  string = "vpx"
)

/*
	@model: WebRTCPipe
	@description:
		model for transmitting webrtc streams
*/
type WebRTCPipe struct {
	PeerConnection *webrtc.PeerConnection
	StreamInstance *StreamInstanceDaemonModel
	ConsumerID     string
	done           chan struct{}
	VideoChan      chan *rtp.Packet
	AudioChan      chan *rtp.Packet
	InputChan      chan []byte
	ICECandidate   string
	isConnected    bool
	isClosed       bool
}

/*
	@function: Open
	@description:
		[1] create new WebRTC peer connection
		[2] create video, audio, input track for newly created connection
		[3] register callbacks of input track and WebRTC connection
*/
func (p *WebRTCPipe) Open(iceServers []string, vCodec string, onICECandidateCallback func(candidate string)) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("Error occurs while open WebRTC Pipe, closed")
			p.Close()
		}
	}()

	// reset WebRTC Pipe if it has been opened
	if p.isConnected {
		log.WithFields(log.Fields{
			"Instance ID": p.StreamInstance.Instanceid,
			"Consumer ID": p.ConsumerID,
		}).Warn("The WebRTC Pipe has been opened, try to reset, wait for a minute")

		p.Close()
		time.Sleep(5 * time.Second)
	}

	// construct webRTC Configuration
	webRTCConfig := webrtc.Configuration{ICEServers: []webrtc.ICEServer{
		{
			URLs: iceServers,
		},
	}}

	// create new connection
	webRTCConnection, err := webrtc.NewPeerConnection(webRTCConfig)
	if err != nil {
		return "", fmt.Errorf("Failed to create new WebRTC peer connection\n")
	}

	// config video codec
	var codec string
	switch vCodec {
	case VCODEC_H264:
		codec = webrtc.MimeTypeH264
	case VCODEC_VPX:
		codec = webrtc.MimeTypeVP8
	default:
		codec = webrtc.MimeTypeVP8
	}

	log.WithFields(log.Fields{
		"Instance ID": p.StreamInstance.Instanceid,
		"Consumer ID": p.ConsumerID,
		"Video Codec": codec,
	}).Info("Choose video codec")

	// create video track and add it to peer connection
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: codec}, "video", "pion")
	if err != nil {
		return "", fmt.Errorf("Failed to create video track: %s\n", err.Error())
	}
	_, err = webRTCConnection.AddTrack(videoTrack)
	if err != nil {
		return "", fmt.Errorf("Failed to add video track to peer connection: %s\n", err.Error())
	}

	// create audio track and add it to peer connection
	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		return "", fmt.Errorf("Failed to create audio track: %s\n", err.Error())
	}
	_, err = webRTCConnection.AddTrack(audioTrack)
	if err != nil {
		return "", fmt.Errorf("Failed to add audio track to peer connection: %s\n", err.Error())
	}

	_, err = webRTCConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	})

	// create data channel and register callbacks of input track
	inputTrack, err := webRTCConnection.CreateDataChannel("input", nil)

	/*
		@callback: OnOpen
		@description:
			log input track is opened
	*/
	inputTrack.OnOpen(func() {
		log.WithFields(log.Fields{
			"Instance ID": p.StreamInstance.Instanceid,
			"Consumer ID": p.ConsumerID,
		}).Info("Input track of WebRTC Pipe is opened")
	})

	/*
		@callback: OnOpen
		@description:
			forward message from input track to input channel of WebRTC Pipe
	*/
	inputTrack.OnMessage(func(msg webrtc.DataChannelMessage) {
		p.InputChan <- msg.Data
	})

	/*
		@callback: OnClose
		@description:
			log input track is closed
	*/
	inputTrack.OnClose(func() {
		log.WithFields(log.Fields{
			"Instance ID": p.StreamInstance.Instanceid,
			"Consumer ID": p.ConsumerID,
		}).Warn("Input track if WebRTC Pipe is closed")
	})

	// register callbacks of WebRTC Connection
	/*
		@callback: OnICECandidate
		@description:
			operations while ice connection state changed
	*/
	webRTCConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.WithFields(log.Fields{
			"Instance ID":          p.StreamInstance.Instanceid,
			"Consumer ID":          p.ConsumerID,
			"ICE Connection State": connectionState.String(),
		}).Warn("ICE connection state has changed")

		// start WebRTC streaming if ice connection state set to connected
		if connectionState == webrtc.ICEConnectionStateConnected {
			go func() {
				p.isConnected = true
				p.StartStreaming(videoTrack, audioTrack)
			}()
		}

		// close pipe if ice connection state set to failed | disconnected | closed
		if connectionState == webrtc.ICEConnectionStateClosed || connectionState == webrtc.ICEConnectionStateFailed || connectionState == webrtc.ICEConnectionStateDisconnected {
			p.Close()
		}
	})

	/*
		@callback: OnICECandidate
		@description:
			operation of onICECandidate
	*/
	webRTCConnection.OnICECandidate(func(iceCandidate *webrtc.ICECandidate) {
		if iceCandidate != nil {
			candidate, err := utils.EncodeBase64(iceCandidate.ToJSON())
			if err != nil {
				log.WithFields(log.Fields{
					"Instance ID":            p.StreamInstance.Instanceid,
					"Consumer ID":            p.ConsumerID,
					"Received ICE Candidate": candidate,
				}).Warn("Failed to encode ICE candidate into base64 string, abandoned")
				return
			}
			onICECandidateCallback(candidate)
		} else {
			onICECandidateCallback("")
		}
	})

	// store peer connection
	p.PeerConnection = webRTCConnection

	// create offer SDP
	offer, err := p.PeerConnection.CreateOffer(nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create offer SDP: %s\n", err.Error())
	}

	// set local description
	err = p.PeerConnection.SetLocalDescription(offer)
	if err != nil {
		return "", fmt.Errorf("Failed to set local description: %s\n", err.Error())
	}
	log.WithFields(log.Fields{
		"Instance ID": p.StreamInstance.Instanceid,
		"Consumer ID": p.ConsumerID,
	}).Info("Set local description")

	// encode offer SDP
	offerSDP, err := utils.EncodeBase64(offer)
	if err != nil {
		return "", fmt.Errorf("Failed to encode offer SDP into base64 string")
	}

	return offerSDP, nil
}

/*
	@function: SetRemoteSDP
	@description:
		set remote sdp while recevied it from the consumer
*/
func (p *WebRTCPipe) SetRemoteSDP(remoteSDP string) error {
	// decode answer sdp
	var answer webrtc.SessionDescription
	err := utils.DecodeBase64(remoteSDP, &answer)
	if err != nil {
		log.WithFields(log.Fields{
			"Instance ID": p.StreamInstance.Instanceid,
			"Consumer ID": p.ConsumerID,
			"error":       err,
		}).Warn("Failed to decode answer SDP from base64, abandoned")
		return err
	}

	// set remote sdp
	err = p.PeerConnection.SetRemoteDescription(answer)
	if err != nil {
		log.WithFields(log.Fields{
			"Instance ID": p.StreamInstance.Instanceid,
			"Consumer ID": p.ConsumerID,
			"error":       err,
		}).Warn("Failed to set remote description in peer connection, abandoned")
		return err
	}

	log.WithFields(log.Fields{
		"Instance ID": p.StreamInstance.Instanceid,
		"Consumer ID": p.ConsumerID,
	}).Info("Set remote description of consumer")

	return nil
}

/*
	@function: AddCandidate
	@description:
		add ice candidate of remote consumer
*/
func (p *WebRTCPipe) AddCandidate(candidate string) error {
	// decode ICE candidate
	var iceCandidate webrtc.ICECandidateInit
	err := utils.DecodeBase64(candidate, &iceCandidate)
	if err != nil {
		log.WithFields(log.Fields{
			"Instance ID": p.StreamInstance.Instanceid,
			"Consumer ID": p.ConsumerID,
			"error":       err,
		}).Warn("Failed to decode ICE candidate from base64, abandoned")
	}

	err = p.PeerConnection.AddICECandidate(iceCandidate)
	if err != nil {
		log.WithFields(log.Fields{
			"Instance ID": p.StreamInstance.Instanceid,
			"Consumer ID": p.ConsumerID,
			"error":       err,
		}).Warn("Failed to add ice candidate in peer connection, abandoned")
		return err
	}

	log.WithFields(log.Fields{
		"Instance ID":   p.StreamInstance.Instanceid,
		"Consumer ID":   p.ConsumerID,
		"ICE Candidate": iceCandidate.Candidate,
	}).Info("Add ice candidate of consumer")

	return nil
}

/*
	@function: StartStreaming
	@description:
		start WebRTC streaming
*/
func (p *WebRTCPipe) StartStreaming(videoTrack *webrtc.TrackLocalStaticRTP, audioTrack *webrtc.TrackLocalStaticRTP) {
	log.WithFields(log.Fields{
		"Instance ID": p.StreamInstance.Instanceid,
		"Consumer ID": p.ConsumerID,
	}).Info("Start streaming")

	// streaming video
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"Instance ID": p.StreamInstance.Instanceid,
					"Consumer ID": p.ConsumerID,
					"error":       err,
				}).Info("Recover from error while streaming video")
			}
		}()

		for packet := range p.VideoChan {
			if err := videoTrack.WriteRTP(packet); err != nil {
				log.WithFields(log.Fields{
					"Instance ID": p.StreamInstance.Instanceid,
					"Consumer ID": p.ConsumerID,
				}).Info("Error occurs while streaming video, panic")
				panic(err)
			}
		}
	}()

	// streaming audio
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"Instance ID": p.StreamInstance.Instanceid,
					"Consumer ID": p.ConsumerID,
					"error":       err,
				}).Info("Recover from error while streaming audio")
			}
		}()

		for packet := range p.AudioChan {
			if err := audioTrack.WriteRTP(packet); err != nil {
				log.WithFields(log.Fields{
					"Instance ID": p.StreamInstance.Instanceid,
					"Consumer ID": p.ConsumerID,
				}).Info("Error occurs while streaming audio, panic")
				panic(err)
			}
		}
	}()
}

/*
	@function: Close
	@description:
		close WebRTC connection and all allocated channels
*/
func (p *WebRTCPipe) Close() {
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"Instance ID": p.StreamInstance.Instanceid,
				"Consumer ID": p.ConsumerID,
			}).Warn("Recover from error since try to close a close channel")
		}
	}()

	// close WebRTC connection
	if p.PeerConnection != nil {
		p.PeerConnection.Close()
		p.PeerConnection = nil
		log.WithFields(log.Fields{
			"Instance ID": p.StreamInstance.Instanceid,
			"Consumer ID": p.ConsumerID,
		}).Info("Close WebRTC Connection")
	}

	// close channels of WebRTC Pipe
	close(p.VideoChan)
	close(p.AudioChan)
	close(p.InputChan)
}
