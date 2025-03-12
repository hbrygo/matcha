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

def test_get_message():
    url = "http://localhost:8181/get_message"
    headers = {"Content-Type": "application/json"}
    tests_passed = 0
    total_tests = 5
    
    # Variables pour stocker les IDs nécessaires aux tests
    chatroom_id = None
    user_id = 1  # Assurez-vous que cet utilisateur existe
    
    def print_test_header(test_name: str):
        print(f"\n{Colors.BLUE}{Colors.BOLD}⚡️ TEST: {test_name}{Colors.RESET}")

    def print_test_result(test_name: str, data, response, expected_status: int):
        nonlocal tests_passed
        success = response.status_code == expected_status
        print(f"\n{Colors.YELLOW}📝 Données envoyées:{Colors.RESET}")
        print(json.dumps(data, indent=2))
        
        if success:
            tests_passed += 1
            print(f"{Colors.GREEN}✅ SUCCÈS: {test_name}")
            print(f"   Status: {response.status_code} (attendu: {expected_status})")
            if response.text:
                try:
                    parsed = json.loads(response.text)
                    print(f"   Réponse: {json.dumps(parsed, indent=2)}{Colors.RESET}")
                except:
                    print(f"   Réponse: {response.text}{Colors.RESET}")
        else:
            print(f"{Colors.RED}❌ ÉCHEC: {test_name}")
            print(f"   Status reçu: {response.status_code}")
            print(f"   Status attendu: {expected_status}")
            print(f"   Réponse: {response.text}{Colors.RESET}")

    # Créer une chatroom et ajouter un message pour les tests
    def setup_test_environment():
        print(f"\n{Colors.BLUE}🔧 Préparation de l'environnement de test...{Colors.RESET}")
        
        # 1. Créer une chatroom
        chatroom_data = {
            "user_ids": [1, 2],  # Ces utilisateurs doivent exister
            "name": "Test Message Room"
        }
        response = requests.post("http://localhost:8181/create_chatroom", 
                                json=chatroom_data, 
                                headers=headers)
        
        if response.status_code != 200:
            print(f"{Colors.RED}❌ Échec de création de la chatroom{Colors.RESET}")
            return None
            
        chatroom_id = response.json().get("chatroom_id")
        
        # 2. Ajouter un message
        message_data = {
            "message": "Ceci est un message de test",
            "chatRoom": chatroom_id,
            "userID": 1
        }
        response = requests.post("http://localhost:8181/new_message",
                               json=message_data,
                               headers=headers)
        
        if response.status_code != 200:
            print(f"{Colors.RED}❌ Échec d'ajout du message{Colors.RESET}")
            return None
            
        print(f"{Colors.GREEN}✅ Environnement préparé avec succès, chatroom ID: {chatroom_id}{Colors.RESET}")
        return chatroom_id

    print(f"{Colors.BLUE}{Colors.BOLD}🚀 DÉBUT DES TESTS GET_MESSAGE{Colors.RESET}\n")
    
    # Préparer l'environnement de test
    chatroom_id = setup_test_environment()
    if not chatroom_id:
        print(f"{Colors.RED}❌ Impossible de continuer les tests sans environnement{Colors.RESET}")
        return
    
    # Test 1: Récupération des messages valide
    print_test_header("Récupération des messages")
    valid_req = {
        "chatRoom": chatroom_id,
        "userID": 1
    }
    response = requests.post(url, json=valid_req, headers=headers)
    print_test_result("Récupération des messages", valid_req, response, 200)
    
    # Test 2: Chatroom invalide
    print_test_header("Chatroom invalide")
    invalid_chatroom = {
        "chatRoom": -1,
        "userID": 1
    }
    response = requests.post(url, json=invalid_chatroom, headers=headers)
    print_test_result("Chatroom invalide", invalid_chatroom, response, 400)
    
    # Test 3: Utilisateur invalide
    print_test_header("Utilisateur invalide")
    invalid_user = {
        "chatRoom": chatroom_id,
        "userID": -1
    }
    response = requests.post(url, json=invalid_user, headers=headers)
    print_test_result("Utilisateur invalide", invalid_user, response, 400)
    
    # Test 4: Utilisateur non membre
    print_test_header("Utilisateur non membre")
    nonmember_user = {
        "chatRoom": chatroom_id,
        "userID": 999  # Cet utilisateur ne devrait pas exister
    }
    response = requests.post(url, json=nonmember_user, headers=headers)
    print_test_result("Utilisateur non membre", nonmember_user, response, 404)
    
    # Test 5: Chatroom inexistante
    print_test_header("Chatroom inexistante")
    nonexistent_chatroom = {
        "chatRoom": 9999,  # Cette chatroom ne devrait pas exister
        "userID": 1
    }
    response = requests.post(url, json=nonexistent_chatroom, headers=headers)
    print_test_result("Chatroom inexistante", nonexistent_chatroom, response, 404)
    
    # Résumé final
    print(f"\n{Colors.BLUE}{Colors.BOLD}📊 RÉSUMÉ DES TESTS{Colors.RESET}")
    if tests_passed == total_tests:
        print(f"{Colors.GREEN}✅ SUCCÈS TOTAL: {tests_passed}/{total_tests} tests réussis{Colors.RESET}")
    else:
        print(f"{Colors.RED}❌ ÉCHEC: {tests_passed}/{total_tests} tests réussis")
        print(f"   {total_tests - tests_passed} test(s) ont échoué{Colors.RESET}")

if __name__ == "__main__":
    test_get_message()