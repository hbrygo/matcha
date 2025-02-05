import requests

def test_create_user():
    url = 'http://localhost:8181/create_user'
    test_data = {
        "nom": "Dupont",
        "prenom": "Jean",
        "email": "jean.dupont@example.com",
        "password": "motdepasse123"
    }

    try:
        response = requests.post(
            url,
            json=test_data,
            headers={'Content-Type': 'application/json'}
        )
        print(f"Status: {response.status_code}")
        if response.text:  # Vérifie si il y a une réponse
            print(f"Response: {response.json()}")
        else:
            print("No response body")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    test_create_user()