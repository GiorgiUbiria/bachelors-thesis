# Unseen Machine Learning: A Multi-Layered Framework for E-Commerce Security, Personalization, and Monitoring

- **Machine Learning methodology and study** (main goal)
- **Periodic retraining and optimization**
- **Request monitoring & auto-actions**
- **Realistic simulations and trend analysis**
- **No Neural Networks** (only traditional ML models)

---

# 🎯 Project Overview

---

## 📜 Main Structure

| Part                       | Purpose                                                           |
| :------------------------- | :---------------------------------------------------------------- |
| **Frontend**               | React app for e-commerce users and Operators                      |
| **Backend**                | Go API + PostgreSQL database                                      |
| **Machine Learning Layer** | Python scripts for training, retraining, predictions, simulations |
| **Automation Layer**       | Cron jobs (or Go-based schedulers) to periodically retrain models |

---

# 🛠 Practical Work Breakdown

---

## 1. 🖥 Frontend (React)

- **User Site**:
  - Product feed
  - Cart & Favorites
  - Personalized Recommendations (update when user behavior changes)
  - Default functionality (Login, Register, Landing Page)

- **Operator Dashboard**:
  - Incoming requests overview (normal, warning, anomaly)
  - Real-time view of the simulated attack data
  - Preemptive actions listing (e.g., "IP blocked", "Session terminated")

**Tech Suggestions:**
- React Query (auto-refreshing dashboards)
- Socket.io / WebSocket (real-time anomalies)
- TailwindCSS (quick UI)

---

## 2. 🖥 Backend (Golang)

- **User Management**: Auth, sessions.
- **Product Management**: CRUD operations.
- **Event Tracking**: Save click, cart, favorites as events.
- **Request Categorization**: Simple rules (normal / warning / anomaly).
- **ML Communication**: Trigger predictions or get recommendations.
- **Operator Actions**: Automatically execute predefined actions when anomaly types are detected.

**Golang Libraries**: 
- `Fiber` for HTTP API.
- `GORM` for PostgreSQL.
- `cron` library for scheduled tasks.

---

## 3. 🧠 Machine Learning System (Python)

---

| ML Task                    | Method                                                              | Purpose                                            |
| :------------------------- | :------------------------------------------------------------------ | :------------------------------------------------- |
| **Anomaly Detection**      | Isolation Forest / One-Class SVM                                    | Detect malicious requests                          |
| **User Behavior Tracking** | Clustering (e.g., KMeans) + Clickstream Analysis                    | Group user activity patterns                       |
| **Recommendations**        | Item-based Collaborative Filtering (k-NN, basic similarity metrics) | Suggest products based on cart, favorites          |
| **Trend Prediction**       | Linear Regression / ARIMA (Time Series)                             | Predict future product demands                     |
| **Manual Implementation**  | Simple custom Logistic Regression / Decision Tree                   | Show theoretical knowledge without using libraries |

---

**Key Points**:

- **Why the method is chosen**: 
  - Explain *classification* vs *regression* where necessary.
  - Compare options (e.g., KMeans vs DBSCAN for clustering users).
  - Personal (manual) implementation of simple models where possible (e.g., own Logistic Regression using NumPy).

- **Periodic Retraining**:
  - Cron jobs or scheduled scripts will automatically retrain models:
    - Example: Every 6 hours, every 1000 events, or based on time+data size.
    - Delete old models or version them if needed.

- **Optimized Chain**:
  - Cache recommendation results per user in the DB or in-memory store (e.g., Redis).
  - Only recompute recommendations when user adds/favorites a product.
  - Async training processes — no blocking user interactions.

---

## 4. 🛡️ Request Monitoring & Attack Simulation

- **Operator View**:
  - Requests displayed in categories (Normal, Warning, Anomaly).
  - When an anomaly is detected:
    - Block IP
    - Notify admin
    - Mark sessions for manual inspection.

- **Attack Simulation**:
  - Write Python script that generates:
    - High volume traffic (simulate DDoS)
    - Suspicious behavior (many cart additions, password guess attempts)
    - SQLi-like payloads.

- Simulated data will be injected into the backend as incoming requests to trigger the detection logic.

---

## 5. 📈 Trend Analysis with Sparse Data

- **Sparse Data Handling**:
  - Mock dataset generation for buys, sells, seasonal trends.
  - Use regression (Linear Regression, ARIMA) to predict next month's best-sellers.

- **Realistic Trends**:
  - "More coats sold in winter"
  - "More electronics sold near holidays"

---

# 📄 Academic Paper Structure

---

