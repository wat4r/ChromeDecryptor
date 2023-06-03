package cookies

import (
	"ChromeDecryptor/pkg/decrypt"
	"ChromeDecryptor/pkg/output"
	"database/sql"
	"fmt"
)

var outputData = `----------------------------
host: %s
path: %s
name: %s
value: %s
expires: %s
lastAccess: %s
lastUpdate: %s
`

func Decrypt(cookiesPath string, key []byte, outputPath string) {
	db, err := sql.Open("sqlite3", cookiesPath)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT host_key, name, encrypted_value, path, expires_utc, last_access_utc, last_update_utc FROM cookies")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var hostKey, name, path string
		var encryptedValue []byte
		var expiresUtc, lastAccessUtc, lastUpdateUtc int64
		rows.Scan(&hostKey, &name, &encryptedValue, &path, &expiresUtc, &lastAccessUtc, &lastUpdateUtc)
		out := fmt.Sprintf(outputData, hostKey, path, name, decrypt.DecryptChrome(key, encryptedValue),
			decrypt.ChromeTimestamp(expiresUtc).Format("2006-01-02 15:04:05"),
			decrypt.ChromeTimestamp(lastAccessUtc).Format("2006-01-02 15:04:05"),
			decrypt.ChromeTimestamp(lastUpdateUtc).Format("2006-01-02 15:04:05"))
		output.Output(outputPath, out)
	}
}
