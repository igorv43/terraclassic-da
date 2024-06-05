package terraclassicda

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	clientx "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"

	tyepstx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/tidwall/gjson"
)
type TxResponse struct {
    Height    int64          `json:"height"`    // The block height where the transaction was included
    TxHash    string         `json:"txhash"`    // The transaction hash
    Codespace string         `json:"codespace"` // Namespace for the error code
    Code      uint32         `json:"code"`      // Response code (0 for success)
    Data      string         `json:"data"`      // Result data (if any)
    RawLog    string         `json:"raw_log"`   // Raw log message
    Logs      ABCIMessageLogs `json:"logs"`      // Logs from the transaction execution
    Info      string         `json:"info"`      // Additional information
    GasWanted int64          `json:"gas_wanted"`// Gas requested
    GasUsed   int64          `json:"gas_used"`  // Gas used
    Tx        Tx             `json:"tx"`        // The transaction itself
    Timestamp string         `json:"timestamp"` // Timestamp of the block
    Events    []Event        `json:"events"`    // Events emitted during the transaction execution
}

// ABCIMessageLogs represents a slice of ABCIMessageLog.
type ABCIMessageLogs []ABCIMessageLog

// ABCIMessageLog defines a structure for logging messages from an ABCI application.
type ABCIMessageLog struct {
    MsgIndex uint32         `json:"msg_index"` // Message index in the transaction
    Log      string         `json:"log"`       // Log message
    Events   StringEvents   `json:"events"`    // Events generated by the message
}

// StringEvents represents a slice of StringEvent.
type StringEvents []StringEvent

// StringEvent defines a structure for string events.
type StringEvent struct {
    Type       string            `json:"type"`       // Type of the event
    Attributes []EventAttribute  `json:"attributes"` // Attributes of the event
}

// EventAttribute defines a structure for event attributes.
type EventAttribute struct {
    Key   string `json:"key"`   // Attribute key
    Value string `json:"value"` // Attribute value
}

// Tx defines a structure for a Cosmos SDK transaction.
type Tx struct {
    // Implementation of the transaction details (e.g., messages, signatures)
}

// Event defines a structure for an event.
type Event struct {
    Type       string            `json:"type"`       // Event type
    Attributes []EventAttribute  `json:"attributes"` // Event attributes
}


type ABCIMessageLogx struct {
	MsgIndex uint32        `json:"msg_index"`
	Events   []StringEvent `json:"events"`
}

type TxResponseLog struct {
	ABCIMessageLogs []ABCIMessageLogx `json:"abcimessage_logs"`
}

type TxRequest struct {
    TxBytes string `json:"tx_bytes"`
    Mode    string `json:"mode"`
}
type SubmitBlob struct {
    Contents Contents `json:"contents"`
}

type Contents struct {
    Action  string `json:"action"`
    BlobID  int    `json:"blob_id"`
    Message string `json:"message"`
}

type SubmitReq struct {
    SubmitBlob SubmitBlob `json:"submit_blob"`
}

