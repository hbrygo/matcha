#import requests
#import json
#import random
#import string
#
#def generate_random_string(length=10):
#    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))
#
#def test_create_user():
#    url = "http://localhost:8181/create_user"
#    headers = {"Content-Type": "application/json"}
#
#    def make_request(data):
#        response = requests.post(url, json=data, headers=headers)
#        print(f"\nTest data: {json.dumps(data, indent=2)}")
#        print(f"Status code: {response.status_code}")
#        print(f"Response: {response.text if response.text else 'No content'}")
#        return response
#
#    # Test 1: Valid user creation
#    valid_user = {
#        "username": f"test_user_{generate_random_string()}",
#        "email": f"test_{generate_random_string()}@example.com",
#        "password": "Test123!@#"
#    }
#    print("\nTest 1: Valid user creation")
#    response = make_request(valid_user)
#    assert response.status_code == 200
#
#    # Test 2: Duplicate email
#    print("\nTest 2: Duplicate email")
#    duplicate_email = valid_user.copy()
#    duplicate_email["username"] = f"test_user_{generate_random_string()}"
#    response = make_request(duplicate_email)
#    assert response.status_code == 409
#
#    # Test 3: Duplicate username
#    print("\nTest 3: Duplicate username")
#    duplicate_username = valid_user.copy()
#    duplicate_username["email"] = f"test_{generate_random_string()}@example.com"
#    response = make_request(duplicate_username)
#    assert response.status_code == 410
#
#    # Test 4: Invalid password
#    print("\nTest 4: Invalid password")
#    weak_password = valid_user.copy()
#    weak_password["password"] = "weak"
#    weak_password["email"] = f"test_{generate_random_string()}@example.com"
#    weak_password["username"] = f"test_user_{generate_random_string()}"
#    response = make_request(weak_password)
#    assert response.status_code == 400
#
#    # Test 5: Invalid email format
#    print("\nTest 5: Invalid email format")
#    invalid_email = valid_user.copy()
#    invalid_email["email"] = "invalid_email"
#    invalid_email["username"] = f"test_user_{generate_random_string()}"
#    response = make_request(invalid_email)
#    assert response.status_code == 400
#
#    print("\nAll tests completed!")
#
#if __name__ == "__main__":
#    test_create_user()

import requests
import json
import random
import string

# Couleurs pour une meilleure lisibilité dans le terminal
class Colors:
    GREEN = '\033[92m'    # Succès
    RED = '\033[91m'      # Erreur
    BLUE = '\033[94m'     # Info
    YELLOW = '\033[93m'   # Avertissement
    BOLD = '\033[1m'      # Texte en gras
    RESET = '\033[0m'     # Reset couleur

def generate_random_string(length=10):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

def print_test_header(test_name: str):
    print(f"\n{Colors.BLUE}{Colors.BOLD}⚡️ TEST: {test_name}{Colors.RESET}")

def print_test_result(test_name: str, success: bool, status_code: int, expected_status: int, response_text: str):
    if success:
        print(f"{Colors.GREEN}✅ SUCCÈS: {test_name}")
        print(f"   Status: {status_code} (attendu: {expected_status})")
        print(f"   Réponse: {response_text}{Colors.RESET}")
    else:
        print(f"{Colors.RED}❌ ÉCHEC: {test_name}")
        print(f"   Status reçu: {status_code}")
        print(f"   Status attendu: {expected_status}")
        print(f"   Réponse: {response_text}{Colors.RESET}")

def test_create_user():
    url = "http://localhost:8181/create_user"
    headers = {"Content-Type": "application/json"}
    tests_passed = 0
    total_tests = 5

    print(f"{Colors.BLUE}{Colors.BOLD}🚀 DÉBUT DES TESTS DE CRÉATION D'UTILISATEUR{Colors.RESET}\n")

    # Test 1: Création d'utilisateur valide
    print_test_header("Création d'utilisateur valide")
    valid_user = {
        "username": f"test_user_{generate_random_string()}",
        "email": f"test_{generate_random_string()}@example.com",
        "password": "Test123!@#"
    }
    response = requests.post(url, json=valid_user, headers=headers)
    success = response.status_code == 200
    print_test_result("Création utilisateur", success, response.status_code, 200, response.text)
    if success:
        tests_passed += 1
        test_data = valid_user  # Sauvegarder pour les tests suivants

    # Test 2: Email en double
    print_test_header("Test email en double")
    duplicate_email = valid_user.copy()
    duplicate_email["username"] = f"test_user_{generate_random_string()}"
    response = requests.post(url, json=duplicate_email, headers=headers)
    success = response.status_code == 409
    print_test_result("Email en double", success, response.status_code, 409, response.text)
    if success:
        tests_passed += 1

    # Test 3: Username en double
    print_test_header("Test username en double")
    duplicate_username = valid_user.copy()
    duplicate_username["email"] = f"test_{generate_random_string()}@example.com"
    response = requests.post(url, json=duplicate_username, headers=headers)
    success = response.status_code == 410
    print_test_result("Username en double", success, response.status_code, 410, response.text)
    if success:
        tests_passed += 1

    # Test 4: Mot de passe faible
    print_test_header("Test mot de passe faible")
    weak_password = {
        "username": f"test_user_{generate_random_string()}",
        "email": f"test_{generate_random_string()}@example.com",
        "password": "weak"
    }
    response = requests.post(url, json=weak_password, headers=headers)
    success = response.status_code == 400
    print_test_result("Mot de passe faible", success, response.status_code, 400, response.text)
    if success:
        tests_passed += 1

    # Test 5: Email invalide
    print_test_header("Test email invalide")
    invalid_email = {
        "username": f"test_user_{generate_random_string()}",
        "email": "invalid_email",
        "password": "Test123!@#"
    }
    response = requests.post(url, json=invalid_email, headers=headers)
    success = response.status_code == 400
    print_test_result("Email invalide", success, response.status_code, 400, response.text)
    if success:
        tests_passed += 1

    # Résumé final
    print(f"\n{Colors.BLUE}{Colors.BOLD}📊 RÉSUMÉ DES TESTS{Colors.RESET}")
    if tests_passed == total_tests:
        print(f"{Colors.GREEN}✅ SUCCÈS TOTAL: {tests_passed}/{total_tests} tests réussis{Colors.RESET}")
    else:
        print(f"{Colors.RED}❌ ÉCHEC: {tests_passed}/{total_tests} tests réussis")
        print(f"   {total_tests - tests_passed} test(s) ont échoué{Colors.RESET}")

if __name__ == "__main__":
    test_create_user()