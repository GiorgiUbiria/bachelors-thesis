import numpy as np
from ml.models.trend import TrendPredictor

# Placeholder: Replace with real data loading
X = np.arange(24).reshape(-1, 1)  # e.g., months
y = np.random.rand(24) * 100  # e.g., sales

model = TrendPredictor()
model.fit(X, y)
model.save('trend_model.joblib')
print('Trend model trained and saved.') 