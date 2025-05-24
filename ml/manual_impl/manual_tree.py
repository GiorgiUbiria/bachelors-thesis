import numpy as np
from typing import Dict, Any, Optional, Union
from collections import Counter

class ManualDecisionTree:
    def __init__(self, max_depth: int = 3, min_samples_split: int = 2, min_samples_leaf: int = 1, 
                 criterion: str = 'gini', max_features: Optional[int] = None):
        """
        Manual implementation of Decision Tree with information gain and pruning
        
        Args:
            max_depth: Maximum depth of the tree
            min_samples_split: Minimum samples required to split a node
            min_samples_leaf: Minimum samples required at a leaf node
            criterion: Splitting criterion ('gini' or 'entropy')
            max_features: Maximum number of features to consider for splitting
        """
        self.max_depth = max_depth
        self.min_samples_split = min_samples_split
        self.min_samples_leaf = min_samples_leaf
        self.criterion = criterion
        self.max_features = max_features
        self.tree = None
        self.feature_importances_ = None
        self.n_features_ = None
        self.n_classes_ = None

    def _gini_impurity(self, y: np.ndarray) -> float:
        """Calculate Gini impurity"""
        if len(y) == 0:
            return 0
        
        _, counts = np.unique(y, return_counts=True)
        probabilities = counts / len(y)
        return 1 - np.sum(probabilities ** 2)

    def _entropy(self, y: np.ndarray) -> float:
        """Calculate entropy"""
        if len(y) == 0:
            return 0
        
        _, counts = np.unique(y, return_counts=True)
        probabilities = counts / len(y)
        # Avoid log(0)
        probabilities = probabilities[probabilities > 0]
        return -np.sum(probabilities * np.log2(probabilities))

    def _calculate_impurity(self, y: np.ndarray) -> float:
        """Calculate impurity based on criterion"""
        if self.criterion == 'gini':
            return self._gini_impurity(y)
        elif self.criterion == 'entropy':
            return self._entropy(y)
        else:
            raise ValueError(f"Unknown criterion: {self.criterion}")

    def _information_gain(self, y: np.ndarray, left_y: np.ndarray, right_y: np.ndarray) -> float:
        """Calculate information gain from a split"""
        if len(y) == 0:
            return 0
        
        # Calculate weighted impurity after split
        n_total = len(y)
        n_left = len(left_y)
        n_right = len(right_y)
        
        if n_left == 0 or n_right == 0:
            return 0
        
        impurity_before = self._calculate_impurity(y)
        impurity_after = (n_left / n_total) * self._calculate_impurity(left_y) + \
                        (n_right / n_total) * self._calculate_impurity(right_y)
        
        return impurity_before - impurity_after

    def _find_best_split(self, X: np.ndarray, y: np.ndarray) -> tuple:
        """Find the best feature and threshold for splitting"""
        best_gain = -1
        best_feature = None
        best_threshold = None
        
        n_features = X.shape[1]
        
        # Select features to consider
        if self.max_features is not None:
            features_to_consider = np.random.choice(n_features, 
                                                  min(self.max_features, n_features), 
                                                  replace=False)
        else:
            features_to_consider = range(n_features)
        
        for feature in features_to_consider:
            # Get unique values for this feature
            unique_values = np.unique(X[:, feature])
            
            # Try each unique value as a threshold
            for i in range(len(unique_values) - 1):
                threshold = (unique_values[i] + unique_values[i + 1]) / 2
                
                # Split data
                left_mask = X[:, feature] <= threshold
                right_mask = ~left_mask
                
                left_y = y[left_mask]
                right_y = y[right_mask]
                
                # Check minimum samples constraint
                if len(left_y) < self.min_samples_leaf or len(right_y) < self.min_samples_leaf:
                    continue
                
                # Calculate information gain
                gain = self._information_gain(y, left_y, right_y)
                
                if gain > best_gain:
                    best_gain = gain
                    best_feature = feature
                    best_threshold = threshold
        
        return best_feature, best_threshold, best_gain

    def _build_tree(self, X: np.ndarray, y: np.ndarray, depth: int = 0) -> Union[Dict[str, Any], int]:
        """Recursively build the decision tree"""
        # Check stopping criteria
        if (depth >= self.max_depth or 
            len(y) < self.min_samples_split or 
            len(np.unique(y)) == 1 or
            len(y) < 2 * self.min_samples_leaf):
            # Return most common class
            return Counter(y).most_common(1)[0][0]
        
        # Find best split
        best_feature, best_threshold, best_gain = self._find_best_split(X, y)
        
        # If no good split found, return leaf
        if best_feature is None or best_gain <= 0:
            return Counter(y).most_common(1)[0][0]
        
        # Split data
        left_mask = X[:, best_feature] <= best_threshold
        right_mask = ~left_mask
        
        # Update feature importance
        if self.feature_importances_ is not None:
            self.feature_importances_[best_feature] += best_gain * len(y) / len(self.y_train_)
        
        # Recursively build subtrees
        left_subtree = self._build_tree(X[left_mask], y[left_mask], depth + 1)
        right_subtree = self._build_tree(X[right_mask], y[right_mask], depth + 1)
        
        return {
            'feature': best_feature,
            'threshold': best_threshold,
            'left': left_subtree,
            'right': right_subtree,
            'gain': best_gain,
            'samples': len(y),
            'impurity': self._calculate_impurity(y)
        }

    def fit(self, X: np.ndarray, y: np.ndarray) -> 'ManualDecisionTree':
        """
        Train the decision tree
        
        Args:
            X: Training features (n_samples, n_features)
            y: Training labels (n_samples,)
        """
        self.n_features_ = X.shape[1]
        self.n_classes_ = len(np.unique(y))
        self.y_train_ = y  # Store for feature importance calculation
        
        # Initialize feature importances
        self.feature_importances_ = np.zeros(self.n_features_)
        
        # Build tree
        self.tree = self._build_tree(X, y)
        
        # Normalize feature importances
        if np.sum(self.feature_importances_) > 0:
            self.feature_importances_ /= np.sum(self.feature_importances_)
        
        return self

    def _predict_one(self, x: np.ndarray, node: Union[Dict[str, Any], int]) -> int:
        """Predict class for a single sample"""
        # If leaf node, return class
        if not isinstance(node, dict):
            return node
        
        # Navigate tree based on feature value
        if x[node['feature']] <= node['threshold']:
            return self._predict_one(x, node['left'])
        else:
            return self._predict_one(x, node['right'])

    def predict(self, X: np.ndarray) -> np.ndarray:
        """Predict classes for multiple samples"""
        if self.tree is None:
            raise ValueError("Model must be fitted before making predictions")
        
        return np.array([self._predict_one(x, self.tree) for x in X])

    def score(self, X: np.ndarray, y: np.ndarray) -> float:
        """Calculate accuracy score"""
        predictions = self.predict(X)
        return np.mean(predictions == y)

    def _prune_node(self, node: Dict[str, Any], X_val: np.ndarray, y_val: np.ndarray) -> Union[Dict[str, Any], int]:
        """Prune a node if it improves validation accuracy"""
        if not isinstance(node, dict):
            return node
        
        # Recursively prune children
        node['left'] = self._prune_node(node['left'], X_val, y_val)
        node['right'] = self._prune_node(node['right'], X_val, y_val)
        
        # If both children are leaves, consider pruning this node
        if (not isinstance(node['left'], dict) and 
            not isinstance(node['right'], dict)):
            
            # Calculate accuracy before pruning
            original_tree = self.tree
            accuracy_before = self.score(X_val, y_val)
            
            # Try pruning (replace with most common class)
            most_common_class = Counter(y_val).most_common(1)[0][0]
            
            # Temporarily replace node with leaf
            temp_tree = self.tree
            self.tree = most_common_class
            accuracy_after = self.score(X_val, y_val)
            
            # Restore original tree
            self.tree = temp_tree
            
            # Prune if accuracy doesn't decrease
            if accuracy_after >= accuracy_before:
                return most_common_class
        
        return node

    def prune(self, X_val: np.ndarray, y_val: np.ndarray) -> 'ManualDecisionTree':
        """
        Prune the tree using validation data
        
        Args:
            X_val: Validation features
            y_val: Validation labels
        """
        if self.tree is None:
            raise ValueError("Model must be fitted before pruning")
        
        self.tree = self._prune_node(self.tree, X_val, y_val)
        return self

    def get_depth(self, node: Union[Dict[str, Any], int] = None) -> int:
        """Get the depth of the tree"""
        if node is None:
            node = self.tree
        
        if not isinstance(node, dict):
            return 0
        
        return 1 + max(self.get_depth(node['left']), self.get_depth(node['right']))

    def get_n_leaves(self, node: Union[Dict[str, Any], int] = None) -> int:
        """Get the number of leaves in the tree"""
        if node is None:
            node = self.tree
        
        if not isinstance(node, dict):
            return 1
        
        return self.get_n_leaves(node['left']) + self.get_n_leaves(node['right'])

    def print_tree(self, node: Union[Dict[str, Any], int] = None, depth: int = 0) -> None:
        """Print the tree structure"""
        if node is None:
            node = self.tree
        
        indent = "  " * depth
        
        if not isinstance(node, dict):
            print(f"{indent}Leaf: class = {node}")
        else:
            print(f"{indent}Feature {node['feature']} <= {node['threshold']:.3f}")
            print(f"{indent}├─ True:")
            self.print_tree(node['left'], depth + 1)
            print(f"{indent}└─ False:")
            self.print_tree(node['right'], depth + 1) 