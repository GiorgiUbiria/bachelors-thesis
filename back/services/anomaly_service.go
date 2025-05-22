package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type AnomalyPredictRequest struct {
	Features [][]float64 `json:"features"`
}

type AnomalyPredictResponse struct {
	Prediction []int `json:"prediction"`
}

var AnomalyAPIURL = "http://localhost:5000/predict" // Set to your Python API URL

var DetectAnomaly = func(features [][]float64) (int, error) {
	requestBody, err := json.Marshal(AnomalyPredictRequest{Features: features})
	if err != nil {
		return 0, err
	}
	resp, err := http.Post(AnomalyAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("anomaly API returned non-200 status")
	}
	var result AnomalyPredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	if len(result.Prediction) == 0 {
		return 0, errors.New("no prediction returned")
	}
	return result.Prediction[0], nil
}
