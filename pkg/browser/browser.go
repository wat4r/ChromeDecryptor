package browser

import (
	"fmt"
	"github.com/wat4r/dpapitk/blob"
	"os"
	"regexp"
)

func ReadEncryptedKey(localStateContext string) string {
	var re = regexp.MustCompile(`(?m)"encrypted_key":"([\w+/=.]+)"`)
	matches := re.FindAllStringSubmatch(localStateContext, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		return matches[0][1]
	}
	return ""
}

func DecryptEncryptedKey(masterKey, encryptedKey []byte) []byte {
	data, err := blob.DecryptWithMasterKey(encryptedKey[5:], masterKey, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	return data
}
