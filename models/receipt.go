package models

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type Receipt struct {
	Id           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var ReceiptList *map[string]int64

func init() {
	ReceiptList = &map[string]int64{}
	*ReceiptList = make(map[string]int64)
}

func SaveReceipt(receipt Receipt) (id string, err error) {
	newReceipt := receipt
	newReceipt.Id = uuid.New().String()
	newReceiptPoints := int64(calculatePoints(&newReceipt))
	(*ReceiptList)[newReceipt.Id] = newReceiptPoints
	return newReceipt.Id, nil
}

func GetPoints(id string) (points int64, ok bool) {
	points, ok = (*ReceiptList)[id]
	return
}

func calculatePoints(receipt *Receipt) (points int) {
	itemCount := len(receipt.Items)

	points = 0
	points += examineRetailerName(receipt.Retailer)
	points += examineTotal(receipt.Total)
	points += examineItemCount(itemCount)
	points += examinePurchaseDate(receipt.PurchaseDate)
	points += examineHour(receipt.PurchaseTime)

	for _, item := range receipt.Items {
		points += examineDescription(item.ShortDescription, item.Price)
	}

	return points
}

// helper functions to make the calculations easier
func examineRetailerName(retailerName string) (points int) {
	points = 0
	for _, char := range retailerName {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points++
		}
	}
	return points
}

func examineTotal(total string) (points int) {

	points = 0
	if isRoundTotal(total) {
		points += 50
	}
	if isMultipleOfQuarter(total) {
		points += 25
	}
	return
}

func isRoundTotal(total string) bool {
	total = strings.TrimSpace(total)
	return strings.HasSuffix(total, ".00")
}

func isMultipleOfQuarter(total string) bool {
	num, err := strconv.ParseFloat(total, 64)
	if err != nil {
		// ADD LOGS
		return false
	}

	return num/0.25 == float64(int(num/0.25))
}

func examineItemCount(count int) (points int) {
	points = 0
	if count%2 == 0 {
		points = (count / 2) * 5
	} else {
		points = (count - 1) / 2 * 5
	}
	return
}

func examinePurchaseDate(purchaseDate string) (points int) {
	points = 0

	date, err := time.Parse("2006-01-02", purchaseDate)
	if err != nil {
		// ADD PROPER LOGS
		return 0
	}

	day := date.Day()
	if day%2 == 1 {
		points += 6
	}
	return
}

func examineHour(purchaseTime string) (points int) {
	//10pts if the time of purchase is after 2pm and before 4pm
	points = 0
	parsedTime, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		// ADD PROPER LOGS
		fmt.Printf("Error parsing hour: %v\n", err)
		return 0
	}

	twoPM := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 14, 0, 0, 0, parsedTime.Location())
	fourPM := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 16, 0, 0, 0, parsedTime.Location())

	if parsedTime.After(twoPM) && parsedTime.Before(fourPM) {
		points += 10

	}
	return
}

func examineDescription(desc string, price string) (points int) {
	points = 0
	trimmedDesc := strings.TrimSpace(desc)
	calcLen := len(trimmedDesc)

	if calcLen%3 == 0 {
		num, err := strconv.ParseFloat(price, 64)
		if err != nil {
			// ADD PROPER LOGS
			fmt.Printf("Error converting desc: %v\n", err)
			return
		}

		points = int(math.Ceil(num * 0.2))
	}
	return
}
