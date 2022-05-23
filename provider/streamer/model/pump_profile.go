package model

import "github.com/robfig/cron"

/*
	@model: PumpProfiler
	@description:
		model for profiling pump
*/
type PumpProfiler struct {
	PerSecondCronTask         *cron.Cron
	VideoPacketCounter        uint64
	AudioPacketCounter        uint64
	VideoByteCounter          uint64
	AudioByteCounter          uint64
	PerSec_VideoPacketCounter uint32
	PerSec_AudioPacketCounter uint32
	PerSec_VideoByteCounter   uint32
	PerSec_AudioByteCounter   uint32
}

/*
	@func: AddVideoPacketCounter
	@description:
		add video packet count to profiler
*/
func (pp *PumpProfiler) AddVideoPacketCounter(adder uint64) {
	pp.VideoPacketCounter = pp.VideoPacketCounter + adder
}

/*
	@func: AddVideoByteCounter
	@description:
		add video byte count to profiler
*/
func (pp *PumpProfiler) AddVideoByteCounter(adder uint64) {
	pp.VideoByteCounter = pp.VideoByteCounter + adder
}

/*
	@func: AddAudioPacketCounter
	@description:
		add audio packet count to profiler
*/
func (pp *PumpProfiler) AddAudioPacketCounter(adder uint64) {
	pp.AudioPacketCounter = pp.AudioPacketCounter + adder
}

/*
	@func: AddAudioByteCounter
	@description:
		add audio byte count to profiler
*/
func (pp *PumpProfiler) AddAudioByteCounter(adder uint64) {
	pp.AudioByteCounter = pp.AudioByteCounter + adder
}

/*
	@func: AddPerSecVideoPacketCounter
	@description:
		add video packet count to profiler (per second)
*/
func (pp *PumpProfiler) AddPerSecVideoPacketCounter(adder uint32) {
	pp.PerSec_VideoPacketCounter = pp.PerSec_VideoPacketCounter + adder
}

/*
	@func: ClearPerSecVideoPacketCounter
	@description:
		clear video packet count (per second)
*/
func (pp *PumpProfiler) ClearPerSecVideoPacketCounter() {
	pp.PerSec_VideoPacketCounter = 0
}

/*
	@func: AddPerSecVideoByteCounter
	@description:
		add video byte count to profiler (per second)
*/
func (pp *PumpProfiler) AddPerSecVideoByteCounter(adder uint32) {
	pp.PerSec_VideoByteCounter = pp.PerSec_VideoByteCounter + adder
}

/*
	@func: ClearPerSecVideoByteCounter
	@description:
		clear video byte count (per second)
*/
func (pp *PumpProfiler) ClearPerSecVideoByteCounter() {
	pp.PerSec_VideoByteCounter = 0
}

/*
	@func: AddPerSecAudioPacketCounter
	@description:
		add audio packet count to profiler (per second)
*/
func (pp *PumpProfiler) AddPerSecAudioPacketCounter(adder uint32) {
	pp.PerSec_AudioPacketCounter = pp.PerSec_AudioPacketCounter + adder
}

/*
	@func: ClearPerSecAudioPacketCounter
	@description:
		clear audio packet count (per second)
*/
func (pp *PumpProfiler) ClearPerSecAudioPacketCounter() {
	pp.PerSec_AudioPacketCounter = 0
}

/*
	@func: AddPerSecAudioByteCounter
	@description:
		add audio byte count to profiler (per second)
*/
func (pp *PumpProfiler) AddPerSecAudioByteCounter(adder uint32) {
	pp.PerSec_AudioByteCounter = pp.PerSec_AudioByteCounter + adder
}

/*
	@func: ClearPerSecAudioByteCounter
	@description:
		clear audio byte count (per second)
*/
func (pp *PumpProfiler) ClearPerSecAudioByteCounter() {
	pp.PerSec_AudioByteCounter = 0
}
