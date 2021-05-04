package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cowin struct {
	Centers []struct {
		CenterID     int    `json:"center_id"`
		Name         string `json:"name"`
		Address      string `json:"address"`
		StateName    string `json:"state_name"`
		DistrictName string `json:"district_name"`
		BlockName    string `json:"block_name"`
		Pincode      int    `json:"pincode"`
		Lat          int    `json:"lat"`
		Long         int    `json:"long"`
		From         string `json:"from"`
		To           string `json:"to"`
		FeeType      string `json:"fee_type"`
		Sessions     []struct {
			SessionID         string   `json:"session_id"`
			Date              string   `json:"date"`
			AvailableCapacity int      `json:"available_capacity"`
			MinAgeLimit       int      `json:"min_age_limit"`
			Vaccine           string   `json:"vaccine"`
			Slots             []string `json:"slots"`
		} `json:"sessions"`
	} `json:"centers"`
}

func main() {

	androidToken := os.Args[1]
	authToken := os.Args[2]

	client := &http.Client{}
	pincode := "382225"   //382225
	minimumAgeLimit := 45 // 18 or 45

	loc, _ := time.LoadLocation("Asia/Kolkata")

	//set timezone,
	now := time.Now().In(loc)

	date := now.Format("02-01-2006")
	fmt.Printf("Calling API... with date : %s, and pincode %s", date, pincode)
	cowinURL := fmt.Sprintf("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByPin?pincode=%s&date=%s", pincode, date)
	req, err := http.NewRequest("GET", cowinURL, nil)
	//req, err := http.NewRequest("GET", "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?district_id=154&date=03-05-2021", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject Cowin
	json.Unmarshal(bodyBytes, &responseObject)
	

	var b bytes.Buffer
	flag := false
	for _, center := range responseObject.Centers {
		for _, session := range center.Sessions {
			if session.AvailableCapacity > 0 && session.MinAgeLimit == minimumAgeLimit {
				b.WriteString(fmt.Sprintln("Name :" + center.Name))
				b.WriteString(fmt.Sprintln("Address : " + center.Address))
				b.WriteString(fmt.Sprintln("Pincode : " + strconv.Itoa(center.Pincode)))
				b.WriteString(fmt.Sprintln("Date : " + session.Date))
				b.WriteString(fmt.Sprintln("Available Capicity : " + strconv.Itoa(session.AvailableCapacity)))
				flag = true
				goto exit_loop

			}
		}

	}
exit_loop:
	if flag {
		fmt.Println("final : " + b.String())
		sendPush(b.String(), androidToken, authToken)
	} else {
		fmt.Println("No Available Capicity in your area!")
	}

}

func sendPush(detail, androidToken, authToken string) {

	url := "https://fcm.googleapis.com/fcm/send"

	payload := strings.NewReader("{\n    \"to\" : \"" + androidToken + "\",\n    \"notification\" : {\n      \"body\" : \"" + detail + "\",\n      \"title\" : \"Vaccination\"\n      }\n}")

	req, error := http.NewRequest("POST", url, payload)

	if error != nil {
		fmt.Println("Error in creating request", error.Error())
		return
	}

	req.Header.Add("authorization", "key="+authToken)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, doError := http.DefaultClient.Do(req)
	if doError != nil {
		fmt.Println("Error :", doError.Error())
		return
	}

	body, readError := ioutil.ReadAll(res.Body)
	if readError != nil {
		fmt.Println("Error :", readError.Error())
		return
	}
	defer res.Body.Close()
	fmt.Println(string(body))

}
