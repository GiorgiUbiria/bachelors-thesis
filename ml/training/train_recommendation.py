import numpy as np
from ml.models.recommendation import ItemBasedRecommender

# Placeholder: Replace with real data loading
user_item_matrix = np.random.randint(0, 2, (100, 50))

model = ItemBasedRecommender()
model.fit(user_item_matrix)
model.save('recommender_model.joblib')
print('Recommendation model trained and saved.') 