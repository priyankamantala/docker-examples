package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/createFile", CreateFile)
	r.Run()
}

//ResponseData ... to send response from server
var ResponseData struct {
	FileData []byte `json:"fileData"`
	CheckSum string `json:"checksum"`
}

//CreateFile ... function to create file of 1 kb
func CreateFile(c *gin.Context) {
	// Function to generate random string to write into the file
	textToWrite := GenerateRandomString()
	log.Println(textToWrite)
	writeErr := ioutil.WriteFile("/serverdata/randomFile.txt", []byte(textToWrite), 0644)
	if writeErr != nil {
		log.Print("Error while writing data to file", writeErr)
		c.JSON(417, gin.H{"message": "Error while writing to file"})
	}
	// reading file passing the path after writing the random text
	filedata, readErr := ioutil.ReadFile("/serverdata/randomFile.txt")
	if readErr != nil {
		log.Fatal("Error while reading file")
		c.JSON(417, gin.H{"message": "Error while reading file"})
	}
	// Checksum generation of file
	checksum := GenerateCheckSum(filedata)
	ResponseData.CheckSum = checksum
	ResponseData.FileData = filedata
	c.JSON(200, ResponseData)
}

//GenerateRandomString will generate random text to write to the file
func GenerateRandomString() string {
	var letterBytes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	bytesToWrite := make([]rune, 1024)
	for i := range bytesToWrite {
		bytesToWrite[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(bytesToWrite)
}

//GenerateCheckSum ... function to generate checksum of file data
func GenerateCheckSum(filedata []byte) string {
	hash := md5.New()
	hashInBytes := hash.Sum(filedata)
	md5HashString := hex.EncodeToString(hashInBytes)
	return md5HashString
}
