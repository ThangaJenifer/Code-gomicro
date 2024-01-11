package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Here we have defined rpcPort 5001
//Excercise 62 in section 7, we have defined RPC payload, RPC type and one method for it
//Excercise 63 we will start the RPC server itself

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	//REGISTER the RPC Server
	//Register publishes the receiver's methods in the DefaultServer.
	err = rpc.Register(new(RPCServer))
	//Run the RPC server or listen for rpc
	go app.rpcListen()

	//We will listen for grpc which is as simple as running it in gorountine
	go app.gRPCListen()

	// start web server
	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}

}

// Create a method for rpcListen for Config type app to start rpc server
func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port ", rpcPort)
	//we declare a varibale listen and err, from standard library net.Listen on tcp and on all ip address port 0.0.0.0: with port rpcPort
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	//need a loop that executes forever, all we say is listen.Accept(), we accept connections
	//we say err != nil continue, just continue dont stop and  start over again and again
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		//We say go rpc.ServConn(rpcConn) to run that in backgorund
		//ServeConn runs the DefaultServer on a single connection. ServeConn blocks, serving the connection until the client hangs up.
		//The caller typically invokes ServeConn in a go statement. ServeConn uses the gob wire format (see package gob) on the connection.
		go rpc.ServeConn(rpcConn)
	}

}

// func (app *Config) serve() {
// 	srv := &http.Server{
// 		Addr: fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}

// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Panic()
// 	}
// }

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
