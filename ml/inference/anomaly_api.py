from fastapi import FastAPI
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