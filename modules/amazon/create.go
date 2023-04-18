package amazon

type Record struct {
	Records []Field `json:records`
}

type Field struct {
	Fiels []amazonOrderInfo `json:fiels`
}

type amazonOrderInfo struct {
	OrderId                   string `json:order_id`
	ShipBy                    string `json:ship_by`
	Item                      string `json:product_name`
	Condition                 string `json:condition`
	SKU                       string `json:sku`
	Quantity                  string `json:quantity`
	OrderDate                 string `json:order_date`
	Price                     string `json:price`
	Tax                       string `json:tax`
	Shipping                  string `json:shipping_fee`
	Promotions                string `json:promotions`
	AmazonFee                 string `json:platform_fee`
	MarketPlaceFacilitatorTax string `json:marketplace_facilitator_tax`
	YourEarning               string `json:earning`
}
