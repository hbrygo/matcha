import requests

def test_db_endpoint():
    try:
        response = requests.get('http://localhost:8181/check-db')
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.json()}")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    test_db_endpoint()