#!/usr/bin/env python3

import requests
import json
import sys
import argparse
from colorama import init, Fore, Style

# Initialiser colorama si disponible
try:
    init()
    USE_COLOR = True
except:
    USE_COLOR = False

def test_get_user_by_username(username, base_url="http://localhost:8181"):
    """
    Teste le handler get_user_by_username avec le nom d'utilisateur fourni
    """
    endpoint = f"{base_url}/get_user_by_username"
    
    # Données de la requête
    data = {
        "username": username
    }
    
    # En-têtes de la requête
    headers = {
        "Content-Type": "application/json"
    }
    
    # Afficher les détails de la requête
    print_info(f"Envoi d'une requête POST à {endpoint}")
    print_info(f"Données: {json.dumps(data, indent=2)}")
    
    try:
        # Envoyer la requête
        response = requests.post(endpoint, json=data, headers=headers)
        
        # Afficher le code de statut et le corps de la réponse
        print_info(f"\nCode de statut: {response.status_code}")
        
        if response.status_code == 200:
            # Succès : afficher les données de l'utilisateur
            user_data = response.json()
            print_success("\nInformations de l'utilisateur:")
            print_success(json.dumps(user_data, indent=2))
            
            # Vérification de la structure de la réponse
            if "user" in user_data:
                user = user_data["user"]
                required_fields = ["uid", "nom", "prenom", "dob", "gender", "preference", "interests", "pictures", "bio"]
                missing_fields = [field for field in required_fields if field not in user]
                
                if missing_fields:
                    print_warning(f"⚠️  Attention: Champs manquants dans la réponse: {', '.join(missing_fields)}")
                else:
                    print_success("✅ Structure de réponse valide: tous les champs requis sont présents")
                    # Afficher l'UID spécifiquement
                    print_success(f"UID de l'utilisateur: {user['uid']}")
            else:
                print_error("❌ Structure de réponse invalide: clé 'user' manquante")
                
        elif response.status_code == 400:
            # Erreur de requête
            print_error("❌ Erreur 400: Données invalides ou manquantes")
            print_error(f"Message: {response.text}")
            
        elif response.status_code == 404:
            # Utilisateur non trouvé
            print_error(f"❌ Erreur 404: Utilisateur '{username}' non trouvé")
            print_error(f"Message: {response.text}")
            
        elif response.status_code == 500:
            # Erreur serveur
            print_error("❌ Erreur 500: Erreur interne du serveur")
            print_error(f"Message: {response.text}")
            
        else:
            # Autre code de statut
            print_error(f"❓ Code de statut inattendu: {response.status_code}")
            print_error(f"Message: {response.text}")
            
    except requests.exceptions.RequestException as e:
        print_error(f"❌ Erreur de connexion: {e}")
        
def print_info(text):
    if USE_COLOR:
        print(f"{Fore.BLUE}{text}{Style.RESET_ALL}")
    else:
        print(text)

def print_success(text):
    if USE_COLOR:
        print(f"{Fore.GREEN}{text}{Style.RESET_ALL}")
    else:
        print(text)

def print_warning(text):
    if USE_COLOR:
        print(f"{Fore.YELLOW}{text}{Style.RESET_ALL}")
    else:
        print(text)

def print_error(text):
    if USE_COLOR:
        print(f"{Fore.RED}{text}{Style.RESET_ALL}")
    else:
        print(text)
        
if __name__ == "__main__":
    # Configuration du parser d'arguments
    parser = argparse.ArgumentParser(description="Test du handler get_user_by_username")
    parser.add_argument("username", help="Nom d'utilisateur à rechercher")
    parser.add_argument("--url", default="http://localhost:8181", help="URL de base de l'API (défaut: http://localhost:8181)")
    
    # Analyser les arguments
    args = parser.parse_args()
    
    # Lancer le test
    test_get_user_by_username(args.username, args.url)