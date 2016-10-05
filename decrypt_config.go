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

}

func main() {
    	hasher := md5.New()
	k := os.Getenv("ASTCONFIG")
    	hasher.Write([]byte(k))
    	fmt.Println(k)
    	key := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println(key)

	file, err := os.Open("conf")
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
	ciphertext := string(data[:count])
	plaintext := decrypt(string(ciphertext), key)
	fmt.Println(string(plaintext))
	buf := bytes.NewBufferString("")
	buf.Write([]byte(plaintext))
	dst := "your_config.json"
	f, _ := os.OpenFile(dst, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)
	f.Write(buf.Bytes())
	defer f.Close()

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
	ciphertext := []byte(cipherstring)
	key := []byte(keystring)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}
