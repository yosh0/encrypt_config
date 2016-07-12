package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"crypto/cipher"
)

func main() {

	k := os.Getenv("ASTCONFIG")

    	hasher := md5.New()
    	hasher.Write([]byte(k))
    	fmt.Println(k)
    	key := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println(key)

	file, err := os.Open("your_config.json")
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
	message := string(data[:count])
	ciphertext := encrypt(message, key)
	writeToFile(ciphertext, "conf")
}

func writeToFile(data, file string) {
	ioutil.WriteFile(file, []byte(data), 400)
}

func readFromFile(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	return data, err
}

func encrypt(plainstring, keystring string) string {
	// Byte array of the string
	plaintext := []byte(plainstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	fmt.Println(block.BlockSize())

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return string(ciphertext)
}
