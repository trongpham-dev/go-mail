package etsy

import (
	"regexp"
	"strconv"
)

type EtsyOrderRecord struct {
	Records []EtsyFieldOrder `json:records`
}

type EtsyFieldOrder struct {
	Fields EtsyOrder `json:fields`
}

type EtsyOrderDetailRecord struct {
	Records []EtsyFieldOrder `json:records`
}

type EtsyFieldOrderDetail struct {
	Fields EtsyOrder `json:fields`
}

type EtsyOrderDetail struct {
	OrderId     string  `json:order_id`
	OrderDate   string  `json:order_date`
	ItemTotal   float32 `json:revenue`
	ShippingFee float32 `json:shipping_fee`
	SalesTax    float32 `json:platform_fee`
	Discount    float32 `json:promotions`
	OrderTotal  float32 `json:earning`
}

type EtsyOrder struct {
	OrderId              string  `json:order_id`
	TransactionId        string  `json:transaction_id`
	ProductName          string  `json:product_name`
	Quantity             uint32  `json:quantity`
	Price                float32 `json:price`
	Personalization      string  `json:omitted;`
	ProductType          string  `json:product_type`
	Personalization_Note string  `json:personalization_note`
}

func ExtractEtsyOrder(t string, rs *EtsyOrder) {
	//TransactionId
	pattern := regexp.MustCompile("Transaction ID:\\s+(\\S+)")
	match := pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		rs.TransactionId = match[1]
	}

	//ProductName
	pattern = regexp.MustCompile("Item:\\s+(.+)")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		rs.ProductName = match[1]
	}

	pattern = regexp.MustCompile("Quantity:\\s+(\\d+)")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		quantity, _ := strconv.Atoi(match[1])
		rs.Quantity = uint32(quantity)
	}

	pattern = regexp.MustCompile("Price:\\s+\\-?\\$(\\d+\\.\\d+)")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		price, _ := strconv.ParseFloat(match[1], 32)
		rs.Price = float32(price)
	}

	pattern = regexp.MustCompile("Personalization:\\s+(\\S+)")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		if match[1] != "" {
			rs.ProductType = "Personalization"
			rs.Personalization_Note = match[1]
		} else {
			rs.ProductType = "Normal"
		}
	}

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
// 	pattern = regexp.MustCompile("Amazon fees:\\s+-\\$(\\d+\\.\regexp
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
