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
import logging
from common.data_loader import (
    load_anomaly_data, load_clustering_data, 
    load_recommendation_data, load_trend_data,
    validate_anomaly_data, validate_clustering_data,
    validate_recommendation_data, validate_trend_data
)

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)

# Configuration
API_URL = os.getenv('API_URL', 'http://api:8080')
API_KEY = os.getenv('API_KEY', '')
MODEL_DIR = 'models'
os.makedirs(MODEL_DIR, exist_ok=True)

# Initialize models
anomaly_detector = IsolationForest(contamination=0.1, random_state=42)
user_clusterer = KMeans(n_clusters=5, random_state=42)
recommendation_model = None
recommendation_matrix = None  # Store the user-item matrix for recommendations

def load_or_create_models():
    """Load existing models or create new ones if they don't exist"""
    global anomaly_detector, user_clusterer, recommendation_model, recommendation_matrix
    
    # Load or create anomaly detector
    try:
        anomaly_detector = joblib.load(f'{MODEL_DIR}/anomaly_detector.joblib')
        logger.info("Loaded existing anomaly detector model")
    except Exception as e:
        logger.info(f"Creating new anomaly detector model: {e}")
        anomaly_detector = IsolationForest(contamination=0.1, random_state=42)
    
    # Load or create user clusterer
    try:
        user_clusterer = joblib.load(f'{MODEL_DIR}/user_clusterer.joblib')
        logger.info("Loaded existing user clustering model")
    except Exception as e:
        logger.info(f"Creating new user clustering model: {e}")
        user_clusterer = KMeans(n_clusters=5, random_state=42)
    
    # Load or create recommendation model
    try:
        recommendation_data = joblib.load(f'{MODEL_DIR}/recommendation_model.joblib')
        recommendation_model = recommendation_data['similarity_matrix']
        recommendation_matrix = recommendation_data['user_item_matrix']
        logger.info("Loaded existing recommendation model")
    except Exception as e:
        logger.info(f"Creating new recommendation model: {e}")
        recommendation_model = None
        recommendation_matrix = None

def save_models():
    """Save all models to disk"""
    try:
        joblib.dump(anomaly_detector, f'{MODEL_DIR}/anomaly_detector.joblib')
        joblib.dump(user_clusterer, f'{MODEL_DIR}/user_clusterer.joblib')
        
        if recommendation_model is not None and recommendation_matrix is not None:
            recommendation_data = {
                'similarity_matrix': recommendation_model,
                'user_item_matrix': recommendation_matrix
            }
            joblib.dump(recommendation_data, f'{MODEL_DIR}/recommendation_model.joblib')
        
        logger.info("Successfully saved all models")
    except Exception as e:
        logger.error(f"Error saving models: {e}")

def _make_authenticated_request(endpoint: str, method: str = 'GET', data: dict = None):
    """Make authenticated API request to backend"""
    try:
        url = f"{API_URL}{endpoint}"
        headers = {}
        
        if API_KEY:
            headers['Authorization'] = f'Bearer {API_KEY}'
        
        if method == 'GET':
            response = requests.get(url, headers=headers, timeout=30)
        elif method == 'POST':
            headers['Content-Type'] = 'application/json'
            response = requests.post(url, headers=headers, json=data, timeout=30)
        else:
            raise ValueError(f"Unsupported HTTP method: {method}")
        
        response.raise_for_status()
        return response.json()
        
    except requests.exceptions.RequestException as e:
        logger.error(f"API request failed for {endpoint}: {e}")
        return None
    except Exception as e:
        logger.error(f"Unexpected error in API request: {e}")
        return None

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    return jsonify({
        "status": "healthy",
        "models_loaded": {
            "anomaly_detector": anomaly_detector is not None,
            "user_clusterer": user_clusterer is not None,
            "recommendation_model": recommendation_model is not None
        },
        "timestamp": datetime.now().isoformat()
    })

@app.route('/analyze/request', methods=['POST'])
def analyze_request():
    """Analyze request for anomalies"""
    try:
        data = request.json
        if not data:
            return jsonify({"error": "No data provided"}), 400
        
        # Extract features
        features = np.array([
            data.get('response_time', 0),
            data.get('request_size', 0),
            data.get('error_count', 0)
        ]).reshape(1, -1)
        
        # Predict anomaly
        prediction = anomaly_detector.predict(features)
        score = anomaly_detector.score_samples(features)
        
        # Calculate anomaly threshold (lower scores indicate anomalies)
        is_anomaly = prediction[0] == -1
        
        return jsonify({
            "is_anomaly": bool(is_anomaly),
            "anomaly_score": float(score[0]),
            "prediction": int(prediction[0])
        })
        
    except Exception as e:
        logger.error(f"Error in anomaly analysis: {e}")
        return jsonify({"error": "Internal server error"}), 500

@app.route('/analyze/user', methods=['POST'])
def analyze_user_behavior():
    """Analyze user behavior patterns"""
    try:
        data = request.json
        if not data:
            return jsonify({"error": "No data provided"}), 400
        
        # Extract user activity features
        features = np.array([
            data.get('login_count', 0),
            data.get('purchase_count', 0),
            data.get('cart_count', 0),
            data.get('favorite_count', 0)
        ]).reshape(1, -1)
        
        # Predict cluster
        cluster = user_clusterer.predict(features)[0]
        
        # Get cluster center for interpretation
        cluster_center = user_clusterer.cluster_centers_[cluster].tolist()
        
        return jsonify({
            "cluster": int(cluster),
            "cluster_center": cluster_center,
            "user_features": features[0].tolist()
        })
        
    except Exception as e:
        logger.error(f"Error in user behavior analysis: {e}")
        return jsonify({"error": "Internal server error"}), 500

@app.route('/recommend', methods=['POST'])
def get_recommendations():
    """Get product recommendations for a user"""
    try:
        data = request.json
        if not data:
            return jsonify({"error": "No data provided"}), 400
            
        user_id = data.get('user_id')
        top_k = data.get('top_k', 5)
        
        if recommendation_model is None or recommendation_matrix is None:
            return jsonify({"error": "Recommendation model not trained"}), 400
        
        # Get user's purchase history from backend
        user_data = _make_authenticated_request(f"/api/users/{user_id}/purchases")
        if user_data is None:
            return jsonify({"error": "Failed to fetch user purchases"}), 400
        
        # Generate recommendations using item-based collaborative filtering
        recommendations = generate_recommendations(user_id, top_k)
        
        return jsonify({
            "user_id": user_id,
            "recommendations": recommendations,
            "top_k": top_k
        })
        
    except Exception as e:
        logger.error(f"Error in recommendation generation: {e}")
        return jsonify({"error": "Internal server error"}), 500

def generate_recommendations(user_id: int, top_k: int = 5):
    """Generate recommendations using item-based collaborative filtering"""
    try:
        if recommendation_matrix is None:
            return []
        
        # Convert user_id to matrix index (assuming user_id maps to row index)
        user_index = user_id % recommendation_matrix.shape[0]
        
        # Get user's ratings
        user_ratings = recommendation_matrix[user_index]
        
        # Calculate scores for all items
        scores = recommendation_model.dot(user_ratings)
        
        # Get items user hasn't rated
        unrated_items = np.where(user_ratings == 0)[0]
        
        if len(unrated_items) == 0:
            return []
        
        # Get scores for unrated items only
        unrated_scores = scores[unrated_items]
        
        # Get top-k recommendations
        top_indices = np.argsort(unrated_scores)[::-1][:top_k]
        recommended_items = unrated_items[top_indices]
        
        # Return item IDs with scores
        recommendations = [
            {
                "item_id": int(item_id),
                "score": float(scores[item_id])
            }
            for item_id in recommended_items
        ]
        
        return recommendations
        
    except Exception as e:
        logger.error(f"Error generating recommendations: {e}")
        return []

@app.route('/train', methods=['POST'])
def train_models_with_versioning():
    """Train all ML models using real data from backend API with versioning support"""
    try:
        logger.info("Starting model training process with versioning")
        
        # Initialize services if not already done
        if scheduler is None or version_manager is None:
            initialize_services()
        
        # Train anomaly detector
        logger.info("Training anomaly detection model")
        anomaly_data = load_anomaly_data()
        if anomaly_data is not None and len(anomaly_data) > 10:
            global anomaly_detector
            new_anomaly_detector = IsolationForest(contamination=0.1, random_state=42)
            new_anomaly_detector.fit(anomaly_data)
            
            # Save with versioning
            if version_manager:
                metadata = {"training_samples": len(anomaly_data), "features": anomaly_data.shape[1]}
                version_id = version_manager.save_model(new_anomaly_detector, "anomaly", metadata)
                logger.info(f"Anomaly detector saved as version {version_id}")
            
            anomaly_detector = new_anomaly_detector
            logger.info(f"Anomaly detector trained with {len(anomaly_data)} samples")
        else:
            logger.warning("Insufficient anomaly data for training")
        
        # Train user clusterer
        logger.info("Training user clustering model")
        clustering_data = load_clustering_data()
        if clustering_data is not None and len(clustering_data) > 5:
            global user_clusterer
            new_user_clusterer = KMeans(n_clusters=min(5, len(clustering_data)), random_state=42)
            new_user_clusterer.fit(clustering_data)
            
            # Save with versioning
            if version_manager:
                metadata = {"training_samples": len(clustering_data), "n_clusters": new_user_clusterer.n_clusters}
                version_id = version_manager.save_model(new_user_clusterer, "clustering", metadata)
                logger.info(f"User clusterer saved as version {version_id}")
            
            user_clusterer = new_user_clusterer
            logger.info(f"User clusterer trained with {len(clustering_data)} samples")
        else:
            logger.warning("Insufficient clustering data for training")
        
        # Train recommendation model
        logger.info("Training recommendation model")
        recommendation_data = load_recommendation_data()
        if recommendation_data is not None:
            user_item_matrix, df = recommendation_data
            if user_item_matrix.size > 0:
                global recommendation_model, recommendation_matrix
                new_recommendation_matrix = user_item_matrix
                new_recommendation_model = cosine_similarity(user_item_matrix.T)
                
                # Save with versioning
                if version_manager:
                    model_data = {
                        'similarity_matrix': new_recommendation_model,
                        'user_item_matrix': new_recommendation_matrix
                    }
                    metadata = {
                        "matrix_shape": user_item_matrix.shape,
                        "n_users": user_item_matrix.shape[0],
                        "n_items": user_item_matrix.shape[1]
                    }
                    version_id = version_manager.save_model(model_data, "recommendation", metadata)
                    logger.info(f"Recommendation model saved as version {version_id}")
                
                recommendation_matrix = new_recommendation_matrix
                recommendation_model = new_recommendation_model
                logger.info(f"Recommendation model trained with matrix shape: {user_item_matrix.shape}")
            else:
                logger.warning("Empty recommendation data")
        else:
            logger.warning("Failed to load recommendation data")
        
        return jsonify({
            "status": "success",
            "message": "All models trained successfully with versioning",
            "timestamp": datetime.now().isoformat(),
            "versioning_enabled": version_manager is not None
        })
    
    except Exception as e:
        logger.error(f"Error in model training: {e}")
        return jsonify({"error": f"Training failed: {str(e)}"}), 500

@app.route('/retrain', methods=['POST'])
def retrain_models():
    """Trigger model retraining"""
    return train_models_with_versioning()

# Manual ML Implementation Endpoints
@app.route('/train/manual-logistic', methods=['POST'])
def train_manual_logistic():
    """Train manual logistic regression model"""
    try:
        from manual_impl.manual_logistic import ManualLogisticRegression
        
        data = request.json
        if not data:
            return jsonify({"error": "No training data provided"}), 400
        
        # Extract features and labels
        X = np.array(data.get('features', []))
        y = np.array(data.get('labels', []))
        
        if len(X) == 0 or len(y) == 0:
            return jsonify({"error": "Empty training data"}), 400
        
        if len(X) != len(y):
            return jsonify({"error": "Features and labels must have same length"}), 400
        
        # Get hyperparameters
        lr = data.get('learning_rate', 0.01)
        n_iter = data.get('n_iterations', 1000)
        regularization = data.get('regularization', None)
        lambda_reg = data.get('lambda_reg', 0.01)
        
        # Train model
        model = ManualLogisticRegression(
            lr=lr, 
            n_iter=n_iter, 
            regularization=regularization, 
            lambda_reg=lambda_reg
        )
        model.fit(X, y)
        
        # Save model
        joblib.dump(model, f'{MODEL_DIR}/manual_logistic.joblib')
        
        return jsonify({
            "status": "success",
            "message": "Manual logistic regression trained successfully",
            "model_info": {
                "learning_rate": lr,
                "n_iterations": n_iter,
                "regularization": regularization,
                "lambda_reg": lambda_reg,
                "n_features": len(X[0]),
                "n_samples": len(X),
                "final_cost": float(model.cost_history[-1]) if model.cost_history else None
            }
        })
        
    except Exception as e:
        logger.error(f"Error training manual logistic regression: {e}")
        return jsonify({"error": f"Training failed: {str(e)}"}), 500

@app.route('/train/manual-tree', methods=['POST'])
def train_manual_tree():
    """Train manual decision tree model"""
    try:
        from manual_impl.manual_tree import ManualDecisionTree
        
        data = request.json
        if not data:
            return jsonify({"error": "No training data provided"}), 400
        
        # Extract features and labels
        X = np.array(data.get('features', []))
        y = np.array(data.get('labels', []))
        
        if len(X) == 0 or len(y) == 0:
            return jsonify({"error": "Empty training data"}), 400
        
        if len(X) != len(y):
            return jsonify({"error": "Features and labels must have same length"}), 400
        
        # Get hyperparameters
        max_depth = data.get('max_depth', 3)
        min_samples_split = data.get('min_samples_split', 2)
        min_samples_leaf = data.get('min_samples_leaf', 1)
        criterion = data.get('criterion', 'gini')
        max_features = data.get('max_features', None)
        
        # Train model
        model = ManualDecisionTree(
            max_depth=max_depth,
            min_samples_split=min_samples_split,
            min_samples_leaf=min_samples_leaf,
            criterion=criterion,
            max_features=max_features
        )
        model.fit(X, y)
        
        # Save model
        joblib.dump(model, f'{MODEL_DIR}/manual_tree.joblib')
        
        return jsonify({
            "status": "success",
            "message": "Manual decision tree trained successfully",
            "model_info": {
                "max_depth": max_depth,
                "min_samples_split": min_samples_split,
                "min_samples_leaf": min_samples_leaf,
                "criterion": criterion,
                "max_features": max_features,
                "n_features": len(X[0]),
                "n_samples": len(X),
                "tree_depth": model.get_depth(),
                "n_leaves": model.get_n_leaves(),
                "feature_importances": model.feature_importances_.tolist()
            }
        })
        
    except Exception as e:
        logger.error(f"Error training manual decision tree: {e}")
        return jsonify({"error": f"Training failed: {str(e)}"}), 500

