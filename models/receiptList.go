package models

import (
	"github.com/google/uuid"
)

type ReceiptList struct {
	receipts map[string]Receipt
}

var SharedReceiptList *ReceiptList

func init() {
	SharedReceiptList = &ReceiptList{
		receipts: make(map[string]Receipt),
	}
}

func (r *ReceiptList) SaveReceipt(receipt Receipt) (id string, err error) {
	receipt.Id = uuid.New().String()
	r.receipts[receipt.Id] = receipt
	return receipt.Id, nil
}

func (r *ReceiptList) RetrievePoints(id string) (points int64, exists bool) {
	if receipt, ok := r.receipts[id]; ok {
		if receipt.IsCalculated == true {
			points = receipt.Points
		} else {
			points = CalculatePoints(receipt)
			receipt.Points = points
			receipt.IsCalculated = true
			r.receipts[id] = receipt
		}
	} else {
		return 0, false
	}
	return points, true
}
