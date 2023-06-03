package main

import (
	"ChromeDecryptor/pkg/browser"
	"ChromeDecryptor/pkg/cookies"
	"ChromeDecryptor/pkg/output"
	"ChromeDecryptor/pkg/password"
	"ChromeDecryptor/pkg/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/goflags"
	"github.com/wat4r/dpapitk/masterkey"
)

type Options struct {
	// Browser
	BrowserFolder string

	// Master key
	MasterKeyFiles  string
	Sid             string
	Password        string
	Hash            string // sha1 or ntlm
	DomainBackupKey string

	// Other
	DecryptPassword bool
	DecryptCookies  bool
	OutputPath      string
}

func main() {
	var opt *Options
	var loginDataPath, cookiesPath string
	var encryptedKey []byte

	logoPrint()
	opt = parseParameter()

	localStatePath := utils.FindFile(opt.BrowserFolder, "Local State")

	if localStatePath == "" {
		fmt.Println("Local State not found!")
		return
	}

	//	encrypted key
	localStateContext := string(utils.ReadFile(localStatePath))
	encryptedKeyBase64 := browser.ReadEncryptedKey(localStateContext)
	encryptedKey = utils.Base64Decode([]byte(encryptedKeyBase64))
	masterKey := decryptMasterKey(opt, encryptedKey)
	key := browser.DecryptEncryptedKey(masterKey, encryptedKey)
	output.Output(opt.OutputPath, fmt.Sprintf("master key: %x\n", key))

	if opt.DecryptPassword && !opt.DecryptCookies {
		loginDataPath = utils.FindFile(opt.BrowserFolder, "Login Data")
		password.Decrypt(loginDataPath, key, opt.OutputPath)
	}

	if opt.DecryptCookies {
		cookiesPath = utils.FindFile(opt.BrowserFolder, "Cookies")
		cookies.Decrypt(cookiesPath, key, opt.OutputPath)
	}

	fmt.Println("\nDone!")
}

func parseParameter() *Options {
	opt := &Options{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("Chrome and Edge decrypt.")

	flagSet.CreateGroup("INPUT", "INPUT",
		flagSet.StringVarP(&opt.BrowserFolder, "browser-folder", "bf", "", "browser folder, include Local State, Login Data, Cookies"),
	)

	flagSet.CreateGroup("DPAPI", "DPAPI",
		flagSet.StringVarP(&opt.MasterKeyFiles, "master-key-files", "mkf", "", "master key files folder"),
		flagSet.StringVarP(&opt.Sid, "sid", "s", "", "user sid"),
		flagSet.StringVarP(&opt.Password, "password", "p", "", "password"),
		flagSet.StringVar(&opt.Hash, "hash", "", "sha1 hash or ntlm hash"),
		flagSet.StringVar(&opt.DomainBackupKey, "pvk", "", "domain backup key (.pvk)"),
	)

	flagSet.CreateGroup("DECRYPT", "DECRYPT",
		flagSet.BoolVarP(&opt.DecryptPassword, "decrypt-password", "dp", true, "decrypt password(default)"),
		flagSet.BoolVarP(&opt.DecryptCookies, "decrypt-cookies", "dc", false, "decrypt cookies"),
	)

	flagSet.CreateGroup("OUTPUT", "OUTPUT",
		flagSet.StringVarP(&opt.OutputPath, "output", "o", "", "output path"),
	)

	if err := flagSet.Parse(); err != nil {
		fmt.Printf("Could not parse flags: %s\n", err)
		os.Exit(0)
	}

	return opt
}

func logoPrint() {
	logoHex := "5f5f5f5f5f5f2020202020202020202020202020202020202020202020202020205f20202020202020202020205f200a7c20205f20205c2020202020204368726f6d652f4564676520202020202020207c207c2020202020202020207c207c0a7c207c207c207c5f5f5f20205f5f5f205f205f5f205f2020205f205f205f5f207c207c5f205f5f5f20205f5f7c207c0a7c207c207c202f205f205c2f205f5f7c20275f5f7c207c207c207c20275f205c7c205f5f2f205f205c2f205f60207c0a7c207c2f202f20205f5f2f20285f5f7c207c20207c207c5f7c207c207c5f29207c207c7c20205f5f2f20285f7c207c0a7c5f5f5f2f205c5f5f5f7c5c5f5f5f7c5f7c2020205c5f5f2c207c202e5f5f2f205c5f5f5c5f5f5f7c5c5f5f2c5f7c0a202020202020202020202020202020202020202020205f5f2f207c207c2020202020202020202020202020202020200a2020202020202020202020202020202020202020207c5f5f5f2f7c5f7c2020202020202076312e302e30202020202020"
	logo := utils.HexToBytes(logoHex)
	fmt.Printf("%s\n\n", logo)
}

func decryptMasterKey(opt *Options, encryptedKey []byte) []byte {
	// Master key file
	var guidMasterKey [16]byte

	err := binary.Read(bytes.NewReader(encryptedKey[29:]), binary.LittleEndian, &guidMasterKey)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	masterKeyFileName := utils.GuidMasterKeyConvert(guidMasterKey)
	output.Output(opt.OutputPath, fmt.Sprintf("master key file name: %s\n", masterKeyFileName))

	masterKeyFilePath := utils.FindFile(opt.MasterKeyFiles, masterKeyFileName)
	if masterKeyFilePath == "" {
		fmt.Printf("master key file `%s` not found!", masterKeyFileName)
		os.Exit(0)
	}

	if opt.Sid == "" {
		absPath, _ := filepath.Abs(masterKeyFilePath)
		folderPath := filepath.Dir(absPath)
		masterKeyFolder := filepath.Base(folderPath)
		if strings.HasPrefix(masterKeyFolder, "S-1-5-21") {
			opt.Sid = masterKeyFolder
		}
	}

	if opt.Sid == "" {
		fmt.Println("user sid not found!")
		os.Exit(0)
	}

	output.Output(opt.OutputPath, fmt.Sprintf("sid: %s\n", opt.Sid))

	masterKeyFileBytes := utils.ReadFile(masterKeyFilePath)
	masterKeyFile := masterkey.InitMasterKeyFile(masterKeyFileBytes)

	if opt.Password != "" {
		masterKeyFile.DecryptWithPassword(opt.Sid, opt.Password)
	} else if opt.Hash != "" {
		masterKeyFile.DecryptWithHash(opt.Sid, opt.Hash)
	} else if opt.DomainBackupKey != "" {
		domainBackupKey := utils.ReadFile(opt.DomainBackupKey)
		masterKeyFile.DecryptWithPvk(domainBackupKey)
	} else {
		fmt.Println("password, hash, domain backup key not found!")
		os.Exit(0)
	}

	if !masterKeyFile.Decrypted {
		fmt.Println("Master key decrypt failed!")
		os.Exit(0)
	}
	return masterKeyFile.Key
}
