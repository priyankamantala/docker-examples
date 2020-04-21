package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Response ... response struct to bind the response from server
type Response struct {
	FileData     []byte `json:"fileData"`
	FileCheckSum string `json:"fileCheckSum"`
}

func main() {
	// get server host and port no
	serverHost := os.Getenv("ServerHost")
	serverPortNo := os.Getenv("serverPortNo")

	// mapping the variables to url to connect to client
	urlPath := "/createFile"
	serverURL := "http://" + serverHost + ":" + serverPortNo + urlPath

	// New client and send request to the server
	client := http.Client{}
	request, reqErr := http.NewRequest("GET", serverURL, nil)
	if reqErr != nil {
		log.Fatal("Error while requesting data from server", reqErr)
	}
	response, respErr := client.Do(request)
	if respErr != nil {
		log.Fatal("Error while recieving reqponse from server", respErr)
	}
	responseBody, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal("Error while getting response body", readErr)
	}
	data := Response{}
	unmarshalErr := json.Unmarshal(responseBody, data)
	if unmarshalErr != nil {
		log.Fatal("Error while unmarshalling the response body from server")
	}
	verification := VerifyCheckSum(data)
	if verification {
		writeErr := ioutil.WriteFile("/clientdata/randomFile.txt", data.FileData, os.ModePerm)
		if writeErr != nil {
			log.Fatal("Error while writing the file recieved from server")
		}
	}
	http.ListenAndServe(":8090", nil)
}

//VerifyCheckSum ... Function to verify the file checksum
func VerifyCheckSum(data Response) bool {

	hash := md5.New()
	hashInBytes := hash.Sum(filedata)
	md5HashString := hex.EncodeToString(hashInBytes)
	if strings.EqualFold(md5HashString, data.FileCheckSum) {
		return true
	}

}
