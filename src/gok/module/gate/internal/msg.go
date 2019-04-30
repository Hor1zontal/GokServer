package internal

import (
	"github.com/gogo/protobuf/proto"
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/network/protobuf"
	//"t08/network/conf"
	"gok/module/gate/conf"
	"gok/module/game"
	"gok/service/msg/protocol"
)

// 使用 Protobuf 代理消息处理器
//var Processor = NewProcessor(NewXXTeaCrypto(conf.Config.SecretKey))
var Processor = NewProcessor()

func init() {
	//注册protobuf处理消息
	Processor.Register(&protocol.C2GS{})
	Processor.Register(&protocol.GS2C{})
	//消息路由
	Processor.SetRouter(&protocol.C2GS{}, game.ChanRPC)
}

func NewProcessor() *MessageProcessor {
	processor := new(MessageProcessor)
	processor.proxy = protobuf.NewProcessor()
	if conf.Config.SecretKey != "" {
		processor.crypto = NewXXTeaCrypto(conf.Config.SecretKey)
	}
	return processor
}

type Crypto interface {
	Encrypt(data []byte) []byte //加密方法
	Decrypt(data []byte) []byte //解密方法
}

type MessageProcessor struct {
	proxy  *protobuf.Processor
	crypto Crypto
}

func (this *MessageProcessor) Register(msg proto.Message) uint16 {
	return this.proxy.Register(msg)
}

func (this *MessageProcessor) SetHandler(msg proto.Message, msgHandler protobuf.MsgHandler) {
	this.proxy.SetHandler(msg, msgHandler)
}

func (this *MessageProcessor) SetRouter(msg proto.Message, msgRouter *chanrpc.Server) {
	this.proxy.SetRouter(msg, msgRouter)
}

func (this *MessageProcessor) Route(msg interface{}, userData interface{}) error {
	return this.proxy.Route(msg, userData)
}

// must goroutine safe
func (this *MessageProcessor) Unmarshal(data []byte) (interface{}, error) {
	if this.crypto != nil {
		//log.Debug("receive data %v - %v", data, base64.StdEncoding.EncodeToString(data))
		data = this.crypto.Decrypt(data)
		//log.Debug("receive decrypt data %v - %v", data, base64.StdEncoding.EncodeToString(data))
		//temp := this.crypto.Decrypt(data)
		//log.Debug("receive Decrypt data %v - %v", temp, base64.StdEncoding.EncodeToString(temp))
	}
	return this.proxy.Unmarshal(data)
}

// must goroutine safe
func (this *MessageProcessor) Marshal(msg interface{}) ([][]byte, error) {
	data, err := this.proxy.Marshal(msg)
	if len(data) != 1 {
		return data, err
	}
	if this.crypto != nil {
		//log.Debug("response data %v - %v", data[0], base64.StdEncoding.EncodeToString(data[0]))
		data[0] = this.crypto.Encrypt(data[0])
		//log.Debug("response encrypt data %v - %v", data[0], base64.StdEncoding.EncodeToString(data[0]))
	}
	return data, err
}
