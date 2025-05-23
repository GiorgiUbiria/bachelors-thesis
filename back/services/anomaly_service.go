package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type AnomalyPredictRequest struct {
	Features [][]float64 `json:"features"`
}

type AnomalyPredictResponse struct {
	Prediction []int `json:"prediction"`
}

type RecommendationRequest struct {
	UserID uint `json:"user_id"`
}

type RecommendationResponse struct {
	Recommendations []uint `json:"recommendations"`
}

type UserBehaviorRequest struct {
	UserID uint `json:"user_id"`
}

type UserBehaviorResponse struct {
	Cluster       int       `json:"cluster"`
	ClusterCenter []float64 `json:"cluster_center"`
}

type MLServiceConfig struct {
	BaseURL      string
	Timeout      time.Duration
	MaxRetries   int
	RetryDelay   time.Duration
	FallbackMode bool
}

var mlConfig = &MLServiceConfig{
	BaseURL:      getEnv("ML_SERVICE_URL", "http://localhost:5000"),
	Timeout:      30 * time.Second,
	MaxRetries:   3,
	RetryDelay:   1 * time.Second,
	FallbackMode: getEnv("ML_FALLBACK_MODE", "true") == "true",
}

// DetectAnomaly detects anomalies in request features with retry logic and fallback
var DetectAnomaly = func(features [][]float64) (int, error) {
	if len(features) == 0 {
		return 0, errors.New("no features provided")
	}

	requestBody, err := json.Marshal(AnomalyPredictRequest{Features: features})
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: mlConfig.Timeout}
	url := fmt.Sprintf("%s/analyze/request", mlConfig.BaseURL)

	var lastErr error
	for attempt := 0; attempt <= mlConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(mlConfig.RetryDelay * time.Duration(attempt))
		}

		resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			lastErr = fmt.Errorf("attempt %d failed: %w", attempt+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("attempt %d: ML service returned status %d", attempt+1, resp.StatusCode)
			continue
		}

		var result AnomalyPredictResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to decode response: %w", attempt+1, err)
			continue
		}

		if len(result.Prediction) == 0 {
			lastErr = fmt.Errorf("attempt %d: no prediction returned", attempt+1)
			continue
		}

		return result.Prediction[0], nil
	}

	// Fallback mechanism
	if mlConfig.FallbackMode {
		return fallbackAnomalyDetection(features[0]), nil
	}

	return 0, fmt.Errorf("all attempts failed, last error: %w", lastErr)
}

// GetRecommendations gets product recommendations for a user
func GetRecommendations(userID uint) ([]uint, error) {
	requestBody, err := json.Marshal(RecommendationRequest{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: mlConfig.Timeout}
	url := fmt.Sprintf("%s/recommend", mlConfig.BaseURL)

	var lastErr error
	for attempt := 0; attempt <= mlConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(mlConfig.RetryDelay * time.Duration(attempt))
		}

		resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			lastErr = fmt.Errorf("attempt %d failed: %w", attempt+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("attempt %d: ML service returned status %d", attempt+1, resp.StatusCode)
			continue
		}

		var result RecommendationResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to decode response: %w", attempt+1, err)
			continue
		}

		return result.Recommendations, nil
	}

	// Fallback mechanism
	if mlConfig.FallbackMode {
		return fallbackRecommendations(userID), nil
	}

	return nil, fmt.Errorf("all attempts failed, last error: %w", lastErr)
}

// AnalyzeUserBehavior analyzes user behavior patterns
func AnalyzeUserBehavior(userID uint) (*UserBehaviorResponse, error) {
	requestBody, err := json.Marshal(UserBehaviorRequest{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: mlConfig.Timeout}
	url := fmt.Sprintf("%s/analyze/user", mlConfig.BaseURL)

	var lastErr error
	for attempt := 0; attempt <= mlConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(mlConfig.RetryDelay * time.Duration(attempt))
		}

		resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			lastErr = fmt.Errorf("attempt %d failed: %w", attempt+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("attempt %d: ML service returned status %d", attempt+1, resp.StatusCode)
			continue
		}

		var result UserBehaviorResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to decode response: %w", attempt+1, err)
			continue
		}

		return &result, nil
	}

	// Fallback mechanism
	if mlConfig.FallbackMode {
		return &UserBehaviorResponse{
			Cluster:       0,
			ClusterCenter: []float64{0, 0, 0, 0},
		}, nil
	}

	return nil, fmt.Errorf("all attempts failed, last error: %w", lastErr)
}

// TriggerMLRetraining triggers retraining of ML models
func TriggerMLRetraining() error {
	client := &http.Client{Timeout: mlConfig.Timeout * 2} // Longer timeout for training
	url := fmt.Sprintf("%s/train", mlConfig.BaseURL)

	var lastErr error
	for attempt := 0; attempt <= mlConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(mlConfig.RetryDelay * time.Duration(attempt))
		}

		resp, err := client.Post(url, "application/json", nil)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d failed: %w", attempt+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("attempt %d: ML service returned status %d", attempt+1, resp.StatusCode)
			continue
		}

		return nil
	}

	return fmt.Errorf("all retraining attempts failed, last error: %w", lastErr)
}

// CheckMLServiceHealth checks if the ML service is available
func CheckMLServiceHealth() error {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/health", mlConfig.BaseURL)

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("ML service health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ML service health check returned status %d", resp.StatusCode)
	}

	return nil
}

// Fallback functions for when ML service is unavailable

// fallbackAnomalyDetection provides simple rule-based anomaly detection
func fallbackAnomalyDetection(features []float64) int {
	if len(features) < 8 {
		return 1 // Normal by default
	}

	responseTime := features[0]
	requestSize := features[1]
	errorCount := features[2]
	statusCode := features[6]

	// Simple rules for anomaly detection
	if responseTime > 5000 { // > 5 seconds
		return -1
	}
	if requestSize > 10000 { // > 10KB
		return -1
	}
	if errorCount > 0 && statusCode >= 400 {
		return -1
	}

	return 1 // Normal
}

// fallbackRecommendations provides simple popularity-based recommendations
func fallbackRecommendations(userID uint) []uint {
	// This would typically query for popular products
	// For now, return some default product IDs
	return []uint{1, 2, 3, 4, 5}
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
