package main

//go run ./cmd/terraclassic-da/main.go
import (
	"context"
	"encoding/json"
	"log"
	"os"

	terraclassicda "github.com/igorv43/terraclassic-da"
)

func main()  {
	ctx := context.Background()
	//terraclassicda.NewTX(ctx)
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
	terraclassicda.NewTerraClassicTX(config,ctx,230000000,2000000)
}