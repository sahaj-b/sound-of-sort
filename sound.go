package main

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

const (
	sampleRate      = beep.SampleRate(24100)
	soundDuration   = 25 * time.Millisecond
	minPitch        = 200.0
	maxPitch        = 1200.0
	minBeepInterval = 5 * time.Millisecond
)

func initAudio() error {
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/20))
	if err != nil {
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}
	return nil
}

func playSine(pitch float64, duration time.Duration, volume *atomic.Uint64) {
	numSamples := sampleRate.N(duration)

	position := 0
	currentVol := math.Float64frombits(volume.Load())

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
			value := float64(currentVol) * envelopeMultiplier * math.Sin(2*math.Pi*pitch*phase/float64(sampleRate))
			samples[i][0] = value
			samples[i][1] = value

			position++
		}
		return count, true
	})

	speaker.Play(streamer)
}
