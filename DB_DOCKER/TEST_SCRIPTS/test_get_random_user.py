#!/usr/bin/env python3
# filepath: /Users/charlesleroy/Documents/ecole_19/cursus/ALUMNI/matcha/matcha/DB_DOCKER/TEST_SCRIPTS/test_get_random_user.py

import requests
import json
import sys
import argparse
from colorama import init, Fore, Style

# Initialiser colorama pour les couleurs dans le terminal
try:
    init()
    USE_COLOR = True
except:
    USE_COLOR = False

def test_get_random_user(uid, base_url="http://localhost:8181"):
    """
    Teste le handler get_random_user avec l'ID utilisateur fourni
    """
    endpoint = f"{base_url}/get_random_user"
    
    # Données de la requête
    data = {
        "uid": uid
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
        
        # Afficher le code de statut
        print_info(f"\nCode de statut: {response.status_code}")
        
        if response.status_code == 200:
            # Succès : afficher les données de l'utilisateur aléatoire
            user_data = response.json()
            print_success("\nProfil utilisateur aléatoire:")
            print_success(json.dumps(user_data, indent=2))
            
            # Vérification de la structure de la réponse
            if "user" in user_data and "uid" in user_data["user"]:
                random_user = user_data["user"]
                print_success(f"\n✅ Utilisateur aléatoire trouvé avec l'ID: {random_user['uid']}")
                
                # Vérifier les champs principaux
                required_fields = ["nom", "prenom", "dob", "interests", "pictures", "bio"]
                missing_fields = [field for field in required_fields if field not in random_user]
                
                if missing_fields:
                    print_warning(f"⚠️  Attention: Champs manquants dans la réponse: {', '.join(missing_fields)}")
                else:
                    print_success("✅ Structure de réponse valide: tous les champs requis sont présents")
            else:
                print_error("❌ Structure de réponse invalide: information utilisateur manquante")
                
        elif response.status_code == 400:
            # Erreur de requête
            print_error("❌ Erreur 400: Données invalides ou manquantes")
            print_error(f"Message: {response.text}")
            
        elif response.status_code == 404:
            # Utilisateur non trouvé ou pas d'utilisateurs compatibles
            print_error("❌ Erreur 404: Utilisateur non trouvé ou aucun utilisateur compatible")
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
    parser = argparse.ArgumentParser(description="Test du handler get_random_user")
    parser.add_argument("uid", type=int, help="ID de l'utilisateur")
    parser.add_argument("--url", default="http://localhost:8181", help="URL de base de l'API (défaut: http://localhost:8181)")
    
    # Analyser les arguments
    args = parser.parse_args()
    
    # Lancer le test
    test_get_random_user(args.uid, args.url)