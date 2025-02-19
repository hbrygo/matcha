#import requests
#import json
#from typing import Dict, Any
#
#def test_login_endpoint():
#    url = "http://localhost:8181/login"
#    headers = {"Content-Type": "application/json"}
#
#    def make_request(payload: Dict[str, Any]) -> requests.Response:
#        return requests.post(url, json=payload, headers=headers)
#
#    def print_test_result(test_name: str, response: requests.Response):
#        print(f"\n=== Test: {test_name} ===")
#        print(f"Status: {response.status_code}")
#        print(f"Response: {response.text if response.text else 'No content'}")
#
#    # Test 1: Login avec Email
#    print("\nTest 1: Login avec Email")
#    email_login = {
#        "Email": "test@exemple.com",
#        "Mot de passe": "test123"
#    }
#    response = make_request(email_login)
#    print_test_result("Login avec Email", response)
#
#    # Test 2: Login avec Username
#    print("\nTest 2: Login avec Username")
#    username_login = {
#        "Username": "testuser",
#        "Mot de passe": "test123"
#    }
#    response = make_request(username_login)
#    print_test_result("Login avec Username", response)
#
#    # Test 3: Email invalide
#    print("\nTest 3: Email invalide")
#    invalid_email = {
#        "Email": "nonexistent@exemple.com",
#        "Mot de passe": "test123"
#    }
#    response = make_request(invalid_email)
#    print_test_result("Email invalide", response)
#    assert response.status_code == 404
#
#    # Test 4: Mot de passe incorrect
#    print("\nTest 4: Mot de passe incorrect")
#    wrong_password = {
#        "Email": "test@exemple.com",
#        "Mot de passe": "wrongpassword"
#    }
#    response = make_request(wrong_password)
#    print_test_result("Mot de passe incorrect", response)
#    assert response.status_code == 401
#
#    # Test 5: Champs manquants
#    print("\nTest 5: Champs manquants")
#    missing_fields = {
#        "Email": "test@exemple.com"
#    }
#    response = make_request(missing_fields)
#    print_test_result("Champs manquants", response)
#    assert response.status_code == 400
#
#    print("\nTous les tests sont terminés!")
#
#if __name__ == "__main__":
#    test_login_endpoint()

import requests
import json
from typing import Dict, Any

# Couleurs pour le terminal
class Colors:
    GREEN = '\033[92m'    # Succès
    RED = '\033[91m'      # Erreur
    BLUE = '\033[94m'     # Info
    BOLD = '\033[1m'      # Gras
    RESET = '\033[0m'     # Reset

def setup_test_user():
    """Crée un utilisateur de test pour les tests de login"""
    url = "http://localhost:8181/create_user"
    test_user = {
        "username": "testuser",
        "email": "test@exemple.com",
        "password": "Test123!@#"
    }
    response = requests.post(
        url, 
        json=test_user,
        headers={"Content-Type": "application/json"}
    )
    return response.status_code == 200

def test_login_endpoint():
    url = "http://localhost:8181/login"
    headers = {"Content-Type": "application/json"}
    tests_passed = 0
    total_tests = 5

    # Configuration initiale
    print(f"{Colors.BLUE}{Colors.BOLD}🔧 CONFIGURATION DES TESTS{Colors.RESET}")
    if not setup_test_user():
        print(f"{Colors.RED}❌ Échec de la création de l'utilisateur de test{Colors.RESET}")
        return
    print(f"{Colors.GREEN}✅ Utilisateur de test créé avec succès{Colors.RESET}\n")


    def make_request(payload: Dict[str, Any]) -> requests.Response:
        return requests.post(url, json=payload, headers=headers)

    def print_test_header(test_name: str):
        print(f"\n{Colors.BLUE}{Colors.BOLD}⚡️ TEST: {test_name}{Colors.RESET}")

    def print_test_result(test_name: str, response: requests.Response, expected_status: int):
        success = response.status_code == expected_status
        if success:
            print(f"{Colors.GREEN}✅ SUCCÈS: {test_name}")
            print(f"   Status: {response.status_code} (attendu: {expected_status})")
            print(f"   Réponse: {response.text}{Colors.RESET}")
            return True
        else:
            print(f"{Colors.RED}❌ ÉCHEC: {test_name}")
            print(f"   Status reçu: {response.status_code}")
            print(f"   Status attendu: {expected_status}")
            print(f"   Réponse: {response.text}{Colors.RESET}")
            return False

    print(f"{Colors.BLUE}{Colors.BOLD}🚀 DÉBUT DES TESTS DE LOGIN{Colors.RESET}\n")

    # Test 1: Login avec Email
    print_test_header("Login avec Email")
    email_login = {
        "email": "test@exemple.com",
        "Mot de password": "Test123!@#"  # Même mot de passe que celui créé
    }
    response = make_request(email_login)
    if print_test_result("Login avec Email", response, 200):
        tests_passed += 1

    # Test 2: Login avec Username
    print_test_header("Login avec Username")
    username_login = {
        "usename": "testuser",
        "password": "Test123!@#"  # Même mot de passe que celui créé
    }
    response = make_request(username_login)
    if print_test_result("Login avec Username", response, 200):
        tests_passed += 1

    # Test 3: Email invalide
    print_test_header("Email invalide")
    invalid_email = {
        "Email": "nonexistent@exemple.com",
        "password": "Test123!@#"
    }
    response = make_request(invalid_email)
    if print_test_result("Email invalide", response, 404):
        tests_passed += 1

    # Test 4: Mot de passe incorrect
    print_test_header("Mot de passe incorrect")
    wrong_password = {
        "Email": "test@exemple.com",
        "password": "wrongpassword"
    }
    response = make_request(wrong_password)
    if print_test_result("Mot de passe incorrect", response, 401):
        tests_passed += 1

    # Test 5: Champs manquants
    print_test_header("Champs manquants")
    missing_fields = {
        "Email": "test@exemple.com"
    }
    response = make_request(missing_fields)
    if print_test_result("Champs manquants", response, 400):
        tests_passed += 1

    # Résumé final
    print(f"\n{Colors.BLUE}{Colors.BOLD}📊 RÉSUMÉ DES TESTS{Colors.RESET}")
    if tests_passed == total_tests:
        print(f"{Colors.GREEN}✅ SUCCÈS TOTAL: {tests_passed}/{total_tests} tests réussis{Colors.RESET}")
    else:
        print(f"{Colors.RED}❌ ÉCHEC: {tests_passed}/{total_tests} tests réussis")
        print(f"   {total_tests - tests_passed} test(s) ont échoué{Colors.RESET}")

if __name__ == "__main__":
    test_login_endpoint()