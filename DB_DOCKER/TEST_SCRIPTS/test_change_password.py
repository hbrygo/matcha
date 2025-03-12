#!/usr/bin/env python3
# filepath: /Users/charlesleroy/Documents/ecole_19/cursus/ALUMNI/matcha/matcha/DB_DOCKER/api/tests/test_change_password.py

import requests
import json
import sys
import argparse
from colorama import init, Fore, Style

# Initialiser colorama pour les couleurs dans le terminal
init()

def test_change_password(user_id, old_password, new_password, base_url="http://localhost:8181"):
    """
    Teste le handler change_password
    """
    endpoint = f"{base_url}/change_password"
    
    # Données de la requête
    data = {
        "userID": user_id,
        "oldPassword": old_password,
        "newPassword": new_password
    }
    
    # En-têtes de la requête
    headers = {
        "Content-Type": "application/json"
    }
    
    # Afficher les détails de la requête
    print(f"{Fore.BLUE}Envoi d'une requête POST à {endpoint}{Style.RESET_ALL}")
    
    # Masquer partiellement les mots de passe pour l'affichage
    display_data = {
        "userID": user_id,
        "oldPassword": old_password[:2] + '*' * (len(old_password) - 2) if old_password else '',
        "newPassword": new_password[:2] + '*' * (len(new_password) - 2) if new_password else ''
    }
    print(f"{Fore.BLUE}Données: {json.dumps(display_data, indent=2)}{Style.RESET_ALL}")
    
    try:
        # Envoyer la requête
        response = requests.post(endpoint, json=data, headers=headers)
        
        # Afficher le code de statut
        print(f"\n{Fore.CYAN}Code de statut: {response.status_code}{Style.RESET_ALL}")
        
        if response.status_code == 200:
            # Succès
            result = response.json()
            print(f"\n{Fore.GREEN}✅ Mot de passe changé avec succès!{Style.RESET_ALL}")
            print(f"{Fore.GREEN}Réponse: {json.dumps(result, indent=2)}{Style.RESET_ALL}")
                
        elif response.status_code == 400:
            # Erreur de validation
            print(f"{Fore.RED}❌ Erreur 400: Données invalides ou manquantes{Style.RESET_ALL}")
            print(f"{Fore.RED}Message: {response.text}{Style.RESET_ALL}")
            
        elif response.status_code == 401:
            # Mot de passe incorrect
            print(f"{Fore.RED}❌ Erreur 401: Ancien mot de passe incorrect{Style.RESET_ALL}")
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
    parser = argparse.ArgumentParser(description="Test du handler change_password")
    parser.add_argument("user_id", type=int, help="ID de l'utilisateur")
    parser.add_argument("old_password", help="Ancien mot de passe")
    parser.add_argument("new_password", help="Nouveau mot de passe")
    parser.add_argument("--url", default="http://localhost:8181", help="URL de base de l'API (défaut: http://localhost:8181)")
    
    # Analyser les arguments
    args = parser.parse_args()
    
    # Lancer le test
    test_change_password(args.user_id, args.old_password, args.new_password, args.url)