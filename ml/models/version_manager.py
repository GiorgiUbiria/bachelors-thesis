import os
import shutil
import joblib
import json
import logging
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple
from pathlib import Path
import numpy as np

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class ModelVersionManager:
    def __init__(self, base_model_dir: str = "models", max_versions: int = 10, retention_days: int = 30):
        """
        Initialize the model version manager
        
        Args:
            base_model_dir: Base directory for storing models
            max_versions: Maximum number of versions to keep per model type
            retention_days: Number of days to retain old models
        """
        self.base_model_dir = Path(base_model_dir)
        self.max_versions = max_versions
        self.retention_days = retention_days
        
        # Create directories
        self.base_model_dir.mkdir(exist_ok=True)
        self.versions_dir = self.base_model_dir / "versions"
        self.versions_dir.mkdir(exist_ok=True)
        self.metadata_dir = self.base_model_dir / "metadata"
        self.metadata_dir.mkdir(exist_ok=True)
        
        logger.info(f"Model Version Manager initialized with base dir: {self.base_model_dir}")

    def _generate_version_id(self) -> str:
        """Generate a timestamp-based version ID"""
        return datetime.now().strftime("%Y%m%d_%H%M%S")

    def _get_model_path(self, model_type: str, version_id: str) -> Path:
        """Get the path for a specific model version"""
        return self.versions_dir / f"{model_type}_{version_id}.joblib"

    def _get_metadata_path(self, model_type: str, version_id: str) -> Path:
        """Get the path for model metadata"""
        return self.metadata_dir / f"{model_type}_{version_id}_metadata.json"

    def _get_current_version_path(self, model_type: str) -> Path:
        """Get the path for the current version symlink"""
        return self.base_model_dir / f"{model_type}_current.joblib"

    def save_model(self, model: Any, model_type: str, metadata: Optional[Dict] = None, 
                   performance_metrics: Optional[Dict] = None) -> str:
        """
        Save a model with versioning
        
        Args:
            model: The model object to save
            model_type: Type of model (e.g., 'anomaly', 'clustering', 'recommendation')
            metadata: Additional metadata about the model
            performance_metrics: Performance metrics for the model
            
        Returns:
            version_id: The version ID of the saved model
        """
        try:
            version_id = self._generate_version_id()
            model_path = self._get_model_path(model_type, version_id)
            metadata_path = self._get_metadata_path(model_type, version_id)
            current_path = self._get_current_version_path(model_type)
            
            # Save the model
            joblib.dump(model, model_path)
            logger.info(f"Saved {model_type} model version {version_id} to {model_path}")
            
            # Prepare metadata
            model_metadata = {
                "model_type": model_type,
                "version_id": version_id,
                "created_at": datetime.now().isoformat(),
                "model_path": str(model_path),
                "performance_metrics": performance_metrics or {},
                "metadata": metadata or {}
            }
            
            # Save metadata
            with open(metadata_path, 'w') as f:
                json.dump(model_metadata, f, indent=2)
            
            # Update current version symlink
            if current_path.exists() or current_path.is_symlink():
                current_path.unlink()
            
            # Create relative symlink
            relative_path = os.path.relpath(model_path, self.base_model_dir)
            current_path.symlink_to(relative_path)
            
            logger.info(f"Updated current {model_type} model to version {version_id}")
            
            # Cleanup old versions
            self._cleanup_old_versions(model_type)
            
            return version_id
            
        except Exception as e:
            logger.error(f"Error saving {model_type} model: {e}")
            raise

    def load_model(self, model_type: str, version_id: Optional[str] = None) -> Tuple[Any, Dict]:
        """
        Load a model by type and version
        
        Args:
            model_type: Type of model to load
            version_id: Specific version to load (None for current)
            
        Returns:
            Tuple of (model, metadata)
        """
        try:
            if version_id is None:
                # Load current version
                current_path = self._get_current_version_path(model_type)
                if not current_path.exists():
                    raise FileNotFoundError(f"No current version found for {model_type} model")
                
                model_path = current_path.resolve()
                # Extract version_id from filename
                version_id = model_path.stem.split('_', 1)[1]
            else:
                model_path = self._get_model_path(model_type, version_id)
                if not model_path.exists():
                    raise FileNotFoundError(f"Model version {version_id} not found for {model_type}")
            
            # Load model
            model = joblib.load(model_path)
            
            # Load metadata
            metadata_path = self._get_metadata_path(model_type, version_id)
            metadata = {}
            if metadata_path.exists():
                with open(metadata_path, 'r') as f:
                    metadata = json.load(f)
            
            logger.info(f"Loaded {model_type} model version {version_id}")
            return model, metadata
            
        except Exception as e:
            logger.error(f"Error loading {model_type} model version {version_id}: {e}")
            raise

    def list_versions(self, model_type: str) -> List[Dict]:
        """
        List all versions of a specific model type
        
        Args:
            model_type: Type of model to list versions for
            
        Returns:
            List of version metadata dictionaries
        """
        try:
            versions = []
            
            # Find all model files for this type
            pattern = f"{model_type}_*.joblib"
            model_files = list(self.versions_dir.glob(pattern))
            
            for model_file in model_files:
                # Extract version_id from filename
                version_id = model_file.stem.split('_', 1)[1]
                
                # Load metadata if available
                metadata_path = self._get_metadata_path(model_type, version_id)
                if metadata_path.exists():
                    with open(metadata_path, 'r') as f:
                        metadata = json.load(f)
                    versions.append(metadata)
                else:
                    # Create basic metadata if file exists but metadata doesn't
                    basic_metadata = {
                        "model_type": model_type,
                        "version_id": version_id,
                        "created_at": datetime.fromtimestamp(model_file.stat().st_mtime).isoformat(),
                        "model_path": str(model_file),
                        "performance_metrics": {},
                        "metadata": {}
                    }
                    versions.append(basic_metadata)
            
            # Sort by creation time (newest first)
            versions.sort(key=lambda x: x["created_at"], reverse=True)
            
            return versions
            
        except Exception as e:
            logger.error(f"Error listing versions for {model_type}: {e}")
            return []

    def rollback_model(self, model_type: str, version_id: str) -> bool:
        """
        Rollback to a specific model version
        
        Args:
            model_type: Type of model to rollback
            version_id: Version to rollback to
            
        Returns:
            True if successful, False otherwise
        """
        try:
            model_path = self._get_model_path(model_type, version_id)
            if not model_path.exists():
                logger.error(f"Cannot rollback: version {version_id} not found for {model_type}")
                return False
            
            current_path = self._get_current_version_path(model_type)
            
            # Remove current symlink
            if current_path.exists() or current_path.is_symlink():
                current_path.unlink()
            
            # Create new symlink to the specified version
            relative_path = os.path.relpath(model_path, self.base_model_dir)
            current_path.symlink_to(relative_path)
            
            logger.info(f"Rolled back {model_type} model to version {version_id}")
            return True
            
        except Exception as e:
            logger.error(f"Error rolling back {model_type} model to version {version_id}: {e}")
            return False

    def delete_version(self, model_type: str, version_id: str) -> bool:
        """
        Delete a specific model version
        
        Args:
            model_type: Type of model
            version_id: Version to delete
            
        Returns:
            True if successful, False otherwise
        """
        try:
            model_path = self._get_model_path(model_type, version_id)
            metadata_path = self._get_metadata_path(model_type, version_id)
            
            # Check if this is the current version
            current_path = self._get_current_version_path(model_type)
            if current_path.exists() and current_path.resolve() == model_path.resolve():
                logger.error(f"Cannot delete current version {version_id} of {model_type} model")
                return False
            
            # Delete model file
            if model_path.exists():
                model_path.unlink()
                logger.info(f"Deleted model file: {model_path}")
            
            # Delete metadata file
            if metadata_path.exists():
                metadata_path.unlink()
                logger.info(f"Deleted metadata file: {metadata_path}")
            
            return True
            
        except Exception as e:
            logger.error(f"Error deleting {model_type} model version {version_id}: {e}")
            return False

    def _cleanup_old_versions(self, model_type: str):
        """Clean up old versions based on retention policy"""
        try:
            versions = self.list_versions(model_type)
            
            # Remove versions beyond max_versions limit
            if len(versions) > self.max_versions:
                versions_to_delete = versions[self.max_versions:]
                for version_metadata in versions_to_delete:
                    version_id = version_metadata["version_id"]
                    logger.info(f"Cleaning up old version {version_id} of {model_type} (max versions exceeded)")
                    self.delete_version(model_type, version_id)
            
            # Remove versions older than retention_days
            cutoff_date = datetime.now() - timedelta(days=self.retention_days)
            for version_metadata in versions:
                created_at = datetime.fromisoformat(version_metadata["created_at"])
                if created_at < cutoff_date:
                    version_id = version_metadata["version_id"]
                    logger.info(f"Cleaning up old version {version_id} of {model_type} (retention period exceeded)")
                    self.delete_version(model_type, version_id)
            
        except Exception as e:
            logger.error(f"Error during cleanup for {model_type}: {e}")

    def update_performance_metrics(self, model_type: str, version_id: str, 
                                 performance_metrics: Dict) -> bool:
        """
        Update performance metrics for a specific model version
        
        Args:
            model_type: Type of model
            version_id: Version to update
            performance_metrics: New performance metrics
            
        Returns:
            True if successful, False otherwise
        """
        try:
            metadata_path = self._get_metadata_path(model_type, version_id)
            
            if not metadata_path.exists():
                logger.error(f"Metadata not found for {model_type} version {version_id}")
                return False
            
            # Load existing metadata
            with open(metadata_path, 'r') as f:
                metadata = json.load(f)
            
            # Update performance metrics
            metadata["performance_metrics"].update(performance_metrics)
            metadata["last_updated"] = datetime.now().isoformat()
            
            # Save updated metadata
            with open(metadata_path, 'w') as f:
                json.dump(metadata, f, indent=2)
            
            logger.info(f"Updated performance metrics for {model_type} version {version_id}")
            return True
            
        except Exception as e:
            logger.error(f"Error updating performance metrics: {e}")
            return False

    def get_best_performing_version(self, model_type: str, metric: str = "accuracy") -> Optional[str]:
        """
        Get the version with the best performance for a specific metric
        
        Args:
            model_type: Type of model
            metric: Performance metric to compare (e.g., 'accuracy', 'f1_score')
            
        Returns:
            Version ID of the best performing model, or None if not found
        """
        try:
            versions = self.list_versions(model_type)
            
            best_version = None
            best_score = -float('inf')
            
            for version_metadata in versions:
                performance_metrics = version_metadata.get("performance_metrics", {})
                if metric in performance_metrics:
                    score = performance_metrics[metric]
                    if score > best_score:
                        best_score = score
                        best_version = version_metadata["version_id"]
            
            if best_version:
                logger.info(f"Best performing {model_type} version for {metric}: {best_version} (score: {best_score})")
            else:
                logger.warning(f"No performance data found for {metric} in {model_type} models")
            
            return best_version
            
        except Exception as e:
            logger.error(f"Error finding best performing version: {e}")
            return None

    def auto_rollback_on_degradation(self, model_type: str, current_metrics: Dict, 
                                   threshold_metric: str = "accuracy", 
                                   degradation_threshold: float = 0.05) -> bool:
        """
        Automatically rollback if current model performance has degraded
        
        Args:
            model_type: Type of model
            current_metrics: Current model performance metrics
            threshold_metric: Metric to check for degradation
            degradation_threshold: Threshold for triggering rollback
            
        Returns:
            True if rollback was performed, False otherwise
        """
        try:
            if threshold_metric not in current_metrics:
                logger.warning(f"Threshold metric {threshold_metric} not found in current metrics")
                return False
            
            current_score = current_metrics[threshold_metric]
            versions = self.list_versions(model_type)
            
            # Find the previous version with performance data
            previous_version = None
            previous_score = None
            
            for version_metadata in versions[1:]:  # Skip current version (index 0)
                performance_metrics = version_metadata.get("performance_metrics", {})
                if threshold_metric in performance_metrics:
                    previous_version = version_metadata["version_id"]
                    previous_score = performance_metrics[threshold_metric]
                    break
            
            if previous_version is None or previous_score is None:
                logger.warning(f"No previous version with {threshold_metric} data found for {model_type}")
                return False
            
            # Check if performance has degraded
            performance_drop = previous_score - current_score
            if performance_drop > degradation_threshold:
                logger.warning(f"Performance degradation detected for {model_type}: "
                             f"{threshold_metric} dropped from {previous_score:.4f} to {current_score:.4f}")
                
                # Perform rollback
                if self.rollback_model(model_type, previous_version):
                    logger.info(f"Auto-rollback successful: {model_type} rolled back to version {previous_version}")
                    return True
                else:
                    logger.error(f"Auto-rollback failed for {model_type}")
                    return False
            else:
                logger.info(f"No significant performance degradation detected for {model_type}")
                return False
            
        except Exception as e:
            logger.error(f"Error in auto-rollback check: {e}")
            return False

    def get_version_summary(self) -> Dict:
        """Get a summary of all model versions"""
        try:
            summary = {}
            
            # Find all unique model types
            model_types = set()
            for model_file in self.versions_dir.glob("*.joblib"):
                model_type = model_file.stem.split('_')[0]
                model_types.add(model_type)
            
            for model_type in model_types:
                versions = self.list_versions(model_type)
                current_path = self._get_current_version_path(model_type)
                current_version = None
                
                if current_path.exists():
                    current_model_path = current_path.resolve()
                    current_version = current_model_path.stem.split('_', 1)[1]
                
                summary[model_type] = {
                    "total_versions": len(versions),
                    "current_version": current_version,
                    "latest_version": versions[0]["version_id"] if versions else None,
                    "oldest_version": versions[-1]["version_id"] if versions else None,
                    "versions": [v["version_id"] for v in versions]
                }
            
            return summary
            
        except Exception as e:
            logger.error(f"Error generating version summary: {e}")
            return {}


# Global version manager instance
version_manager_instance = None

def get_version_manager() -> ModelVersionManager:
    """Get or create the global version manager instance"""
    global version_manager_instance
    if version_manager_instance is None:
        model_dir = os.getenv('MODEL_DIR', 'models')
        max_versions = int(os.getenv('MAX_MODEL_VERSIONS', '10'))
        retention_days = int(os.getenv('MODEL_RETENTION_DAYS', '30'))
        version_manager_instance = ModelVersionManager(model_dir, max_versions, retention_days)
    return version_manager_instance 