@app.route('/predict/manual-logistic', methods=['POST'])
def predict_manual_logistic():
    """Make predictions using manual logistic regression"""
    try:
        # Load model
        model = joblib.load(f'{MODEL_DIR}/manual_logistic.joblib')
        
        data = request.json
        if not data:
            return jsonify({"error": "No prediction data provided"}), 400
        
        # Extract features
        X = np.array(data.get('features', []))
        if len(X) == 0:
            return jsonify({"error": "Empty prediction data"}), 400
        
        # Make predictions
        predictions = model.predict(X)
        probabilities = model.predict_proba(X)
        
        return jsonify({
            "predictions": predictions.tolist(),
            "probabilities": probabilities.tolist(),
            "n_samples": len(X)
        })
        
    except FileNotFoundError:
        return jsonify({"error": "Manual logistic regression model not found. Train it first."}), 404
    except Exception as e:
        logger.error(f"Error in manual logistic regression prediction: {e}")
        return jsonify({"error": f"Prediction failed: {str(e)}"}), 500

@app.route('/predict/manual-tree', methods=['POST'])
def predict_manual_tree():
    """Make predictions using manual decision tree"""
    try:
        # Load model
        model = joblib.load(f'{MODEL_DIR}/manual_tree.joblib')
        
        data = request.json
        if not data:
            return jsonify({"error": "No prediction data provided"}), 400
        
        # Extract features
        X = np.array(data.get('features', []))
        if len(X) == 0:
            return jsonify({"error": "Empty prediction data"}), 400
        
        # Make predictions
        predictions = model.predict(X)
        
        return jsonify({
            "predictions": predictions.tolist(),
            "n_samples": len(X),
            "feature_importances": model.feature_importances_.tolist()
        })
        
    except FileNotFoundError:
        return jsonify({"error": "Manual decision tree model not found. Train it first."}), 404
    except Exception as e:
        logger.error(f"Error in manual decision tree prediction: {e}")
        return jsonify({"error": f"Prediction failed: {str(e)}"}), 500

