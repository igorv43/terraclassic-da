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
