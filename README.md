# E-Commerce Platform with ML Integration

A modern e-commerce platform that leverages machine learning for security, personalization, and trend analysis.

## 🏗️ Project Structure

```
.
├── front/                 # React + Vite frontend
├── back/                  # Go + Fiber backend
├── ml/                    # Python ML service
├── docs/                  # Documentation
├── nginx/                 # Nginx configuration
├── docker-compose.yaml    # Docker orchestration
└── README.md             # This file
```

## 🚀 Features

- **Frontend**: React application with Vite
  - User interface for shopping
  - Operator dashboard for monitoring
  - Real-time updates

- **Backend**: Go + Fiber API
  - User management
  - Product management
  - Event tracking
  - Request categorization
  - **Automated and integration tests for request analysis**

- **ML Service**: Python-based
  - Anomaly detection
  - User behavior tracking
  - Product recommendations
  - Trend analysis

## 🛠️ Setup

1. **Prerequisites**
   - Docker and Docker Compose
   - Node.js (for local frontend development)
   - Go (for local backend development)
   - Python 3.11+ (for local ML service development)

2. **Environment Setup**
   ```bash
   # Clone the repository
   git clone <repository-url>
   cd <repository-name>

   # Start all services
   docker-compose up -d
   ```

3. **Access the Application**
   - Frontend: http://localhost
   - Backend API: http://localhost/api
   - ML Service: http://localhost:5000

## 🔧 Development

### Frontend Development
```bash
cd front
npm install
npm run dev
```

### Backend Development
```bash
cd back
go mod download
go run main.go
```

### ML Service Development
```bash
cd ml
python -m venv venv
source venv/bin/activate  # or `venv\Scripts\activate` on Windows
pip install -r requirements.txt
python app.py
```

## 📊 ML Models

The project implements several ML models:

1. **Anomaly Detection**
   - Isolation Forest for request anomaly detection
   - Real-time monitoring and alerts

2. **User Behavior Tracking**
   - Clustering for user activity patterns
   - Clickstream analysis

3. **Recommendations**
   - Item-based Collaborative Filtering
   - Periodic retraining

4. **Trend Analysis**
   - Time series analysis
   - Seasonal trend prediction

## 🛡️ Request Analysis & Security Workflow

- **All incoming requests** are logged and analyzed by the backend.
- **Features** are extracted and sent to the ML service for anomaly detection.
- **Anomalous requests** are flagged and the IP is automatically banned for a period.
- **All requests** are available for analytics and operator review in the dashboard.
- **Automated Go tests** and a **Python attack simulation script** ensure the workflow is robust.

## 🧪 Automated Testing & Attack Simulation

### Backend Unit/Integration Tests
- Located in `back/routes/handlers/request_log_handlers_test.go`
- **How to run:**
  ```bash
  cd back
  go test ./routes/handlers
  ```
- **Covers:**
  - Normal request logging
  - Anomaly detection and IP banning
  - Banned IP cannot make requests
  - Analytics endpoint returns correct data

### Attack Simulation Script
- Located at `ml/simulation/attack_simulation.py`
- **How to run:**
  ```bash
  cd ml/simulation
  pip install requests
  python attack_simulation.py
  ```
- **What it does:**
  - Simulates normal, rapid, and anomalous requests
  - Demonstrates anomaly detection and IP banning
  - Prints results for each step

## 📈 Interpreting Results
- **Test output**: Go test output will show pass/fail for each scenario.
- **Simulation output**: Python script prints request results and ban status.
- **Dashboard**: View real-time request logs and anomalies in the operator dashboard (Frontend > Analytics > Requests).

## 📝 Documentation

Detailed documentation can be found in the `docs/` directory:
- Architecture overview
- API documentation
- ML model documentation
- Deployment guide

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details. 