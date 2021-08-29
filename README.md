# RTSPtoRTMP

RTSP Stream to YouTube or Twitch

full native! not use ffmpeg or gstreamer

 
### Download Source

1. Download source
   ```bash 
   $ git clone https://github.com/deepch/RTSPtoRTMP 
   ```
3. CD to Directory
   ```bash
    $ cd RTSPtoRTMP/
   ```
4. Test Run
   ```bash
    $ GO111MODULE=on go run *.go
   ```
5. Edit config.json


## Configuration

### Edit file config.json

format:

```bash
{
  "streams": {
    "H264_AAC": {
      "on_demand": false,
      "url": "YOU_RTSP_CAMERA_URL",
      "broadcast": {
        "enable": true,
        "url": "rtmp://a.rtmp.youtube.com/live2/YOU_YOUTUBE_KEY"
      }
    }
  }
}

```

## Limitations

Video Codecs Supported: H264 all profiles

Audio Codecs Supported: AAC or NONE

## Test

CPU usage 0.2% one core cpu intel core i7 / stream

## Team

Deepch - https://github.com/deepch streaming developer

Dmitry - https://github.com/vdalex25 web developer

## Other Example

Examples of working with video on golang

- [RTSPtoWeb](https://github.com/deepch/RTSPtoWeb)
- [RTSPtoWebRTC](https://github.com/deepch/RTSPtoWebRTC)
- [RTSPtoWSMP4f](https://github.com/deepch/RTSPtoWSMP4f)
- [RTSPtoImage](https://github.com/deepch/RTSPtoImage)
- [RTSPtoHLS](https://github.com/deepch/RTSPtoHLS)
- [RTSPtoHLSLL](https://github.com/deepch/RTSPtoHLSLL)
- [RTSPtoRTMP](https://github.com/deepch/RTSPtoRTMP)
- 
[![paypal.me/AndreySemochkin](https://ionicabizau.github.io/badges/paypal.svg)](https://www.paypal.me/AndreySemochkin) - You can make one-time donations via PayPal. I'll probably buy a ~~coffee~~ tea. :tea: