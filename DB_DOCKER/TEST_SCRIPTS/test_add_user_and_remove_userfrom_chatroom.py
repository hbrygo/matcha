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

def test_chatroom_user_management():
    base_url = "http://localhost:8181"
    headers = {"Content-Type": "application/json"}
    
    # Variables globales pour suivre les tests
    total_tests = 0
    passed_tests = 0

    def print_test_header(endpoint: str, test_name: str):
        print(f"\n{Colors.BLUE}{Colors.BOLD}⚡️ TEST {endpoint}: {test_name}{Colors.RESET}")

    def print_test_result(test_name: str, data, response, expected_status: int):
        nonlocal total_tests, passed_tests
        total_tests += 1
        
        success = response.status_code == expected_status
        print(f"\n{Colors.YELLOW}📝 Données envoyées:{Colors.RESET}")
        print(json.dumps(data, indent=2))
        
        if success:
            passed_tests += 1
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

    def create_test_chatroom():
        # Créer une chatroom pour les tests
        print(f"{Colors.BLUE}🔧 Création d'une chatroom de test...{Colors.RESET}")
        chatroom_data = {
            "user_ids": [1, 2],  # Assurez-vous que ces utilisateurs existent
            "name": "Test Chatroom",
            "is_group": True
        }
        response = requests.post(f"{base_url}/create_chatroom", json=chatroom_data, headers=headers)
        
        if response.status_code == 200:
            chatroom_id = response.json().get("chatroom_id")
            print(f"{Colors.GREEN}✅ Chatroom créée avec succès, ID: {chatroom_id}{Colors.RESET}")
            return chatroom_id
        else:
            print(f"{Colors.RED}❌ Échec de création de la chatroom: {response.text}{Colors.RESET}")
            return None

    def test_add_user_to_chatroom(chatroom_id):
        print(f"\n{Colors.BLUE}{Colors.BOLD}📌 TESTS ADD_USER_TO_CHATROOM{Colors.RESET}")
        
        # Test 1: Ajouter un utilisateur à un chatroom existant
        print_test_header("ADD", "Ajout d'un utilisateur à un chatroom")
        data = {
            "chatroom_id": chatroom_id,
            "user_id": 3  # Assurez-vous que cet utilisateur existe
        }
        response = requests.post(f"{base_url}/add_user_to_chatroom", json=data, headers=headers)
        print_test_result("Ajout d'utilisateur", data, response, 200)
        
        # Test 2: Ajouter un utilisateur déjà présent
        print_test_header("ADD", "Ajout d'un utilisateur déjà présent")
        response = requests.post(f"{base_url}/add_user_to_chatroom", json=data, headers=headers)
        print_test_result("Utilisateur déjà présent", data, response, 400)
        
        # Test 3: Ajouter un utilisateur inexistant
        print_test_header("ADD", "Ajout d'un utilisateur inexistant")
        data = {
            "chatroom_id": chatroom_id,
            "user_id": 9999  # Utilisateur qui ne devrait pas exister
        }
        response = requests.post(f"{base_url}/add_user_to_chatroom", json=data, headers=headers)
        print_test_result("Utilisateur inexistant", data, response, 404)
        
        # Test 4: Ajouter à un chatroom inexistant
        print_test_header("ADD", "Ajout à un chatroom inexistant")
        data = {
            "chatroom_id": 9999,  # Chatroom qui ne devrait pas exister
            "user_id": 1
        }
        response = requests.post(f"{base_url}/add_user_to_chatroom", json=data, headers=headers)
        print_test_result("Chatroom inexistant", data, response, 404)
        
        # Test 5: Données invalides
        print_test_header("ADD", "Données invalides")
        data = {
            "chatroom_id": -1,
            "user_id": -1
        }
        response = requests.post(f"{base_url}/add_user_to_chatroom", json=data, headers=headers)
        print_test_result("Données invalides", data, response, 400)

    def test_remove_user_from_chatroom(chatroom_id):
        print(f"\n{Colors.BLUE}{Colors.BOLD}📌 TESTS REMOVE_USER_FROM_CHATROOM{Colors.RESET}")
        
        # Test 1: Retirer un utilisateur d'un chatroom
        print_test_header("REMOVE", "Retrait d'un utilisateur d'un chatroom")
        data = {
            "chatroom_id": chatroom_id,
            "user_id": 3  # L'utilisateur que nous avons ajouté précédemment
        }
        response = requests.post(f"{base_url}/remove_user_from_chatroom", json=data, headers=headers)
        print_test_result("Retrait d'utilisateur", data, response, 200)
        
        # Test 2: Retirer un utilisateur déjà retiré
        print_test_header("REMOVE", "Retrait d'un utilisateur absent")
        response = requests.post(f"{base_url}/remove_user_from_chatroom", json=data, headers=headers)
        print_test_result("Utilisateur absent", data, response, 400)
        
        # Test 3: Retirer d'un chatroom inexistant
        print_test_header("REMOVE", "Retrait d'un chatroom inexistant")
        data = {
            "chatroom_id": 9999,  # Chatroom qui ne devrait pas exister
            "user_id": 1
        }
        response = requests.post(f"{base_url}/remove_user_from_chatroom", json=data, headers=headers)
        print_test_result("Chatroom inexistant", data, response, 404)
        
        # Test 4: Données invalides
        print_test_header("REMOVE", "Données invalides")
        data = {
            "chatroom_id": -1,
            "user_id": -1
        }
        response = requests.post(f"{base_url}/remove_user_from_chatroom", json=data, headers=headers)
        print_test_result("Données invalides", data, response, 400)

    # Exécution des tests
    print(f"{Colors.BLUE}{Colors.BOLD}🚀 DÉBUT DES TESTS DE GESTION DES UTILISATEURS DANS LES CHATROOMS{Colors.RESET}")
    
    # Créer une chatroom pour les tests
    chatroom_id = create_test_chatroom()
    
    if chatroom_id:
        # Tester l'ajout d'utilisateurs
        test_add_user_to_chatroom(chatroom_id)
        
        # Tester le retrait d'utilisateurs
        test_remove_user_from_chatroom(chatroom_id)
        
        # Résumé final
        print(f"\n{Colors.BLUE}{Colors.BOLD}📊 RÉSUMÉ DES TESTS{Colors.RESET}")
        if passed_tests == total_tests:
            print(f"{Colors.GREEN}✅ SUCCÈS TOTAL: {passed_tests}/{total_tests} tests réussis{Colors.RESET}")
        else:
            print(f"{Colors.RED}❌ ÉCHEC: {passed_tests}/{total_tests} tests réussis")
            print(f"   {total_tests - passed_tests} test(s) ont échoué{Colors.RESET}")
    else:
        print(f"{Colors.RED}❌ Impossible de lancer les tests sans chatroom de test{Colors.RESET}")

if __name__ == "__main__":
    test_chatroom_user_management()