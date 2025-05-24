# Load data from the DB or CSV (reusable)
import numpy as np
import pandas as pd
import requests
import os
from typing import Tuple, Optional, Dict, Any
import logging
from datetime import datetime

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# API Configuration
API_URL = os.getenv('API_URL', 'http://api:8080')
API_TIMEOUT = int(os.getenv('API_TIMEOUT', '30'))

def _make_api_request(endpoint: str, params: Optional[Dict] = None) -> Optional[Dict[Any, Any]]:
    """Make API request with error handling and retries"""
    try:
        url = f"{API_URL}{endpoint}"
        logger.info(f"Making API request to: {url}")
        
        response = requests.get(url, params=params, timeout=API_TIMEOUT)
        response.raise_for_status()
        
        data = response.json()
        logger.info(f"Successfully fetched data from {endpoint}")
        return data
        
    except requests.exceptions.Timeout:
        logger.error(f"Timeout error for endpoint {endpoint}")
        return None
    except requests.exceptions.ConnectionError:
        logger.error(f"Connection error for endpoint {endpoint}")
        return None
    except requests.exceptions.HTTPError as e:
        logger.error(f"HTTP error for endpoint {endpoint}: {e}")
        return None
    except Exception as e:
        logger.error(f"Unexpected error for endpoint {endpoint}: {e}")
        return None

def load_anomaly_data() -> Optional[np.ndarray]:
    """Load and return data for anomaly detection from backend API."""
    data = _make_api_request('/api/requests/logs')
    
    if data is None:
        logger.warning("Failed to fetch anomaly data, using fallback mock data")
        return np.random.randn(1000, 3)  # Reduced features to match expected format
    
    try:
        # Validate and extract features
        request_logs = data.get('request_logs', [])
        if not request_logs:
            logger.warning("No request logs found, using mock data")
            return np.random.randn(1000, 3)
        
        features = []
        for log in request_logs:
            # Extract features: response_time, request_size, error_count
            feature_row = [
                float(log.get('response_time', 0)),
                float(log.get('request_size', 0)),
                float(log.get('error_count', 0))
            ]
            features.append(feature_row)
        
        if len(features) < 10:  # Need minimum samples for training
            logger.warning("Insufficient data samples, using mock data")
            return np.random.randn(1000, 3)
            
        logger.info(f"Loaded {len(features)} anomaly detection samples")
        return np.array(features)
        
    except Exception as e:
        logger.error(f"Error processing anomaly data: {e}")
        return np.random.randn(1000, 3)

def load_clustering_data() -> Optional[np.ndarray]:
    """Load and return data for user clustering from backend API."""
    data = _make_api_request('/api/users/activity-stats')
    
    if data is None:
        logger.warning("Failed to fetch clustering data, using fallback mock data")
        return np.random.rand(500, 4)  # 4 features for user activity
    
    try:
        # Validate and extract features
        user_activities = data.get('user_activities', [])
        if not user_activities:
            logger.warning("No user activities found, using mock data")
            return np.random.rand(500, 4)
        
        features = []
        for activity in user_activities:
            # Extract features: login_count, purchase_count, cart_count, favorite_count
            feature_row = [
                float(activity.get('login_count', 0)),
                float(activity.get('purchase_count', 0)),
                float(activity.get('cart_count', 0)),
                float(activity.get('favorite_count', 0))
            ]
            features.append(feature_row)
        
        if len(features) < 5:  # Need minimum samples for clustering
            logger.warning("Insufficient user activity samples, using mock data")
            return np.random.rand(500, 4)
            
        logger.info(f"Loaded {len(features)} user clustering samples")
        return np.array(features)
        
    except Exception as e:
        logger.error(f"Error processing clustering data: {e}")
        return np.random.rand(500, 4)

def load_recommendation_data() -> Optional[Tuple[np.ndarray, pd.DataFrame]]:
    """Load and return user-item matrix for recommendations from backend API."""
    data = _make_api_request('/api/purchases/user-item-matrix')
    
    if data is None:
        logger.warning("Failed to fetch recommendation data, using fallback mock data")
        return np.random.randint(0, 2, (100, 50)), pd.DataFrame()
    
    try:
        # Validate and extract purchase data
        purchases = data.get('purchases', [])
        if not purchases:
            logger.warning("No purchase data found, using mock data")
            return np.random.randint(0, 2, (100, 50)), pd.DataFrame()
        
        # Convert to DataFrame for easier manipulation
        df = pd.DataFrame(purchases)
        
        # Validate required columns
        required_cols = ['user_id', 'product_id', 'quantity']
        if not all(col in df.columns for col in required_cols):
            logger.warning("Missing required columns in purchase data, using mock data")
            return np.random.randint(0, 2, (100, 50)), pd.DataFrame()
        
        # Create user-item matrix
        user_item_matrix = df.pivot_table(
            index='user_id',
            columns='product_id',
            values='quantity',
            fill_value=0,
            aggfunc='sum'
        )
        
        if user_item_matrix.empty or user_item_matrix.shape[0] < 5:
            logger.warning("Insufficient purchase data, using mock data")
            return np.random.randint(0, 2, (100, 50)), pd.DataFrame()
            
        logger.info(f"Loaded user-item matrix: {user_item_matrix.shape}")
        return user_item_matrix.values, user_item_matrix
        
    except Exception as e:
        logger.error(f"Error processing recommendation data: {e}")
        return np.random.randint(0, 2, (100, 50)), pd.DataFrame()

