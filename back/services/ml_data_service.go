package services

import (
	"time"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
)

// MLTrainingData represents the complete training data structure for ML models
type MLTrainingData struct {
	RequestLogs    []RequestLogFeatures   `json:"request_logs"`
	UserActivities []UserActivityFeatures `json:"user_activities"`
	Purchases      []PurchaseFeatures     `json:"purchases"`
	SalesTrends    []SalesTrendFeatures   `json:"sales_trends"`
	UserItemMatrix [][]float64            `json:"user_item_matrix"`
	Timestamp      time.Time              `json:"timestamp"`
}

// RequestLogFeatures represents features extracted from request logs for anomaly detection
type RequestLogFeatures struct {
	ResponseTime float64 `json:"response_time"`
	RequestSize  float64 `json:"request_size"`
	ErrorCount   float64 `json:"error_count"`
	IPHash       float64 `json:"ip_hash"`
	MethodCode   float64 `json:"method_code"`
	PathLength   float64 `json:"path_length"`
	StatusCode   float64 `json:"status_code"`
	UserAgentLen float64 `json:"user_agent_len"`
}

// UserActivityFeatures represents aggregated user activity features for clustering
type UserActivityFeatures struct {
	UserID        uint    `json:"user_id"`
	LoginCount    float64 `json:"login_count"`
	PurchaseCount float64 `json:"purchase_count"`
	CartCount     float64 `json:"cart_count"`
	FavoriteCount float64 `json:"favorite_count"`
	ViewCount     float64 `json:"view_count"`
	ClickCount    float64 `json:"click_count"`
	SearchCount   float64 `json:"search_count"`
	SessionCount  float64 `json:"session_count"`
	AvgSessionDur float64 `json:"avg_session_duration"`
}

// PurchaseFeatures represents purchase data for recommendation training
type PurchaseFeatures struct {
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	Rating    float64 `json:"rating,omitempty"`
}

// SalesTrendFeatures represents sales data for trend analysis
type SalesTrendFeatures struct {
	Date      time.Time `json:"date"`
	ProductID uint      `json:"product_id"`
	Category  string    `json:"category"`
	Sales     float64   `json:"sales"`
	Revenue   float64   `json:"revenue"`
	Units     float64   `json:"units"`
}

// GetMLTrainingData aggregates all data needed for ML model training
func GetMLTrainingData() (*MLTrainingData, error) {
	data := &MLTrainingData{
		Timestamp: time.Now(),
	}

	// Get request log features
	requestFeatures, err := getRequestLogFeatures()
	if err != nil {
		return nil, err
	}
	data.RequestLogs = requestFeatures

	// Get user activity features
	userFeatures, err := getUserActivityFeatures()
	if err != nil {
		return nil, err
	}
	data.UserActivities = userFeatures

	// Get purchase features
	purchaseFeatures, err := getPurchaseFeatures()
	if err != nil {
		return nil, err
	}
	data.Purchases = purchaseFeatures

	// Get sales trend features
	trendFeatures, err := getSalesTrendFeatures()
	if err != nil {
		return nil, err
	}
	data.SalesTrends = trendFeatures

	// Generate user-item matrix
	userItemMatrix, err := generateUserItemMatrix()
	if err != nil {
		return nil, err
	}
	data.UserItemMatrix = userItemMatrix

	return data, nil
}

// getRequestLogFeatures extracts features from request logs for anomaly detection
func getRequestLogFeatures() ([]RequestLogFeatures, error) {
	var logs []models.RequestLog

	// Get logs from last 7 days for training
	since := time.Now().AddDate(0, 0, -7)
	result := config.DB.Where("created_at >= ?", since).Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	features := make([]RequestLogFeatures, len(logs))
	for i, log := range logs {
		features[i] = RequestLogFeatures{
			ResponseTime: log.ResponseTime,
			RequestSize:  float64(len(log.Path) + len(log.UserAgent)), // Approximate request size
			ErrorCount:   getErrorCount(log.Status),
			IPHash:       hashIP(log.IP),
			MethodCode:   methodToCode(log.Method),
			PathLength:   float64(len(log.Path)),
			StatusCode:   float64(log.Status),
			UserAgentLen: float64(len(log.UserAgent)),
		}
	}

	return features, nil
}

// getUserActivityFeatures aggregates user activity data for clustering
func getUserActivityFeatures() ([]UserActivityFeatures, error) {
	var users []models.User
	result := config.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	features := make([]UserActivityFeatures, len(users))
	since := time.Now().AddDate(0, 0, -30) // Last 30 days

	for i, user := range users {
		// Count different activity types
		var activities []models.UserActivity
		config.DB.Where("user_id = ? AND created_at >= ?", user.ID, since).Find(&activities)

		activityCounts := make(map[string]int)
		for _, activity := range activities {
			activityCounts[activity.Type]++
		}

		// Count purchases
		var purchaseCount int64
		config.DB.Model(&models.Order{}).Where("user_id = ? AND created_at >= ?", user.ID, since).Count(&purchaseCount)

		// Count cart items
		var cartCount int64
		config.DB.Model(&models.CartItem{}).
			Joins("JOIN carts ON cart_items.cart_id = carts.id").
			Where("carts.user_id = ?", user.ID).Count(&cartCount)

		// Count favorites
		var favoriteCount int64
		config.DB.Model(&models.Favorite{}).Where("user_id = ?", user.ID).Count(&favoriteCount)

		features[i] = UserActivityFeatures{
			UserID:        user.ID,
			LoginCount:    float64(activityCounts["login"]),
			PurchaseCount: float64(purchaseCount),
			CartCount:     float64(cartCount),
			FavoriteCount: float64(favoriteCount),
			ViewCount:     float64(activityCounts["view"]),
			ClickCount:    float64(activityCounts["click"]),
			SearchCount:   float64(activityCounts["search"]),
			SessionCount:  float64(len(getUniqueSessions(activities))),
			AvgSessionDur: calculateAvgSessionDuration(activities),
		}
	}

	return features, nil
}

