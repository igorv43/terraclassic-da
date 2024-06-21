package terraclassicda

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"log"
	"net/http"

	"fmt"

	clientx "github.com/cosmos/cosmos-sdk/client"
	"github.com/rollkit/go-da"
	"github.com/tidwall/gjson"
)

// SubmitRequest represents a request to submit data.
type SubmitRequest struct {
	Data string `json:"data"`
}

// SubmitResponse represents the response after submitting data.
type SubmitResponse struct {
	BlockNumber      uint32 `json:"block_number"`
	BlockHash        string `json:"block_hash"`
	TransactionHash  string `json:"hash"`
	TransactionIndex uint32 `json:"index"`
}

// BlocksResponse represents the structure of a response containing blocks information.
// type BlocksResponse struct {
// 	BlockNumber      uint32             `json:"block_number"`
// 	DataTransactions []DataTransactions `json:"data_transactions"`
// }

// DataTransactions represents data transactions within the block.
type DataContractTransactions struct {
	BlockNumber              uint32             `json:"terra_block_number"`
	PreviousBlockNumber      uint32             `json:"terra_previous_block"`
	Data                     string             `json:"data"`
}
type DataContact struct {
	Data      []DataContractTransactions `json:"data"`
	
}
type DataContractTransactionsModel struct {
	BlockNumber         uint32 `json:"terra_block_number"`
	PreviousBlockNumber uint32 `json:"terra_previous_block"`
	Data                []string `json:"data"`
}
type DataContactModel struct {
	Data DataContractTransactionsModel `json:"data"`
}
// Config represents the configuration structure.
type Config struct {
	AppID             string `json:"app_ID"`
	LcURL             string `json:"lc_url"`
	GRPCServerAddress string `json:"grpc_server_address"`
	PrivateKeyHex 	  string `json:"private_key_hex"`
	FromAddress       string `json:"from_address"`
	ContractAddress   string `json:"contract_address"`
	RestURL           string `json:"rest_url"`
	FcdURL            string  `json:"fcd_url"`
	
}

type GetBlobByBlockRequest struct {
	GetBlobByBlock struct {
		TerraBlockNumber string `json:"terra_block_number"`
	} `json:"get_blob_by_block"`
}

// BlockURL represents the URL pattern for retrieving data and extrinsic information
const BlockURL = "/cosmwasm/wasm/v1/contract/"

// BLOCK_NOT_FOUND represents the string indicating that a block is not found.
const BLOCK_NOT_FOUND = "\"Not found\""

// PROCESSING_BLOCK represents the string indicating that a block is still being processed.
const PROCESSING_BLOCK = "\"Processing block\""

// TerraClassicDA implements the avail backend for the DA interface
type TerraClassicDA struct {
	config Config
	ctx    context.Context
}

// NewTerraClassicDA returns an instance of AvailDA
// func NewTerraClassicDA(config Config, ctx context.Context) *TerraClassicDA {
// 	return &TerraClassicDA{
// 		ctx:    ctx,
// 		config: Config{LcURL: config.LcURL, AppID: config.AppID,PrivateKeyHex: config.PrivateKeyHex, FromAddress: config.FromAddress,ContractAddress: config.ContractAddress,RestURL: config.RestURL},
// 	}
// }
// func NewTerraClassicDA(config Config, ctx context.Context) *TerraClassicDA {
// 	return &TerraClassicDA{
// 		ctx:    ctx,
// 		config: Config{LcURL: config.LcURL, AppID: config.AppID,PrivateKeyHex: config.PrivateKeyHex, FromAddress: config.FromAddress,ContractAddress: config.ContractAddress,RestURL: config.RestURL},
// 	}
// }
func NewTerraClassicDA(opts ...func(*TerraClassicDA) *TerraClassicDA) *TerraClassicDA {
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	data, err := os.ReadFile("config.json")
		if err != nil {
			log.Fatalln("Error reading config file:", err)
		}
		// Parse the configuration data into a Config struct
		var config Config
		err = json.Unmarshal(data, &config)
		if err != nil {
			log.Fatalln("Error parsing config file:", err)
		}
	da := &TerraClassicDA{
		ctx:    ctx,
		config: Config{LcURL: config.LcURL, AppID: config.AppID,PrivateKeyHex: config.PrivateKeyHex, FromAddress: config.FromAddress,ContractAddress: config.ContractAddress,RestURL: config.RestURL,FcdURL: config.FcdURL},
	}
	for _, f := range opts {
		da = f(da)
	}
	//da.pubKey, da.privKey, _ = ed25519.GenerateKey(rand.Reader)
	return da
}

