package model

import (
	"container/ring"
	"fmt"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/pion/rtp"
)

/*
	@constant:
		VIDEO_LISTENER_RING_LENGTH
		AUDIO_LISTENER_RING_LENGTH
	@description:
		length of ring buffer of audio and video listener
*/
const (
	VIDEO_LISTENER_RING_LENGTH int = 100
	AUDIO_LISTENER_RING_LENGTH int = 100
)

/*
	@model: Websocket
	@description:
		model for hijacking instance streams
*/
type WebRTCStreamer struct {
	StreamInstance  *StreamInstanceDaemonModel
	VideoListener   *net.UDPConn
	AudioListener   *net.UDPConn
	wineConn        *net.TCPConn
	VideoStreamSSRC uint32
	AudioStreamSSRC uint32
	VideoStream     chan *rtp.Packet
	AudioStream     chan *rtp.Packet
	Hub             []*WebRTCPipe
}

/*
	@func: CreateVideoListener
	@description:
		create UDP listened on video stream
*/
func (s *WebRTCStreamer) CreateVideoListener() error {
	// obtain listen metadata
	videoRTCPort, _ := strconv.Atoi(s.StreamInstance.VideoRTCPort)

	log.WithFields(log.Fields{
		"Video RTC Port": s.StreamInstance.VideoRTCPort,
	}).Info("Try to create video listener")

	// obtain listen
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: videoRTCPort})
	if err != nil {
		return fmt.Errorf("Failed to obtain listener of the video stream")
	}

	// listen for a single RTP packet to determine the SSRC of video stream
	inboundRTPPacket := make([]byte, 4096)
	n, _, err := listener.ReadFromUDP(inboundRTPPacket)
	if err != nil {
		return fmt.Errorf("Failed to listen on video stream")
	}

	// unmarshal the incoming packet
	packet := &rtp.Packet{}
	if err = packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
		return fmt.Errorf("Failed to unmarshal RTP packet received from video stream")
	}

	// record in model
	s.VideoListener = listener
	s.VideoStreamSSRC = packet.SSRC

	return nil
}

/*
	@func: CreateAudioListener
	@description:
		create UDP listened on audio stream
*/
func (s *WebRTCStreamer) CreateAudioListener() error {
	// obtain listen metadata
	audioRTCPort, _ := strconv.Atoi(s.StreamInstance.AudioRTCPort)

	log.WithFields(log.Fields{
		"Video RTC Port": s.StreamInstance.AudioRTCPort,
	}).Info("Try to create audio listener")

	// obtain listen
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: audioRTCPort})
	if err != nil {
		return fmt.Errorf("Failed to obtain listener of the audio stream")
	}

	// listen for a single RTP packet to determine the SSRC of audio stream
	inboundRTPPacket := make([]byte, 4096)
	n, _, err := listener.ReadFromUDP(inboundRTPPacket)
	if err != nil {
		return fmt.Errorf("Failed to listen on audio stream")
	}

	// unmarshal the incoming packet
	packet := &rtp.Packet{}
	if err = packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
		return fmt.Errorf("Failed to unmarshal RTP packet received from audio stream")
	}

	// record in model
	s.AudioListener = listener
	s.AudioStreamSSRC = packet.SSRC

	return nil
}

/*
	@func: ListenVideoStream
	@description:
		start a goroutine to listen on video stream
*/
func (s *WebRTCStreamer) ListenVideoStream() {
	go func() {
		// defer the closure of video stream listener
		defer func() {
			s.VideoListener.Close()
			log.WithFields(log.Fields{
				"Stream Instance ID": s.StreamInstance.Instanceid,
			}).Warn("WebRTCStreamer stopped listen to the video stream from the instance")
		}()

		// initialize a ring buffer
		ringBuffer := ring.New(VIDEO_LISTENER_RING_LENGTH)
		for i := 0; i < VIDEO_LISTENER_RING_LENGTH; i++ {
			ringBuffer.Value = make([]byte, 1500)
			ringBuffer = ringBuffer.Next()
		}

		// streaming loop
		for {
			inboundRTPPacket := ringBuffer.Value.([]byte)
			ringBuffer = ringBuffer.Next()

			n, _, err := s.VideoListener.ReadFrom(inboundRTPPacket)
			if err != nil {
				log.WithFields(log.Fields{
					"Stream Instance ID": s.StreamInstance.Instanceid,
					"error":              err.Error(),
				}).Warn("Error occurs while fetching video stream, continued")
				continue
			}

			packet := &rtp.Packet{}
			if err := packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
				log.WithFields(log.Fields{
					"Stream Instance ID": s.StreamInstance.Instanceid,
					"error":              err.Error(),
				}).Warn("Error occurs while unmarshal UDP datagram of video stream into RTP Packet, continued")
				continue
			}

			s.VideoStream <- packet
		}
	}()
}

