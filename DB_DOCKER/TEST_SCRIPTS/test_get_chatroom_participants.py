#!/usr/bin/env python3
# filepath: /Users/charlesleroy/Documents/ecole_19/cursus/ALUMNI/matcha/matcha/DB_DOCKER/TEST_SCRIPTS/test_get_chatroom_participants.py

import requests
import json
import sys
import argparse
from colorama import init, Fore, Style

# Initialiser colorama pour les couleurs dans le terminal
init()

def test_get_chatroom_participants(chatroom_id, base_url="http://localhost:8181"):
    """
    Teste le handler get_chatroom_participants avec l'ID de chatroom fourni
    """
    endpoint = f"{base_url}/get_chatroom_participants"
    
    # Données de la requête
    data = {
        "chatRoomID": chatroom_id
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
        
        # Afficher le code de statut
        print(f"\n{Fore.CYAN}Code de statut: {response.status_code}{Style.RESET_ALL}")
        
        if response.status_code == 200:
            # Succès : afficher les participants
            result = response.json()
            print(f"\n{Fore.GREEN}Liste des participants de la chatroom {chatroom_id}:{Style.RESET_ALL}")
            print(json.dumps(result, indent=2))
            
            # Vérifier que la réponse contient la liste des participants
            if "participants" in result:
                participants = result["participants"]
                print(f"\n{Fore.GREEN}✅ Structure de réponse valide{Style.RESET_ALL}")
                print(f"{Fore.GREEN}Nombre de participants: {len(participants)}{Style.RESET_ALL}")
                
                # Afficher les IDs des participants
                if len(participants) > 0:
                    print(f"{Fore.GREEN}IDs des participants: {', '.join(map(str, participants))}{Style.RESET_ALL}")
                else:
                    print(f"{Fore.YELLOW}⚠️ Aucun participant trouvé pour cette chatroom{Style.RESET_ALL}")
            else:
                print(f"{Fore.RED}❌ Structure de réponse invalide: clé 'participants' manquante{Style.RESET_ALL}")
                
        elif response.status_code == 400:
            # Erreur de validation
            print(f"{Fore.RED}❌ Erreur 400: Données invalides ou manquantes{Style.RESET_ALL}")
            print(f"{Fore.RED}Message: {response.text}{Style.RESET_ALL}")
            
        elif response.status_code == 404:
            # Chatroom non trouvée
            print(f"{Fore.RED}❌ Erreur 404: Chatroom avec ID {chatroom_id} non trouvée{Style.RESET_ALL}")
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
    parser = argparse.ArgumentParser(description="Test du handler get_chatroom_participants")
    parser.add_argument("chatroom_id", type=int, help="ID de la chatroom")
    parser.add_argument("--url", default="http://localhost:8181", help="URL de base de l'API (défaut: http://localhost:8181)")
    
    # Analyser les arguments
    args = parser.parse_args()
    
    # Lancer le test
    test_get_chatroom_participants(args.chatroom_id, args.url)