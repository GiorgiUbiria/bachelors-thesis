# Load data from the DB or CSV (reusable)
import numpy as np

def load_anomaly_data():
    """Load and return data for anomaly detection (replace with real DB/CSV logic)."""
    return np.random.randn(1000, 10)

def load_clustering_data():
    """Load and return data for user clustering (replace with real DB/CSV logic)."""
    return np.random.rand(500, 5)

def load_recommendation_data():
    """Load and return user-item matrix for recommendations (replace with real DB/CSV logic)."""
    return np.random.randint(0, 2, (100, 50))

def load_trend_data():
    """Load and return data for trend analysis (replace with real DB/CSV logic)."""
    X = np.arange(24).reshape(-1, 1)
    y = np.random.rand(24) * 100
    return X, y