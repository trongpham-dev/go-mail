package etsy

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-mail/component"
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
	OrderId  string `json:"order_id"`
	ShopName string `json:"shop_name"`
	// CustMail      string  `json:"customer_email"`
	// Address       string  `json:"address"`
	TransactionId        string  `json:"transaction_id"`
	OrderDate            string  `json:"date"`
	ProductName          string  `json:"product_name"`
	Quantity             uint32  `json:"quantity"`
	Price                float32 `json:"price"`
	Personalization      string  `json:"-"`
	ProductType          string  `json:"-"`
	Personalization_Note string  `json:"personalization_note"`
}

type EtsyOrderShippingRecord struct {
	Records []EtsyFieldOrderShipping `json:"records"`
}

type EtsyFieldOrderShipping struct {
	Fields EtsyOrderShipping `json:"fields"`
}

type EtsyOrderShipping struct {
	OrderId  string `json:"order_id"`
	Email    string `json:"email"`
	CustMail string `json:"customer_email"`
	CustName string `json:"customer_name"`
	Road     string `json:"road"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
	Country  string `json:"country"`
}

func NewEtsyOrderShippingRecord(e []EtsyFieldOrderShipping) *EtsyOrderShippingRecord {
	return &EtsyOrderShippingRecord{
		Records: e,
	}
}

func NewEtsyFieldOrderShipping(e EtsyOrderShipping) EtsyFieldOrderShipping {
	return EtsyFieldOrderShipping{
		Fields: e,
	}
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
			return
		}
		rs.OrderTotal = float32(dollarValue)
	}

}

type Response struct {
	Code int `json:"code"`
	// Msg   string `json:"msg"`
	// Error struct {
	//     Message string `json:"message"`
	//     LogID   string `json:"log_id"`
	// } `json:"error"`
}

func CreateEtsyOrder(appCtx component.AppContext, r *EtsyOrderRecord) error {
	var res *http.Response
	var response Response
	for i := 1; i <= 3; i++ {
		appTkn, err := appCtx.GetAppToken().GetAppAccessToken()

		if err != nil {
			return err
		}
		appCtx.SetAppToken(appTkn)

		client := &http.Client{}
		postBody, _ := json.Marshal(r)
		responseBody := bytes.NewBuffer(postBody)
		req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/bitable/v1/apps/KhcHb8CvtajCzUsTNBYlzxEtgId/tables/tbls3bFafW446965/records/batch_create", responseBody)
		req.Header.Set("Authorization", "Bearer "+appCtx.GetAppToken().AppAccessToken)
		req.Header.Add("Content-Type", "application/json")

		if err != nil {
			return err
		}

		res, err = client.Do(req)

		if err != nil {
			return err
		}

		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return err
		}

		if response.Code == 0 {
			break
		}
		defer res.Body.Close()
	}

	if response.Code != 0 {

		return errors.New("can not create order!")
	}
	return nil
}

func CreateEtsyOrderDetail(appCtx component.AppContext, r *EtsyOrderDetailRecord) error {
	var res *http.Response
	var response Response
	for i := 1; i <= 3; i++ {
		appTkn, err := appCtx.GetAppToken().GetAppAccessToken()

		if err != nil {
			return err
		}
		appCtx.SetAppToken(appTkn)

		client := &http.Client{}
		postBody, _ := json.Marshal(r)
		responseBody := bytes.NewBuffer(postBody)
		req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/bitable/v1/apps/KhcHb8CvtajCzUsTNBYlzxEtgId/tables/tblwPkRqnAVoNwc8/records/batch_create", responseBody)
		req.Header.Set("Authorization", "Bearer "+appCtx.GetAppToken().AppAccessToken)
		req.Header.Add("Content-Type", "application/json")

		if err != nil {
			return err
		}

		res, err = client.Do(req)

		if err != nil {
			log.Println(err)
			return err
		}

		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return err
		}

		if response.Code == 0 {
			break
		}
	}
	defer res.Body.Close()

	if response.Code != 0 {

		return errors.New("can not create order!")
	}
	return nil
}

func CreateEtsyOrderShipping(appCtx component.AppContext, r *EtsyOrderShippingRecord) error {
	var res *http.Response
	var response Response
	for i := 1; i <= 3; i++ {
		appTkn, err := appCtx.GetAppToken().GetAppAccessToken()

		if err != nil {
			return err
		}
		appCtx.SetAppToken(appTkn)

		client := &http.Client{}
		postBody, _ := json.Marshal(r)
		responseBody := bytes.NewBuffer(postBody)
		req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/bitable/v1/apps/KhcHb8CvtajCzUsTNBYlzxEtgId/tables/tblFiG4gWcGG5a9A/records/batch_create", responseBody)
		req.Header.Set("Authorization", "Bearer "+appCtx.GetAppToken().AppAccessToken)
		req.Header.Add("Content-Type", "application/json")

		if err != nil {
			return err
		}

		res, err = client.Do(req)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		var response Response
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return err
		}

		if response.Code == 0 {
			break
		}
	}

	defer res.Body.Close()

	if response.Code != 0 {

		return errors.New("can not create order!")
	}

	return nil
}
