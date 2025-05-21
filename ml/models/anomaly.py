import joblib
from sklearn.ensemble import IsolationForest
from sklearn.svm import OneClassSVM
import numpy as np

class AnomalyDetector:
    def __init__(self, method='isolation_forest', **kwargs):
        if method == 'isolation_forest':
            self.model = IsolationForest(**kwargs)
        elif method == 'one_class_svm':
            self.model = OneClassSVM(**kwargs)
        else:
            raise ValueError('Unknown method')
        self.method = method

    def fit(self, X):
        self.model.fit(X)

    def predict(self, X):
        return self.model.predict(X)

    def save(self, path):
        joblib.dump(self.model, path)

    def load(self, path):
        self.model = joblib.load(path) 