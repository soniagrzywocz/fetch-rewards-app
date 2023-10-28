package api

import (
	"encoding/json"
	"fetch/models"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
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

func ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
	// ADD PROPER ERROR MESSAGES
	var receipt models.Receipt
	if err := ParseJSONRequest(r, &receipt); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	id, err := models.SaveReceipt(receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Id string `json:"id"`
	}{
		Id: id,
	}

	jsonResponse, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		fmt.Printf("Error marshaling response")
	}
	w.Header().Set("Content-Type", "application/json")
	// Respond with the created receipt's ID
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("New Receipt processed")
}

func GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	receiptId := vars["id"]

	points, exists := models.GetPoints(receiptId)
	if exists {
		response := struct {
			Points int64 `json:"points"`
		}{
			Points: points,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Error marshaling response")
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		response := struct {
			Description string
		}{
			Description: "No receipt found for that id",
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Error marshaling response")
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
	}
}

func ParseJSONRequest(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}
