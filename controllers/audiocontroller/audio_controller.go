package audiocontroller

import (
	"net/http"
	"streaming/util"
)

const audioDir = "media/audio"

// ffmpeg command
// ffmpeg -i BachGavotteShort.mp3 -c:a libmp3lame -b:a 128k -map 0:0 -f segment -segment_time 10 -segment_list outputlist.m3u8 -segment_format mpegts output%03d.ts

// HandleAudio handles and serves http requests for /audio
func HandleAudio(res http.ResponseWriter, req *http.Request) http.Handler {
	var head string
	head, req.URL.Path = util.ShiftPath(req.URL.Path)
	var handler http.Handler

	switch head {
	// sample url for testing
	// http://localhost:8080/audio/playback/bachgavotteshort/outputlist.m3u8
	case "playback":
		handler = handlePlayback(res, req)
	}
	return handler
}

func handlePlayback(res http.ResponseWriter, req *http.Request) http.Handler {
	return http.FileServer(http.Dir(audioDir))
}
