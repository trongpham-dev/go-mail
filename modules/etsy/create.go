package etsy

type EtsyOrderDetail struct {
	OrderId     string `json:order_id`
	OrderDate   string `json:order_date`
	ItemTotal   string `json:revenue`
	ShippingFee string `json:shipping_fee`
	SalesTax    string `json:platform_fee`
	Discount    string `json:promotions`
	OrderTotal  string `json:earning`
}

type EtsyOrder struct {
	OrderId              string `json:order_id`
	TransactionId        string `json:transaction_id`
	ProductName          string `json:product_name`
	Quantity             string `json:quantity`
	Price                string `json:price`
	Personalization      string `json:personalization`
	ProductType          string `json:product_type`
	Personalization_Note string `json:personalization_note`
}

// func FindOrderInfor(t string, rs *EtsyOrder)  {

// 	//extracting OrderId
// 	pattern := regexp.MustCompile("Order ID:\\s+(\\S+)")
// 	rs.OrderId = pattern.FindString(t)
// 	rs.OrderId = rs.OrderId[10:len(rs.OrderId)]

// 	//extracting Ship By
// 	pattern = regexp.MustCompile("Ship by:\\s+([\\d/]+)")
// 	rs.ShipBy = pattern.FindString(t)
// 	rs.ShipBy = rs.ShipBy[9:len(rs.ShipBy)]

// 	//extracting item name
// 	pattern = regexp.MustCompile("Item:\\s+(\\S+)")
// 	rs.Item = pattern.FindString(t)
// 	rs.Item = rs.Item[6:len(rs.Item)]

// 	//extracting condition
// 	pattern = regexp.MustCompile("Condition:\\s+(\\S+)")
// 	rs.Condition = pattern.FindString(t)
// 	rs.Condition = rs.Condition[11:len(rs.Condition)]

// 	//extracting SKU
// 	pattern = regexp.MustCompile("SKU:\\s+(\\S+)")
// 	rs.SKU = pattern.FindString(t)
// 	rs.SKU = rs.SKU[5:len(rs.SKU)]

// 	//extracting quantity
// 	pattern = regexp.MustCompile("Quantity:\\s+(\\d+)")
// 	rs.Quantity = pattern.FindString(t)
// 	rs.Quantity = rs.Quantity[10:len(rs.Quantity)]

// 	//extracting orderdate
// 	pattern = regexp.MustCompile("Order date:\\s+([\\d/]+)")
// 	rs.OrderDate = pattern.FindString(t)
// 	rs.OrderDate = rs.OrderDate[12:len(rs.OrderDate)]

// 	//extracting price
// 	pattern = regexp.MustCompile("Price:\\s+\\$(\\d+\\.\\d+)")
// 	rs.Price = pattern.FindString(t)
// 	rs.Price = rs.Price[7:len(rs.Price)]

// 	//extracting Tax
// 	pattern = regexp.MustCompile("Tax:\\s+\\$(\\d+\\.\\d+)")
// 	rs.Tax = pattern.FindString(t)
// 	rs.Tax = rs.Tax[5:len(rs.Tax)]

// 	//extracting Shipping
// 	pattern = regexp.MustCompile("Shipping:\\s+\\$(\\d+\\.\\d+)")
// 	rs.Shipping = pattern.FindString(t)
// 	rs.Shipping = rs.Shipping[10:len(rs.Shipping)]

// 	//extracting Promotion
// 	pattern = regexp.MustCompile("Promotions:\\s+-\\$(\\d+\\.\\d+)")
// 	rs.Promotions = pattern.FindString(t)
// 	rs.Promotions = rs.Promotions[12:len(rs.Promotions)]

// 	//extracting Amazregexp
// 	pattern = regexp.MustCompile("Amazon fees:\\s+-\\$(\\d+\\.\\d+)")
// 	rs.AmazonFee = pattern.FindString(t)
// 	rs.AmazonFee = rs.AmazonFee[13:len(rs.AmazonFee)]

// 	//extracting Marketplace Facilitator Tax
// 	pattern = regexp.MustCompile("Marketplace Facilitator Tax:\\s+-\\$(\\d+\\.\\d+)")
// 	rs.MarketPlaceFacilitatorTax = pattern.FindString(t)
// 	rs.MarketPlaceFacilitatorTax = rs.MarketPlaceFacilitatorTax[29:len(rs.MarketPlaceFacilitatorTax)]

// 	//extracting Your earnings
// 	pattern = regexp.MustCompile("Your earnings:\\s+\\$(\\d+\\.\\d+)")
// 	rs.YourEarning = pattern.FindString(t)
// 	rs.YourEarning = rs.YourEarning[15:len(rs.YourEarning)]

// 	return &rs
// }