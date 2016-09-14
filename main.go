package main

import (
	json "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
)

var counter int = 1

func printlog(msg string, increase bool) {
	if increase {
		counter++
		log.Print("------------------------------------")
	}
	log.Print(strconv.Itoa(counter) + ". - " + msg)
}

func main() {
	if len(os.Args) > 1 || StartTesting {
		printlog("ENV", false)
		for _, ts := range os.Environ() {
			printlog(ts, false)
			printlog(fmt.Sprintf("%v: %v", ts, os.Getenv(ts)), false)
		}

		printlog(fmt.Sprintf("%v: %v", "IP_ADDRESS", GatewayHost), false)
		printlog(fmt.Sprintf("%v: %v", "PORT", GatewayPort), false)
		printlog("SMOKING TESTING SERVER", false)
		smokeTest(10, 3)
		printlog("STARTING TESTS", false)
		startTests()
	} else {
		for true {

		}
	}

}

func smokeTest(attempts int, sleep int64) {
	passed := false
	for total := 0; total < attempts; total++ {
		_, err := sendRequest("GET", "/", defaultHeaders, "")
		if err == nil {
			passed = true
			break
		}
		time.Sleep(time.Second * time.Duration(sleep))
	}
	if !passed {
		panic(fmt.Sprintf("DIDNT PASS THE SMOKE TEST AFTER %d attempts", attempts))
	}

}

func startTests() {
	//Running the tests now

	//Signup
	printlog("Attempting signup", true)
	username, err := signup()

	if err != nil {
		printlog("Signup failed: "+err.Error(), false)
		panic(err)
	}
	printlog("Signup success: "+username, false)

	//Login
	printlog("Attempting login", true)
	username, err = login(username)
	if err != nil {
		printlog("Login failed: "+err.Error(), false)
		panic(err)
	}

	printlog("Login success: "+username, false)

	//List tasks

	//Creates tasks

	//Completes tasks
}

func sendRequest(method string, endpoint string, headers map[string]string, data string) (*http.Response, error) {
	client := &http.Client{}

	fullURL := GatewayHost + GatewayPort + endpoint

	req, err := http.NewRequest(method, fullURL, strings.NewReader(data))
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

const password string = "123456"

var defaultHeaders map[string]string = map[string]string{"Content-type": "application/json"}

var token *Token

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupUser struct {
	Name string `json:"name"`
	LoginUser
}

type Token struct {
	Token     string    `json:"token"`
	UserId    string    `json:"userId"`
	Ttl       int       `json:"ttl"`
	CreatedOn time.Time `json:"createdOn"`
}

func signup() (string, error) {
	username := randomdata.Adjective() + randomdata.FirstName(randomdata.Male)
	user := SignupUser{
		Name: randomdata.FullName(randomdata.Male),
		LoginUser: LoginUser{
			Username: username,
			Password: password,
		},
	}

	data, err := json.Marshal(user)
	if err != nil {
		return username, err
	}
	printlog("Will send "+string(data), false)

	res, err := sendRequest("POST", "/users/signup", defaultHeaders, string(data))
	if err != nil {
		return username, err
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&token)
	if err != nil {
		return username, err
	}

	return username, nil
}

func login(username string) (string, error) {
	user := LoginUser{
		Username: username,
		Password: password,
	}
	data, err := json.Marshal(user)
	if err != nil {
		return username, err
	}
	res, err := sendRequest("POST", "/users/login", defaultHeaders, string(data))
	if err != nil {
		return username, err
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&token)
	if err != nil {
		return username, err
	}

	return username, nil

}
