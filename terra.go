package terraclassicda

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"

// 	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
// 	clientx "github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/client/tx"
// 	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
// 	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
// 	sdktypes "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/types/tx/signing"
// 	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
// 	"github.com/tidwall/gjson"

// 	//authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
// 	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
// 	//rpchttp "github.com/tendermint/tendermint/rpc/client/http"
// )

// func NewTX2(){

// 	const (
// 		nodeURL         = "http://localhost:1317" // URL do nó Cosmos REST
// 		nodeURIRPC        = "tcp://localhost:26657"
// 		chainID         = "localterra"
// 		denom           = "uluna"
// 		privateKeyHex   = "21a5a38c18761a6225ba032dbf398d98595aefaac2b5ace8c18f7a476939e64e" // Chave privada em formato hexadecimal
// 		fromAddress     = "terra1dcegyrekltswvyy0xy69ydgxn9x8x32zdtapd8"           // Endereço do remetente
// 		contractAddress = "terra14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9ssrc8au"           // Endereço do contrato inteligente
// 		feeAmount       = "2000000"
// 		gasLimit        = "200000000"
// 	)

// 	//fmt.Printf("Transaction broadcast response TX JSON: ")
// 	config := sdktypes.GetConfig()
// 	config.SetBech32PrefixForAccount("terra", "terrapub")
// 	config.Seal()
// //wasm/v1beta1/tx
//  //sender, _ := sdktypes.AccAddressFromBech32("terra1mecnwyxx7zjwjrpfc0as7wc2f47kqgpp7lvmh7")
//  fromAddr, err := sdktypes.AccAddressFromBech32(fromAddress)
//  if err != nil {
// 	log.Fatalf("failed to AccAddressFromBech32 message: %v", err)
// }
// // contractAddress, _ := types.AccAddressFromBech32("cosmos1contractaddress")
// // msgData := []byte("dados da mensagem")

// // // Criar uma nova mensagem para executar um contrato

// executeMsg := map[string]interface{}{
// 	"increment": map[string]interface{}{},
// }
// executeMsgBytes, err := json.Marshal(executeMsg)
// if err != nil {
// 	log.Fatalf("failed to marshal execute message: %v", err)
// }
// //rawMsg := wasmtypes.RawContractMessage(executeMsgBytes)

// msg := wasmtypes.MsgExecuteContract{
// 	Sender:   fromAddr.String(),
// 	Contract: contractAddress,
// 	Msg:      executeMsgBytes,
// 	Funds:    sdktypes.Coins{}, // Corrigido para o tipo sdk.Coins
// }

// // Configurar client.Context
// clientCtx := clientx.Context{}.

// WithChainID(chainID).
// WithNodeURI(nodeURIRPC)

// // Configurar cliente RPC
// rpcClient, err := rpchttp.New(nodeURIRPC, "/websocket")
// if err != nil {
// log.Fatalf("failed to create RPC client: %v", err)
// }
// clientCtx = clientCtx.WithClient(rpcClient)

// 	txBuilder := clientCtx.TxConfig.NewTxBuilder()

// 	err = txBuilder.SetMsgs(&msg)
//     if err != nil {
// 		log.Fatalf("Failed to create message: %v", err)
//     }

// 	txBuilder.SetGasLimit(200000) // Set the gas limit
//     txBuilder.SetFeeAmount(sdktypes.NewCoins(sdktypes.NewInt64Coin("uluna", 230000000) ))
//     txBuilder.SetMemo("")
//     txBuilder.SetTimeoutHeight(0)
// 	privKeyBytes, err := hex.DecodeString(privateKeyHex)
// 	if err != nil {
// 		log.Fatalf("failed to decode private key: %v", err)
// 	}
// 	//priv1, _, addr1 := testdata.KeyTestPubAddr()
// 	//privKey := secp256k1.PrivKey(privKeyBytes)
// 	var privKey cryptotypes.PrivKey = secp256k1.GenPrivKeyFromSecret(privKeyBytes)
// 	pubKey := privKey.PubKey()

// 	// Obter sequência da conta e número da conta
// 	accountResp, err := http.Get(fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts/%s", nodeURL, fromAddress))
// 	if err != nil {
// 		log.Fatalf("failed to get account info: %v", err)
// 	}
// 	defer accountResp.Body.Close()

// 	accountBody, err := ioutil.ReadAll(accountResp.Body)
// 	if err != nil {
// 		log.Fatalf("failed to read account response body: %v", err)
// 	}
// 	//priv1, _, addr1 := testdata.KeyTestPubAddr()

// 	sequence := gjson.Get(string(accountBody), "account.sequence").Uint()
// 	accountNumber := gjson.Get(string(accountBody), "account.account_number").Uint()
// 	var sigsV2 []signing.SignatureV2
// 	sigV2 := signing.SignatureV2{
// 		PubKey: pubKey,
// 		Data: &signing.SingleSignatureData{
// 			SignMode:  signing.SignMode(clientCtx.TxConfig.SignModeHandler().DefaultMode()),
// 			Signature: nil,
// 		},
// 		Sequence:  sequence,
// 	}
// 	sigsV2 = append(sigsV2,sigV2)
// 	err = txBuilder.SetSignatures(sigsV2...)
//     if err != nil {
// 		log.Fatalf("Erro na assinatura: %v", err)
//     }

// 	sigsV2 = []signing.SignatureV2{}

// 	var signerData = xauthsigning.SignerData{
// 		ChainID:       chainID,
// 		AccountNumber: accountNumber,
// 		Sequence:      sequence,
// 	}

// 	sigV2, err =  tx.SignWithPrivKey(
// 		signing.SignMode(clientCtx.TxConfig.SignModeHandler().DefaultMode()), signerData,
// 		txBuilder, privKey, clientCtx.TxConfig, sequence)
// 	if err != nil {
// 		log.Fatalf("Erro na assinatura: %v", err)
// 	}

// 	sigsV2 = append(sigsV2, sigV2)

// 	// Generated Protobuf-encoded bytes.
//     txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
//     if err != nil {
// 		log.Fatalf(".GetTx: %v", err)
//     }
// 	fmt.Printf("Transaction broadcast response TX: %v\n", string(txBytes))
//     // Generate a JSON string.
//     txJSONBytes, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
//     if err != nil {
//         log.Fatalf("TxJSONEncoder .GetTx: %v", err)
//     }
// 	fmt.Printf("Transaction broadcast response TX JSON: %v\n", string(txJSONBytes))

// }