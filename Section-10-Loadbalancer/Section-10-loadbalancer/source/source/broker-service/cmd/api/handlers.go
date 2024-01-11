package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)

}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)

}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

// Created a type exactly matches in RPC server end with Name and Data
type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	//need to get the rpc client, Go makes it remearkable easy by using rpc package and Dial
	//Dial using tcp and name of microservice in docker-compose port is 5001
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//once we have client rpc need to create a payload, payload type should be exactly matching one that Remote RPC server expects to get
	//Create var rpcPayload of type RPCPayload, populate l LogPayload values into it
	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}
	//we want to get some kind of result back by creating var result storing response of remote rpc call
	var result string
	//Clinet rpc call we say exactly want to call that is RPCServer which is type created at server side and LogInfo function func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error in rpc.go in logger-service/cmd/api
	//Any method Want to expose to rpc on the server end should be exported
	//Second parameter is data rpcPayload which logpayload data is populated
	//Last thing is response from server which is stored as reference to result variable
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//This will write JSON back to end-user
	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// we need to create a grpc client in the broker handler to communicate with log service GRPC server for grpc responses
// Creating LogViaGRPC handler func with response writter and request
func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	//This handler will receive JSON payload so declare a var requestPayload of type RequestPayload
	var requestPayload RequestPayload
	//we read posting JSON into that variable requestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//connect to our grpc server in logservice using conn var and grpc.Dial and dial logger service address and listen port
	//next parameter is credintails as GRPC needs valid creds to connect to it. Now we dont need creds as running in own docker cluster. still we gonna pass something grpc.WithTransportCredentials(insecure.NewCredentials()) which is empty credientails
	//Last thing is what are options, grpc.WithBlock()
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()
	//All above code is for connection now we need client
	//Create var c for client asssign it by calling logs package where generated our code NewLogServiceClient and pass conn we got

	c := logs.NewLogServiceClient(conn)
	//need a context with timeout as we use it in next step using context.Background available and just a second timeout as grpc is fast transport
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//last step is calling the WriteLog func, it requires context and reference to logs.LogRequest
	//Inside of logs.LogRequest, we create a reference log entry requestPayload.Log.Name and requestPayload.Log.Data
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//if we get pass this control, it means we have written log using grpc so needs to send response back
	//below response we send it back to frontend
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)

	//Last thing we need to do is go to routes.go and create a route for this handler
}
