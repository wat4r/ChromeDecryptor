package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func Base64Decode(b64Data []byte) []byte {
	rawDecodedText, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		panic(err)
	}
	return rawDecodedText
}

func GuidMasterKeyConvert(guidMasterKey [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		ReverseBytes(guidMasterKey[:4]),
		ReverseBytes(guidMasterKey[4:6]),
		ReverseBytes(guidMasterKey[6:8]),
		guidMasterKey[8:10],
		guidMasterKey[10:])
}

func ReverseBytes(data []byte) []byte {
	for i := 0; i < len(data)/2; i++ {
		j := len(data) - i - 1
		data[i], data[j] = data[j], data[i]
	}
	return data
}

func HexToBytes(hexData string) []byte {
	byteData, _ := hex.DecodeString(hexData)
	return byteData
}
