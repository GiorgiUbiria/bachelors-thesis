# Define expected input/output formats

# Anomaly Detection
ANOMALY_INPUT_SCHEMA = {
    'features': 'np.ndarray, shape (n_samples, n_features)'
}
ANOMALY_OUTPUT_SCHEMA = {
    'predictions': 'np.ndarray, shape (n_samples,); values: 1 (normal), -1 (anomaly)'
}

# Clustering
CLUSTERING_INPUT_SCHEMA = {
    'features': 'np.ndarray, shape (n_samples, n_features)'
}
CLUSTERING_OUTPUT_SCHEMA = {
    'labels': 'np.ndarray, shape (n_samples,); cluster indices'
}

# Recommendation
RECOMMENDATION_INPUT_SCHEMA = {
    'user_index': 'int',
    'user_item_matrix': 'np.ndarray, shape (n_users, n_items)'
}
RECOMMENDATION_OUTPUT_SCHEMA = {
    'recommended_items': 'np.ndarray, shape (top_k,); item indices'
}

# Trend Analysis
TREND_INPUT_SCHEMA = {
    'X': 'np.ndarray, shape (n_samples, n_features)'
}
TREND_OUTPUT_SCHEMA = {
    'y_pred': 'np.ndarray, shape (n_samples,)'
}