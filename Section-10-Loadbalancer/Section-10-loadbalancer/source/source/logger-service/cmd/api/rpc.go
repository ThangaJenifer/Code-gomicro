package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

// RPCPayload is the type for data we receive from RPC
// RPCPayload is type of data we are going to receive for any methods tied to rpc server struct type
// Create same type in broker-service handlers.go matching it above logItemViaRPC() func for rpc calls
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo writes our payload to mongo
// This will take rpcpayload type and resp pointer to string , can return a potential error
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	//to write our payload to mongo, it needs a mongodb collection database and collection name
	collection := client.Database("logs").Collection("logs")
	//Inserting our payload into mongoDB using data package Logentry type of model Name, data and Time as mentioned
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	// resp is the message sent back to the RPC caller
	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}
