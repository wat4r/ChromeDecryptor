package output

import (
	"ChromeDecryptor/pkg/utils"
	"fmt"
)

func Output(outputPath string, outputString string) {
	if outputPath != "" {
		utils.WriteOutput(outputPath, []byte(outputString))
	}
	fmt.Print(outputString)
}
