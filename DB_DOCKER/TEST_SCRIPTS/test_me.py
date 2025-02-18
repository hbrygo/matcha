#import requests
#import json
#from typing import Dict, Any
#
#def test_me_endpoint():
#    url = "http://localhost:8181/me"
#    headers = {"Content-Type": "application/json"}
#
#    def make_request(payload: Dict[str, Any]) -> requests.Response:
#        return requests.post(url, json=payload, headers=headers)
#
#    def print_response(test_name: str, response: requests.Response):
#        print(f"\n=== {test_name} ===")
#        print(f"Status: {response.status_code}")
#        print(f"Response: {response.text if response.text else 'No content'}")
#
#    # Test 1: Requête valide
#    print("\nTest 1: Requête valide")
#    response = make_request({"uid": 1})
#    print_response("Requête valide", response)
#    assert response.status_code in [200, 404], "Le status devrait être 200 ou 404"
#
#    # Test 2: UID invalide
#    print("\nTest 2: UID invalide")
#    response = make_request({"uid": -1})
#    print_response("UID invalide", response)
#    assert response.status_code == 400, "Le status devrait être 400"
#
#    # Test 3: UID manquant
#    print("\nTest 3: UID manquant")
#    response = make_request({})
#    print_response("UID manquant", response)
#    assert response.status_code == 400, "Le status devrait être 400"
#
#    # Test 4: Format JSON invalide
#    print("\nTest 4: Format JSON invalide")
#    response = requests.post(url, data="invalid json", headers=headers)
#    print_response("Format JSON invalide", response)
#    assert response.status_code == 400, "Le status devrait être 400"
#
#    print("\nTous les tests sont terminés!")
#
#if __name__ == "__main__":
#    test_me_endpoint()

import requests
import json
from typing import Dict, Any

# Couleurs pour une meilleure lisibilité
class Colors:
    GREEN = '\033[92m'    # Succès
    RED = '\033[91m'      # Erreur
    BLUE = '\033[94m'     # Info
    BOLD = '\033[1m'      # Texte en gras
    RESET = '\033[0m'     # Reset couleur

def test_me_endpoint():
    url = "http://localhost:8181/me"
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

    print(f"{Colors.BLUE}{Colors.BOLD}🚀 DÉBUT DES TESTS ME ENDPOINT{Colors.RESET}\n")

    # Test 1: Récupération de mes informations (utilisateur existant)
    print_test_header("Récupération de mes informations")
    response = requests.post(url, 
                           json={"uid": 1},  # Premier utilisateur créé
                           headers=headers)
    if print_test_result("Mes informations", response, 200):
        tests_passed += 1
        try:
            data = response.json()
            if "user" in data:
                print(f"{Colors.GREEN}   ✓ Format de réponse correct{Colors.RESET}")
                print(f"{Colors.BLUE}   📝 Informations reçues:{Colors.RESET}")
                print(json.dumps(data, indent=2))
            else:
                print(f"{Colors.RED}   ✗ Format de réponse incorrect{Colors.RESET}")
        except json.JSONDecodeError:
            print(f"{Colors.RED}   ✗ Réponse non JSON{Colors.RESET}")

    # Test 2: UID invalide
    print_test_header("UID invalide")
    response = requests.post(url,
                           json={"uid": -1},
                           headers=headers)
    if print_test_result("UID invalide", response, 400):
        tests_passed += 1

    # Test 3: UID manquant
    print_test_header("UID manquant")
    response = requests.post(url,
                           json={},
                           headers=headers)
    if print_test_result("UID manquant", response, 400):
        tests_passed += 1

    # Test 4: Format JSON invalide
    print_test_header("Format JSON invalide")
    response = requests.post(url,
                           data="invalid json",
                           headers=headers)
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
    test_me_endpoint()