# ML Layer - Machine Learning Service

This directory contains the machine learning service implementation for the e-commerce platform, providing anomaly detection, user behavior analysis, recommendation systems, and trend prediction capabilities.

## 🚀 Features Implemented

### ✅ Phase 1: Critical Infrastructure (COMPLETED)

#### 1. Data Integration & Pipeline
- **Real API Integration**: Replaced mock data loaders with actual backend API calls
- **Data Validation**: Comprehensive input validation and schema checking
- **Fallback Mechanisms**: Graceful degradation when backend APIs are unavailable
- **Error Handling**: Robust error handling with logging and retry logic

#### 2. Manual ML Implementations
- **Manual Logistic Regression**: Complete implementation with regularization and feature scaling
- **Manual Decision Tree**: Full implementation with information gain, pruning, and feature importance
- **API Endpoints**: Exposed manual implementations via REST API
- **Comparison Framework**: Side-by-side comparison with sklearn implementations

#### 3. Periodic Retraining System
- **Automated Scheduler**: APScheduler-based retraining with configurable schedules
- **Model Versioning**: Timestamp-based model versioning with rollback capabilities
- **Performance Monitoring**: Automatic rollback on performance degradation
- **Manual Triggers**: API endpoints for manual retraining

#### 4. Model Version Management
- **Version Control**: Complete model lifecycle management
- **Metadata Tracking**: Performance metrics and training metadata storage
- **Rollback System**: Automatic and manual rollback capabilities
- **Cleanup Policies**: Automatic cleanup of old model versions

## 📁 Directory Structure

```
ml/
├── app.py                      # Main Flask application
├── requirements.txt            # Python dependencies
├── Dockerfile                  # Docker configuration
├── test_ml_service.py         # Test suite
├── README.md                  # This file
├── common/
│   ├── config.py              # Configuration management
│   ├── data_loader.py         # Real API data loading
│   ├── schema.py              # Data schemas
│   └── utils.py               # Utility functions
├── models/
│   ├── anomaly.py             # Anomaly detection models
│   ├── clustering.py          # User clustering models
│   ├── recommendation.py      # Recommendation models
│   ├── trend.py               # Trend prediction models
│   └── version_manager.py     # Model versioning system
├── manual_impl/
│   ├── manual_logistic.py     # Manual logistic regression
│   └── manual_tree.py         # Manual decision tree
├── training/
│   ├── scheduler.py           # Automated retraining scheduler
│   ├── train_anomaly.py       # Anomaly model training
│   ├── train_clustering.py    # Clustering model training
│   ├── train_recommendation.py # Recommendation model training
│   └── train_trend.py         # Trend model training
├── inference/
│   └── anomaly_api.py         # FastAPI inference service
└── simulation/
    └── attack_simulation.py   # Attack simulation scripts
```

## 🔧 Installation & Setup

### Prerequisites
- Python 3.11+
- Docker (optional)
- Backend API service running

### Local Installation

1. **Install dependencies:**
```bash
cd ml
pip install -r requirements.txt
```

2. **Set environment variables:**
```bash
export API_URL="http://localhost:8080"  # Backend API URL
export API_KEY="your-api-key"           # Optional API key
export MODEL_DIR="models"               # Model storage directory
```

3. **Run the service:**
```bash
python app.py
```

### Docker Installation

```bash
cd ml
docker build -t ml-service .
docker run -p 5000:5000 -e API_URL="http://host.docker.internal:8080" ml-service
```

## 📚 API Documentation

### Core ML Endpoints

#### Health Check
```http
GET /health
```
Returns service health status and loaded models information.

#### Train Models
```http
POST /train
```
Trains all ML models using real data from backend API with versioning support.

#### Anomaly Detection
```http
POST /analyze/request
Content-Type: application/json

{
  "response_time": 0.1,
  "request_size": 1024,
  "error_count": 0
}
```

#### User Behavior Analysis
```http
POST /analyze/user
Content-Type: application/json

{
  "login_count": 10,
  "purchase_count": 5,
  "cart_count": 8,
  "favorite_count": 12
}
```

#### Recommendations
```http
POST /recommend
Content-Type: application/json

{
  "user_id": 123,
  "top_k": 5
}
```

### Manual ML Implementation Endpoints

#### Train Manual Logistic Regression
```http
POST /train/manual-logistic
Content-Type: application/json

{
  "features": [[1, 2, 3], [4, 5, 6]],
  "labels": [0, 1],
  "learning_rate": 0.01,
  "n_iterations": 1000,
  "regularization": "l2",
  "lambda_reg": 0.01
}
```

#### Train Manual Decision Tree
```http
POST /train/manual-tree
Content-Type: application/json

{
  "features": [[1, 2, 3], [4, 5, 6]],
  "labels": [0, 1],
  "max_depth": 5,
  "criterion": "gini"
}
```

#### Compare Models
```http
POST /compare/models
Content-Type: application/json

{
  "features": [[1, 2, 3], [4, 5, 6]],
  "labels": [0, 1],
  "test_size": 0.2
}
```

### Scheduler Management Endpoints

#### Get Scheduler Status
```http
GET /scheduler/status
```

