import requests
import json

# Couleurs pour une meilleure lisibilité
class Colors:
    GREEN = '\033[92m'    # Succès
    RED = '\033[91m'      # Erreur
    BLUE = '\033[94m'     # Info
    YELLOW = '\033[93m'   # Avertissement
    BOLD = '\033[1m'      # Texte en gras
    RESET = '\033[0m'     # Reset couleur

def test_create_chatroom():
    url = "http://localhost:8181/create_chatroom"
    headers = {"Content-Type": "application/json"}
    tests_passed = 0
    total_tests = 5

    def print_test_header(test_name: str):
        print(f"\n{Colors.BLUE}{Colors.BOLD}⚡️ TEST: {test_name}{Colors.RESET}")

    def print_test_result(test_name: str, data, response, expected_status: int):
        success = response.status_code == expected_status
        print(f"\n{Colors.YELLOW}📝 Données envoyées:{Colors.RESET}")
        print(json.dumps(data, indent=2))
        
        if success:
            print(f"{Colors.GREEN}✅ SUCCÈS: {test_name}")
            print(f"   Status: {response.status_code} (attendu: {expected_status})")
            if response.text:
                try:
                    parsed = json.loads(response.text)
                    print(f"   Réponse: {json.dumps(parsed, indent=2)}{Colors.RESET}")
                except:
                    print(f"   Réponse: {response.text}{Colors.RESET}")
            return True
        else:
            print(f"{Colors.RED}❌ ÉCHEC: {test_name}")
            print(f"   Status reçu: {response.status_code}")
            print(f"   Status attendu: {expected_status}")
            print(f"   Réponse: {response.text}{Colors.RESET}")
            return False

    print(f"{Colors.BLUE}{Colors.BOLD}🚀 DÉBUT DES TESTS CREATE_CHATROOM{Colors.RESET}\n")

    # Test 1: Création d'une conversation privée
    print_test_header("Création d'une conversation privée")
    private_chat = {
        "user_ids": [1, 2],  # Assurez-vous que ces utilisateurs existent
        "name": "Private Discussion"
    }
    response = requests.post(url, json=private_chat, headers=headers)
    if print_test_result("Conversation privée", private_chat, response, 200):
        tests_passed += 1

    # Test 2: Création d'un groupe
    print_test_header("Création d'un groupe")
    group_chat = {
        "user_ids": [1, 2, 3],  # Assurez-vous que ces utilisateurs existent
        "name": "Group Chat"
    }
    response = requests.post(url, json=group_chat, headers=headers)
    if print_test_result("Groupe", group_chat, response, 200):
        tests_passed += 1

    # Test 3: Aucun utilisateur spécifié
    print_test_header("Aucun utilisateur spécifié")
    no_users = {
        "user_ids": [],
        "name": "Empty Group",
        "is_group": True
    }
    response = requests.post(url, json=no_users, headers=headers)
    if print_test_result("Sans utilisateurs", no_users, response, 400):
        tests_passed += 1

    # Test 4: Conversation privée avec un seul utilisateur
    print_test_header("Conversation privée avec un utilisateur")
    single_user = {
        "user_ids": [1],
        "name": "Single User Chat",
        "is_group": False
    }
    response = requests.post(url, json=single_user, headers=headers)
    if print_test_result("Un seul utilisateur", single_user, response, 400):
        tests_passed += 1

    # Test 5: Utilisateur inexistant
    print_test_header("Utilisateur inexistant")
    nonexistent_user = {
        "user_ids": [999, 1000],  # Ces utilisateurs ne devraient pas exister
        "name": "Nonexistent Users",
        "is_group": False
    }
    response = requests.post(url, json=nonexistent_user, headers=headers)
    if print_test_result("Utilisateur inexistant", nonexistent_user, response, 404):
        tests_passed += 1

    # Résumé final
    print(f"\n{Colors.BLUE}{Colors.BOLD}📊 RÉSUMÉ DES TESTS{Colors.RESET}")
    if tests_passed == total_tests:
        print(f"{Colors.GREEN}✅ SUCCÈS TOTAL: {tests_passed}/{total_tests} tests réussis{Colors.RESET}")
    else:
        print(f"{Colors.RED}❌ ÉCHEC: {tests_passed}/{total_tests} tests réussis")
        print(f"   {total_tests - tests_passed} test(s) ont échoué{Colors.RESET}")

if __name__ == "__main__":
    test_create_chatroom()