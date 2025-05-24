#!/usr/bin/env python3
"""
Test script for ML service functionality
"""

import requests
import json
import numpy as np
from sklearn.datasets import make_classification, make_blobs
import time

# Configuration
ML_SERVICE_URL = "http://localhost:5000"

def test_health_check():
    """Test health check endpoint"""
    print("Testing health check...")
    try:
        response = requests.get(f"{ML_SERVICE_URL}/health")
        if response.status_code == 200:
            data = response.json()
            print(f"✓ Health check passed: {data}")
            return True
        else:
            print(f"✗ Health check failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Health check error: {e}")
        return False

def test_manual_logistic_regression():
    """Test manual logistic regression implementation"""
    print("\nTesting manual logistic regression...")
    try:
        # Generate synthetic data
        X, y = make_classification(n_samples=1000, n_features=4, n_classes=2, random_state=42)
        
        # Train model
        train_data = {
            "features": X.tolist(),
            "labels": y.tolist(),
            "learning_rate": 0.01,
            "n_iterations": 500,
            "regularization": "l2",
            "lambda_reg": 0.01
        }
        
        response = requests.post(f"{ML_SERVICE_URL}/train/manual-logistic", json=train_data)
        if response.status_code == 200:
            result = response.json()
            print(f"✓ Manual logistic regression training: {result['message']}")
            
            # Test prediction
            test_X = X[:10]  # Use first 10 samples for testing
            pred_data = {"features": test_X.tolist()}
            
            pred_response = requests.post(f"{ML_SERVICE_URL}/predict/manual-logistic", json=pred_data)
            if pred_response.status_code == 200:
                pred_result = pred_response.json()
                print(f"✓ Manual logistic regression prediction: {len(pred_result['predictions'])} predictions made")
                return True
            else:
                print(f"✗ Manual logistic regression prediction failed: {pred_response.status_code}")
                return False
        else:
            print(f"✗ Manual logistic regression training failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Manual logistic regression error: {e}")
        return False

def test_manual_decision_tree():
    """Test manual decision tree implementation"""
    print("\nTesting manual decision tree...")
    try:
        # Generate synthetic data
        X, y = make_classification(n_samples=500, n_features=4, n_classes=3, random_state=42)
        
        # Train model
        train_data = {
            "features": X.tolist(),
            "labels": y.tolist(),
            "max_depth": 5,
            "min_samples_split": 5,
            "criterion": "gini"
        }
        
        response = requests.post(f"{ML_SERVICE_URL}/train/manual-tree", json=train_data)
        if response.status_code == 200:
            result = response.json()
            print(f"✓ Manual decision tree training: {result['message']}")
            
            # Test prediction
            test_X = X[:10]  # Use first 10 samples for testing
            pred_data = {"features": test_X.tolist()}
            
            pred_response = requests.post(f"{ML_SERVICE_URL}/predict/manual-tree", json=pred_data)
            if pred_response.status_code == 200:
                pred_result = pred_response.json()
                print(f"✓ Manual decision tree prediction: {len(pred_result['predictions'])} predictions made")
                return True
            else:
                print(f"✗ Manual decision tree prediction failed: {pred_response.status_code}")
                return False
        else:
            print(f"✗ Manual decision tree training failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Manual decision tree error: {e}")
        return False

def test_model_comparison():
    """Test model comparison functionality"""
    print("\nTesting model comparison...")
    try:
        # Generate synthetic data
        X, y = make_classification(n_samples=1000, n_features=4, n_classes=2, random_state=42)
        
        comparison_data = {
            "features": X.tolist(),
            "labels": y.tolist(),
            "test_size": 0.2
        }
        
        response = requests.post(f"{ML_SERVICE_URL}/compare/models", json=comparison_data)
        if response.status_code == 200:
            result = response.json()
            print(f"✓ Model comparison completed")
            
            # Print comparison results
            comparison_results = result['comparison_results']
            print(f"  Logistic Regression - Manual: {comparison_results['logistic_regression']['manual']['accuracy']:.4f}")
            print(f"  Logistic Regression - Sklearn: {comparison_results['logistic_regression']['sklearn']['accuracy']:.4f}")
            print(f"  Decision Tree - Manual: {comparison_results['decision_tree']['manual']['accuracy']:.4f}")
            print(f"  Decision Tree - Sklearn: {comparison_results['decision_tree']['sklearn']['accuracy']:.4f}")
            
            return True
        else:
            print(f"✗ Model comparison failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Model comparison error: {e}")
        return False

