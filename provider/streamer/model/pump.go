package model

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/robfig/cron"
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
	VIDEO_LISTENER_RING_LENGTH int = 300
	AUDIO_LISTENER_RING_LENGTH int = 300
	VIDEO_PUMP_CHANNEL_LENGTH  int = 300
	AUDIO_PUMP_CHANNEL_LENGTH  int = 300
)

const (
	INPUT_EVENT_WS_TYPE_KEYDOWN   = "KEYDOWN"
	INPUT_EVENT_WS_TYPE_KEYUP     = "KEYUP"
	INPUT_EVENT_WS_TYPE_MOUSEDOWN = "MOUSEDOWN"
	INPUT_EVENT_WS_TYPE_MOUSEUP   = "MOUSEUP"
	INPUT_EVENT_WS_TYPE_MOUSEMOVE = "MOUSEMOVE"
)

/*
	@model: Pump
	@description:
		model for hijacking instance streams
*/
type Pump struct {
	StreamInstance  *StreamInstanceDaemonModel
	PumpProfiler    PumpProfiler
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
	@func: PerSecondProfiling
	@description:
		per-second profiling
*/
func (p *Pump) PerSecondProfiling() {
	p.PumpProfiler.PerSecondCronTask = cron.New()
	p.PumpProfiler.PerSecondCronTask.AddFunc("*/1 * * * * * ", func() {
		// log overall statistics
		log.WithFields(log.Fields{
			"Instance ID":        p.StreamInstance.Instanceid,
			"Video Packet Count": p.PumpProfiler.VideoPacketCounter,
			"Audio Packet Count": p.PumpProfiler.AudioPacketCounter,
			"Video Byte Count":   p.PumpProfiler.VideoByteCounter,
			"Audio Byte Count":   p.PumpProfiler.AudioByteCounter,
		}).Info("Pump overall profiling")

		// log per second statistics
		log.WithFields(log.Fields{
			"Instance ID":       p.StreamInstance.Instanceid,
			"Video Packet Rate": p.PumpProfiler.PerSec_VideoPacketCounter,
			"Audio Packet Rate": p.PumpProfiler.PerSec_AudioPacketCounter,
			"Video Byte Rate":   p.PumpProfiler.PerSec_VideoByteCounter,
			"Audio Byte Rate":   p.PumpProfiler.PerSec_AudioByteCounter,
		}).Info("Pump per second profiling")

		// clear per second counter
		p.PumpProfiler.ClearPerSecVideoPacketCounter()
		p.PumpProfiler.ClearPerSecVideoByteCounter()
		p.PumpProfiler.ClearPerSecAudioPacketCounter()
		p.PumpProfiler.ClearPerSecAudioByteCounter()
	})

	p.PumpProfiler.PerSecondCronTask.Start()
}

/*
	@func: CreateVideoListener
	@description:
		create UDP listened on video stream
*/
func (s *Pump) CreateVideoListener() error {
	// obtain listen metadata
	videoRTCPort, _ := strconv.Atoi(s.StreamInstance.VideoRTCPort)

	log.WithFields(log.Fields{
		"Video RTC Port": s.StreamInstance.VideoRTCPort,
	}).Info("Try to create video listener")

	// obtain listen
	listener, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: videoRTCPort})
	if err != nil {
		return fmt.Errorf("failed to obtain listener of the video stream, %s", err.Error())
	}

	// listen for a single RTP packet to determine the SSRC of video stream
	inboundRTPPacket := make([]byte, 4096)
	n, _, err := listener.ReadFromUDP(inboundRTPPacket)
	if err != nil {
		return fmt.Errorf("failed to listen on video stream")
	}

	// unmarshal the incoming packet
	packet := &rtp.Packet{}
	if err = packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
		return fmt.Errorf("failed to unmarshal RTP packet received from video stream")
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
func (s *Pump) CreateAudioListener() error {
	// obtain listen metadata
	audioRTCPort, _ := strconv.Atoi(s.StreamInstance.AudioRTCPort)

	log.WithFields(log.Fields{
		"Video RTC Port": s.StreamInstance.AudioRTCPort,
	}).Info("Try to create audio listener")

	// obtain listen
	listener, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: audioRTCPort})
	if err != nil {
		return fmt.Errorf("failed to obtain listener of the audio stream, %s", err.Error())
	}

	// listen for a single RTP packet to determine the SSRC of audio stream
	inboundRTPPacket := make([]byte, 4096)
	n, _, err := listener.ReadFromUDP(inboundRTPPacket)
	if err != nil {
		return fmt.Errorf("failed to listen on audio stream")
	}

	// unmarshal the incoming packet
	packet := &rtp.Packet{}
	if err = packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
		return fmt.Errorf("failed to unmarshal RTP packet received from audio stream")
	}

	// record in model
	s.AudioListener = listener
	s.AudioStreamSSRC = packet.SSRC

	return nil
}