func NewTerraClassicTX(configtx Config,ctx context.Context, feeAmount uint64, gasLimit uint64){
	const pathTx  = "/cosmos/tx/v1beta1/txs"
	
		nodeURL         := configtx.LcURL // URL do nó Cosmos REST
		chainID         := configtx.AppID
		denom           := "uluna"
		privateKeyHex   := configtx.PrivateKeyHex // Chave privada em formato hexadecimal
		fromAddress     := configtx.FromAddress           // Endereço do remetente
		contractAddress := configtx.ContractAddress         // Endereço do contrato inteligente
		// feeAmount       := 230000000
		// gasLimit        := 2000000
	
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount("terra", "terrapub")
	config.Seal()


	// executeMsg := map[string]interface{}{
	// 	"increment": map[string]interface{}{},
	// }
	var executeMsg SubmitReq

    // Atribuição de valores aos campos da estrutura
    executeMsg.SubmitBlob.Contents.Action = "submit"
    executeMsg.SubmitBlob.Contents.BlobID = 0
    executeMsg.SubmitBlob.Contents.Message = "VGhpcyBpcyBibG9iICMy"

	executeMsgBytes, err := json.Marshal(executeMsg)
	if err != nil {
		log.Fatalf("failed to marshal execute message: %v", err)
	}

	msg := wasmtypes.MsgExecuteContract{
		Sender:   fromAddress,
		Contract: contractAddress,
		Msg:      executeMsgBytes,
		Funds:    sdktypes.Coins{}, // Corrigido para o tipo sdk.Coins
	}


    clientCtx := clientx.Context{}.
      
        WithChainID(chainID).
        WithNodeURI(nodeURL)

    // Consultar o nó para obter informações de estado
    clientCtx = clientCtx.WithNodeURI(nodeURL)
    rpcClient, err := clientx.NewClientFromNode(nodeURL)
    if err != nil {
        log.Fatalf("failed to create RPC client: %v", err)
    }
    clientCtx = clientCtx.WithClient(rpcClient)

    // Carregar configurações de transações
    marshaler := codec.NewProtoCodec(clientCtx.InterfaceRegistry)
	
    txConfig := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)
    clientCtx = clientCtx.WithTxConfig(txConfig)
	txBuilder:= clientCtx.TxConfig.NewTxBuilder()
	
	err = txBuilder.SetMsgs(&msg)
    if err != nil {
		log.Fatalf("Failed to create message: %v", err)
    }

	txBuilder.SetGasLimit(uint64(gasLimit)) // Set the gas limit
    txBuilder.SetFeeAmount(sdktypes.NewCoins(sdktypes.NewInt64Coin(denom, int64(feeAmount)  ) )) 
    txBuilder.SetMemo("")
    txBuilder.SetTimeoutHeight(0)
	
	
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

	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("failed to decode private key: %v", err)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	
	pubKey := privKey.PubKey()



	var sigsV2 []signing.SignatureV2
	sigV2 := signing.SignatureV2{
		PubKey: pubKey,
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode(clientCtx.TxConfig.SignModeHandler().DefaultMode()),
			Signature: nil,
		},
		Sequence:  sequence,
	}
	sigsV2 = append(sigsV2,sigV2)
	err = txBuilder.SetSignatures(sigsV2...)
    if err != nil {
		log.Fatalf("Erro na assinatura: %v", err)
    }

	
	sigsV2 = []signing.SignatureV2{}
	var signerData = xauthsigning.SignerData{
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
	
	
	sigV2, err =  tx.SignWithPrivKey(
		signing.SignMode(clientCtx.TxConfig.SignModeHandler().DefaultMode()), signerData,
		txBuilder,&privKey, clientCtx.TxConfig, sequence)
	if err != nil {
		log.Fatalf("Erro na assinatura: %v", err)
	}

	sigsV2 = append(sigsV2, sigV2)
	err = txBuilder.SetSignatures(sigsV2...)
    if err != nil {
        log.Fatalf("erro na assinatura : %v", err)
    }
	// Generated Protobuf-encoded bytes.
    txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
    if err != nil {
		log.Fatalf(".GetTx: %v", err)
    }
	//fmt.Printf("Transaction broadcast response TX: %v\n", string(txBytes))
    // Generate a JSON string.
    // txJSONBytes, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
    // if err != nil {
    //     log.Fatalf("TxJSONEncoder .GetTx: %v", err)
    // }
	// fmt.Printf("Transaction broadcast response TX JSON: %v\n", string(txJSONBytes))

	// base64Request := base64.StdEncoding.EncodeToString(txJSONBytes)
	// var txRequest  TxRequest
	// txRequest.Mode="BROADCAST_MODE_SYNC"
	// txRequest.TxBytes = string(base64Request)
	// requestBody, err := json.Marshal(txRequest)
	// if err != nil {
	// 	log.Fatalf("erro na strutura txRequest: %v", err)
	// }
	// response, err := http.Post(nodeURL+pathTx, "application/json", bytes.NewBuffer(requestBody))
	// if err != nil {
	// 	log.Fatalf("erro na  Post: %v", err)
	// }

	// defer func() {
	// 	err = response.Body.Close()
	// 	if err != nil {
	// 		log.Println("error closing response body", err)
	// 	}
	// }()

	
	// responseData, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	log.Fatalf("erro na Body - Post: %v", err)
	// }
	// log.Println("sucesso",string(responseData))
	grpcConn,_ := grpc.Dial(
        "127.0.0.1:9090", // Or your gRPC server address.
        grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
    )
    defer grpcConn.Close()

    // Broadcast the tx via gRPC. We create a new client for the Protobuf Tx
    // service.
    txClient := tyepstx.NewServiceClient(grpcConn)
    // We then call the BroadcastTx method on this client.
    grpcRes, err := txClient.BroadcastTx(
        ctx,
        &tyepstx.BroadcastTxRequest{
            Mode:    tyepstx.BroadcastMode_BROADCAST_MODE_BLOCK,
            TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
        },
    )
    if err != nil {
		log.Fatalf("erro na transação: %v", err)
    }
    
    fmt.Println(grpcRes.TxResponse) // Should be `0` if the tx is successful
	TxResponseBytes, err := json.Marshal(grpcRes.TxResponse)
	if err != nil {
		log.Fatalf("failed to marshal TxResponseBytes message: %v", err)
	}
	var txResponse TxResponse
	err = json.Unmarshal(TxResponseBytes, &txResponse)
	if err != nil {
		log.Fatalln("Error parsing TxResponse:", err)
	}
	// TxResponseLogBytes, err := json.Marshal(txResponse.RawLog)
	// if err != nil {
	// 	log.Fatalf("failed to marshal TxResponseLogBytes message: %v", err)
	// }
	var logs []ABCIMessageLog
	// var txResponseLog TxResponseLog
	err = json.Unmarshal([]byte(txResponse.RawLog), &logs)
	if err != nil {
		log.Fatalln("Error parsing TxResponseLog:", err)
	}
	terra_block_number := findAttributeByKeyName(logs, "terra_block_number")
	if terra_block_number != nil {
		fmt.Printf("Valor da chave '%s': %s\n", "terra_block_number", terra_block_number.Value)
	} else {
		fmt.Printf("Chave '%s' não encontrada\n", "terra_block_number")
	}
	// blob_id := findAttributeByKeyName(txResponseLog, "blob_id")
	// terra_previous_block := findAttributeByKeyName(txResponseLog, "terra_previous_block")
	// blob_count := findAttributeByKeyName(txResponseLog, "blob_count")

	fmt.Println(string(terra_block_number.Value))
}
// Função personalizada para encontrar um EventAttribute com uma chave específica
func findAttributeByKeyName(logs []ABCIMessageLog, keyName string) *EventAttribute {
	for _, msgLog := range logs {
		for _, event := range msgLog.Events {
			for _, attr := range event.Attributes {
				if attr.Key == keyName {
					return &attr
				}
			}
		}
	}
	return nil // Retorna nil se não encontrar
}