package terraclassicda

//https://github.com/terra-money/terra.go/blob/master/tx/txbuilder.go
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
	tyepstx "github.com/cosmos/cosmos-sdk/types/tx"
	"google.golang.org/grpc"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	//banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/tidwall/gjson"
)

func NewTX4(ctx context.Context){
	
	const (
		nodeURL         = "http://localhost:1317" // URL do nó Cosmos REST
		chainID         = "localterra"
		denom           = "uluna"
		privateKeyHex   = "21a5a38c18761a6225ba032dbf398d98595aefaac2b5ace8c18f7a476939e64e" // Chave privada em formato hexadecimal
		fromAddress     = "terra1dcegyrekltswvyy0xy69ydgxn9x8x32zdtapd8"           // Endereço do remetente
		contractAddress = "terra14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9ssrc8au"           // Endereço do contrato inteligente
		feeAmount       = "2000000"
		gasLimit        = "200000000"
	)
//wasm/v1beta1/tx
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount("terra", "terrapub")
	config.Seal()


executeMsg := map[string]interface{}{
	"increment": map[string]interface{}{},
}
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

	txBuilder.SetGasLimit(200000) // Set the gas limit
    txBuilder.SetFeeAmount(sdktypes.NewCoins(sdktypes.NewInt64Coin("uluna", 230000000) )) 
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
	//priv1, _, addr1 := testdata.KeyTestPubAddr()
	
	sequence := gjson.Get(string(accountBody), "account.sequence").Uint()
	accountNumber := gjson.Get(string(accountBody), "account.account_number").Uint()

	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("failed to decode private key: %v", err)
	}


	var privKey cryptotypes.PrivKey = secp256k1.GenPrivKeyFromSecret(privKeyBytes)
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
		txBuilder, privKey, clientCtx.TxConfig, sequence)
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
	fmt.Printf("Transaction broadcast response TX: %v\n", string(txBytes))
    // Generate a JSON string.
    txJSONBytes, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
    if err != nil {
        log.Fatalf("TxJSONEncoder .GetTx: %v", err)
    }
	fmt.Printf("Transaction broadcast response TX JSON: %v\n", string(txJSONBytes))

	// Create a connection to the gRPC server.
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
            Mode:    tyepstx.BroadcastMode_BROADCAST_MODE_SYNC,
            TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
        },
    )
    if err != nil {
		log.Fatalf("erro na transação: %v", err)
    }

    fmt.Println(grpcRes.TxResponse) // Should be `0` if the tx is successful

}
