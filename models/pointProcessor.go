package models

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func CalculatePoints(receipt Receipt) int64 {
	itemCount := len(receipt.Items)

	points := 0
	points += examineRetailerName(receipt.Retailer)
	points += examineTotal(receipt.Total)
	points += examineItemCount(itemCount)
	points += examinePurchaseDate(receipt.PurchaseDate)
	points += examineHour(receipt.PurchaseTime)

	for _, item := range receipt.Items {
		points += examineDescription(item.ShortDescription, item.Price)
	}

	calcPoints := int64(points)
	return calcPoints
}

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
		log.Printf("isMultipleOfQuarter: error parsing total: %v\n", err)
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
		log.Printf("examinePurchaseDate: error parsing purchaseDate: %v\n", err)
		return 0
	}

	day := date.Day()
	if day%2 == 1 {
		points += 6
	}
	return
}

func examineHour(purchaseTime string) (points int) {
	points = 0
	parsedTime, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		log.Printf("ExamineHour error parsing purchaseTime: %v\n", err)
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
			log.Printf("ExamineDescription: error parsing description:  %v\n", err)
			return
		}

		points = int(math.Ceil(num * 0.2))
	}
	return
}
