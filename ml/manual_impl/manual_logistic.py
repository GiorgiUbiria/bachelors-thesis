import numpy as np
from typing import Optional

class ManualLogisticRegression:
    def __init__(self, lr: float = 0.01, n_iter: int = 1000, regularization: Optional[str] = None, lambda_reg: float = 0.01):
        """
        Manual implementation of Logistic Regression with regularization
        
        Args:
            lr: Learning rate
            n_iter: Number of iterations
            regularization: Type of regularization ('l1', 'l2', or None)
            lambda_reg: Regularization strength
        """
        self.lr = lr
        self.n_iter = n_iter
        self.regularization = regularization
        self.lambda_reg = lambda_reg
        self.weights = None
        self.bias = None
        self.cost_history = []
        self.feature_means = None
        self.feature_stds = None

    def _normalize_features(self, X: np.ndarray, fit: bool = False) -> np.ndarray:
        """Normalize features using z-score normalization"""
        if fit:
            self.feature_means = np.mean(X, axis=0)
            self.feature_stds = np.std(X, axis=0)
            # Avoid division by zero
            self.feature_stds = np.where(self.feature_stds == 0, 1, self.feature_stds)
        
        return (X - self.feature_means) / self.feature_stds

    def _sigmoid(self, x: np.ndarray) -> np.ndarray:
        """Sigmoid activation function with numerical stability"""
        # Clip x to prevent overflow
        x = np.clip(x, -500, 500)
        return 1 / (1 + np.exp(-x))

    def _compute_cost(self, y_true: np.ndarray, y_pred: np.ndarray) -> float:
        """Compute logistic regression cost with regularization"""
        # Avoid log(0) by adding small epsilon
        epsilon = 1e-15
        y_pred = np.clip(y_pred, epsilon, 1 - epsilon)
        
        # Binary cross-entropy loss
        cost = -np.mean(y_true * np.log(y_pred) + (1 - y_true) * np.log(1 - y_pred))
        
        # Add regularization
        if self.regularization == 'l1':
            cost += self.lambda_reg * np.sum(np.abs(self.weights))
        elif self.regularization == 'l2':
            cost += self.lambda_reg * np.sum(self.weights ** 2)
        
        return cost

    def fit(self, X: np.ndarray, y: np.ndarray) -> 'ManualLogisticRegression':
        """
        Train the logistic regression model
        
        Args:
            X: Training features (n_samples, n_features)
            y: Training labels (n_samples,)
        """
        # Normalize features
        X_norm = self._normalize_features(X, fit=True)
        
        n_samples, n_features = X_norm.shape
        
        # Initialize parameters
        self.weights = np.random.normal(0, 0.01, n_features)
        self.bias = 0
        self.cost_history = []
        
        # Training loop
        for i in range(self.n_iter):
            # Forward pass
            linear_model = np.dot(X_norm, self.weights) + self.bias
            y_pred = self._sigmoid(linear_model)
            
            # Compute cost
            cost = self._compute_cost(y, y_pred)
            self.cost_history.append(cost)
            
            # Compute gradients
            dw = (1 / n_samples) * np.dot(X_norm.T, (y_pred - y))
            db = (1 / n_samples) * np.sum(y_pred - y)
            
            # Add regularization to weight gradients
            if self.regularization == 'l1':
                dw += self.lambda_reg * np.sign(self.weights)
            elif self.regularization == 'l2':
                dw += 2 * self.lambda_reg * self.weights
            
            # Update parameters
            self.weights -= self.lr * dw
            self.bias -= self.lr * db
            
            # Early stopping if cost doesn't improve
            if i > 10 and abs(self.cost_history[-1] - self.cost_history[-2]) < 1e-8:
                print(f"Early stopping at iteration {i}")
                break
        
        return self

    def predict_proba(self, X: np.ndarray) -> np.ndarray:
        """Predict class probabilities"""
        X_norm = self._normalize_features(X, fit=False)
        linear_model = np.dot(X_norm, self.weights) + self.bias
        return self._sigmoid(linear_model)

    def predict(self, X: np.ndarray) -> np.ndarray:
        """Predict binary classes"""
        probabilities = self.predict_proba(X)
        return np.where(probabilities > 0.5, 1, 0)

    def score(self, X: np.ndarray, y: np.ndarray) -> float:
        """Calculate accuracy score"""
        predictions = self.predict(X)
        return np.mean(predictions == y)

    def get_feature_importance(self) -> np.ndarray:
        """Get feature importance based on absolute weights"""
        if self.weights is None:
            raise ValueError("Model must be fitted before getting feature importance")
        return np.abs(self.weights) 