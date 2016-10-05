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
    	hasher := md5.New()
	k := os.Getenv("ASTCONFIG")
    	hasher.Write([]byte(k))
    	fmt.Println(k)
    	key := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println(key)
	file, err := os.Open("your_config.json")
   	if err != nil {
     		fmt.Println(err)
     		os.Exit(1)
   	}
	data := make([]byte, 20000)
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
	plaintext := []byte(plainstring)
	key := []byte(keystring)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	fmt.Println(block.BlockSize())
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return string(ciphertext)
}
