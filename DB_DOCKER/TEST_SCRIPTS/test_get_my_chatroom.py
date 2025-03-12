#!/usr/bin/env python3
# filepath: /Users/charlesleroy/Documents/ecole_19/cursus/ALUMNI/matcha/matcha/DB_DOCKER/api/tests/test_get_my_chatroom.py

import requests
import json
import sys
import argparse
from colorama import init, Fore, Style

# Initialiser colorama pour les couleurs dans le terminal
init()

def test_get_my_chatroom(user_id, base_url="http://localhost:8181"):
    """
    Teste le handler get_my_chatroom avec l'ID utilisateur fourni
    """
    endpoint = f"{base_url}/get_my_chatroom"
    
    # Données de la requête
    data = {
        "userID": user_id
    }
    
    # En-têtes de la requête
    headers = {
        "Content-Type": "application/json"
    }
    
    # Afficher les détails de la requête
    print(f"{Fore.BLUE}Envoi d'une requête POST à {endpoint}{Style.RESET_ALL}")
    print(f"{Fore.BLUE}Données: {json.dumps(data, indent=2)}{Style.RESET_ALL}")
    
    try:
        # Envoyer la requête
        response = requests.post(endpoint, json=data, headers=headers)
        
        # Afficher le code de statut et le corps de la réponse
        print(f"\n{Fore.CYAN}Code de statut: {response.status_code}{Style.RESET_ALL}")
        
        if response.status_code == 200:
            # Succès : afficher les données des chatrooms
            chatroom_data = response.json()
            print(f"\n{Fore.GREEN}Chatrooms de l'utilisateur:{Style.RESET_ALL}")
            print(json.dumps(chatroom_data, indent=2))
            
            # Vérification de la structure de la réponse
            if "user" in chatroom_data and "chatRoomID" in chatroom_data["user"]:
                chatrooms = chatroom_data["user"]["chatRoomID"]
                print(f"\n{Fore.GREEN}✅ Structure de réponse valide{Style.RESET_ALL}")
                print(f"{Fore.GREEN}Nombre de chatrooms trouvées: {len(chatrooms)}{Style.RESET_ALL}")
                
                # Afficher les IDs des chatrooms
                if len(chatrooms) > 0:
                    print(f"{Fore.GREEN}IDs des chatrooms: {', '.join(map(str, chatrooms))}{Style.RESET_ALL}")
                else:
                    print(f"{Fore.YELLOW}⚠️ Aucune chatroom trouvée pour cet utilisateur{Style.RESET_ALL}")
            else:
                print(f"{Fore.RED}❌ Structure de réponse invalide: format inattendu{Style.RESET_ALL}")
                
        elif response.status_code == 400:
            # Erreur de requête
            print(f"{Fore.RED}❌ Erreur 400: Données invalides ou manquantes{Style.RESET_ALL}")
            print(f"{Fore.RED}Message: {response.text}{Style.RESET_ALL}")
            
        elif response.status_code == 404:
            # Utilisateur non trouvé
            print(f"{Fore.RED}❌ Erreur 404: Utilisateur avec ID {user_id} non trouvé{Style.RESET_ALL}")
            print(f"{Fore.RED}Message: {response.text}{Style.RESET_ALL}")
            
        elif response.status_code == 500:
            # Erreur serveur
            print(f"{Fore.RED}❌ Erreur 500: Erreur interne du serveur{Style.RESET_ALL}")
            print(f"{Fore.RED}Message: {response.text}{Style.RESET_ALL}")
            
        else:
            # Autre code de statut
            print(f"{Fore.RED}❓ Code de statut inattendu: {response.status_code}{Style.RESET_ALL}")
            print(f"{Fore.RED}Message: {response.text}{Style.RESET_ALL}")
            
    except requests.exceptions.RequestException as e:
        print(f"{Fore.RED}❌ Erreur de connexion: {e}{Style.RESET_ALL}")
        
if __name__ == "__main__":
    # Configuration du parser d'arguments
    parser = argparse.ArgumentParser(description="Test du handler get_my_chatroom")
    parser.add_argument("user_id", type=int, help="ID de l'utilisateur")
    parser.add_argument("--url", default="http://localhost:8181", help="URL de base de l'API (défaut: http://localhost:8181)")
    
    # Analyser les arguments
    args = parser.parse_args()
    
    # Lancer le test
    test_get_my_chatroom(args.user_id, args.url)