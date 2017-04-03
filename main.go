package main

import (
	json "encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
)

var counter int = 1

func printlog(msg string, increase bool) {
	if increase {
		counter++
		println("------------------------------------")
	}
	println(counter, "-", msg)
}

func getEnvName(name, proto, port string) string {
	return fmt.Sprintf("%s_PORT_%s_%s", name, port, proto)
}

var FinalURL string

func main() {

	all := os.Args
	println("args")
	println(all)
	name := "GATEWAY"
	proto := "HTTP"
	port := "80"
	if len(all) > 1 {
		name = strings.ToUpper(all[1])
	}
	if len(all) > 2 {
		proto = strings.ToUpper(all[2])
	}
	if len(all) > 3 {
		port = strings.ToUpper(all[3])
	}
	shouldEx := false
	if len(all) > 4 {
		shouldEx = true
	}

	printlog("ENV", false)
	for _, ts := range os.Environ() {
		printlog(ts, false)
	}

	envName := getEnvName(name, proto, port)
	println(" will use this one here: ", envName)

	if proto != "HTTP" {
		println("Changing the protocol because this is a http server...")
		FinalURL = "http://" + os.Getenv(envName+"_ADDR")
	} else {
		FinalURL = os.Getenv(envName)
		if port == "80" {
			FinalURL = strings.Replace(FinalURL, ":80", "", -1)
		}
	}

	// printlog(fmt.Sprintf("%v: %v", "IP_ADDRESS", GatewayHost), false)
	// printlog(fmt.Sprintf("%v: %v", "PORT", GatewayPort), false)
	// printlog(fmt.Sprintf("%v: %v", "ENDPOINT", Endpoint), false)

	printlog(fmt.Sprintf("%v: %v", "FinalURL", FinalURL), false)
	printlog("SMOKING TESTING SERVER", false)
	smokeTest(10, 5)
	printlog("STARTING TESTS", false)
	if shouldEx {
		println("will quit by command", shouldEx)
		os.Exit(0)
	}
	startTests()
}

func smokeTest(attempts int, sleep int64) {
	passed := false
	for total := 0; total < attempts; total++ {
		_, err := sendRequest("GET", "/health", defaultHeaders, "")
		if err == nil {

			passed = true
			break
		} else {
			fmt.Println("Smoke test error ", err)
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
	printlog("Checking tasks list, should be empty", true)
	tasks, err := getTaskList()
	if err != nil {
		printlog("Task list failed: "+err.Error(), false)
		panic(err)
	}

	if tasks == nil || len(tasks) != 0 {
		printlog(fmt.Sprintf("Invalid list. Should be empty but not nil: %v", tasks), false)
		panic(errors.New("invalid list"))
	}

	//Creates tasks
	printlog("Creating a few tasks...", true)
	newTasks := make([]*Task, 3)
	for i := range newTasks {
		name := randomdata.Street()
		printlog(fmt.Sprintf("Creating task: %s", name), false)
		newTasks[i], err = createTask(name)
		if err != nil {
			printlog(fmt.Sprintf("Creating task err: %v", err), false)
			panic(err)
		}
		if newTasks[i] == nil {
			printlog(fmt.Sprintf("Creating task issue: Task is nil"), false)
			panic(errors.New("task is new"))
		}
		if newTasks[i].Name != name {
			printlog(fmt.Sprintf("Creating task issue: Task name is not the same: %s != %s", name, newTasks[i].Name), false)
			panic(errors.New("task name is not the same"))
		}
	}

	printlog("Fetching the list to compare again...", true)
	tasks, err = getTaskList()
	if err != nil {
		printlog("Task list failed: "+err.Error(), false)
		panic(err)
	}

	if len(tasks) != len(newTasks) {
		printlog(fmt.Sprintf("Creating task issue: Task list size is not the same: %v != %v", tasks, newTasks), false)
		panic(errors.New("list size is not the same"))
	}

	for _, task := range newTasks {
		found := false
		for _, t := range tasks {
			if task.ID == t.ID {
				found = true
				if task.Name != t.Name || task.Done != t.Done {
					printlog(fmt.Sprintf("Tasks are different: %v != %v", task, t), false)
					panic(errors.New("task not the same"))
				}
				// preparing for update
				task.Name = randomdata.Street()
				t.Name = task.Name
				t.Done = true
				task.Done = true
			}
		}
		if !found {
			printlog(fmt.Sprintf("Task not found: %v", task), false)
			panic(errors.New("list not found"))
		}
	}

	// Updating tasks
	printlog("Updating tasks...", true)
	for i, task := range tasks {
		tasks[i], err = updateTask(task.ID, task.Name, task.Done)
		if err != nil {
			printlog("Task update failed: "+err.Error(), false)
			panic(err)
		}
	}

	printlog("Fetching the list to compare again...", true)
	tasks, err = getTaskList()
	if err != nil {
		printlog("Task list failed: "+err.Error(), false)
		panic(err)
	}

	if len(tasks) != len(newTasks) {
		printlog(fmt.Sprintf("Creating task issue: Task list size is not the same: %v != %v", tasks, newTasks), false)
		panic(errors.New("list size is not the same"))
	}

	for _, task := range newTasks {
		found := false
		for _, t := range tasks {
			if task.ID == t.ID {
				found = true
				if task.Name != t.Name || task.Done != t.Done {
					printlog(fmt.Sprintf("Tasks are different: %v != %v", task, t), false)
					panic(errors.New("task not the same"))
				}
				// preparing for update
				task.Name = randomdata.Street()
				t.Name = task.Name
				t.Done = true
				task.Done = true
			}
		}
		if !found {
			printlog(fmt.Sprintf("Task not found: %v", task), false)
			panic(errors.New("list not found"))
		}
	}
}

func sendRequest(method string, endpoint string, headers map[string]string, data string) (*http.Response, error) {
	client := &http.Client{}

	fullURL := GatewayHost + GatewayPort + endpoint
	if FinalURL != "" {
		fullURL = FinalURL + endpoint
	} else if Endpoint != "" {
		fullURL = Endpoint + endpoint
	}
	fmt.Println("[", method, "]", fullURL)

	req, err := http.NewRequest(method, fullURL, strings.NewReader(data))
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	if token != nil && len(token.Token) > 0 {
		req.Header.Add("Authorization", token.Token)
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

	println("token", token)

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

	println("token", token)

	return username, nil

}

type Task struct {
	ID        string    `json:"_id"`
	UserID    string    `json:"userid"`
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"createdon"`
	Done      bool      `json:"done"`
}

func getTaskList() ([]*Task, error) {

	res, err := sendRequest("GET", "/tasks", defaultHeaders, "")
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func createTask(name string) (*Task, error) {
	task := &Task{Name: name}
	data, err := json.Marshal(task)
	res, err := sendRequest(http.MethodPost, "/tasks", defaultHeaders, string(data))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func updateTask(id, name string, done bool) (*Task, error) {
	task := &Task{ID: id, Name: name, Done: done}
	data, err := json.Marshal(task)
	res, err := sendRequest(http.MethodPut, "/tasks", defaultHeaders, string(data))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