| Section                                             | Details                                                                                                                                                                                                                                                       |
| :-------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **1. Introduction**                                 | Importance of ML in modern applications, daily life examples                                                                                                                                                                                                  |
| **2. Literature Review**                            | Overview of ML methods (classification, regression, clustering), real-world ML uses                                                                                                                                                                           |
| **3. Problem Statement**                            | Why ML is critical for modern e-commerce systems                                                                                                                                                                                                              |
| **4. Methodology**                                  | <ul><li>Explain each method you chose (why classification, why regression, etc.)</li><li>Compare alternative methods (e.g., Random Forest vs Decision Tree)</li><li>Show manual (simplified) implementation of Logistic Regression or Decision Tree</li></ul> |
| **5. System Architecture**                          | Diagram and description of React, Go API, Python ML connection                                                                                                                                                                                                |
| **6. Implementation**                               | Details on Frontend, Backend, ML System                                                                                                                                                                                                                       |
| **7. Training and Optimization**                    | How periodic training happens, when models update, performance optimizations                                                                                                                                                                                  |
| **8. Attack Simulation and Request Categorization** | How anomaly detection and categorization works, how you simulate attacks                                                                                                                                                                                      |
| **9. Results**                                      | Screenshots, graphs (training accuracy, anomaly detection rates, trend prediction accuracy)                                                                                                                                                                   |
| **10. Discussion and Challenges**                   | Problems encountered and how you overcame them                                                                                                                                                                                                                |
| **11. Conclusion and Future Work**                  | How system could be improved (real-time updates, better scaling)                                                                                                                                                                                              |
| **References**                                      | Libraries, papers, methods used                                                                                                                                                                                                                               |

---

# 📦 Deliverables

- **Web Application** (React, Go API, Python ML services)
- **Scripts** for:
  - Attack Simulation
  - Mock Sparse Data Generation
- **Paper (~30 pages)** with:
  - Manual simple ML algorithm implementations
  - Comparisons and critical analysis
- **GitHub Repo** with clear separation:
  ```
  /frontend
  /backend
  /ml
  /docs
  ```

---

# 🗺 Suggested Timeline (adjusted)

| Week | Task                                                              |
| :--- | :---------------------------------------------------------------- |
| 1    | Set up React frontend, Go backend, PostgreSQL database            |
| 2    | Implement basic product browsing, user system, cart, favorites    |
| 3    | Implement event tracking (clicks, cart actions)                   |
| 4    | Implement request categorization rules, operator dashboard        |
| 5    | Develop basic anomaly detection script and API                    |
| 6    | Develop basic recommendation system with periodic retraining      |
| 7    | Trend analysis model, sparse data seeding script                  |
| 8    | Attack simulation script, integrate with dashboard                |
| 9    | Optimization passes (cache results, asynchronous processing)      |
| 10   | Write academic paper, polish final codebase, prepare presentation |

---

# 🧠 Notes

- Don't overcomplicate ML models; prioritize **interpretability and justification**.
- Show that you **manually understand and can implement** basic ML techniques.
- Keep **periodic retraining light and scheduled** — no re-training triggered after every user event unless necessary.

---


# Other prompts

---

I have a bachelor's academic paper to prepare and the project to write for it.

The project should be a web application (an e-commerce) platform, where several processes are automated using machine learning.

These processes are: Incoming request anomaly detection, visitor count and activity time + user click rate + user cart/favourite options for statistical analysis and future trends prediction + recommendation system on the user's feed.

I myself envision this project to be separated into 3 parts:
1. Frontend - something in React
2. Backend - API in Golang + Databases
3. Machine Learning - Scripts in Python + trained machine learning models

---

---

1. Main focus is the Machine Learning itself. I should provide its applications in real life, but also expand on following topics: Why I chose the method (regression, classification, etc.), Compare options (Different classifiers, etc.), Manual/Personal implementation (simplified of course) of the method.

2. Models should be periodically trained, meaning the trained models should not stay stale and based on realistic criteria, be fed with new data and trained accordingly. This task should be automated.

3. Important part is optimization. Because there's a chain of request-responses & operations, everything still should run smoothly. For example, when the user favourites a new thing or adds a new thing to the cart, his recommendations should update accordingly.

4. Request analysis will be provided to the user with role Operator, that will have his own dashboard of all the incoming requests that will be classified (not ml) in different categories (normal, warning, etc.). When an anomaly is detected, based on what kind of an anomaly it is, preemptive measures should be automatically taken by the program. There also should exist a script that will simulate "real-wordly" attack/issue scenarios for visual representation.

5. "Trends analysis" will be used to regress the information about what kind of products are most requested, least requested, what to invest in (since the data is sparse for the later, we will include a script that will seed the database with a mock data about buys and sells for different seasons, etc.).

Lastly, the project won't use Neural Networks.
---