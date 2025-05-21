# Centralized config (e.g., file paths, model params)
ANOMALY_MODEL_PATH = 'anomaly_model.joblib'
CLUSTERING_MODEL_PATH = 'clustering_model.joblib'
RECOMMENDER_MODEL_PATH = 'recommender_model.joblib'
TREND_MODEL_PATH = 'trend_model.joblib'

# Example model parameters
def get_anomaly_params():
    return {'n_estimators': 100}

def get_clustering_params():
    return {'n_clusters': 3}