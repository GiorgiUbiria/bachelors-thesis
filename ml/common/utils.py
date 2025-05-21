# Shared helper functions (e.g., evaluation, logging)

def accuracy(y_true, y_pred):
    return (y_true == y_pred).mean()

def log(msg):
    print(f'[ML LOG] {msg}')