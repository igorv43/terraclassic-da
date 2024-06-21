package terraclassicda

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type DataContractTransactionsx struct {
	BlockNumber         uint32 `json:"terra_block_number"`
	PreviousBlockNumber uint32 `json:"terra_previous_block"`
	Data                string `json:"data"`
}
type DataContactx struct {
	Data []DataContractTransactionsx `json:"data"`
}

type DataContractTransactionsModelx struct {
	BlockNumber         uint32 `json:"terra_block_number"`
	PreviousBlockNumber uint32 `json:"terra_previous_block"`
	Data                []string `json:"data"`
}
type DataContactModelx struct {
	Data DataContractTransactionsModelx `json:"data"`
}
type GetBlobByBlockRequestx struct {
	GetBlobByBlock struct {
		TerraBlockNumber int `json:"terra_block_number"`
	} `json:"get_blob_by_block"`
}
func Ver() {
	contract_address:="terra1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrquka9l6"
	id := 23162
	requestData := GetBlobByBlockRequestx{
		GetBlobByBlock: struct {
			TerraBlockNumber int `json:"terra_block_number"`
		}{
			TerraBlockNumber: int(id),
		},
	}

	// Converter a estrutura da solicitação em JSON
	requestDataBase64, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Erro ao criar o corpo da solicitação:", err)

	}
	base64Request := base64.StdEncoding.EncodeToString(requestDataBase64)
	blocksURL := "http://localhost:1317"+BlockURL+contract_address+"/smart/"+base64Request
	parsedURL, err := url.Parse(blocksURL)
	if err != nil {
		log.Println("error 1", err)

	}
	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	log.Println("URL ", parsedURL.String())
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
	log.Println("teste:",string(responseData))
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
		//goto Loop
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

	 log.Println("sucesso  0",dataContactx.Data)
	// for _, dataTransaction := range blocksObject.Data {
	// 	decodeStr, _ := base64.StdEncoding.DecodeString(dataTransaction.Data)
	// 	log.Println("sucesso ", string(decodeStr))

	// 	//blobs = append(blobs, []byte(string(decodeStr)))
	// }

}
func Ver2(){
	contract_address:="terra1wdz7f49letx7fs58yke57tkmn24ffzxfj8hqmvafuzh5aaevzy7qgkterx"
	id := 18103462
	blocksURL := "https://fcd.terra-classic.hexxagon.io/v1/txs?block="+fmt.Sprint(id)+"&limit=10"
	parsedURL, err := url.Parse(blocksURL)
	if err != nil {
		log.Println("error 1", err)

	}
	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	
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
	if err != nil {
		log.Println("error 3", err)
	}
	var blocksObject Response
	if string(responseData) == BLOCK_NOT_FOUND {
		log.Println("sucesso BLOCK_NOT_FOUND")
		blocksObject = Response{Txs: []Tx{}}
	} else if string(responseData) == PROCESSING_BLOCK {
		log.Println("sucesso PROCESSING_BLOCK")
		time.Sleep(10 * time.Second)
		//goto Loop
	} else {
		err = json.Unmarshal(responseData, &blocksObject)
		if err != nil {
			log.Println("error 4", err)
		}
	}
	
	fmt.Println("limit:", blocksObject.Limit)
	for _, dataTx := range blocksObject.Txs {
		for _,msgTx := range dataTx.Tx.Value.Msg {
			 if(msgTx.Value.Contract ==contract_address)  {
				parsedMsg, err := parseMsg(msgTx.Value.Msg)
				if err != nil {
					fmt.Println("Error parsing msg:", err)
					continue
				}
				switch vData := parsedMsg.(type) {
				case SubmitBlob:
					
					fmt.Printf("Parsed SubmitBlob: %+v\n", vData)
					
				// case ClaimReward:
				// 	fmt.Printf("Parsed ClaimReward: %+v\n", v)
				default:
					fmt.Println("Unknown message type")
				}
			 } 
		}
		// dataTx.Tx.Value.
		// decodeStr, _ := base64.StdEncoding.DecodeString(dataTransaction.Data)
		// log.Println("sucesso ", string(decodeStr))

		//blobs = append(blobs, []byte(string(decodeStr)))
	}
}
func parseMsg(msg map[string]interface{}) (interface{}, error) {
	if submitBlob, ok := msg["submit_blob"]; ok {
		submitBlobData, err := json.Marshal(submitBlob)
		if err != nil {
			return nil, err
		}
		var sb SubmitBlob
		err = json.Unmarshal(submitBlobData, &sb)
		if err != nil {
			return nil, err
		}
		return sb, nil
	}
	// if claimReward, ok := msg["claim_reward"]; ok {
	// 	claimRewardData, err := json.Marshal(claimReward)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	var cr ClaimReward
	// 	err = json.Unmarshal(claimRewardData, &cr)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return cr, nil
	// }
	return nil, fmt.Errorf("unknown message type")
}