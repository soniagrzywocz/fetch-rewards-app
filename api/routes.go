package api

import (
	"encoding/json"
	"fetch/models"
	"io"
	"log"
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
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		log.Printf("The receipt format is invalid")
		return
	}

	id, err := models.SaveReceipt(receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Panicf("Error saving the receipt: %v", err)
		return
	}

	response := struct {
		Id string `json:"id"`
	}{
		Id: id,
	}

	jsonResponse, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		log.Printf("ProcessReceiptHandler: error marshaling response")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	// fmt.Printf("New Receipt processed")
	log.Printf("New Receipt processed with an ID: %v\n", id)
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
			log.Printf("GetPointsHandler: error marshaling response in a happy path")
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		log.Printf("GET request for receiptID:  %v ,points:  %d", receiptId, points)
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
			log.Printf("GetPointsHandler: error marshaling response in a negative path")
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
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
