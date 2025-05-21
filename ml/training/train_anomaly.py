import numpy as np
from ml.models.anomaly import AnomalyDetector

# Placeholder: Replace with real data loading
X = np.random.randn(1000, 10)

model = AnomalyDetector(method='isolation_forest', n_estimators=100)
model.fit(X)
model.save('anomaly_model.joblib')
print('Anomaly model trained and saved.') 