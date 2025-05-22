from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import numpy as np
import joblib
import os

MODEL_PATH = os.getenv('ANOMALY_MODEL_PATH', 'anomaly_model.joblib')

class PredictRequest(BaseModel):
    features: list[list[float]]

class PredictResponse(BaseModel):
    prediction: list[int]

app = FastAPI()

model = joblib.load(MODEL_PATH)

@app.post('/predict', response_model=PredictResponse)
def predict(req: PredictRequest):
    X = np.array(req.features)
    pred = model.predict(X)
    return PredictResponse(prediction=pred.tolist())

@app.post('/retrain')
def retrain(req: PredictRequest):
    try:
        X = np.array(req.features)
        # Re-instantiate the model (optionally, allow method/params via request)
        from ml.models.anomaly import AnomalyDetector
        new_model = AnomalyDetector(method='isolation_forest', n_estimators=100)
        new_model.fit(X)
        new_model.save(MODEL_PATH)
        global model
        model = joblib.load(MODEL_PATH)
        return {"status": "retrained", "n_samples": len(X)}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) 