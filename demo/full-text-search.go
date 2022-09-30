package main

import (
	"context"
	"fmt"
	client "github.com/zinclabs/sdk-go-zincsearch"
)

func main() {
	zincHandler := getZincHandler()

	matchQuery := *client.NewMetaMatchQuery()
	matchQuery.SetQuery("wlop")
	subQuery := *client.NewMetaQuery()
	subQuery.SetMatch(map[string]client.MetaMatchQuery{
		"Prompt": matchQuery,
	})
	boolQuery := *client.NewMetaBoolQuery()
	boolQuery.SetShould([]client.MetaQuery{subQuery})
	queryQuery := *client.NewMetaQuery()
	queryQuery.SetBool(boolQuery)
	zincHandler.Search(queryQuery)

	//zincHandler.Bulk([]map[string]interface{}{
	//	{"Prompt": "hyper detailed ultra sharp of a beautiful fractal. trending on artstation, golden, delicate, facing camera, hyper realism, 1 4 5 0, engraving, ultra realistic, 8 k"},
	//	{"Prompt": "steampunk cybernetic biomechanical cecropia moth with wings, 3 d model, very coherent symmetrical artwork, unreal engine realistic render, 8 k, micro detail, intricate, elegant, highly detailed, centered, digital painting, artstation, smooth, sharp focus, illustration, artgerm, tomasz alen kopera, wlop"},
	//})
	//
	//zincHandler.Update("1S8x9go1zWM", map[string]any{
	//	"Prompt": "steampunk cybernetic biomechanical cecropia moth with wings, 3 d model, very coherent symmetrical artwork, unreal engine realistic render, 8 k, micro detail, intricate, elegant, highly detailed, centered, digital painting, artstation, smooth, sharp focus, illustration, artgerm, tomasz alen kopera, wlop",
	//})
	//zincHandler.Delete("1S8x9go1zWM")
}

type ZincHandler struct {
	apiClient *client.APIClient
	ctx       context.Context
	index     string
}

const (
	//zincURL = "https://zinc.wenuts.top"
	zincURL  = "http://zinc.wenuts.top"
	userName = "liuwei"
	password = "liuwei123"
	index    = "prompt"
)

func getZincHandler() *ZincHandler {
	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{
		client.ServerConfiguration{
			URL: zincURL,
		},
	}

	return &ZincHandler{
		apiClient: client.NewAPIClient(configuration),
		ctx: context.WithValue(context.Background(), client.ContextBasicAuth, client.BasicAuth{
			UserName: userName,
			Password: password,
		}),
		index: index,
	}
}

func (h *ZincHandler) Search(queryQuery client.MetaQuery) {
	query := *client.NewMetaZincQuery()
	query.SetSize(10) // default: 10
	query.SetQuery(queryQuery)

	resp, r, err := h.apiClient.Search.Search(h.ctx, h.index).Query(query).Execute()
	if err != nil {
		fmt.Printf("Error when calling `SearchApi.Search``: %v\n", err)
		fmt.Printf("Full HTTP response: %v\n", r)
	}

	fmt.Println("Total num:", resp.Hits.Total.GetValue())
	for _, hit := range resp.Hits.Hits {
		fmt.Println(hit.Id, hit.GetTimestamp(), hit.GetSource())
	}
}

func (h *ZincHandler) Bulk(document []map[string]any) {
	query := client.NewMetaJSONIngest()
	query.SetIndex(h.index)
	query.SetRecords(document)

	resp, r, err := h.apiClient.Document.Bulkv2(h.ctx).Query(*query).Execute()
	if err != nil {
		fmt.Printf("Error when calling `Document.Bulkv2``: %v\n", err)
		fmt.Printf("Full HTTP response: %v\n", r)
	}
	// response from `Bulkv2`: MetaHTTPResponseRecordCount
	fmt.Printf("Response from `Document.Bulkv2`: %v\n", resp.GetRecordCount())
}

func (h *ZincHandler) Update(id string, document map[string]any) {
	resp, r, err := h.apiClient.Document.Update(h.ctx, h.index, id).Document(document).Execute()
	if err != nil {
		fmt.Printf("Error when calling `Document.Update``: %v\n", err)
		fmt.Printf("Full HTTP response: %v\n", r)
	}
	// response from `Update`: MetaHTTPResponseID
	fmt.Printf("Response from `Document.Update`: %v\n", resp.GetId())
}

func (h *ZincHandler) Delete(id string) {
	resp, r, err := h.apiClient.Document.Delete(h.ctx, h.index, id).Execute()
	if err != nil {
		fmt.Printf("Error when calling `Document.Delete``: %v\n", err)
		fmt.Printf("Full HTTP response: %v\n", r)
	}
	// response from `Delete`: MetaHTTPResponseDocument
	fmt.Printf("Response from `Document.Delete`: %v\n", resp.GetId())
}
