package etsy

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type EtsyOrderRecord struct {
	Records []EtsyFieldOrder `json:"records"`
}

type EtsyFieldOrder struct {
	Fields EtsyOrder `json:"fields"`
}

type EtsyOrderDetailRecord struct {
	Records []EtsyFieldOrderDetail `json:"records"`
}

type EtsyFieldOrderDetail struct {
	Fields EtsyOrderDetail `json:"fields"`
}

type EtsyOrderDetail struct {
	OrderId     string  `json:"order_id"`
	OrderDate   string  `json:"date"`
	ItemTotal   float32 `json:"revenue"`
	ShippingFee float32 `json:"shipping_fee"`
	Tax         float32 `json:"tax"`
	SalesTax    float32 `json:"platform_fee"`
	Discount    float32 `json:"promotions"`
	OrderTotal  float32 `json:"earning"`
}

type EtsyOrder struct {
	OrderId              string  `json:"order_id"`
	Email                string  `json:"email"`
	CustMail			string `json:"customer_email"`
	Address				string `json:"address"`
	TransactionId        string  `json:"transaction_id"`
	OrderDate            string  `json:"date"`
	ProductName          string  `json:"product_name"`
	Quantity             uint32  `json:"quantity"`
	Price                float32 `json:"price"`
	Personalization      string  `json:"-"`
	ProductType          string  `json:"product_type"`
	Personalization_Note string  `json:"personalization_note"`
}

type appToken struct {
	appAccessToken string `json:"app_access_token"`
	Code int 	`json:"code"`
	expire         int    `json:"expire"`
	Message string `json:"msg"`
	ternantAccessToken string `json:"tenant_access_token"`
}

type appInfo struct {
	appId     string `json:"app_id"`
	appSecret string `json:"app_secret"`
}

func NewEtsyOrderDetailRecord(e []EtsyFieldOrderDetail) *EtsyOrderDetailRecord {
	return &EtsyOrderDetailRecord{
		Records: e,
	}
}

func NewEtsyFieldOrderDetail(e EtsyOrderDetail) EtsyFieldOrderDetail {
	return EtsyFieldOrderDetail{
		Fields: e,
	}
}

func NewEtsyOrderRecord(e []EtsyFieldOrder) *EtsyOrderRecord {
	return &EtsyOrderRecord{
		Records: e,
	}
}

func NewEtsyFieldOrder(e EtsyOrder) EtsyFieldOrder {
	return EtsyFieldOrder{
		Fields: e,
	}
}

func ExtractEtsyOrder(t string, rs *EtsyOrder) {
	//TransactionId
	pattern := regexp.MustCompile("Transaction ID:.*?(\\S+)")
	match := pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		rs.TransactionId = match[1]
	}

	//ProductName
	pattern = regexp.MustCompile("Item:\\s+(.*?)\n")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		rs.ProductName = match[1]
	}

	pattern = regexp.MustCompile("Quantity:.*?(\\d+)")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		quantity, _ := strconv.Atoi(match[1])
		rs.Quantity = uint32(quantity)
	}

	pattern = regexp.MustCompile("(?:^|\\n)Price:.*?\\b(?:[A-Z]{2}\\$|\\$|\\bUSD\\s*)?([\\d,]+(?:\\.\\d{2})?)\\s*\\$?\\b")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 1 {
		dollarStr := match[1]
		dollarStr = regexp.MustCompile(",").ReplaceAllString(dollarStr, "")
		dollarValue, err := strconv.ParseFloat(dollarStr, 64)
		if err != nil {
			log.Fatal(err)
			return
		}
		rs.Price = float32(dollarValue)
	}

	pattern = regexp.MustCompile("Personalization:.*?(.+)")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 0 {
		if strings.TrimSpace(match[1]) != "" {
			rs.ProductType = "Personalization"
			rs.Personalization_Note = match[1]
		} else {
			rs.ProductType = "Normal"
		}
	} else {
		rs.ProductType = "Normal"
	}

}

