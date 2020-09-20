package hsutil

var (
	lastDeviceScale float64
)

func GetLastDeviceScale() float64 {
	return lastDeviceScale
}

func SetDeviceScale(scale float64) {
	lastDeviceScale = scale
}

func ScaleToDevice(x int) int {
	return int(lastDeviceScale * float64(x))
}

func UnscaleFromDevice(x int) int {
	return int(lastDeviceScale / float64(x))
}