@app.route('/compare/models', methods=['POST'])
def compare_models():
    """Compare manual implementations with sklearn implementations"""
    try:
        from sklearn.linear_model import LogisticRegression
        from sklearn.tree import DecisionTreeClassifier
        from sklearn.model_selection import train_test_split
        from sklearn.metrics import accuracy_score, precision_score, recall_score, f1_score
        from manual_impl.manual_logistic import ManualLogisticRegression
        from manual_impl.manual_tree import ManualDecisionTree
        import time
        
        data = request.json
        if not data:
            return jsonify({"error": "No comparison data provided"}), 400
        
        # Extract features and labels
        X = np.array(data.get('features', []))
        y = np.array(data.get('labels', []))
        
        if len(X) == 0 or len(y) == 0:
            return jsonify({"error": "Empty comparison data"}), 400
        
        # Split data
        test_size = data.get('test_size', 0.2)
        X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=test_size, random_state=42)
        
        results = {}
        
        # Compare Logistic Regression
        logger.info("Comparing Logistic Regression implementations")
        
        # Manual implementation
        start_time = time.time()
        manual_lr = ManualLogisticRegression(lr=0.01, n_iter=1000)
        manual_lr.fit(X_train, y_train)
        manual_lr_time = time.time() - start_time
        
        manual_lr_pred = manual_lr.predict(X_test)
        manual_lr_accuracy = accuracy_score(y_test, manual_lr_pred)
        
        # Sklearn implementation
        start_time = time.time()
        sklearn_lr = LogisticRegression(max_iter=1000, random_state=42)
        sklearn_lr.fit(X_train, y_train)
        sklearn_lr_time = time.time() - start_time
        
        sklearn_lr_pred = sklearn_lr.predict(X_test)
        sklearn_lr_accuracy = accuracy_score(y_test, sklearn_lr_pred)
        
        results['logistic_regression'] = {
            'manual': {
                'accuracy': float(manual_lr_accuracy),
                'training_time': float(manual_lr_time),
                'final_cost': float(manual_lr.cost_history[-1]) if manual_lr.cost_history else None
            },
            'sklearn': {
                'accuracy': float(sklearn_lr_accuracy),
                'training_time': float(sklearn_lr_time)
            }
        }
        
        # Compare Decision Tree
        logger.info("Comparing Decision Tree implementations")
        
        # Manual implementation
        start_time = time.time()
        manual_tree = ManualDecisionTree(max_depth=5, criterion='gini')
        manual_tree.fit(X_train, y_train)
        manual_tree_time = time.time() - start_time
        
        manual_tree_pred = manual_tree.predict(X_test)
        manual_tree_accuracy = accuracy_score(y_test, manual_tree_pred)
        
        # Sklearn implementation
        start_time = time.time()
        sklearn_tree = DecisionTreeClassifier(max_depth=5, criterion='gini', random_state=42)
        sklearn_tree.fit(X_train, y_train)
        sklearn_tree_time = time.time() - start_time
        
        sklearn_tree_pred = sklearn_tree.predict(X_test)
        sklearn_tree_accuracy = accuracy_score(y_test, sklearn_tree_pred)
        
        results['decision_tree'] = {
            'manual': {
                'accuracy': float(manual_tree_accuracy),
                'training_time': float(manual_tree_time),
                'tree_depth': manual_tree.get_depth(),
                'n_leaves': manual_tree.get_n_leaves()
            },
            'sklearn': {
                'accuracy': float(sklearn_tree_accuracy),
                'training_time': float(sklearn_tree_time),
                'tree_depth': sklearn_tree.get_depth(),
                'n_leaves': sklearn_tree.get_n_leaves()
            }
        }
        
        # Overall comparison
        results['summary'] = {
            'dataset_info': {
                'n_samples': len(X),
                'n_features': X.shape[1],
                'n_classes': len(np.unique(y)),
                'test_size': test_size
            },
            'performance_comparison': {
                'logistic_regression_accuracy_diff': float(manual_lr_accuracy - sklearn_lr_accuracy),
                'decision_tree_accuracy_diff': float(manual_tree_accuracy - sklearn_tree_accuracy),
                'logistic_regression_time_ratio': float(manual_lr_time / sklearn_lr_time) if sklearn_lr_time > 0 else None,
                'decision_tree_time_ratio': float(manual_tree_time / sklearn_tree_time) if sklearn_tree_time > 0 else None
            }
        }
        
        return jsonify({
            "status": "success",
            "comparison_results": results
        })
        
    except Exception as e:
        logger.error(f"Error in model comparison: {e}")
        return jsonify({"error": f"Comparison failed: {str(e)}"}), 500

if __name__ == '__main__':
    load_or_create_models()
    app.run(host='0.0.0.0', port=5000)

# Scheduler and Version Management Integration
from training.scheduler import get_scheduler
from models.version_manager import get_version_manager

# Initialize scheduler and version manager
scheduler = None
version_manager = None

def initialize_services():
    """Initialize scheduler and version manager services"""
    global scheduler, version_manager
    try:
        scheduler = get_scheduler()
        version_manager = get_version_manager()
        logger.info("Scheduler and version manager initialized successfully")
    except Exception as e:
        logger.error(f"Error initializing services: {e}")

# Scheduler Management Endpoints
@app.route('/scheduler/status', methods=['GET'])
def get_scheduler_status():
    """Get scheduler status and job information"""
    try:
        if scheduler is None:
            return jsonify({"error": "Scheduler not initialized"}), 500
        
        status = scheduler.get_schedule_status()
        return jsonify({
            "status": "success",
            "scheduler_status": status
        })
        
    except Exception as e:
        logger.error(f"Error getting scheduler status: {e}")
        return jsonify({"error": f"Failed to get scheduler status: {str(e)}"}), 500

