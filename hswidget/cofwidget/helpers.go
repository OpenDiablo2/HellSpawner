package cofwidget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"
)

// this likely needs to be a method of d2cof.COF
func speedToFPS(speed int) float64 {
	const (
		baseFPS      = 25
		speedDivisor = 256
	)

	fps := baseFPS * (float64(speed) / speedDivisor)
	if fps == 0 {
		fps = baseFPS
	}

	return fps
}

// this should also probably be a method of COF
func calculateDuration(cof *d2cof.COF) float64 {
	const (
		milliseconds = 1000
	)

	frameDelay := milliseconds / speedToFPS(cof.Speed)

	return float64(cof.FramesPerDirection) * frameDelay
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
