import os
import logging
import requests
from datetime import datetime, timedelta
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.cron import CronTrigger
from apscheduler.triggers.interval import IntervalTrigger
import atexit
from typing import Optional

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class MLRetrainingScheduler:
    def __init__(self, ml_service_url: str = "http://localhost:5000", api_key: Optional[str] = None):
        """
        Initialize the ML retraining scheduler
        
        Args:
            ml_service_url: URL of the ML service
            api_key: Optional API key for authentication
        """
        self.ml_service_url = ml_service_url
        self.api_key = api_key
        self.scheduler = BackgroundScheduler()
        self.last_retrain_times = {}
        
        # Configure scheduler
        self.scheduler.start()
        atexit.register(lambda: self.scheduler.shutdown())
        
        logger.info("ML Retraining Scheduler initialized")

    def _make_retrain_request(self, model_type: str = "all") -> bool:
        """Make a retraining request to the ML service"""
        try:
            url = f"{self.ml_service_url}/train"
            headers = {}
            
            if self.api_key:
                headers['Authorization'] = f'Bearer {self.api_key}'
            
            data = {"model_type": model_type}
            
            logger.info(f"Triggering retraining for {model_type} models")
            response = requests.post(url, json=data, headers=headers, timeout=300)  # 5 min timeout
            
            if response.status_code == 200:
                logger.info(f"Successfully triggered retraining for {model_type} models")
                self.last_retrain_times[model_type] = datetime.now()
                return True
            else:
                logger.error(f"Failed to retrain {model_type} models: {response.status_code} - {response.text}")
                return False
                
        except requests.exceptions.Timeout:
            logger.error(f"Timeout while retraining {model_type} models")
            return False
        except requests.exceptions.ConnectionError:
            logger.error(f"Connection error while retraining {model_type} models")
            return False
        except Exception as e:
            logger.error(f"Unexpected error while retraining {model_type} models: {e}")
            return False

    def _check_data_volume_threshold(self, model_type: str, threshold: int = 1000) -> bool:
        """Check if data volume threshold is reached for retraining"""
        try:
            # This would typically check the database for new data since last retrain
            # For now, we'll implement a simple check
            url = f"{self.ml_service_url}/health"
            response = requests.get(url, timeout=30)
            
            if response.status_code == 200:
                # In a real implementation, this would check actual data counts
                # For now, we'll simulate based on time since last retrain
                last_retrain = self.last_retrain_times.get(model_type)
                if last_retrain is None:
                    return True  # Never retrained before
                
                hours_since_retrain = (datetime.now() - last_retrain).total_seconds() / 3600
                # Simulate data accumulation: assume 100 samples per hour
                estimated_new_samples = int(hours_since_retrain * 100)
                
                return estimated_new_samples >= threshold
            
            return False
            
        except Exception as e:
            logger.error(f"Error checking data volume for {model_type}: {e}")
            return False

    def _check_model_performance(self, model_type: str, min_accuracy: float = 0.85) -> bool:
        """Check if model performance has degraded below threshold"""
        try:
            # This would typically evaluate model performance on validation data
            # For now, we'll implement a simple simulation
            url = f"{self.ml_service_url}/health"
            response = requests.get(url, timeout=30)
            
            if response.status_code == 200:
                # In a real implementation, this would check actual model metrics
                # For now, we'll simulate performance degradation over time
                last_retrain = self.last_retrain_times.get(model_type)
                if last_retrain is None:
                    return True  # Never retrained before
                
                days_since_retrain = (datetime.now() - last_retrain).days
                # Simulate performance degradation: 1% per day
                simulated_accuracy = max(0.5, 0.95 - (days_since_retrain * 0.01))
                
                return simulated_accuracy < min_accuracy
            
            return False
            
        except Exception as e:
            logger.error(f"Error checking performance for {model_type}: {e}")
            return False

    def retrain_anomaly_models(self):
        """Retrain anomaly detection models"""
        logger.info("Scheduled anomaly model retraining triggered")
        
        # Check if retraining is needed based on data volume or performance
        if (self._check_data_volume_threshold("anomaly", threshold=500) or 
            self._check_model_performance("anomaly", min_accuracy=0.85)):
            
            success = self._make_retrain_request("anomaly")
            if success:
                logger.info("Anomaly models retrained successfully")
            else:
                logger.error("Failed to retrain anomaly models")
        else:
            logger.info("Anomaly model retraining skipped - thresholds not met")

    def retrain_clustering_models(self):
        """Retrain user clustering models"""
        logger.info("Scheduled clustering model retraining triggered")
        
        # Check if retraining is needed
        if (self._check_data_volume_threshold("clustering", threshold=200) or 
            self._check_model_performance("clustering", min_accuracy=0.80)):
            
            success = self._make_retrain_request("clustering")
            if success:
                logger.info("Clustering models retrained successfully")
            else:
                logger.error("Failed to retrain clustering models")
        else:
            logger.info("Clustering model retraining skipped - thresholds not met")

    def retrain_recommendation_models(self):
        """Retrain recommendation models"""
        logger.info("Scheduled recommendation model retraining triggered")
        
        # Check if retraining is needed
        if (self._check_data_volume_threshold("recommendation", threshold=100) or 
            self._check_model_performance("recommendation", min_accuracy=0.75)):
            
            success = self._make_retrain_request("recommendation")
            if success:
                logger.info("Recommendation models retrained successfully")
            else:
                logger.error("Failed to retrain recommendation models")
        else:
            logger.info("Recommendation model retraining skipped - thresholds not met")

    def retrain_trend_models(self):
        """Retrain trend prediction models"""
        logger.info("Scheduled trend model retraining triggered")
        
        # Check if retraining is needed
        if (self._check_data_volume_threshold("trend", threshold=50) or 
            self._check_model_performance("trend", min_accuracy=0.70)):
            
            success = self._make_retrain_request("trend")
            if success:
                logger.info("Trend models retrained successfully")
            else:
                logger.error("Failed to retrain trend models")
        else:
            logger.info("Trend model retraining skipped - thresholds not met")

    def retrain_all_models(self):
        """Retrain all models"""
        logger.info("Scheduled full model retraining triggered")
        success = self._make_retrain_request("all")
        if success:
            logger.info("All models retrained successfully")
        else:
            logger.error("Failed to retrain all models")

    def setup_schedules(self):
        """Setup all retraining schedules"""
        
        # Hourly anomaly model retraining (every hour)
        self.scheduler.add_job(
            func=self.retrain_anomaly_models,
            trigger=IntervalTrigger(hours=1),
            id='anomaly_retrain',
            name='Anomaly Model Retraining',
            replace_existing=True
        )
        logger.info("Scheduled anomaly model retraining every hour")
        
        # Daily user clustering retraining (at 2 AM)
        self.scheduler.add_job(
            func=self.retrain_clustering_models,
            trigger=CronTrigger(hour=2, minute=0),
            id='clustering_retrain',
            name='Clustering Model Retraining',
            replace_existing=True
        )
        logger.info("Scheduled clustering model retraining daily at 2 AM")
        
        # Weekly recommendation model retraining (Sundays at 3 AM)
        self.scheduler.add_job(
            func=self.retrain_recommendation_models,
            trigger=CronTrigger(day_of_week=6, hour=3, minute=0),  # Sunday = 6
            id='recommendation_retrain',
            name='Recommendation Model Retraining',
            replace_existing=True
        )
        logger.info("Scheduled recommendation model retraining weekly on Sundays at 3 AM")
        
        # Monthly trend prediction retraining (1st of month at 4 AM)
        self.scheduler.add_job(
            func=self.retrain_trend_models,
            trigger=CronTrigger(day=1, hour=4, minute=0),
            id='trend_retrain',
            name='Trend Model Retraining',
            replace_existing=True
        )
        logger.info("Scheduled trend model retraining monthly on 1st at 4 AM")
        
        # Emergency full retraining (can be triggered manually)
        self.scheduler.add_job(
            func=self.retrain_all_models,
            trigger=CronTrigger(day_of_week=0, hour=5, minute=0),  # Monday = 0
            id='full_retrain',
            name='Full Model Retraining',
            replace_existing=True
        )
        logger.info("Scheduled full model retraining weekly on Mondays at 5 AM")

    def trigger_manual_retrain(self, model_type: str = "all"):
        """Manually trigger retraining for specific model type"""
        logger.info(f"Manual retraining triggered for {model_type}")
        
        if model_type == "anomaly":
            self.retrain_anomaly_models()
        elif model_type == "clustering":
            self.retrain_clustering_models()
        elif model_type == "recommendation":
            self.retrain_recommendation_models()
        elif model_type == "trend":
            self.retrain_trend_models()
        elif model_type == "all":
            self.retrain_all_models()
        else:
            logger.error(f"Unknown model type for manual retraining: {model_type}")

    def get_schedule_status(self) -> dict:
        """Get current schedule status and next run times"""
        jobs = self.scheduler.get_jobs()
        status = {
            "scheduler_running": self.scheduler.running,
            "jobs": []
        }
        
        for job in jobs:
            job_info = {
                "id": job.id,
                "name": job.name,
                "next_run_time": job.next_run_time.isoformat() if job.next_run_time else None,
                "trigger": str(job.trigger)
            }
            status["jobs"].append(job_info)
        
        status["last_retrain_times"] = {
            model_type: time.isoformat() if time else None
            for model_type, time in self.last_retrain_times.items()
        }
        
        return status

    def pause_schedule(self, job_id: str):
        """Pause a specific scheduled job"""
        try:
            self.scheduler.pause_job(job_id)
            logger.info(f"Paused scheduled job: {job_id}")
        except Exception as e:
            logger.error(f"Failed to pause job {job_id}: {e}")

    def resume_schedule(self, job_id: str):
        """Resume a specific scheduled job"""
        try:
            self.scheduler.resume_job(job_id)
            logger.info(f"Resumed scheduled job: {job_id}")
        except Exception as e:
            logger.error(f"Failed to resume job {job_id}: {e}")

    def shutdown(self):
        """Shutdown the scheduler"""
        logger.info("Shutting down ML retraining scheduler")
        self.scheduler.shutdown()


# Global scheduler instance
scheduler_instance = None

def get_scheduler() -> MLRetrainingScheduler:
    """Get or create the global scheduler instance"""
    global scheduler_instance
    if scheduler_instance is None:
        ml_service_url = os.getenv('ML_SERVICE_URL', 'http://localhost:5000')
        api_key = os.getenv('ML_API_KEY', None)
        scheduler_instance = MLRetrainingScheduler(ml_service_url, api_key)
        scheduler_instance.setup_schedules()
    return scheduler_instance

def start_scheduler():
    """Start the retraining scheduler"""
    scheduler = get_scheduler()
    logger.info("ML retraining scheduler started")
    return scheduler

if __name__ == "__main__":
    # Start scheduler when run directly
    scheduler = start_scheduler()
    
    # Keep the script running
    try:
        import time
        while True:
            time.sleep(60)  # Sleep for 1 minute
    except KeyboardInterrupt:
        logger.info("Received interrupt signal, shutting down...")
        scheduler.shutdown() 