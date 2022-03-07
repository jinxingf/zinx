package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// data len(4 byte) + id(4byte)  = 8 byte
	return 8
}
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil

}
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// read binary from io by IOReader
	dataBuf := bytes.NewReader(binaryData)

	// unpack header(len and id)

	msg := &Message{}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// check if dataLen is bigger than maxLen or not
	if utils.GlobalConf.MaxPackageSize > 0 && msg.DataLen > utils.GlobalConf.MaxPackageSize {
		return nil, errors.New("too large message package")
	}

	return msg, nil
}