@app.route('/scheduler/trigger', methods=['POST'])
def trigger_manual_retrain():
    """Manually trigger model retraining"""
    try:
        if scheduler is None:
            return jsonify({"error": "Scheduler not initialized"}), 500
        
        data = request.json or {}
        model_type = data.get('model_type', 'all')
        
        scheduler.trigger_manual_retrain(model_type)
        
        return jsonify({
            "status": "success",
            "message": f"Manual retraining triggered for {model_type} models"
        })
        
    except Exception as e:
        logger.error(f"Error triggering manual retrain: {e}")
        return jsonify({"error": f"Failed to trigger retraining: {str(e)}"}), 500

@app.route('/scheduler/pause/<job_id>', methods=['POST'])
def pause_scheduled_job(job_id: str):
    """Pause a specific scheduled job"""
    try:
        if scheduler is None:
            return jsonify({"error": "Scheduler not initialized"}), 500
        
        scheduler.pause_schedule(job_id)
        
        return jsonify({
            "status": "success",
            "message": f"Job {job_id} paused successfully"
        })
        
    except Exception as e:
        logger.error(f"Error pausing job {job_id}: {e}")
        return jsonify({"error": f"Failed to pause job: {str(e)}"}), 500

@app.route('/scheduler/resume/<job_id>', methods=['POST'])
def resume_scheduled_job(job_id: str):
    """Resume a specific scheduled job"""
    try:
        if scheduler is None:
            return jsonify({"error": "Scheduler not initialized"}), 500
        
        scheduler.resume_schedule(job_id)
        
        return jsonify({
            "status": "success",
            "message": f"Job {job_id} resumed successfully"
        })
        
    except Exception as e:
        logger.error(f"Error resuming job {job_id}: {e}")
        return jsonify({"error": f"Failed to resume job: {str(e)}"}), 500

# Model Version Management Endpoints
@app.route('/models/versions', methods=['GET'])
def list_model_versions():
    """List all model versions"""
    try:
        if version_manager is None:
            return jsonify({"error": "Version manager not initialized"}), 500
        
        model_type = request.args.get('model_type')
        
        if model_type:
            versions = version_manager.list_versions(model_type)
            return jsonify({
                "status": "success",
                "model_type": model_type,
                "versions": versions
            })
        else:
            summary = version_manager.get_version_summary()
            return jsonify({
                "status": "success",
                "version_summary": summary
            })
        
    except Exception as e:
        logger.error(f"Error listing model versions: {e}")
        return jsonify({"error": f"Failed to list versions: {str(e)}"}), 500

@app.route('/models/rollback', methods=['POST'])
def rollback_model():
    """Rollback to a specific model version"""
    try:
        if version_manager is None:
            return jsonify({"error": "Version manager not initialized"}), 500
        
        data = request.json
        if not data:
            return jsonify({"error": "No rollback data provided"}), 400
        
        model_type = data.get('model_type')
        version_id = data.get('version_id')
        
        if not model_type or not version_id:
            return jsonify({"error": "model_type and version_id are required"}), 400
        
        success = version_manager.rollback_model(model_type, version_id)
        
        if success:
            # Reload the model in the application
            global anomaly_detector, user_clusterer, recommendation_model, recommendation_matrix
            
            if model_type == 'anomaly':
                anomaly_detector, _ = version_manager.load_model('anomaly')
            elif model_type == 'clustering':
                user_clusterer, _ = version_manager.load_model('clustering')
            elif model_type == 'recommendation':
                model_data, _ = version_manager.load_model('recommendation')
                recommendation_model = model_data['similarity_matrix']
                recommendation_matrix = model_data['user_item_matrix']
            
            return jsonify({
                "status": "success",
                "message": f"Successfully rolled back {model_type} to version {version_id}"
            })
        else:
            return jsonify({"error": "Rollback failed"}), 500
        
    except Exception as e:
        logger.error(f"Error rolling back model: {e}")
        return jsonify({"error": f"Rollback failed: {str(e)}"}), 500

