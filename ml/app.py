from flask import Flask, request, jsonify
from flask_cors import CORS
import numpy as np
import pandas as pd
from sklearn.ensemble import IsolationForest
from sklearn.cluster import KMeans
from sklearn.metrics.pairwise import cosine_similarity
import joblib
import os
from datetime import datetime, timedelta
import requests
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

app = Flask(__name__)
CORS(app)

# Configuration
API_URL = os.getenv('API_URL', 'http://api:8080')
MODEL_DIR = 'models'
os.makedirs(MODEL_DIR, exist_ok=True)

# Initialize models
anomaly_detector = IsolationForest(contamination=0.1, random_state=42)
user_clusterer = KMeans(n_clusters=5, random_state=42)
recommendation_model = None

def load_or_create_models():
    """Load existing models or create new ones if they don't exist"""
    global anomaly_detector, user_clusterer, recommendation_model
    
    # Load or create anomaly detector
    try:
        anomaly_detector = joblib.load(f'{MODEL_DIR}/anomaly_detector.joblib')
    except:
        anomaly_detector = IsolationForest(contamination=0.1, random_state=42)
    
    # Load or create user clusterer
    try:
        user_clusterer = joblib.load(f'{MODEL_DIR}/user_clusterer.joblib')
    except:
        user_clusterer = KMeans(n_clusters=5, random_state=42)
    
    # Load or create recommendation model
    try:
        recommendation_model = joblib.load(f'{MODEL_DIR}/recommendation_model.joblib')
    except:
        recommendation_model = None

def save_models():
    """Save all models to disk"""
    joblib.dump(anomaly_detector, f'{MODEL_DIR}/anomaly_detector.joblib')
    joblib.dump(user_clusterer, f'{MODEL_DIR}/user_clusterer.joblib')
    if recommendation_model is not None:
        joblib.dump(recommendation_model, f'{MODEL_DIR}/recommendation_model.joblib')

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    return jsonify({"status": "healthy"})

@app.route('/analyze/request', methods=['POST'])
def analyze_request():
    """Analyze request for anomalies"""
    data = request.json
    
    # Extract features
    features = np.array([
        data.get('response_time', 0),
        data.get('request_size', 0),
        data.get('error_count', 0)
    ]).reshape(1, -1)
    
    # Predict anomaly
    score = anomaly_detector.score_samples(features)
    is_anomaly = score < np.percentile(anomaly_detector.score_samples(features), 10)
    
    return jsonify({
        "is_anomaly": bool(is_anomaly),
        "anomaly_score": float(score[0])
    })

@app.route('/analyze/user', methods=['POST'])
def analyze_user_behavior():
    """Analyze user behavior patterns"""
    data = request.json
    
    # Extract user activity features
    features = np.array([
        data.get('login_count', 0),
        data.get('purchase_count', 0),
        data.get('cart_count', 0),
        data.get('favorite_count', 0)
    ]).reshape(1, -1)
    
    # Predict cluster
    cluster = user_clusterer.predict(features)[0]
    
    return jsonify({
        "cluster": int(cluster),
        "cluster_center": user_clusterer.cluster_centers_[cluster].tolist()
    })

@app.route('/recommend', methods=['POST'])
def get_recommendations():
    """Get product recommendations for a user"""
    data = request.json
    user_id = data.get('user_id')
    
    if recommendation_model is None:
        return jsonify({"error": "Recommendation model not trained"}), 400
    
    # Get user's purchase history
    response = requests.get(f"{API_URL}/api/users/{user_id}/purchases")
    if response.status_code != 200:
        return jsonify({"error": "Failed to fetch user purchases"}), 400
    
    purchases = response.json()
    
    # Generate recommendations
    recommendations = recommendation_model.predict(purchases)
    
    return jsonify({
        "recommendations": recommendations.tolist()
    })

@app.route('/train', methods=['POST'])
def train_models():
    """Train all ML models"""
    try:
        # Get training data from API
        response = requests.get(f"{API_URL}/api/ml/training-data")
        if response.status_code != 200:
            return jsonify({"error": "Failed to fetch training data"}), 400
        
        data = response.json()
        
        # Train anomaly detector
        request_features = np.array([
            [d['response_time'], d['request_size'], d['error_count']]
            for d in data['request_logs']
        ])
        anomaly_detector.fit(request_features)
        
        # Train user clusterer
        user_features = np.array([
            [d['login_count'], d['purchase_count'], d['cart_count'], d['favorite_count']]
            for d in data['user_activities']
        ])
        user_clusterer.fit(user_features)
        
        # Train recommendation model
        global recommendation_model
        recommendation_model = train_recommendation_model(data['purchases'])
        
        # Save models
        save_models()
        
        return jsonify({"status": "success"})
    
    except Exception as e:
        return jsonify({"error": str(e)}), 500

def train_recommendation_model(purchases):
    """Train a simple collaborative filtering model"""
    # Create user-item matrix
    user_item_matrix = pd.DataFrame(purchases).pivot(
        index='user_id',
        columns='product_id',
        values='quantity'
    ).fillna(0)
    
    # Calculate item similarity
    item_similarity = cosine_similarity(user_item_matrix.T)
    
    return item_similarity

if __name__ == '__main__':
    load_or_create_models()
    app.run(host='0.0.0.0', port=5000) 