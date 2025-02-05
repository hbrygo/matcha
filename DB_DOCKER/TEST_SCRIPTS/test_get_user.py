import requests

def test_get_user():
    url = 'http://localhost:8181/get_user'
    
    # Test 1: Recherche par email et password (utilisateur test)
    print("\nTest 1: Recherche par email/password")
    test_data1 = {
        "email": "test@example.com",
        "password": "testpassword"
    }
    
    try:
        response = requests.post(url, json=test_data1)
        print(f"Status: {response.status_code}")
        print(f"Response: {response.json()}")
    except Exception as e:
        print(f"Error: {e}")
    
    # Test 2: Recherche par UID (assumant que l'utilisateur test a l'UID 1)
    print("\nTest 2: Recherche par UID")
    test_data2 = {
        "uid": 1
    }
    
    try:
        response = requests.post(url, json=test_data2)
        print(f"Status: {response.status_code}")
        print(f"Response: {response.json()}")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    test_get_user()