def load_trend_data() -> Optional[Tuple[np.ndarray, np.ndarray]]:
    """Load and return data for trend analysis from backend API."""
    data = _make_api_request('/api/analytics/sales-trends')
    
    if data is None:
        logger.warning("Failed to fetch trend data, using fallback mock data")
        X = np.arange(24).reshape(-1, 1)
        y = np.random.rand(24) * 100
        return X, y
    
    try:
        # Validate and extract trend data
        sales_trends = data.get('sales_trends', [])
        if not sales_trends:
            logger.warning("No sales trend data found, using mock data")
            X = np.arange(24).reshape(-1, 1)
            y = np.random.rand(24) * 100
            return X, y
        
        # Convert to DataFrame for easier manipulation
        df = pd.DataFrame(sales_trends)
        
        # Validate required columns
        required_cols = ['date', 'sales_count']
        if not all(col in df.columns for col in required_cols):
            logger.warning("Missing required columns in trend data, using mock data")
            X = np.arange(24).reshape(-1, 1)
            y = np.random.rand(24) * 100
            return X, y
        
        # Sort by date and prepare features
        df['date'] = pd.to_datetime(df['date'])
        df = df.sort_values('date')
        
        # Create time-based features (days since first date)
        first_date = df['date'].min()
        df['days_since_start'] = (df['date'] - first_date).dt.days
        
        X = df['days_since_start'].values.reshape(-1, 1)
        y = df['sales_count'].values
        
        if len(X) < 5:  # Need minimum samples for trend analysis
            logger.warning("Insufficient trend data samples, using mock data")
            X = np.arange(24).reshape(-1, 1)
            y = np.random.rand(24) * 100
            return X, y
            
        logger.info(f"Loaded {len(X)} trend analysis samples")
        return X, y
        
    except Exception as e:
        logger.error(f"Error processing trend data: {e}")
        X = np.arange(24).reshape(-1, 1)
        y = np.random.rand(24) * 100
        return X, y

def validate_anomaly_data(data: Dict[str, Any]) -> bool:
    """Validate incoming request log data format"""
    required_fields = ['response_time', 'request_size', 'error_count']
    
    if not isinstance(data, dict):
        return False
    
    request_logs = data.get('request_logs', [])
    if not isinstance(request_logs, list) or not request_logs:
        return False
    
    # Check first few samples for required fields
    for log in request_logs[:5]:
        if not all(field in log for field in required_fields):
            return False
        
        # Check data types
        try:
            float(log['response_time'])
            float(log['request_size'])
            float(log['error_count'])
        except (ValueError, TypeError):
            return False
    
    return True

def validate_clustering_data(data: Dict[str, Any]) -> bool:
    """Validate user activity data format"""
    required_fields = ['login_count', 'purchase_count', 'cart_count', 'favorite_count']
    
    if not isinstance(data, dict):
        return False
    
    user_activities = data.get('user_activities', [])
    if not isinstance(user_activities, list) or not user_activities:
        return False
    
    # Check first few samples for required fields
    for activity in user_activities[:5]:
        if not all(field in activity for field in required_fields):
            return False
        
        # Check data types
        try:
            float(activity['login_count'])
            float(activity['purchase_count'])
            float(activity['cart_count'])
            float(activity['favorite_count'])
        except (ValueError, TypeError):
            return False
    
    return True

def validate_recommendation_data(data: Dict[str, Any]) -> bool:
    """Validate purchase data format"""
    required_fields = ['user_id', 'product_id', 'quantity']
    
    if not isinstance(data, dict):
        return False
    
    purchases = data.get('purchases', [])
    if not isinstance(purchases, list) or not purchases:
        return False
    
    # Check first few samples for required fields
    for purchase in purchases[:5]:
        if not all(field in purchase for field in required_fields):
            return False
        
        # Check data types
        try:
            int(purchase['user_id'])
            int(purchase['product_id'])
            float(purchase['quantity'])
        except (ValueError, TypeError):
            return False
    
    return True

def validate_trend_data(data: Dict[str, Any]) -> bool:
    """Validate sales trend data format"""
    required_fields = ['date', 'sales_count']
    
    if not isinstance(data, dict):
        return False
    
    sales_trends = data.get('sales_trends', [])
    if not isinstance(sales_trends, list) or not sales_trends:
        return False
    
    # Check first few samples for required fields
    for trend in sales_trends[:5]:
        if not all(field in trend for field in required_fields):
            return False
        
        # Check data types
        try:
            pd.to_datetime(trend['date'])
            float(trend['sales_count'])
        except (ValueError, TypeError):
            return False
    
    return True