@app.route('/models/performance', methods=['POST'])
def update_model_performance():
    """Update performance metrics for a model version"""
    try:
        if version_manager is None:
            return jsonify({"error": "Version manager not initialized"}), 500
        
        data = request.json
        if not data:
            return jsonify({"error": "No performance data provided"}), 400
        
        model_type = data.get('model_type')
        version_id = data.get('version_id')
        performance_metrics = data.get('performance_metrics', {})
        
        if not model_type or not version_id:
            return jsonify({"error": "model_type and version_id are required"}), 400
        
        success = version_manager.update_performance_metrics(model_type, version_id, performance_metrics)
        
        if success:
            return jsonify({
                "status": "success",
                "message": f"Performance metrics updated for {model_type} version {version_id}"
            })
        else:
            return jsonify({"error": "Failed to update performance metrics"}), 500
        
    except Exception as e:
        logger.error(f"Error updating performance metrics: {e}")
        return jsonify({"error": f"Failed to update metrics: {str(e)}"}), 500

@app.route('/models/best', methods=['GET'])
def get_best_model_version():
    """Get the best performing version for a model type"""
    try:
        if version_manager is None:
            return jsonify({"error": "Version manager not initialized"}), 500
        
        model_type = request.args.get('model_type')
        metric = request.args.get('metric', 'accuracy')
        
        if not model_type:
            return jsonify({"error": "model_type parameter is required"}), 400
        
        best_version = version_manager.get_best_performing_version(model_type, metric)
        
        if best_version:
            return jsonify({
                "status": "success",
                "model_type": model_type,
                "metric": metric,
                "best_version": best_version
            })
        else:
            return jsonify({
                "status": "success",
                "message": f"No performance data found for {metric} in {model_type} models"
            })
        
    except Exception as e:
        logger.error(f"Error finding best model version: {e}")
        return jsonify({"error": f"Failed to find best version: {str(e)}"}), 500

@app.route('/models/auto-rollback', methods=['POST'])
def check_auto_rollback():
    """Check for performance degradation and auto-rollback if needed"""
    try:
        if version_manager is None:
            return jsonify({"error": "Version manager not initialized"}), 500
        
        data = request.json
        if not data:
            return jsonify({"error": "No performance data provided"}), 400
        
        model_type = data.get('model_type')
        current_metrics = data.get('current_metrics', {})
        threshold_metric = data.get('threshold_metric', 'accuracy')
        degradation_threshold = data.get('degradation_threshold', 0.05)
        
        if not model_type:
            return jsonify({"error": "model_type is required"}), 400
        
        rollback_performed = version_manager.auto_rollback_on_degradation(
            model_type, current_metrics, threshold_metric, degradation_threshold
        )
        
        if rollback_performed:
            # Reload the model in the application
            global anomaly_detector, user_clusterer, recommendation_model, recommendation_matrix
            
            if model_type == 'anomaly':
                anomaly_detector, _ = version_manager.load_model('anomaly')
            elif model_type == 'clustering':
                user_clusterer, _ = version_manager.load_model('clustering')
            elif model_type == 'recommendation':
                model_data, _ = version_manager.load_model('recommendation')
                recommendation_model = model_data['similarity_matrix']
                recommendation_matrix = model_data['user_item_matrix']
            
            return jsonify({
                "status": "success",
                "rollback_performed": True,
                "message": f"Auto-rollback performed for {model_type} due to performance degradation"
            })
        else:
            return jsonify({
                "status": "success",
                "rollback_performed": False,
                "message": f"No rollback needed for {model_type}"
            })
        
    except Exception as e:
        logger.error(f"Error in auto-rollback check: {e}")
        return jsonify({"error": f"Auto-rollback check failed: {str(e)}"}), 500

if __name__ == '__main__':
    # Initialize services
    initialize_services()
    load_or_create_models()
    app.run(host='0.0.0.0', port=5000) 