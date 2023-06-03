# ChromeDecryptor
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Release](https://img.shields.io/github/release/wat4r/ChromeDecryptor)](https://github.com/wat4r/ChromeDecryptor/releases)


## Introduction
ChromeDecryptor is a tool for offline decryption of Chrome/Edge browser passwords and cookies.


## Features
 - Support offline decryption
 - Support multiple ways of decrypting the Master key (password, hash, domain backup key)
 - Support all browsers with Chromium core


## Usage
```sh
chromeDecryptor -h
```
This will display help for the tool. Here are all the switches it supports.

```yaml
Usage:
  chromeDecryptor.exe [flags]

Flags:
INPUT:
   -bf, -browser-folder string  browser folder, include Local State, Login Data, Cookies

DPAPI:
   -mkf, -master-key-files string  master key files folder
   -s, -sid string                 user sid
   -p, -password string            password
   -hash string                    sha1 hash or ntlm hash
   -pvk string                     domain backup key (.pvk)

DECRYPT:
   -dp, -decrypt-password  decrypt password(default) (default true)
   -dc, -decrypt-cookies   decrypt cookies

OUTPUT:
   -o, -output string  output path
```

## Example
### Directory structure
```yaml
+---chromeFolder
|       Cookies
|       Local State
|       Login Data
|
+---masterKeyFolder
        +---S-1-5-21-1099483827-325504281-218701502-1001
        |       cd8b97f2-c875-4e9e-85d8-bb8d73954e50
        |       Preferred
```


### Decrypt password
```sh
chromeDecryptor -bf ./chromeFolder -mkf ./masterKeyFolder -p password123 -o pwd.txt
```

### Decrypt cookies
```sh
chromeDecryptor -bf ./chromeFolder -mkf ./masterKeyFolder -p password123 -dc -o cookies.txt
```


## License
This project is licensed under the [Apache 2.0 license](LICENSE).


## Contact
If you have any issues or feature requests, please contact us. PR is welcomed.
 - https://github.com/wat4r/ChromeDecryptor/issues

