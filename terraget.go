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
	Data DataContractTransactionsx `json:"data"`
}
type GetBlobByBlockRequestx struct {
	GetBlobByBlock struct {
		TerraBlockNumber int `json:"terra_block_number"`
	} `json:"get_blob_by_block"`
}
func Ver() {
	contract_address:="terra1wdz7f49letx7fs58yke57tkmn24ffzxfj8hqmvafuzh5aaevzy7qgkterx"
	id := 18103462
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
	blocksURL := "https://terra-classic-lcd.publicnode.com"+BlockURL+contract_address+"/smart/"+base64Request
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
	if err != nil {
		log.Println("error 3", err)
	}
	var blocksObject DataContactx
	if string(responseData) == BLOCK_NOT_FOUND {
		log.Println("sucesso BLOCK_NOT_FOUND")
		blocksObject = DataContactx{Data: DataContractTransactionsx{}}
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
	log.Println("sucesso  0",blocksObject.Data.Data)
	// for _, dataTransaction := range blocksObject.Data {
	// 	decodeStr, _ := base64.StdEncoding.DecodeString(dataTransaction.Data)
	// 	log.Println("sucesso ", string(decodeStr))

	// 	//blobs = append(blobs, []byte(string(decodeStr)))
	// }

}
