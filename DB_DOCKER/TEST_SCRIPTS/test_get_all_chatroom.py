#!/usr/bin/env python3
# filepath: /Users/charlesleroy/Documents/ecole_19/cursus/ALUMNI/matcha/matcha/DB_DOCKER/api/tests/test_get_all_chatroom.py

import requests
import json
import sys
import argparse
from colorama import init, Fore, Style

# Initialiser colorama pour les couleurs dans le terminal
init()

def test_get_all_chatroom(base_url="http://localhost:8181"):
    """
    Teste le handler get_all_chatroom qui récupère toutes les chatrooms existantes
    """
    endpoint = f"{base_url}/get_all_chatroom"
    
    # En-têtes de la requête
    headers = {
        "Content-Type": "application/json"
    }
    
    # Afficher les détails de la requête
    print(f"{Fore.BLUE}Envoi d'une requête GET à {endpoint}{Style.RESET_ALL}")
    
    try:
        # Envoyer la requête GET (pas de body car c'est GET)
        response = requests.get(endpoint, headers=headers)
        
        # Afficher le code de statut
        print(f"\n{Fore.CYAN}Code de statut: {response.status_code}{Style.RESET_ALL}")
        
        if response.status_code == 200:
            # Succès : afficher les données des chatrooms
            data = response.json()
            print(f"\n{Fore.GREEN}Liste des chatrooms existantes:{Style.RESET_ALL}")
            print(json.dumps(data, indent=2))
            
            # Vérification de la structure de la réponse
            if "chatRoom" in data and "chatRoomID" in data["chatRoom"]:
                chatrooms = data["chatRoom"]["chatRoomID"]
                print(f"\n{Fore.GREEN}✅ Structure de réponse valide{Style.RESET_ALL}")
                print(f"{Fore.GREEN}Nombre de chatrooms trouvées: {len(chatrooms)}{Style.RESET_ALL}")
                
                # Afficher les IDs des chatrooms
                if len(chatrooms) > 0:
                    print(f"{Fore.GREEN}IDs des chatrooms: {', '.join(map(str, chatrooms))}{Style.RESET_ALL}")
                else:
                    print(f"{Fore.YELLOW}⚠️ Aucune chatroom existante dans le système{Style.RESET_ALL}")
            else:
                print(f"{Fore.RED}❌ Structure de réponse invalide: format inattendu{Style.RESET_ALL}")
                
        elif response.status_code == 405:
            # Méthode non autorisée
            print(f"{Fore.RED}❌ Erreur 405: Méthode non autorisée, seul GET est accepté{Style.RESET_ALL}")
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
    parser = argparse.ArgumentParser(description="Test du handler get_all_chatroom")
    parser.add_argument("--url", default="http://localhost:8181", help="URL de base de l'API (défaut: http://localhost:8181)")
    
    # Analyser les arguments
    args = parser.parse_args()
    
    # Lancer le test
    test_get_all_chatroom(args.url)