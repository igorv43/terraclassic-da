package main

import (
	"context"

	terraclassicda "github.com/igorv43/terraclassic-da"
)

func main()  {
	ctx := context.Background()
	terraclassicda.NewTX3(ctx)
}