/*
	@func: ListenAudioStream
	@description:
		start a goroutine to listen on audio stream
*/
func (s *WebRTCStreamer) ListenAudioStream() {
	go func() {
		// defer the closure of audio stream listener
		defer func() {
			s.AudioListener.Close()
			log.WithFields(log.Fields{
				"Stream Instance ID": s.StreamInstance.Instanceid,
			}).Warn("WebRTCStreamer stopped listen to the audio stream from the instance")
		}()

		// initialize a ring buffer
		ringBuffer := ring.New(AUDIO_LISTENER_RING_LENGTH)
		for i := 0; i < AUDIO_LISTENER_RING_LENGTH; i++ {
			ringBuffer.Value = make([]byte, 1500)
			ringBuffer = ringBuffer.Next()
		}

		// streaming loop
		for {
			inboundRTPPacket := ringBuffer.Value.([]byte)
			ringBuffer = ringBuffer.Next()

			n, _, err := s.AudioListener.ReadFrom(inboundRTPPacket)
			if err != nil {
				log.WithFields(log.Fields{
					"Stream Instance ID": s.StreamInstance.Instanceid,
					"error":              err.Error(),
				}).Warn("Error occurs while fetching audio stream, continued")
				continue
			}

			packet := &rtp.Packet{}
			if err := packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
				log.WithFields(log.Fields{
					"Stream Instance ID": s.StreamInstance.Instanceid,
					"error":              err.Error(),
				}).Warn("Error occurs while unmarshal UDP datagram of audio stream into RTP Packet, continued")
				continue
			}

			s.AudioStream <- packet
		}
	}()
}

/*
	@func: AddWebRTCPipe
	@description:
		add a new WebRTC pipe to the hub
*/
func (s *WebRTCStreamer) AddWebRTCPipe(p *WebRTCPipe) {
	s.Hub = append(s.Hub, p)
}

/*
	@func: Discharge
	@description:
		discharge local streams to loaded WebRTCPipe in the hub
*/
func (s *WebRTCStreamer) Discharge() {
	// discharge video stream
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.WithFields(log.Fields{
					"Instance ID": s.StreamInstance.Instanceid,
					"error":       r,
				}).Warn("Recover from discharging video stream, error occurs since sent to closed video stream channel")
			}
		}()

		for packet := range s.VideoStream {
			for pipeIndex, pipe := range s.Hub {
				select {
				case <-pipe.done:
					s.Hub = append(s.Hub[:pipeIndex], s.Hub[pipeIndex+1:]...)
					close(pipe.VideoChan)
					close(pipe.AudioChan)
				case pipe.VideoChan <- packet:
				}
			}
		}
	}()

	// discharge audio stream
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.WithFields(log.Fields{
					"Instance ID": s.StreamInstance.Instanceid,
					"error":       r,
				}).Warn("Recover from discharging audio stream, error occurs since sent to closed audio stream channel")
			}
		}()

		for packet := range s.AudioStream {
			for pipeIndex, pipe := range s.Hub {
				select {
				case <-pipe.done:
					s.Hub = append(s.Hub[:pipeIndex], s.Hub[pipeIndex+1:]...)
					close(pipe.VideoChan)
					close(pipe.AudioChan)
				case pipe.AudioChan <- packet:
				}
			}
		}
	}()
}
