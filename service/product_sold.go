package service

type ProductSoldRequest struct {
	ProductId  string `json:"product_id"`
	OrderId    string `json:"order_id"`
	ExternalId string `json:"external_id"`
	Quantity   string `json:"quantity"`
}

func ProductSold(request ProductSoldRequest) error {
	//get order
	//get date by timestamp
	//find break by date
	//add order there
	return nil
}