func ExtractEtsyOrderDetail(t string, rs *EtsyOrderDetail) {
	//TransactionId
	pattern := regexp.MustCompile("(?:^|\\n)Item total:.*?\\b(?:[A-Z]{2}\\$|\\$|\\bUSD\\s*)?([\\d,]+(?:\\.\\d{2})?)\\s*\\$?\\b")
	match := pattern.FindStringSubmatch(t)
	if len(match) > 1 {
		dollarStr := match[1]
		dollarStr = regexp.MustCompile(",").ReplaceAllString(dollarStr, "")
		dollarValue, err := strconv.ParseFloat(dollarStr, 64)
		if err != nil {
			log.Fatal(err)
			return
		}
		rs.ItemTotal = float32(dollarValue)
	}

	//ProductName
	pattern = regexp.MustCompile("(?:^|\\n)Discount:.*?\\b(?:[A-Z]{2}\\$|\\$|\\bUSD\\s*)?([\\d,]+(?:\\.\\d{2})?)\\s*\\$?\\b")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 1 {
		dollarStr := match[1]
		dollarStr = regexp.MustCompile(",").ReplaceAllString(dollarStr, "")
		dollarValue, err := strconv.ParseFloat(dollarStr, 64)
		if err != nil {
			log.Fatal(err)
			return
		}
		rs.Discount = float32(dollarValue)
	}

	pattern = regexp.MustCompile("(?:^|\\n)Shipping:.*?\\b(?:[A-Z]{2}\\$|\\$|\\bUSD\\s*)?([\\d,]+(?:\\.\\d{2})?)\\s*\\$?\\b")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 1 {
		dollarStr := match[1]
		dollarStr = regexp.MustCompile(",").ReplaceAllString(dollarStr, "")
		dollarValue, err := strconv.ParseFloat(dollarStr, 64)
		if err != nil {
			log.Fatal(err)
			return
		}
		rs.ShippingFee = float32(dollarValue)
	}

	pattern = regexp.MustCompile("(?:^|\\n)Tax:.*?\\b(?:[A-Z]{2}\\$|\\$|\\bUSD\\s*)?([\\d,]+(?:\\.\\d{2})?)\\s*\\$?\\b")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 1 {
		dollarStr := match[1]
		dollarStr = regexp.MustCompile(",").ReplaceAllString(dollarStr, "")
		dollarValue, err := strconv.ParseFloat(dollarStr, 64)
		if err != nil {
			log.Fatal(err)
			return
		}
		rs.Tax = float32(dollarValue)
	}

	pattern = regexp.MustCompile("(?:^|\\n)Sales tax:.*?\\b(?:[A-Z]{2}\\$|\\$|\\bUSD\\s*)?([\\d,]+(?:\\.\\d{2})?)\\s*\\$?\\b")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 1 {
		dollarStr := match[1]
		dollarStr = regexp.MustCompile(",").ReplaceAllString(dollarStr, "")
		dollarValue, err := strconv.ParseFloat(dollarStr, 64)
		if err != nil {
			log.Fatal(err)
			return
		}
		rs.SalesTax = float32(dollarValue)
	}

	pattern = regexp.MustCompile("(?:^|\\n)Order total:.*?\\b(?:[A-Z]{2}\\$|\\$|\\bUSD\\s*)?([\\d,]+(?:\\.\\d{2})?)\\s*\\$?\\b")
	match = pattern.FindStringSubmatch(t)
	if len(match) > 1 {
		dollarStr := match[1]
		dollarStr = regexp.MustCompile(",").ReplaceAllString(dollarStr, "")
		dollarValue, err := strconv.ParseFloat(dollarStr, 64)
		if err != nil {
			log.Fatal(err)
			return
		}
		rs.OrderTotal = float32(dollarValue)
	}

}

func GetAppAccessToken(a *appToken) error {
	for {
		appTkn := appInfo{
			appId:     "cli_a4b0a37dd8f8d02f",
			appSecret: "ziCKGTkVuprRLpoV17rrzcaCkjZV5lBq",
		}

		client := &http.Client{}
		postBody, _ := json.Marshal(appTkn)
		responseBody := bytes.NewBuffer(postBody)
		req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/auth/v3/app_access_token/internal", responseBody)
		req.Header.Add("Content-Type", "application/json")

		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := io.ReadAll(res.Body)
		json.Unmarshal(body, a)

		if a.expire > 1820 {
			break
		}
		defer res.Body.Close()
	}
	return nil
}

func CreateEtsyOrder(r *EtsyOrderRecord) error {
	// appToken := appToken{}
	// err := GetAppAccessToken(&appToken)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	client := &http.Client{}
	postBody, _ := json.Marshal(r)
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/bitable/v1/apps/KhcHb8CvtajCzUsTNBYlzxEtgId/tables/tbls3bFafW446965/records/batch_create", responseBody)
	req.Header.Set("Authorization", "Bearer t-g206587yQLLAGY5OVDOXKDCWQDKOEF3P7AG6KUVL")
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer res.Body.Close()

	return nil
}

func CreateEtsyOrderDetail(r *EtsyOrderDetailRecord) error {
	// appToken := appToken{}
	// err := GetAppAccessToken(&appToken)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	client := &http.Client{}
	postBody, _ := json.Marshal(r)
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/bitable/v1/apps/KhcHb8CvtajCzUsTNBYlzxEtgId/tables/tblwPkRqnAVoNwc8/records/batch_create", responseBody)
	req.Header.Set("Authorization", "Bearer t-g206587yQLLAGY5OVDOXKDCWQDKOEF3P7AG6KUVL")
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer res.Body.Close()

	return nil
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
