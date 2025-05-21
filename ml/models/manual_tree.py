import numpy as np

class ManualDecisionTree:
    def __init__(self, max_depth=3):
        self.max_depth = max_depth
        self.tree = None

    def fit(self, X, y):
        self.tree = self._build_tree(X, y, depth=0)

    def predict(self, X):
        return np.array([self._predict_one(x, self.tree) for x in X])

    def _build_tree(self, X, y, depth):
        if depth >= self.max_depth or len(set(y)) == 1:
            return np.bincount(y).argmax()
        best_feat = 0  # For demo, always split on first feature
        thresh = np.median(X[:, best_feat])
        left = X[:, best_feat] < thresh
        right = ~left
        return {
            'feature': best_feat,
            'thresh': thresh,
            'left': self._build_tree(X[left], y[left], depth+1),
            'right': self._build_tree(X[right], y[right], depth+1)
        }

    def _predict_one(self, x, node):
        if not isinstance(node, dict):
            return node
        if x[node['feature']] < node['thresh']:
            return self._predict_one(x, node['left'])
        else:
            return self._predict_one(x, node['right']) 