/*
	@func: CreateInputSimulator
	@description:
		construct connection to the applicatio instance container
*/
func (s *Pump) CreateInputSimulator() error {
	// resolve TCP address
	addressString := fmt.Sprintf("0.0.0.0:%s", s.StreamInstance.InputPort)
	address, err := net.ResolveTCPAddr("tcp4", addressString)
	if err != nil {
		return fmt.Errorf("failed to resolve TCP address: %s", addressString)
	}

	// start listen
	connection, err := net.ListenTCP("tcp4", address)
	if err != nil {
		return fmt.Errorf("failed to listen on TCP address: %s", addressString)
	}

	// accept connection within goroutine
	go func() {
		conn, err := connection.AcceptTCP()
		if err != nil {
			log.WithFields(log.Fields{
				"Instance ID": s.StreamInstance.Instanceid,
				"error":       err.Error(),
			}).Warn("Failed to accept tcp connection with instance container")
		}

		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(10 * time.Second)
		s.wineConn = conn

		log.WithFields(log.Fields{
			"Instance ID": s.StreamInstance.Instanceid,
		}).Warn("Accept tcp connection with instance container for input track")
	}()

	return nil
}

/*
	@func: ListenVideoStream
	@description:
		start a goroutine to listen on video stream
*/
func (s *Pump) ListenVideoStream() {
	go func() {
		// defer the closure of video stream listener
		defer func() {
			s.VideoListener.Close()
			log.WithFields(log.Fields{
				"Stream Instance ID": s.StreamInstance.Instanceid,
			}).Warn("Pump stopped listen to the video stream from the instance")
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

			if ENABLE_PUPMP_PROFILING {
				// profile (overall)
				s.PumpProfiler.AddVideoPacketCounter(1)
				s.PumpProfiler.AddVideoByteCounter(uint64(n))

				// profile (per second)
				s.PumpProfiler.AddPerSecVideoPacketCounter(1)
				s.PumpProfiler.AddPerSecVideoByteCounter(uint32(n))
			}
		}
	}()
}

/*
	@func: ListenAudioStream
	@description:
		start a goroutine to listen on audio stream
*/
func (s *Pump) ListenAudioStream() {
	go func() {
		// defer the closure of audio stream listener
		defer func() {
			s.AudioListener.Close()
			log.WithFields(log.Fields{
				"Stream Instance ID": s.StreamInstance.Instanceid,
			}).Warn("Pump stopped listen to the audio stream from the instance")
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

			if ENABLE_PUPMP_PROFILING {
				// profile (overall)
				s.PumpProfiler.AddAudioPacketCounter(1)
				s.PumpProfiler.AddAudioByteCounter(uint64(n))

				// profile (per second)
				s.PumpProfiler.AddPerSecAudioPacketCounter(1)
				s.PumpProfiler.AddPerSecAudioByteCounter(uint32(n))
			}
		}
	}()
}

/*
	@func: AddWebRTCPipe
	@description:
		add a new WebRTC pipe to the hub
*/
func (s *Pump) AddWebRTCPipe(p *WebRTCPipe) {
	s.Hub = append(s.Hub, p)
}

