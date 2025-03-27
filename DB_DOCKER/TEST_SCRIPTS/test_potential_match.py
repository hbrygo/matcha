#!/usr/bin/env python3
# filepath: /Users/charlesleroy/Documents/ecole_19/cursus/ALUMNI/matcha/matcha/DB_DOCKER/TEST_SCRIPTS/test_potential_match.py

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

def test_potential_match(my_uid, other_uid, base_url="http://localhost:8181"):
    """
    Teste le handler potential_match avec les IDs utilisateurs fournis
    """
    endpoint = f"{base_url}/potential_match"
    
    # Données de la requête
    data = {
        "my_uid": my_uid,
        "other_uid": other_uid
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
            # Succès : afficher le résultat du match
            result = response.json()
            print_success("\nRésultat du match potentiel:")
            print_success(json.dumps(result, indent=2))
            
            # Vérifier si c'est un match mutuel
            if "message" in result and "It's a match" in result["message"]:
                print_success("\n💖 C'est un match! Les deux utilisateurs se sont mutuellement likés.")
            else:
                print_success("\n👍 Like enregistré avec succès, en attente de réciprocité.")
                
        elif response.status_code == 400:
            # Erreur de requête
            print_error("❌ Erreur 400: Données invalides ou manquantes")
            print_error(f"Message: {response.text}")
            
        elif response.status_code == 404:
            # Utilisateur non trouvé
            print_error("❌ Erreur 404: Un ou plusieurs utilisateurs non trouvés")
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
    parser = argparse.ArgumentParser(description="Test du handler potential_match")
    parser.add_argument("my_uid", type=int, help="ID de l'utilisateur qui effectue le like")
    parser.add_argument("other_uid", type=int, help="ID de l'utilisateur qui est liké")
    parser.add_argument("--url", default="http://localhost:8181", help="URL de base de l'API (défaut: http://localhost:8181)")
    
    # Analyser les arguments
    args = parser.parse_args()
    
    # Lancer le test
    test_potential_match(args.my_uid, args.other_uid, args.url)