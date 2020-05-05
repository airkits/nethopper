package transport

const (
	// HeaderType define JSONHeader,CSPBHeader,SSPBHeader
	HeaderType = iota
	//HeaderTypeWSJSON type WS json
	HeaderTypeWSJSON
	//HeaderTypeWSPB type WS protobuf message
	HeaderTypeWSPB
	//HeaderTypeGRPCPB type GRPC protobuf message
	HeaderTypeGRPCPB
)

//IMessage message interface
type IMessage interface {
	GetID() uint32
	GetCmd() string
	GetMsgType() uint32
	GetSeq() uint32
}

// //NewMessage create new message with headerType
// func NewMessage(headerType int32, codec codec.Codec) *Message {
// 	return &Message{HeaderType: headerType, codec: codec}
// }

// //Message message struct
// type Message struct {
// 	HeaderType int32
// 	Header     interface{}
// 	Body       IBody
// 	codec      codec.Codec
// }

// //CreateHeader create message header
// func (m *Message) CreateHeader() interface{} {
// 	switch m.HeaderType {
// 	case HeaderTypeWSJSON:
// 		{
// 			return &json.Header{}
// 		}
// 	case HeaderTypeWSPB:
// 		{
// 			return &cs.Header{}
// 		}
// 	case HeaderTypeGRPCPB:
// 		{
// 			return &ss.Header{}
// 		}
// 	}
// 	return nil
// }

// //NewHeader create header with params
// func (m *Message) NewHeader(id int32, cmd string, msgType int32) interface{} {
// 	switch m.HeaderType {
// 	case HeaderTypeWSJSON:
// 		{
// 			return &json.Header{ID: id, Cmd: cmd, MsgType: msgType}
// 		}
// 	case HeaderTypeWSPB:
// 		{
// 			return &cs.Header{ID: id, Cmd: cmd, MsgType: msgType}
// 		}
// 	case HeaderTypeGRPCPB:
// 		{
// 			return &ss.Header{ID: id, Cmd: cmd, MsgType: msgType}
// 		}
// 	}
// 	return nil
// }

// //Codec get message codec
// func (m *Message) Codec() codec.Codec {
// 	return m.codec
// }

// //DecodeHeader head decode
// func (m *Message) DecodeHeader(payload []byte) error {

// 	head := m.CreateHeader()
// 	if err := m.codec.Unmarshal(payload, head, nil); err != nil {
// 		return err
// 	}
// 	m.Header = head
// 	return nil
// }

// //Encode encode message
// func (m *Message) Encode() ([]byte, error) {
// 	if m.Header == nil {
// 		return nil, errors.New("message head is null")
// 	}

// 	var err error
// 	if err = m.EncodeBody(); err != nil {
// 		return nil, err
// 	}
// 	var payload []byte
// 	if payload, err = m.codec.Marshal(m.Header, nil); err != nil {
// 		return nil, err
// 	}
// 	return payload, nil
// }

// //encodeWSJSONBody encode body
// func (m *Message) encodeWSJSONBody(payload []byte) error {
// 	(m.Header).(*json.Header).Payload = string(payload)
// 	return nil
// }

// //encodeWSPBBody encode body
// func (m *Message) encodeWSPBBody(payload []byte) error {
// 	(m.Header).(*cs.Header).Payload = payload
// 	return nil
// }

// //encodeGRPCPBBody encode body
// func (m *Message) encodeGRPCPBBody(payload []byte) error {
// 	(m.Header).(*ss.Header).Payload = payload
// 	return nil
// }

// //EncodeBody encode body
// func (m *Message) EncodeBody() error {
// 	if m.Header == nil {
// 		return errors.New("message head is null")
// 	}
// 	var payload []byte
// 	var err error
// 	if payload, err = m.codec.Marshal(m.Body, nil); err != nil {
// 		return err
// 	}
// 	switch m.HeaderType {
// 	case HeaderTypeWSJSON:
// 		{
// 			return m.encodeWSJSONBody(payload)
// 		}
// 	case HeaderTypeWSPB:
// 		{
// 			return m.encodeWSPBBody(payload)
// 		}
// 	case HeaderTypeGRPCPB:
// 		{
// 			return m.encodeGRPCPBBody(payload)
// 		}
// 	default:
// 		{
// 			return errors.New("unknown codec type")
// 		}
// 	}
// }
