package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// TXOutput 交易输出
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock 锁定输出
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address[:])
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey 检查输出是否以指定密钥锁定
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput 创建新交易输出
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// TXOutputs 交易输出数组
type TXOutputs struct {
	Outputs []TXOutput
}

// Serialize 序列化输出数组
func (outs *TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs 反序列化输出数组
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)

	if err != nil {
		log.Panic(err)
	}

	return outputs
}
