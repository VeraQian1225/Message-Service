package main

import (
	"fmt"
	"net/http"
	"github.com/messagebird/go-rest-api"
	"math/big"
	"crypto/rand"
	"log"
	"strconv"
	_ "debug/dwarf"
	_ "github.com/shirou/gopsutil/docker"
)

var info = make(map[string]string)

func sendSMS(writer http.ResponseWriter, request *http.Request) {
	s := sixDigits()
	ss := strconv.Itoa(int(s))
	log.Printf("%06d\n", s)
	Phone := request.URL.Query().Get("Phone")
	Bodymsg := "嗨Hi, Vera! Verification code is " + ss + ". Please input the code in 15mins, otherwise you need to request a new code."
	params := &messagebird.MessageParams{Reference: "MyReference"}
	params.DataCoding="unicode"
	getMsgInfo(Phone, Bodymsg, params)
	fmt.Fprintf(writer, "Hello World, %s!", request.URL.Query().Get("Phone"), params)
}

func getMsgInfo(Phone string, Bodymsg string, params *messagebird.MessageParams){
	client := messagebird.New("")
	message, err := client.NewMessage(
		"1111",
		[]string{Phone},
		Bodymsg,
		params)
	if err != nil{
		fmt.Println(err)
	}
	//info[message.ID] = Bodymsg
	fmt.Println(message.ID)
	fmt.Println(message.Reference)
	fmt.Println(message.Recipients)
	fmt.Println(message)

}

func sendVM(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Vera is calling, %s!", request.URL.Query().Get("Phone"))
	params := &messagebird.VoiceMessageParams{Reference: "MyReference"}
	client := messagebird.New("ViJzmNPRl7BQwTMU7QBoYwhjx")

	message, err := client.NewVoiceMessage(
		[]string{request.URL.Query().Get("Phone")},
		"Hello World",
		params)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(message)
}

func feedback(writer http.ResponseWriter, request *http.Request) {
	Phone := request.URL.Query().Get("recipient")
	status := request.URL.Query().Get("status")
	Id := request.URL.Query().Get("id")
	Bodymsg := info[Id]
	params := &messagebird.MessageParams{Reference: "MyReference"}
	if status != "delivered" {
		writer.WriteHeader(500)
		getMsgInfo(Phone, Bodymsg, params)
	}
}

func main() {
	http.HandleFunc("/sendsms", sendSMS)
	http.HandleFunc("/sendVM", sendVM)
	http.HandleFunc("/processfeedback/script", feedback)
	http.ListenAndServe(":8080", nil)
}

func sixDigits() int64 {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	return n.Int64()
}