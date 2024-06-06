package terraclassicda

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"sync"
	"time"

	"io"
	"log"
	"net/http"
	"net/url"

	"fmt"

	"github.com/rollkit/go-da"
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
func NewTerraClassicDA(config Config, ctx context.Context) *TerraClassicDA {
	return &TerraClassicDA{
		ctx:    ctx,
		config: Config{LcURL: config.LcURL, AppID: config.AppID,PrivateKeyHex: config.PrivateKeyHex, FromAddress: config.FromAddress,ContractAddress: config.ContractAddress,RestURL: config.RestURL},
	}
}

var _ da.DA = &TerraClassicDA{}

// MaxBlobSize returns the max blob size
func (c *TerraClassicDA) MaxBlobSize(ctx context.Context) (uint64, error) {
	var maxBlobSize uint64 = 64 * 64 * 500
	return maxBlobSize, nil
}

// Submit each blob to avail data availability layer
func (c *TerraClassicDA) Submit(ctx context.Context, daBlobs []da.Blob, gasPrice float64, namespace da.Namespace) ([]da.ID, error) {
	resultChan := make(chan SubmitResponse, len(daBlobs))
	errorChan := make(chan error, len(daBlobs))

	var wg sync.WaitGroup

	var mu sync.Mutex

	for _, blob := range daBlobs {
		wg.Add(1)

		// Start a goroutine for each blob
		go func(blob da.Blob) {
			defer wg.Done()
			encodedBlob := base64.StdEncoding.EncodeToString(blob)
			blobID:= 1
			var submitResponsetipo = NewTerraClassicTX(c.config,ctx ,encodedBlob,blobID , 900860000, 2000000)
			requestBody, err := json.Marshal(submitResponsetipo)
			if err != nil {
				errorChan <- err
				return
			}
			var submitResponse SubmitResponse
			err = json.Unmarshal(requestBody, &submitResponse)
			if err != nil {
				errorChan <- err
				return
			}

			// Acquire the mutex before updating slices
			//fmt.Println("pega dados para comparar: ",string(requestBody))
			mu.Lock()
			resultChan <- SubmitResponse{
				BlockNumber:      submitResponse.BlockNumber,
				BlockHash:        submitResponse.BlockHash,
				TransactionHash:  submitResponse.TransactionHash,
				TransactionIndex: submitResponse.TransactionIndex,
			}
			mu.Unlock()
           
		}(blob)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results from channels
	var ids []da.ID

	for result := range resultChan {
		ids = append(ids, makeID(result.BlockNumber))
		
	}

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
	Loop:
		blockNumber := binary.BigEndian.Uint32(id)
		requestData := GetBlobByBlockRequestx{
			GetBlobByBlock: struct {
				TerraBlockNumber int `json:"terra_block_number"`
			}{
				TerraBlockNumber: int(blockNumber),
			},
		}
	
		// Converter a estrutura da solicitação em JSON
		requestDataBase64, err := json.Marshal(requestData)
		if err != nil {
			fmt.Println("Erro ao criar o corpo da solicitação:", err)
	
		}
		base64Request := base64.StdEncoding.EncodeToString(requestDataBase64)
		blocksURL := c.config.RestURL+BlockURL+c.config.ContractAddress+"/smart/"+base64Request
		parsedURL, err := url.Parse(blocksURL)
		if err != nil {
			log.Println("error 1", err)
	
		}
		req, err := http.NewRequest("GET", parsedURL.String(), nil)
		//log.Println("URL ", parsedURL.String())
		if err != nil {
			log.Println("error 2", err)
		}
		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			log.Println("error 3", err)
		}
		defer func() {
			err = response.Body.Close()
			if err != nil {
				log.Println("error closing response body", err)
			}
		}()
		responseData, err := io.ReadAll(response.Body)
		//log.Println("teste:",string(responseData))
		if err != nil {
			log.Println("error 3", err)
		}
		var blocksObject DataContactModel
		if string(responseData) == BLOCK_NOT_FOUND {
			log.Println("sucesso BLOCK_NOT_FOUND")
			blocksObject = DataContactModel{Data: DataContractTransactionsModel{}}
		} else if string(responseData) == PROCESSING_BLOCK {
			log.Println("sucesso PROCESSING_BLOCK")
			time.Sleep(10 * time.Second)
			goto Loop
		} else {
			err = json.Unmarshal(responseData, &blocksObject)
			if err != nil {
				log.Println("error 4", err)
			}
		}
		var dataContactx DataContact 
		var listDataContactx  []DataContractTransactions
		for _, msData := range blocksObject.Data.Data { 
			listDataContactx = append(listDataContactx, DataContractTransactions{
				BlockNumber: blocksObject.Data.BlockNumber        ,
				PreviousBlockNumber: blocksObject.Data.PreviousBlockNumber,
				Data:                msData,
			})
		}
		dataContactx.Data =listDataContactx
		
		for _, dataTransaction := range dataContactx.Data {
			decodeStr, _ := base64.StdEncoding.DecodeString(dataTransaction.Data)
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