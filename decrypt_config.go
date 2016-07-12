package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"crypto/cipher"
	"encoding/json"
)

type Config struct  {
	Tg Tg
}

type Tg struct {
	Rcp []string
	Path string
}

func main() {

	k := os.Getenv("ASTCONFIG")

    	hasher := md5.New()
    	hasher.Write([]byte(k))
    	fmt.Println(k)
    	key := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println(key)

	file, err := os.Open("astconf")
   	if err != nil {
     		fmt.Println(err)
     		os.Exit(1)
   	}
	data := make([]byte, 10000)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
        fmt.Println(count)
	ciphertext := string(data[:count])
	plaintext := decrypt(string(ciphertext), key)
	fmt.Println(string(plaintext))
	respByte := bytes.NewReader(plaintext)
	decoder := json.NewDecoder(io.Reader(respByte))
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("CONF START")
	fmt.Println(conf)
	fmt.Println("CONF END")
}

func decrypt(cipherstring string, keystring string) []byte {
	// Byte array of the string
	ciphertext := []byte(cipherstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

//	return string(ciphertext)
	return ciphertext
}
