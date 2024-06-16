package pdu

import (
	"encoding/json"
	"fmt"
)

const (
	TYPE_DATA  = 0
	TYPE_JOIN  = 1
	TYPE_LEAVE = 2
	TYPE_ACK   = 3

	MAX_PDU_SIZE = 1024
)

type PDU struct {
	Mtype uint8  `json:"mtype"`
	Len   uint32 `json:"len"`
	Data  []byte `json:"data"`
}

func MakePduBuffer() []byte {
	return make([]byte, MAX_PDU_SIZE)
}

func NewPDU(mtype uint8, data []byte) *PDU {
	return &PDU{
		Mtype: mtype,
		Len:   uint32(len(data)),
		Data:  data,
	}
}

func (pdu *PDU) GetTypeAsString() string {
	switch pdu.Mtype {
	case TYPE_DATA:
		return "***DATA"
	case TYPE_JOIN:
		return "****JOIN"
	case TYPE_LEAVE:
		return "****LEAVE"
	case TYPE_ACK:
		return "****ACK"
	default:
		return "UNKNOWN"
	}
}

func (pdu *PDU) ToJsonString() string {
	jsonData, err := json.MarshalIndent(pdu, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "{}"
	}

	return string(jsonData)
}

func PduFromBytes(raw []byte) (*PDU, error) {
	pdu := &PDU{}
	err := json.Unmarshal(raw, pdu)
	return pdu, err
}

func (pdu *PDU) PduToBytes() ([]byte, error) {
	return json.Marshal(pdu)
}
