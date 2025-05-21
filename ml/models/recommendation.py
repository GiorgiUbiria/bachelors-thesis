import joblib
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

class ItemBasedRecommender:
    def __init__(self):
        self.item_similarity = None
        self.user_item_matrix = None

    def fit(self, user_item_matrix):
        self.user_item_matrix = user_item_matrix
        self.item_similarity = cosine_similarity(user_item_matrix.T)

    def recommend(self, user_index, top_k=5):
        user_ratings = self.user_item_matrix[user_index]
        scores = self.item_similarity.dot(user_ratings)
        recommended = np.argsort(scores)[::-1]
        return recommended[:top_k]

    def save(self, path):
        joblib.dump((self.item_similarity, self.user_item_matrix), path)

    def load(self, path):
        self.item_similarity, self.user_item_matrix = joblib.load(path) 