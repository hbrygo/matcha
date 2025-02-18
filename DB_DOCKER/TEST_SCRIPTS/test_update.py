#import requests
#import json
#
#def test_update_endpoint():
#    # URL de l'endpoint
#    url = "http://localhost:8181/update"
#    
#    # Données de test
#    test_data = {
#        "uid": 1,
#        "nom": "Test",
#        "prenom": "User",
#        "dob": "1990-01-01",
#        "gender": "M",
#        "preference": "F",
#        "interests": ["sport", "musique", "cinéma"],
#        "pictures": ["path/to/pic1.jpg", "path/to/pic2.jpg"],
#        "bio": "Une bio de test"
#    }
#
#    # Headers
#    headers = {
#        "Content-Type": "application/json"
#    }
#
#    # Test cases
#    def run_test(data, expected_status):
#        response = requests.post(url, json=data, headers=headers)
#        print(f"\nTest avec données: {json.dumps(data, indent=2)}")
#        print(f"Status attendu: {expected_status}")
#        print(f"Status reçu: {response.status_code}")
#        print(f"Réponse: {json.dumps(response.json(), indent=2) if response.text else 'No content'}")
#        return response.status_code == expected_status
#
#    # Test 1: Mise à jour valide
#    print("\nTest 1: Mise à jour valide")
#    assert run_test(test_data, 200)
#
#    # Test 2: UID invalide
#    print("\nTest 2: UID invalide")
#    invalid_data = test_data.copy()
#    invalid_data["uid"] = -1
#    assert run_test(invalid_data, 400)
#
#    # Test 3: Champs manquants
#    print("\nTest 3: Champs manquants")
#    missing_data = test_data.copy()
#    del missing_data["nom"]
#    assert run_test(missing_data, 400)
#
#    # Test 4: UID inexistant
#    print("\nTest 4: UID inexistant")
#    nonexistent_data = test_data.copy()
#    nonexistent_data["uid"] = 9999
#    assert run_test(nonexistent_data, 404)
#
#if __name__ == "__main__":
#    test_update_endpoint()

import requests
import json
from typing import Dict, Any

# Couleurs pour une meilleure lisibilité
class Colors:
    GREEN = '\033[92m'    # Succès
    RED = '\033[91m'      # Erreur
    BLUE = '\033[94m'     # Info
    YELLOW = '\033[93m'   # Avertissement
    BOLD = '\033[1m'      # Texte en gras
    RESET = '\033[0m'     # Reset couleur

def test_update_endpoint():
    url = "http://localhost:8181/update"
    headers = {"Content-Type": "application/json"}
    tests_passed = 0
    total_tests = 4

    # Données de test pour la mise à jour
    test_data = {
        "uid": 1,  # UID de l'utilisateur créé précédemment
        "nom": "Dupont",
        "prenom": "Jean",
        "dob": "1990-01-01",
        "gender": "M",
        "preference": "F",
        "interests": ["sport", "musique", "lecture"],
        "pictures": ["profile1.jpg", "profile2.jpg"],
        "bio": "Une bio mise à jour pour les tests"
    }

    def print_test_header(test_name: str):
        print(f"\n{Colors.BLUE}{Colors.BOLD}⚡️ TEST: {test_name}{Colors.RESET}")

    def print_test_result(test_name: str, data: Dict[str, Any], response: requests.Response, expected_status: int):
        success = response.status_code == expected_status
        print(f"\n{Colors.YELLOW}📝 Données envoyées:{Colors.RESET}")
        print(json.dumps(data, indent=2))
        
        if success:
            print(f"{Colors.GREEN}✅ SUCCÈS: {test_name}")
            print(f"   Status: {response.status_code} (attendu: {expected_status})")
            if response.text:
                print(f"   Réponse: {response.text}{Colors.RESET}")
            return True
        else:
            print(f"{Colors.RED}❌ ÉCHEC: {test_name}")
            print(f"   Status reçu: {response.status_code}")
            print(f"   Status attendu: {expected_status}")
            print(f"   Réponse: {response.text}{Colors.RESET}")
            return False

    print(f"{Colors.BLUE}{Colors.BOLD}🚀 DÉBUT DES TESTS UPDATE{Colors.RESET}\n")

    # Test 1: Mise à jour valide
    print_test_header("Mise à jour du profil")
    response = requests.post(url, json=test_data, headers=headers)
    if print_test_result("Mise à jour complète", test_data, response, 200):
        tests_passed += 1

    # Test 2: UID invalide
    print_test_header("UID invalide")
    invalid_data = test_data.copy()
    invalid_data["uid"] = -1
    response = requests.post(url, json=invalid_data, headers=headers)
    if print_test_result("UID invalide", invalid_data, response, 400):
        tests_passed += 1

    # Test 3: Champs manquants
    print_test_header("Champs manquants")
    missing_data = test_data.copy()
    del missing_data["nom"]
    response = requests.post(url, json=missing_data, headers=headers)
    if print_test_result("Données incomplètes", missing_data, response, 400):
        tests_passed += 1

    # Test 4: UID inexistant
    print_test_header("UID inexistant")
    nonexistent_data = test_data.copy()
    nonexistent_data["uid"] = 9999
    response = requests.post(url, json=nonexistent_data, headers=headers)
    if print_test_result("Utilisateur inexistant", nonexistent_data, response, 404):
        tests_passed += 1

    # Résumé final
    print(f"\n{Colors.BLUE}{Colors.BOLD}📊 RÉSUMÉ DES TESTS{Colors.RESET}")
    if tests_passed == total_tests:
        print(f"{Colors.GREEN}✅ SUCCÈS TOTAL: {tests_passed}/{total_tests} tests réussis{Colors.RESET}")
    else:
        print(f"{Colors.RED}❌ ÉCHEC: {tests_passed}/{total_tests} tests réussis")
        print(f"   {total_tests - tests_passed} test(s) ont échoué{Colors.RESET}")

if __name__ == "__main__":
    test_update_endpoint()