import requests
import time

API_URL = 'http://localhost:8080/log-request'


def send_request(ip=None, headers=None):
    session = requests.Session()
    if ip:
        session.headers.update({'X-Forwarded-For': ip})
    if headers:
        session.headers.update(headers)
    try:
        resp = session.post(API_URL)
        print(f"Request from {ip or 'default'}: {resp.status_code} {resp.text}")
        return resp
    except Exception as e:
        print(f"Error for {ip}: {e}")
        return None

def simulate_normal():
    print("--- Normal Requests ---")
    for i in range(3):
        send_request(ip=f"10.0.0.{i+1}")
        time.sleep(1)

def simulate_rapid():
    print("--- Rapid Requests (Potential Anomaly) ---")
    ip = "20.0.0.1"
    for i in range(5):
        send_request(ip=ip)
        time.sleep(0.1)

def simulate_anomalous():
    print("--- Anomalous Requests (Malformed) ---")
    ip = "30.0.0.1"
    # Simulate strange user agent or path
    for i in range(3):
        send_request(ip=ip, headers={"User-Agent": "sqlmap/1.0"})
        time.sleep(0.2)

def test_ban():
    print("--- Test Banned IP ---")
    ip = "40.0.0.1"
    # Simulate anomaly to trigger ban
    send_request(ip=ip, headers={"User-Agent": "attack-bot"})
    # Immediately try again
    resp = send_request(ip=ip)
    if resp is not None and resp.status_code == 403:
        print(f"IP {ip} is banned as expected.")
    else:
        print(f"IP {ip} is NOT banned (unexpected).")

def main():
    simulate_normal()
    simulate_rapid()
    simulate_anomalous()
    test_ban()

if __name__ == "__main__":
    main() 