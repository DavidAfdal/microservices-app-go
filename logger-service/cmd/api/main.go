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

const (
	webPort   = "80"
	rpcPort   = "5001"
	monggoUrl = "mongodb://mongo:27017"
	gRpcPort  = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconect

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)

	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	} ()

	app := Config{
		Models: data.New(client),
	}

	// register the rpc Server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()


    
	log.Println("Starting service on port: ", webPort)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
    
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}


}

// func (app *Config) serve() {
// 	srv := &http.Server{
// 		Addr: fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}
    
// 	err := srv.ListenAndServe()

// 	if err != nil {
// 		log.Panic(err)
// 	}
// }


func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port", rpcPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))

	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()

		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(monggoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "david",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("error connecting to Mongo : ", err)
		return nil, err
	}
	

	log.Println("Connenct To Mongo")

	return c, nil
}