/*
	@func: Discharge
	@description:
		discharge local streams to loaded WebRTCPipe in the hub
*/
func (s *Pump) Discharge() {
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

func (s *Pump) HarvestInput(pipe *WebRTCPipe) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.WithFields(log.Fields{
					"Instance ID": s.StreamInstance.Instanceid,
					"Consumer ID": pipe.ConsumerID,
					"error":       r,
				}).Warn("Recover from discharging audio stream, error occurs since read from closed input data channel")
			}
		}()

		type keyboardEventData struct {
			KeyCode int `json:"key_code"`
		}

		type mouseEventData struct {
			IsLeft byte    `json:"is_left"`
			X      float32 `json:"x"`
			Y      float32 `json:"y"`
			Width  float32 `json:"width"`
			Height float32 `json:"height"`
		}

		for rawInput := range pipe.InputChan {
			// parse packet
			inputEvent := &WSPacket{}
			err := json.Unmarshal(rawInput, &inputEvent)
			if err != nil {
				log.WithFields(log.Fields{
					"Instance ID": s.StreamInstance.Instanceid,
					"Consumer ID": pipe.ConsumerID,
					"error":       err.Error(),
				}).Warn("Failed to parse packet from input track, abandoned")
				continue
			}

			switch inputEvent.PacketType {
			/*
				@case: detect input event as keyup
			*/
			case INPUT_EVENT_WS_TYPE_KEYUP:
				p := &keyboardEventData{}
				json.Unmarshal([]byte(inputEvent.Data), &p)
				vmKeyMsg := fmt.Sprintf("K%d,%b|", p.KeyCode, 0)
				_, err := s.wineConn.Write([]byte(vmKeyMsg))
				if err != nil {
					log.WithFields(log.Fields{
						"Given Keyboard Code": p.KeyCode,
						"Instance ID":         s.StreamInstance.Instanceid,
						"Consumer ID":         pipe.ConsumerID,
						"error":               err.Error(),
					}).Warn("Failed to pass keyup event to instance container")
				}

			/*
				@case: detect input event as keydown
			*/
			case INPUT_EVENT_WS_TYPE_KEYDOWN:
				p := &keyboardEventData{}
				json.Unmarshal([]byte(inputEvent.Data), &p)
				vmKeyMsg := fmt.Sprintf("K%d,%b|", p.KeyCode, 1)
				_, err := s.wineConn.Write([]byte(vmKeyMsg))
				if err != nil {
					log.WithFields(log.Fields{
						"Given Keyboard Code": p.KeyCode,
						"Instance ID":         s.StreamInstance.Instanceid,
						"Consumer ID":         pipe.ConsumerID,
						"error":               err.Error(),
					}).Warn("Failed to pass keyup event to instance container")
				}

			/*
				@case: detect input event as mouseup
			*/
			case INPUT_EVENT_WS_TYPE_MOUSEUP:
				p := &mouseEventData{}
				json.Unmarshal([]byte(inputEvent.Data), &p)
				p.X = p.X * float32(s.StreamInstance.ScreenWidth) / p.Width
				p.Y = p.Y * float32(s.StreamInstance.ScreenHeight) / p.Height

				vmKeyMsg := fmt.Sprintf("M%d,%d,%f,%f,%f,%f|", p.IsLeft, 2, p.X, p.Y, p.Width, p.Height)
				_, err := s.wineConn.Write([]byte(vmKeyMsg))
				if err != nil {
					log.WithFields(log.Fields{
						"Instance ID": s.StreamInstance.Instanceid,
						"Consumer ID": pipe.ConsumerID,
						"error":       err.Error(),
					}).Warn("Failed to pass mouseup event to instance container")
				}

			/*
				@case: detect input event as mousedown
			*/
			case INPUT_EVENT_WS_TYPE_MOUSEDOWN:
				p := &mouseEventData{}
				json.Unmarshal([]byte(inputEvent.Data), &p)
				p.X = p.X * float32(s.StreamInstance.ScreenWidth) / p.Width
				p.Y = p.Y * float32(s.StreamInstance.ScreenHeight) / p.Height

				vmKeyMsg := fmt.Sprintf("M%d,%d,%f,%f,%f,%f|", p.IsLeft, 1, p.X, p.Y, p.Width, p.Height)
				_, err := s.wineConn.Write([]byte(vmKeyMsg))
				if err != nil {
					log.WithFields(log.Fields{
						"Instance ID": s.StreamInstance.Instanceid,
						"Consumer ID": pipe.ConsumerID,
						"error":       err.Error(),
					}).Warn("Failed to pass mousedown event to instance container")
				}

			/*
				@case: detect input event as mousemove
			*/
			case INPUT_EVENT_WS_TYPE_MOUSEMOVE:
				p := &mouseEventData{}
				json.Unmarshal([]byte(inputEvent.Data), &p)
				p.X = p.X * float32(s.StreamInstance.ScreenWidth) / p.Width
				p.Y = p.Y * float32(s.StreamInstance.ScreenHeight) / p.Height

				vmKeyMsg := fmt.Sprintf("M%d,%d,%f,%f,%f,%f|", p.IsLeft, 0, p.X, p.Y, p.Width, p.Height)
				_, err := s.wineConn.Write([]byte(vmKeyMsg))
				if err != nil {
					log.WithFields(log.Fields{
						"Instance ID": s.StreamInstance.Instanceid,
						"Consumer ID": pipe.ConsumerID,
						"error":       err.Error(),
					}).Warn("Failed to pass mousemove event to instance container")
				}
			}
		}
	}()
}
