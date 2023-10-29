package models

type Receipt struct {
	Id           string `json:"id"`
	Retailer     string `json:"retailer" validate:"required"`
	PurchaseDate string `json:"purchaseDate" validate:"required"`
	PurchaseTime string `json:"purchaseTime" validate:"required"`
	Items        []Item `json:"items"`
	Total        string `json:"total" validate:"required"`
	Points       int64
	IsCalculated bool
}

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required"`
	Price            string `json:"price" validate:"required"`
}

func ProcessReceipt(receipt Receipt) (id string, err error) {

	receiptId, saveErr := SharedReceiptList.SaveReceipt(receipt)
	if saveErr != nil {
		return "", saveErr
	}

	return receiptId, nil
}
