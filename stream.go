package main

import (
	"errors"
	"log"
	"time"

	"github.com/deepch/vdk/av"

	"github.com/deepch/vdk/format/rtmp"

	"github.com/deepch/vdk/format/rtspv2"
)

var (
	ErrorStreamExitNoVideoOnStream = errors.New("Stream Exit No Video On Stream")
	ErrorStreamExitRtspDisconnect  = errors.New("Stream Exit Rtsp Disconnect")
	ErrorStreamExitNoViewer        = errors.New("Stream Exit On Demand No Viewer")
)

func serveStreams() {
	for k, v := range Config.Streams {
		if v.OnDemand {
			log.Println("OnDemand not supported")
			v.OnDemand = false
		}
		if !v.OnDemand {
			go RTSPWorkerLoop(k, v.URL, v.OnDemand)
		}
	}
}

func RTSPWorkerLoop(name, url string, OnDemand bool) {
	defer Config.RunUnlock(name)
	for {
		log.Println(name, "Stream Try Connect")
		err := RTSPWorker(name, url, OnDemand)
		if err != nil {
			log.Println(err)
		}
		if OnDemand && !Config.HasViewer(name) {
			log.Println(name, ErrorStreamExitNoViewer)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func RTSPWorker(name, url string, OnDemand bool) error {
	keyTest := time.NewTimer(20 * time.Second)
	clientTest := time.NewTimer(20 * time.Second)
	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: url, DisableAudio: false, DialTimeout: 3 * time.Second, ReadWriteTimeout: 3 * time.Second, Debug: false})
	if err != nil {
		return err
	}
	defer RTSPClient.Close()
	var connRTMP *rtmp.Conn
	broadCastEnable, broadCastURL := Config.bcConfig(name)
	if broadCastEnable {
		connRTMP, err = rtmp.DialTimeout(broadCastURL, 5*time.Second)
		if err != nil {
			return err
		}
		defer connRTMP.Close()
	}
	if len(RTSPClient.CodecData) < 2 {
		AudioStreamFake := NewFakeAudio()
		RTSPClient.CodecData, err = AudioStreamFake.Start(RTSPClient.CodecData, RTSPClient.OutgoingPacketQueue)
		defer AudioStreamFake.Close()
	} else if RTSPClient.CodecData[1].Type() != av.AAC {
		log.Fatalln("Only AAC or Disable Audio")
	}
	if broadCastEnable {
		err = connRTMP.WriteHeader(RTSPClient.CodecData)
		if err != nil {
			return err
		}
	}
	if RTSPClient.CodecData != nil {
		Config.coAd(name, RTSPClient.CodecData)
	}
	var start bool
	var lineA time.Duration
	var lineV time.Duration
	for {
		select {
		case <-clientTest.C:
			if OnDemand && !Config.HasViewer(name) {
				return ErrorStreamExitNoViewer
			}
		case <-keyTest.C:
			return ErrorStreamExitNoVideoOnStream
		case signals := <-RTSPClient.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				Config.coAd(name, RTSPClient.CodecData)
			case rtspv2.SignalStreamRTPStop:
				return ErrorStreamExitRtspDisconnect
			}
		case packetAV := <-RTSPClient.OutgoingPacketQueue:
			if packetAV.IsKeyFrame {
				keyTest.Reset(20 * time.Second)
			}
			if packetAV.IsKeyFrame && !start {
				start = true
			}
			if packetAV.Idx == 0 {
				lineV += packetAV.Duration
				packetAV.Time = lineV
			}
			if packetAV.Idx == 1 {
				lineA += packetAV.Duration
				packetAV.Time = lineA
			}
			if broadCastEnable && start {
				if err = connRTMP.WritePacket(*packetAV); err != nil {
					return err
				}
			}
		}
	}
}
