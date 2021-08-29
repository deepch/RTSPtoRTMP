package main

import (
	"time"

	"github.com/deepch/vdk/codec/aacparser"

	"github.com/deepch/vdk/av"
)

type FakeAAC struct {
	closer              chan bool
	OutgoingPacketQueue string
}

func NewFakeAudio() *FakeAAC {
	return &FakeAAC{closer: make(chan bool)}
}

func (obj *FakeAAC) Start(CodecData []av.CodecData, OutgoingPacketQueue chan *av.Packet) ([]av.CodecData, error) {
	go func() {
		ticker := time.NewTicker(128 * time.Millisecond)
		for {
			select {
			case <-obj.closer:
				return
			case <-ticker.C:
				OutgoingPacketQueue <- &av.Packet{Duration: 128 * time.Millisecond, Idx: 1, CompositionTime: 1 * time.Millisecond, Data: []byte{1, 2, 52, 24, 67, 243, 0, 69, 107, 122, 169, 68, 50, 47, 46, 147, 26, 82, 161, 42, 80, 119, 97, 77, 199, 109, 20, 8, 26, 99, 53, 19, 85, 4, 186, 150, 69, 121, 208, 72, 42, 74, 224, 144, 172, 132, 146, 70, 246, 74, 13, 98, 76, 226, 5, 214, 169, 73, 18, 154, 33, 43, 11, 196, 2, 132, 43, 2, 240, 10, 14, 162, 64, 149, 243, 217, 96, 225, 85, 233, 90, 11, 19, 41, 67, 32, 5, 45, 5, 74, 72, 138, 215, 107, 33, 96, 49, 66, 85, 152, 150, 222, 87, 151, 89, 46, 80, 1, 172, 66, 150, 150, 150, 150, 150, 150, 150, 150, 150, 151, 7}}

			}
		}
	}()
	audioCodec, err := aacparser.NewCodecDataFromMPEG4AudioConfig(aacparser.MPEG4AudioConfig{SampleRate: 8000, SampleRateIndex: 11, ObjectType: 2, ChannelLayout: 1, ChannelConfig: 1})
	if err != nil {
		return nil, err
	}
	return append(CodecData, audioCodec), nil
}

func (obj *FakeAAC) Close() {
	obj.closer <- true
}