#### Trigger Manual Retraining
```http
POST /scheduler/trigger
Content-Type: application/json

{
  "model_type": "anomaly"  // or "clustering", "recommendation", "trend", "all"
}
```

#### Pause/Resume Jobs
```http
POST /scheduler/pause/{job_id}
POST /scheduler/resume/{job_id}
```

### Model Version Management Endpoints

#### List Model Versions
```http
GET /models/versions?model_type=anomaly
```

#### Rollback Model
```http
POST /models/rollback
Content-Type: application/json

{
  "model_type": "anomaly",
  "version_id": "20241201_143022"
}
```

#### Update Performance Metrics
```http
POST /models/performance
Content-Type: application/json

{
  "model_type": "anomaly",
  "version_id": "20241201_143022",
  "performance_metrics": {
    "accuracy": 0.95,
    "precision": 0.92,
    "recall": 0.88
  }
}
```

#### Get Best Model Version
```http
GET /models/best?model_type=anomaly&metric=accuracy
```

#### Auto-Rollback Check
```http
POST /models/auto-rollback
Content-Type: application/json

{
  "model_type": "anomaly",
  "current_metrics": {"accuracy": 0.80},
  "threshold_metric": "accuracy",
  "degradation_threshold": 0.05
}
```

## 🔄 Retraining Schedules

The automated retraining system runs on the following schedule:

- **Anomaly Detection**: Every hour (if data/performance thresholds met)
- **User Clustering**: Daily at 2 AM
- **Recommendations**: Weekly on Sundays at 3 AM
- **Trend Prediction**: Monthly on 1st at 4 AM
- **Full Retraining**: Weekly on Mondays at 5 AM

## 🧪 Testing

Run the comprehensive test suite:

```bash
python test_ml_service.py
```

The test suite covers:
- Health checks
- Manual ML implementations
- Model comparison
- Scheduler functionality
- Model versioning
- Anomaly detection
- User clustering

## 📊 Model Performance

### Manual vs Sklearn Comparison

The manual implementations are designed to match sklearn performance:

| Algorithm | Manual Accuracy | Sklearn Accuracy | Performance Ratio |
|-----------|----------------|------------------|-------------------|
| Logistic Regression | ~0.85-0.90 | ~0.85-0.90 | ~1.0x |
| Decision Tree | ~0.80-0.85 | ~0.80-0.85 | ~1.0x |

### Features Implemented

#### Manual Logistic Regression
- ✅ Sigmoid activation with numerical stability
- ✅ L1 and L2 regularization
- ✅ Feature normalization
- ✅ Early stopping
- ✅ Cost history tracking
- ✅ Feature importance

#### Manual Decision Tree
- ✅ Gini impurity and entropy criteria
- ✅ Information gain calculation
- ✅ Pruning functionality
- ✅ Feature importance
- ✅ Tree visualization
- ✅ Configurable hyperparameters

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `API_URL` | `http://api:8080` | Backend API URL |
| `API_KEY` | `""` | Optional API authentication key |
| `MODEL_DIR` | `models` | Model storage directory |
| `MAX_MODEL_VERSIONS` | `10` | Maximum versions per model type |
| `MODEL_RETENTION_DAYS` | `30` | Days to retain old models |
| `ML_SERVICE_URL` | `http://localhost:5000` | ML service URL for scheduler |

### Retraining Thresholds

| Model Type | Data Threshold | Performance Threshold |
|------------|----------------|----------------------|
| Anomaly | 500 samples | 85% accuracy |
| Clustering | 200 samples | 80% accuracy |
| Recommendation | 100 samples | 75% accuracy |
| Trend | 50 samples | 70% accuracy |

## 🚀 Next Steps (Phase 2)

### Advanced Features (Planned)
- [ ] Ensemble anomaly detection (Isolation Forest + One-Class SVM)
- [ ] Real-time streaming anomaly detection
- [ ] Advanced user behavior sequence analysis
- [ ] Hybrid recommendation systems
- [ ] ARIMA and Prophet time series models
- [ ] A/B testing framework
- [ ] Model performance monitoring dashboard

### Academic Support
- [ ] Experimental comparison documentation
- [ ] ROC curves and performance visualizations
- [ ] Methodology documentation for thesis
- [ ] Benchmark results generation

## 📝 Logging

The service provides comprehensive logging:
- Model training events
- API requests and responses
- Scheduler job execution
- Version management operations
- Error tracking and debugging

Logs are structured and include timestamps, log levels, and contextual information.

## 🔒 Security

- Input validation for all endpoints
- Error handling without information leakage
- Optional API key authentication
- Request size limits
- Rate limiting (when integrated with backend)

## 🤝 Integration

The ML service integrates with:
- **Backend API**: For real training data
- **Database**: Via backend API endpoints
- **Frontend**: Via REST API calls
- **Monitoring**: Through health checks and metrics

## 📈 Monitoring

Key metrics tracked:
- Model training success/failure rates
- Prediction latency
- Model accuracy over time
- API response times
- Scheduler job execution status
- Model version usage statistics

---

**Status**: Phase 1 Complete ✅
**Next Phase**: Advanced Features & Academic Documentation
**Last Updated**: December 2024 