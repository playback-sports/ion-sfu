package buffer

import (
	"math"
)

// Abs returns the absolute value of x.
func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

type resolution struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

func (r resolution) getWidthDiff(w int64) int64 {
	return Abs(r.Width - w)
}

func (r resolution) getHeightDiff(h int64) int64 {
	return Abs(r.Height - h)
}

type targetBitrate struct {
	Resolution  resolution
	LowRate     uint64
	MidRate     uint64
	MidHighRate uint64
	HighRate    uint64
}

const (
	RESOLUTION_144p = iota
	RESOLUTION_240p
	RESOLUTION_360p
	RESOLUTION_480p
	RESOLUTION_540p
	RESOLUTION_720p
	RESOLUTION_1080p
	RESOLUTION_1440p
	RESOLUTION_2160p
)

var targetBitrates = map[int]targetBitrate{
	RESOLUTION_144p: {
		Resolution: resolution{Width: 176, Height: 144},
		LowRate:    80 * 1000,  // 80 Kbps
		MidRate:    100 * 1000, // 100 Kbps
		HighRate:   120 * 1000, // 120 Kbps
	},
	RESOLUTION_240p: {
		Resolution: resolution{Width: 426, Height: 240},
		LowRate:    250 * 1000, // 250 Kbps
		MidRate:    400 * 1000, // 400 Kbps
		HighRate:   700 * 1000, // 700 Kbps
	},
	RESOLUTION_360p: {
		Resolution: resolution{Width: 640, Height: 360},
		LowRate:    500 * 1000,  // 500 Kbps
		MidRate:    800 * 1000,  // 800 Kbps
		HighRate:   1400 * 1000, // 1.4 Mbps
	},
	RESOLUTION_480p: {
		Resolution: resolution{Width: 854, Height: 480},
		LowRate:    750 * 1000,  // 750 Kbps
		MidRate:    1200 * 1000, // 1.2 Mbps
		HighRate:   2100 * 1000, // 2.1 Mbps
	},
	RESOLUTION_540p: {
		Resolution:  resolution{Width: 960, Height: 540},
		LowRate:     1125 * 1000, // 1.125 Mbps
		MidRate:     1800 * 1000, // 1.8 Mbps
		MidHighRate: 2700 * 1000, // 1.8 Mbps
		HighRate:    3150 * 1000, // 3.15 Mbps
	},
	RESOLUTION_720p: {
		Resolution:  resolution{Width: 1280, Height: 720},
		LowRate:     1500 * 1000, // 1.5 Mbps
		MidRate:     2400 * 1000, // 2.4 Mbps
		MidHighRate: 3800 * 1000, // 3.8 Mbps
		HighRate:    4200 * 1000, // 4.2 Mbps
	},
	RESOLUTION_1080p: {
		Resolution:  resolution{Width: 1920, Height: 1080},
		LowRate:     3000 * 1000, // 3 Mbps
		MidRate:     4800 * 1000, // 4.8 Mbps
		MidHighRate: 6200 * 1000, // 6.2 Mbps
		HighRate:    8400 * 1000, // 8.4 Mbps
	},
	RESOLUTION_1440p: {
		Resolution: resolution{Width: 2560, Height: 1440},
		LowRate:    6000 * 1000,  // 6 Mbps
		MidRate:    10400 * 1000, // 10.4 Mbps
		HighRate:   18200 * 1000, // 18.2 Mbps
	},
	RESOLUTION_2160p: {
		Resolution: resolution{Width: 3840, Height: 2160},
		LowRate:    10000 * 1000, // 10 Mbps
		MidRate:    16000 * 1000, // 16 Mbps
		HighRate:   28000 * 1000, // 28 Mbps
	},
}

func getMinIndex(values []int64) (i int) {
	m := int64(math.MaxInt64)
	for curr, v := range values {
		if v < m {
			i = curr
			m = v
		}
	}
	return i
}

func getTargetBitrateForResolution(width, height int64) targetBitrate {
	if width == 0 && height == 0 {
		return targetBitrates[RESOLUTION_540p]
	}

	wDiffs := make([]int64, 0, len(targetBitrates))
	hDiffs := make([]int64, 0, len(targetBitrates))

	// Swap dimensions for vertical video
	if width < height {
		width, height = height, width
	}

	for i := 0; i < len(targetBitrates); i++ {
		wDiffs = append(wDiffs, targetBitrates[i].Resolution.getWidthDiff(width))
		hDiffs = append(hDiffs, targetBitrates[i].Resolution.getHeightDiff(height))
	}

	closestResolutionWidth := getMinIndex(wDiffs)
	closestResolutionHeight := getMinIndex(hDiffs)

	targetBitrateResolution := closestResolutionHeight
	if closestResolutionWidth < closestResolutionHeight {
		targetBitrateResolution = closestResolutionWidth
	}

	return targetBitrates[targetBitrateResolution]
}
