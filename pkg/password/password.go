package password

import (
	"ChromeDecryptor/pkg/decrypt"
	"ChromeDecryptor/pkg/output"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var outputData = `----------------------------
url: %s
username: %s
password: %s
created: %s
modified: %s
`

func Decrypt(loginDataPath string, key []byte, outputPath string) {
	db, err := sql.Open("sqlite3", loginDataPath)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT origin_url, username_value, password_value, date_created, date_password_modified FROM logins")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var originUrl, userNameValue string
		var passwordValue []byte
		var dataCreated, datePasswordModified int64

		rows.Scan(&originUrl, &userNameValue, &passwordValue, &dataCreated, &datePasswordModified)
		out := fmt.Sprintf(outputData, originUrl, userNameValue,
			decrypt.DecryptChrome(key, passwordValue),
			decrypt.ChromeTimestamp(dataCreated).Format("2006-01-02 15:04:05"),
			decrypt.ChromeTimestamp(datePasswordModified).Format("2006-01-02 15:04:05"))
		output.Output(outputPath, out)
	}
}
