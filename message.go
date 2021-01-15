package delay_server

const (
	StringEncoder = iota + 1
	ByteEncoder
)

type (
	// kafka kv序列化类型
	Encoder int32
	// DelayMessage 延迟消息体
	DelayMessage struct {
		DelayTime int64

		Topic   string
		Brokers []string

		KEncoder    Encoder
		VEncoder    Encoder
		RawMsgKey   []byte
		RawMsgValue []byte
	}
)
