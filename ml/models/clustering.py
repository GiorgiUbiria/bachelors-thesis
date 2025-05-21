import joblib
from sklearn.cluster import KMeans
import numpy as np

class UserClustering:
    def __init__(self, n_clusters=3, **kwargs):
        self.model = KMeans(n_clusters=n_clusters, **kwargs)

    def fit(self, X):
        self.model.fit(X)

    def predict(self, X):
        return self.model.predict(X)

    def save(self, path):
        joblib.dump(self.model, path)

    def load(self, path):
        self.model = joblib.load(path) 