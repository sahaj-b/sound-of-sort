package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

const (
	sampleRate    = beep.SampleRate(24100)
	volume        = 0.1
	soundDuration = 30 * time.Millisecond
	minPitch      = 200.0
	maxPitch      = 1200.0
)

var (
	minVal = 1
	maxVal = 100
)

func setArrBounds(min, max int) {
	minVal = min
	maxVal = max
}

func initAudio() {
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/20))
	if err != nil {
		fmt.Println("Error initializing speaker:", err)
		return
	}
}

func playSine(pitch float64, duration time.Duration) {
	numSamples := sampleRate.N(duration)

	position := 0

	streamer := beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		if position >= numSamples {
			return 0, false
		}

		// how many samples to generate in this batch
		count := len(samples)
		if position+count > numSamples {
			count = numSamples - position
		}

		for i := 0; i < count; i++ {
			envelopeMultiplier := 1.0
			progress := float64(position) / float64(numSamples)

			const fadeInDuration = 0.2
			if progress < fadeInDuration {
				envelopeMultiplier = progress / fadeInDuration
			}

			const fadeOutDuration = 0.2
			if progress > (1.0 - fadeOutDuration) {
				envelopeMultiplier = (1.0 - progress) / fadeOutDuration
			}

			phase := float64(position)
			value := volume * envelopeMultiplier * math.Sin(2*math.Pi*pitch*phase/float64(sampleRate))
			samples[i][0] = value
			samples[i][1] = value

			position++
		}
		return count, true
	})

	// Just play the damn sound. No waiting.
	speaker.Play(streamer)
}

func playBeepArr(val int) {
	// Map the array value to the pitch, not the index. This sounds better.
	pitch := minPitch + (maxPitch-minPitch)*float64(val-minVal)/float64(maxVal-minVal)
	playSine(pitch, soundDuration)
}