// getPurchaseFeatures extracts purchase data for recommendation training
func getPurchaseFeatures() ([]PurchaseFeatures, error) {
	var orders []models.Order
	result := config.DB.Preload("Items").Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	var features []PurchaseFeatures
	for _, order := range orders {
		for _, item := range order.Items {
			features = append(features, PurchaseFeatures{
				UserID:    order.UserID,
				ProductID: item.ProductID,
				Quantity:  float64(item.Quantity),
				Price:     item.Price,
				Rating:    5.0, // Default rating, can be enhanced with actual ratings
			})
		}
	}

	return features, nil
}

// getSalesTrendFeatures extracts sales data for trend analysis
func getSalesTrendFeatures() ([]SalesTrendFeatures, error) {
	var features []SalesTrendFeatures

	// Get sales data grouped by date and product
	rows, err := config.DB.Raw(`
		SELECT 
			DATE(orders.created_at) as date,
			order_items.product_id,
			products.category,
			COUNT(*) as sales,
			SUM(order_items.price * order_items.quantity) as revenue,
			SUM(order_items.quantity) as units
		FROM orders 
		JOIN order_items ON orders.id = order_items.order_id
		JOIN products ON order_items.product_id = products.id
		WHERE orders.created_at >= ?
		GROUP BY DATE(orders.created_at), order_items.product_id, products.category
		ORDER BY date DESC
	`, time.Now().AddDate(0, -6, 0)).Rows() // Last 6 months

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var feature SalesTrendFeatures
		err := rows.Scan(
			&feature.Date,
			&feature.ProductID,
			&feature.Category,
			&feature.Sales,
			&feature.Revenue,
			&feature.Units,
		)
		if err != nil {
			continue
		}
		features = append(features, feature)
	}

	return features, nil
}

// generateUserItemMatrix creates a user-item interaction matrix for collaborative filtering
func generateUserItemMatrix() ([][]float64, error) {
	// Get all users and products
	var users []models.User
	var products []models.Product

	config.DB.Find(&users)
	config.DB.Find(&products)

	if len(users) == 0 || len(products) == 0 {
		return [][]float64{}, nil
	}

	// Create user-product mapping
	userMap := make(map[uint]int)
	productMap := make(map[uint]int)

	for i, user := range users {
		userMap[user.ID] = i
	}
	for i, product := range products {
		productMap[product.ID] = i
	}

	// Initialize matrix
	matrix := make([][]float64, len(users))
	for i := range matrix {
		matrix[i] = make([]float64, len(products))
	}

	// Fill matrix with purchase data
	var orders []models.Order
	config.DB.Preload("Items").Find(&orders)

	for _, order := range orders {
		userIdx, userExists := userMap[order.UserID]
		if !userExists {
			continue
		}

		for _, item := range order.Items {
			productIdx, productExists := productMap[item.ProductID]
			if !productExists {
				continue
			}

			// Use quantity as interaction strength
			matrix[userIdx][productIdx] += float64(item.Quantity)
		}
	}

	return matrix, nil
}

// Helper functions
func getErrorCount(status int) float64 {
	if status >= 400 {
		return 1.0
	}
	return 0.0
}

func hashIP(ip string) float64 {
	// Simple hash function for IP addresses
	hash := 0.0
	for _, char := range ip {
		hash = hash*31 + float64(char)
	}
	return hash
}

func methodToCode(method string) float64 {
	switch method {
	case "GET":
		return 1.0
	case "POST":
		return 2.0
	case "PUT":
		return 3.0
	case "DELETE":
		return 4.0
	case "PATCH":
		return 5.0
	default:
		return 0.0
	}
}

func getUniqueSessions(activities []models.UserActivity) []string {
	sessions := make(map[string]bool)
	for _, activity := range activities {
		if activity.SessionID != "" {
			sessions[activity.SessionID] = true
		}
	}

	result := make([]string, 0, len(sessions))
	for session := range sessions {
		result = append(result, session)
	}
	return result
}

func calculateAvgSessionDuration(activities []models.UserActivity) float64 {
	if len(activities) == 0 {
		return 0.0
	}

	sessionTimes := make(map[string][]time.Time)
	for _, activity := range activities {
		if activity.SessionID != "" {
			sessionTimes[activity.SessionID] = append(sessionTimes[activity.SessionID], activity.CreatedAt)
		}
	}

	totalDuration := 0.0
	sessionCount := 0

	for _, times := range sessionTimes {
		if len(times) > 1 {
			// Sort times and calculate duration
			minTime := times[0]
			maxTime := times[0]
			for _, t := range times {
				if t.Before(minTime) {
					minTime = t
				}
				if t.After(maxTime) {
					maxTime = t
				}
			}
			duration := maxTime.Sub(minTime).Minutes()
			totalDuration += duration
			sessionCount++
		}
	}

	if sessionCount == 0 {
		return 0.0
	}
	return totalDuration / float64(sessionCount)
}
