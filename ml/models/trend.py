import joblib
from sklearn.linear_model import LinearRegression
import numpy as np

class TrendPredictor:
    def __init__(self, **kwargs):
        self.model = LinearRegression(**kwargs)

    def fit(self, X, y):
        self.model.fit(X, y)

    def predict(self, X):
        return self.model.predict(X)

    def save(self, path):
        joblib.dump(self.model, path)

    def load(self, path):
        self.model = joblib.load(path) 