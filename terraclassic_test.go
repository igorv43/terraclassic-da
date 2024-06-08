package terraclassicda_test

import (
	"sync"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	terraclassicda "github.com/igorv43/terraclassic-da"
	dummy "github.com/igorv43/terraclassic-da/test"
)

func TestTerraClassicDA(t *testing.T) {
	// config := terraclassicda.Config{
	// 	AppID: "localterra",
	// 	LcURL: "http://localhost:9000/v2",
	// 	RestURL:"http://localhost:1317",
	// 	PrivateKeyHex:"21a5a38c18761a6225ba032dbf398d98595aefaac2b5ace8c18f7a476939e64e",
	// 	FromAddress:"terra1dcegyrekltswvyy0xy69ydgxn9x8x32zdtapd8",
	// 	ContractAddress:"terra1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrquka9l6",

	// }
configx := sdktypes.GetConfig()
	configx.SetBech32PrefixForAccount("terra", "terrapub")
	configx.Seal()
	// nodeURL         := configtx.RestURL // URL do nó Cosmos REST
	// chainID         := configtx.AppID
	// denom           := "uluna"
	// privateKeyHex   := configtx.PrivateKeyHex // Chave privada em formato hexadecimal
	// fromAddress     := configtx.FromAddress           // Endereço do remetente
	// contractAddress := configtx.ContractAddress         // Endereço do contrato inteligente
	//ctx := context.Background()

	da := terraclassicda.NewTerraClassicDA()

	var wg sync.WaitGroup
	wg.Add(1)

	// Start the mock server in a separate goroutine
	go func() {
		defer wg.Done()
		dummy.StartMockServer()
	}()

	// Wait for the mock server to start
	wg.Wait()

	dummy.RunDATestSuite(t, da)
}