def test_scheduler_status():
    """Test scheduler status endpoint"""
    print("\nTesting scheduler status...")
    try:
        response = requests.get(f"{ML_SERVICE_URL}/scheduler/status")
        if response.status_code == 200:
            result = response.json()
            print(f"✓ Scheduler status retrieved")
            scheduler_status = result['scheduler_status']
            print(f"  Scheduler running: {scheduler_status['scheduler_running']}")
            print(f"  Number of jobs: {len(scheduler_status['jobs'])}")
            return True
        else:
            print(f"✗ Scheduler status failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Scheduler status error: {e}")
        return False

def test_model_versions():
    """Test model versioning functionality"""
    print("\nTesting model versioning...")
    try:
        response = requests.get(f"{ML_SERVICE_URL}/models/versions")
        if response.status_code == 200:
            result = response.json()
            print(f"✓ Model versions retrieved")
            version_summary = result['version_summary']
            for model_type, info in version_summary.items():
                print(f"  {model_type}: {info['total_versions']} versions, current: {info['current_version']}")
            return True
        else:
            print(f"✗ Model versions failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Model versions error: {e}")
        return False

def test_anomaly_detection():
    """Test anomaly detection with sample data"""
    print("\nTesting anomaly detection...")
    try:
        # Test with normal request
        normal_request = {
            "response_time": 0.1,
            "request_size": 1024,
            "error_count": 0
        }
        
        response = requests.post(f"{ML_SERVICE_URL}/analyze/request", json=normal_request)
        if response.status_code == 200:
            result = response.json()
            print(f"✓ Normal request analysis: anomaly={result['is_anomaly']}, score={result['anomaly_score']:.4f}")
            
            # Test with suspicious request
            suspicious_request = {
                "response_time": 10.0,
                "request_size": 1000000,
                "error_count": 5
            }
            
            response = requests.post(f"{ML_SERVICE_URL}/analyze/request", json=suspicious_request)
            if response.status_code == 200:
                result = response.json()
                print(f"✓ Suspicious request analysis: anomaly={result['is_anomaly']}, score={result['anomaly_score']:.4f}")
                return True
            else:
                print(f"✗ Suspicious request analysis failed: {response.status_code}")
                return False
        else:
            print(f"✗ Normal request analysis failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Anomaly detection error: {e}")
        return False

def test_user_clustering():
    """Test user behavior clustering"""
    print("\nTesting user clustering...")
    try:
        # Test with sample user behavior
        user_behavior = {
            "login_count": 10,
            "purchase_count": 5,
            "cart_count": 8,
            "favorite_count": 12
        }
        
        response = requests.post(f"{ML_SERVICE_URL}/analyze/user", json=user_behavior)
        if response.status_code == 200:
            result = response.json()
            print(f"✓ User clustering: cluster={result['cluster']}")
            return True
        else:
            print(f"✗ User clustering failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ User clustering error: {e}")
        return False

def run_all_tests():
    """Run all tests"""
    print("=" * 50)
    print("ML Service Test Suite")
    print("=" * 50)
    
    tests = [
        test_health_check,
        test_manual_logistic_regression,
        test_manual_decision_tree,
        test_model_comparison,
        test_scheduler_status,
        test_model_versions,
        test_anomaly_detection,
        test_user_clustering
    ]
    
    passed = 0
    total = len(tests)
    
    for test in tests:
        if test():
            passed += 1
        time.sleep(1)  # Small delay between tests
    
    print("\n" + "=" * 50)
    print(f"Test Results: {passed}/{total} tests passed")
    print("=" * 50)
    
    if passed == total:
        print("🎉 All tests passed!")
        return True
    else:
        print("❌ Some tests failed!")
        return False

if __name__ == "__main__":
    run_all_tests() 