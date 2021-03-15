package hsstate

type State interface {
	Dispose()
	Encode() []byte
	Decode([]byte)
}
