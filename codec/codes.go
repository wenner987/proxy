package codec

const magicNum uint64 = 0xFFFFFAAAAAAFFFFF
const headerLen = 16

// magic: magic number proto tag
// Version: proto version
// ProtoLen: proto length
// From: Msg origin
type ProtoHeader struct {
	Magic uint64
	Version uint16
	ProtoLen uint32
	From uint16
	Dst uint16
}

type Proto struct {
	Header ProtoHeader
	Context []byte
}