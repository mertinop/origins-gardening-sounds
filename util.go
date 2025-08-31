package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

const (
	AssetServerBaseURL = "https://github.com/mertinop/origins-fishing-sounds/raw/refs/heads/main/assets/"
)

func playSound(filename string) {
	resp, err := http.Get(AssetServerBaseURL + filename)
	if err != nil {
		log.Println("Error downloading sound file:", err)
		return
	}
	defer resp.Body.Close()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "sound-*.wav")
	if err != nil {
		log.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Copy the downloaded data to the temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		log.Println("Error writing to temporary file:", err)
		return
	}
	tmpFile.Close()

	// Open the temporary file for playing
	f, err := os.Open(tmpFile.Name())
	log.Printf("Playing sound: %s", tmpFile.Name())
	if err != nil {
		log.Println("Error opening temporary sound file:", err)
		return
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Println("Error decoding WAV file:", err)
		return
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Println("Error initializing speaker:", err)
		return
	}

	done := make(chan bool)
	ctrl := &beep.Ctrl{Streamer: beep.Seq(streamer, beep.Callback(func() {
		done <- true
	}))}

	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   1,
		Silent:   false,
	}
	speaker.Play(volume)

	<-done
}