var _ da.DA = &TerraClassicDA{}

// MaxBlobSize returns the max blob size
func (c *TerraClassicDA) MaxBlobSize(ctx context.Context) (uint64, error) {
	var maxBlobSize uint64 = 64 * 64 * 250
	return maxBlobSize, nil
}
func   ClientTX() (clientx.Context,uint64,uint64, error){
	data, err := os.ReadFile("config.json")
		if err != nil {
			log.Fatalln("Error reading config file:", err)
		}
		// Parse the configuration data into a Config struct
		var config Config
		err = json.Unmarshal(data, &config)
		if err != nil {
			log.Fatalln("Error parsing config file:", err)
		}
	nodeURL         := config.RestURL // URL do nó Cosmos REST
	chainID         := config.AppID
	fromAddress     := config.FromAddress
	clientCtx := clientx.Context{}. WithChainID(chainID). WithNodeURI(nodeURL)

    // Consultar o nó para obter informações de estado
    clientCtx = clientCtx.WithNodeURI(nodeURL)
    rpcClient, err := clientx.NewClientFromNode(nodeURL)
    if err != nil {
        log.Fatalf("failed to create RPC client: %v", err)
		//return nil, err
    }
    clientCtx = clientCtx.WithClient(rpcClient)
	// Obter sequência da conta e número da conta
	accountResp, err := http.Get(fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts/%s", nodeURL, fromAddress))
	if err != nil {
		log.Fatalf("failed to get account info: %v", err)
	}
	defer accountResp.Body.Close()

	accountBody, err := ioutil.ReadAll(accountResp.Body)
	if err != nil {
		log.Fatalf("failed to read account response body: %v", err)
	}
	
	
	sequence := gjson.Get(string(accountBody), "account.sequence").Uint()
	accountNumber := gjson.Get(string(accountBody), "account.account_number").Uint()
	return clientCtx, sequence,accountNumber, err
}
func isSequenceMismatchError(err error) bool {
    // Verifique se o erro é devido à incompatibilidade de sequência
    return strings.Contains(err.Error(), "account sequence mismatch")
}
// Submit each blob to avail data availability layer
func (c *TerraClassicDA) Submit(ctx context.Context, daBlobs []da.Blob, gasPrice float64, namespace da.Namespace) ([]da.ID, error) {
	resultChan := make(chan SubmitResponse, len(daBlobs))
	errorChan := make(chan error, len(daBlobs))

	//var wg sync.WaitGroup

	//var mu sync.Mutex
    

	for id, blob := range daBlobs {
		//wg.Add(1)
		
		// Start a goroutine for each blob
	//	go func(blob da.Blob) {
			//defer wg.Done()
			 clientCtx,sequence,accountNumber,err := ClientTX()
			if err != nil {
				log.Fatalf("Failed to create message: %v", err)
			}
			
			encodedBlob := base64.StdEncoding.EncodeToString(blob)
			blobID:= id
			 submitResponsetipo,err := NewTerraClassicTX(clientCtx,sequence,accountNumber,c.config,ctx ,encodedBlob,blobID , 900860000, 300000000)
			if err != nil {
				errorChan <- err
				fmt.Println(err)
				//return
			}
			
			requestBody, err := json.Marshal(submitResponsetipo)
			//fmt.Println("pega dados para comparar: ",submitResponsetipo)
			if err != nil {
				errorChan <- err
				//return
				fmt.Println(err)
			}
			var submitResponse SubmitResponse
			err = json.Unmarshal(requestBody, &submitResponse)
			if err != nil {
			
				errorChan <- err
				//return
				fmt.Println(err)
			}

			// Acquire the mutex before updating slices
			//fmt.Println("pega dados para comparar: ",string(requestBody))
			//mu.Lock()
			resultChan <- SubmitResponse{
				BlockNumber:      submitResponse.BlockNumber,
				BlockHash:        submitResponse.BlockHash,
				TransactionHash:  submitResponse.TransactionHash,
				TransactionIndex: submitResponse.TransactionIndex,
			}
			//mu.Unlock()
           
		//}(blob)
		time.Sleep(15 * time.Second)
	}

	//go func() {
		//wg.Wait()
		close(resultChan)
		close(errorChan)
	//}()

	// Collect results from channels
	var ids []da.ID

	for result := range resultChan {
		ids = append(ids, makeID(result.BlockNumber))
		
	}
	// for err := range errorChan {
		
	// 	if err == nil {
    //         fmt.Printf("Transaction successful")
    //         break
    //     } else if isSequenceMismatchError(err) {
    //         fmt.Println("Sequence mismatch, retrying...")
    //        // time.Sleep(10 * time.Second) // Atraso antes de tentar novamente
    //        // continue
    //     } else {
	// 		time.Sleep(10 * time.Second) // Atraso antes de tentar novamente
    //        // continue
	// 	   fmt.Println("erro geral: ",err)
    //        // panic(err)
    //     }
	// }
	// Check for errors
	if err := <-errorChan; err != nil {
		return nil, err
	}

	fmt.Println("successfully submitted blobs to terraclassic")
	for _, id := range ids {
		blockNumber := binary.BigEndian.Uint32(id)
		fmt.Println("ids blob terraclassic",blockNumber)
	}
	return ids, nil
}

// Get returns Blob for each given ID, or an error
func (c *TerraClassicDA) Get(ctx context.Context, ids []da.ID, namespace da.Namespace) ([]da.Blob, error) {
	
	
	var blobs [][]byte
	//var blockNumber uint32
	for _, id := range ids {
	
		blockNumber := binary.BigEndian.Uint32(id)
		
	
		dataBlobs,err:= GetBlock(blockNumber,c.config)
		if err != nil {
			//log.Println("error 1", err)
			return nil, err
	
		}
		for _, data := range dataBlobs {
			decodeStr, _ := base64.StdEncoding.DecodeString(data)
			blobs = append(blobs, []byte(string(decodeStr)))
		}
		
	}
	return blobs, nil
}

// GetIDs returns the ID
func (c *TerraClassicDA) GetIDs(ctx context.Context, height uint64, namespace da.Namespace) ([]da.ID, error) {
	// todo:currently returning height as ID, need to extend avail-light api
	heightAsUint32 := uint32(height)
	ids := make([]byte, 8)
	binary.BigEndian.PutUint32(ids, heightAsUint32)
	return [][]byte{ids}, nil
}

// GetProofs returns inclusion Proofs for Blobs specified by their IDs.
func (c *TerraClassicDA) GetProofs(ctx context.Context, ids []da.ID, namespace da.Namespace) ([]da.Proof, error) {
	// TODO: add transaction hash to ID, so we can use it for proofs here
	var proofs []da.Proof
	for _, id := range ids {
		proofs = append(proofs, makeProofs(string(id)))
	}
	return proofs, nil
}

// Commit creates a Commitment for each given Blob.
func (c *TerraClassicDA) Commit(ctx context.Context, daBlobs []da.Blob, namespace da.Namespace) ([]da.Commitment, error) {
	return nil, nil
}

// Validate validates Commitments against the corresponding Proofs
func (c *TerraClassicDA) Validate(ctx context.Context, ids []da.ID, daProofs []da.Proof, namespace da.Namespace) ([]bool, error) {
	return nil, nil
}

func makeID(blockNumber uint32) da.ID {
	// IDs are not unique in rollkit context and that this has to be improved
	id := make([]byte, 8)
	binary.BigEndian.PutUint32(id, blockNumber)
	return id
}

func makeProofs(proofs string) da.ID {
	return []byte(proofs)
}