package api

import (
	"encoding/json"
	"fetch/models"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
)

func SetRoutes(r *mux.Router) {
	r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := parseJSONRequest(r, &receipt); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if vErr := validateReq(receipt); vErr != nil {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		log.Printf("The receipt format is invalid")
		return
	}

	id, err := models.ProcessReceipt(receipt)
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
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	log.Printf("New Receipt processed with an ID: %v\n", id)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	receiptId := vars["id"]

	points, exists := models.SharedReceiptList.RetrievePoints(receiptId)
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
			return
		}
		log.Printf("GET request for receiptID:  %v ,points:  %d", receiptId, points)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		notFoundError := "No receipt found for that id"
		jsonResponse, err := json.Marshal(notFoundError)
		if err != nil {
			log.Printf("GetPointsHandler: error marshaling response in a negative path")
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
	}
}

func parseJSONRequest(r *http.Request, v interface{}) error {
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

func validateReq(receipt models.Receipt) error {
	validate := validator.New()

	if err := validate.Struct(receipt); err != nil {
		return err
	}
	return nil
}
