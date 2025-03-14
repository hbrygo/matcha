#import requests
#import json
#
#def test_get_user():
#    # API endpoint configuration
#    url = "http://localhost:8181/get_user"
#    headers = {"Content-Type": "application/json"}
#
#    def print_test(name, response):
#        print(f"\n=== Test: {name} ===")
#        print(f"Status: {response.status_code}")
#        print(f"Response: {response.text if response.text else 'No content'}")
#
#    # Test 1: Existing user request
#    print("\nTest 1: Existing user")
#    response = requests.post(url, 
#                           headers=headers,
#                           json={"uid": 1})
#    print_test("Existing user", response)
#    assert response.status_code in [200, 404]
#
#    # Test 2: Invalid UID
#    print("\nTest 2: Invalid UID")
#    response = requests.post(url, 
#                           headers=headers,
#                           json={"uid": -1})
#    print_test("Invalid UID", response)
#    assert response.status_code == 400
#
#    # Test 3: Missing UID
#    print("\nTest 3: Missing UID")
#    response = requests.post(url, 
#                           headers=headers,
#                           json={})
#    print_test("Missing UID", response)
#    assert response.status_code == 400
#
#    # Test 4: Invalid JSON format
#    print("\nTest 4: Invalid JSON format")
#    response = requests.post(url, 
#                           headers=headers,
#                           data="invalid json")
#    print_test("Invalid JSON format", response)
#    assert response.status_code == 400
#
#    print("\nAll tests completed!")
#
#if __name__ == "__main__":
#    test_get_user()

import requests
import json

# Couleurs pour le terminal
class Colors:
    GREEN = '\033[92m'    # Succès
    RED = '\033[91m'      # Erreur
    BLUE = '\033[94m'     # Info
    BOLD = '\033[1m'      # Gras
    RESET = '\033[0m'     # Reset

def test_get_user():
    url = "http://localhost:8181/get_user"
    headers = {"Content-Type": "application/json"}
    tests_passed = 0
    total_tests = 4

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

    print(f"{Colors.BLUE}{Colors.BOLD}🚀 DÉBUT DES TESTS GET_USER{Colors.RESET}\n")

    # Test 1: Utilisateur existant (créé précédemment)
    print_test_header("Récupération utilisateur existant")
    response = requests.post(url, 
                           headers=headers,
                           json={"uid": 1})  # Premier utilisateur créé
    if print_test_result("Utilisateur existant", response, 200):
        tests_passed += 1
        # Vérification du contenu de la réponse
        try:
            data = response.json()
            if "user" in data:
                print(f"{Colors.GREEN}   ✓ Format de réponse correct{Colors.RESET}")
            else:
                print(f"{Colors.RED}   ✗ Format de réponse incorrect{Colors.RESET}")
        except json.JSONDecodeError:
            print(f"{Colors.RED}   ✗ Réponse non JSON{Colors.RESET}")

    # Test 2: UID invalide
    print_test_header("UID invalide")
    response = requests.post(url, 
                           headers=headers,
                           json={"uid": -1})
    if print_test_result("UID invalide", response, 400):
        tests_passed += 1

    # Test 3: UID manquant
    print_test_header("UID manquant")
    response = requests.post(url, 
                           headers=headers,
                           json={})
    if print_test_result("UID manquant", response, 400):
        tests_passed += 1

    # Test 4: Format JSON invalide
    print_test_header("Format JSON invalide")
    response = requests.post(url, 
                           headers=headers,
                           data="invalid json")
    if print_test_result("Format JSON invalide", response, 400):
        tests_passed += 1

    # Résumé final
    print(f"\n{Colors.BLUE}{Colors.BOLD}📊 RÉSUMÉ DES TESTS{Colors.RESET}")
    if tests_passed == total_tests:
        print(f"{Colors.GREEN}✅ SUCCÈS TOTAL: {tests_passed}/{total_tests} tests réussis{Colors.RESET}")
    else:
        print(f"{Colors.RED}❌ ÉCHEC: {tests_passed}/{total_tests} tests réussis")
        print(f"   {total_tests - tests_passed} test(s) ont échoué{Colors.RESET}")

if __name__ == "__main__":
    test_get_user()