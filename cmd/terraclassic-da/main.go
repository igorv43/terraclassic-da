// package main

// //go run ./cmd/terraclassic-da/main.go
// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"log"
// 	"net"
// 	"os"

// 	sdktypes "github.com/cosmos/cosmos-sdk/types"
// 	terraclassicda "github.com/igorv43/terraclassic-da"
// 	proxy "github.com/rollkit/go-da/proxy/grpc"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func main()  {
// 	ctx := context.Background()
// 	configx := sdktypes.GetConfig()
// 	configx.SetBech32PrefixForAccount("terra", "terrapub")
// 	configx.Seal()
// 	data, err := os.ReadFile("config.json")
// 	if err != nil {
// 		log.Fatalln("Error reading config file:", err)
// 	}
// 	// Parse the configuration data into a Config struct
// 	var config terraclassicda.Config
// 	err = json.Unmarshal(data, &config)
// 	if err != nil {
// 		log.Fatalln("Error parsing config file:", err)
// 	}
// 	da :=terraclassicda.NewTerraClassicDA(config,ctx)

//		srv := proxy.NewServer(da, grpc.Creds(insecure.NewCredentials()))
//		lis, err := net.Listen("tcp", config.GRPCServerAddress)
//		if err != nil {
//			log.Fatalln("failed to create network listener:", err)
//		}
//		log.Println("serving terraclassic-da over gRPC on:", lis.Addr())
//		err = srv.Serve(lis)
//		if !errors.Is(err, grpc.ErrServerStopped) {
//			log.Fatalln("gRPC server stopped with error:", err)
//		}
//		//terraclassicda.Ver()
//	}
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	terraclassicda "github.com/igorv43/terraclassic-da"
	proxy "github.com/rollkit/go-da/proxy/jsonrpc"
	//goDATest "github.com/rollkit/go-da/test"
)

const (
	defaultHost = "localhost"
	defaultPort = "7980"
)

func main() {
	
	var (
		host      string
		port      string
		listenAll bool
	)
	flag.StringVar(&port, "port", defaultPort, "listening port")
	flag.StringVar(&host, "host", defaultHost, "listening address")
	flag.BoolVar(&listenAll, "listen-all", false, "listen on all network interfaces (0.0.0.0) instead of just localhost")
	flag.Parse()

	if listenAll {
		host = "0.0.0.0"
	}
	//ctx := context.Background()
	configx := sdktypes.GetConfig()
	configx.SetBech32PrefixForAccount("terra", "terrapub")
	configx.Seal()
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Error reading config file:", err)
	}
	// Parse the configuration data into a Config struct
	var config terraclassicda.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("Error parsing config file:", err)
	}
	da := terraclassicda.NewTerraClassicDA()
	srv := proxy.NewServer(host, port, da)
	log.Printf("Listening on: %s:%s", host, port)
	if err := srv.Start(context.Background()); err != nil {
		log.Fatal("error while serving:", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)
	<-interrupt
	fmt.Println("\nCtrl+C pressed. Exiting...")
	os.Exit(0)
}

