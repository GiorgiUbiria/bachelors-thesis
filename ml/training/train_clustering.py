import numpy as np
from ml.models.clustering import UserClustering

# Placeholder: Replace with real data loading
X = np.random.rand(500, 5)

model = UserClustering(n_clusters=3)
model.fit(X)
model.save('clustering_model.joblib')
print('Clustering model trained and saved.') 