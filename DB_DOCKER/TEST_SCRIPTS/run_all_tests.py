#!/usr/bin/env python3
# filepath: /Users/charlesleroy/Documents/ecole_19/cursus/ALUMNI/matcha/matcha/DB_DOCKER/run_all_tests.py

import os
import subprocess
import time
import sys
from datetime import datetime

# Couleurs pour le terminal
class Colors:
    GREEN = '\033[92m'    # Succès
    RED = '\033[91m'      # Erreur
    BLUE = '\033[94m'     # Info
    YELLOW = '\033[93m'   # Avertissement
    BOLD = '\033[1m'      # Texte en gras
    RESET = '\033[0m'     # Reset couleur

# Configuration
# Déterminer le chemin absolu de la racine du projet
TEST_DIR = os.path.dirname(os.path.abspath(__file__))
total_tests = 0
passed_tests = 0
failed_tests = 0

def print_separator():
    print(f"{Colors.BLUE}{Colors.BOLD}================================================={Colors.RESET}")

def print_test_header(test_name):
    print_separator()
    print(f"{Colors.BLUE}{Colors.BOLD}🚀 EXÉCUTION: {test_name}{Colors.RESET}")
    print_separator()

def run_test(script_name, test_name, *args):
    global total_tests, passed_tests, failed_tests
    
    print_test_header(test_name)
    
    script_path = os.path.join(TEST_DIR, script_name)
    if not os.path.exists(script_path):
        print(f"{Colors.RED}❌ ERREUR: Script non trouvé: {script_path}{Colors.RESET}")
        failed_tests += 1
        total_tests += 1
        return
    
    start_time = time.time()
    try:
        cmd = [sys.executable, script_path]
        cmd.extend(args)
        process = subprocess.run(cmd, check=True)
        exit_code = process.returncode
    except subprocess.CalledProcessError as e:
        exit_code = e.returncode
    except Exception as e:
        print(f"{Colors.RED}❌ ERREUR D'EXÉCUTION: {str(e)}{Colors.RESET}")
        exit_code = 1
    end_time = time.time()
    duration = round(end_time - start_time, 2)
    
    # Afficher le résultat
    if exit_code == 0:
        print(f"\n{Colors.GREEN}{Colors.BOLD}✅ SUCCÈS: {test_name} ({duration}s){Colors.RESET}")
        passed_tests += 1
    else:
        print(f"\n{Colors.RED}{Colors.BOLD}❌ ÉCHEC: {test_name} ({duration}s) - Code: {exit_code}{Colors.RESET}")
        failed_tests += 1
    total_tests += 1
    
    return exit_code

def create_test_user():
    """Crée un utilisateur de test pour les tests qui en ont besoin"""
    print_test_header("CRÉATION D'UN UTILISATEUR DE TEST")
    script_path = os.path.join(TEST_DIR, "test_create_user.py")
    try:
        subprocess.run([sys.executable, script_path], check=True)
    except:
        # Ignorer l'échec, l'utilisateur peut déjà exister
        pass

def main():
    start_time = time.time()
    
    # En-tête principal
    print(f"{Colors.BLUE}{Colors.BOLD}")
    print("=================================================")
    print("       TESTS AUTOMATISÉS DE L'API MATCHA         ")
    print("=================================================")
    print(f"{Colors.RESET}")
    
    # Vérifier si le répertoire de test existe
    if not os.path.isdir(TEST_DIR):
        print(f"{Colors.RED}{Colors.BOLD}Erreur: Dossier de test introuvable: {TEST_DIR}{Colors.RESET}")
        return 1
    
    # Créer un utilisateur pour les tests
    create_test_user()
    
    # Exécuter tous les tests sans arguments
    run_test("test_create_user.py", "Création d'utilisateur")
    run_test("test_login.py", "Login d'utilisateur")
    run_test("test_create_chatroom.py", "Création de chatroom")
    run_test("test_get_all_chatroom.py", "Récupération de toutes les chatrooms")
    run_test("test_add_user_and_remove_userfrom_chatroom.py", "Gestion des utilisateurs dans les chatrooms")
    run_test("test_new_message.py", "Envoi de message")
    run_test("test_get_messages.py", "Récupération des messages")
    run_test("test_get_user.py", "Récupération d'utilisateur par ID")
    run_test("test_me.py", "Récupération de mon profil")
    run_test("test_update.py", "Mise à jour de profil")
    
    # Tests avec des arguments spécifiques
    run_test("test_get_my_chatroom.py", "Récupération de mes chatrooms", "1")
    run_test("test_get_user_by_username.py", "Récupération d'utilisateur par username", "testuser1")
    run_test("test_change_password.py", "Changement de mot de passe", "1", "Test123!@#", "NewPass456!@#")
    run_test("test_get_chatroom_participants.py", "Récupération des participants d'une chatroom", "1")

    # Nouveaux tests à ajouter
    run_test("test_get_random_user.py", "Récupération d'un utilisateur aléatoire", "1")
    run_test("test_potential_match.py", "Création d'un match potentiel", "1", "2")

    # Résumé final
    end_time = time.time()
    total_duration = round(end_time - start_time, 2)
    
    print(f"\n\n{Colors.BLUE}{Colors.BOLD}================================================={Colors.RESET}")
    print(f"{Colors.BOLD}📊 RÉSUMÉ DES TESTS{Colors.RESET}")
    print(f"{Colors.BLUE}{Colors.BOLD}================================================={Colors.RESET}")
    print(f"⏱️  Durée totale: {total_duration} secondes")
    print(f"🧪 Tests exécutés: {total_tests}")
    print(f"{Colors.GREEN}✅ Tests réussis: {passed_tests}{Colors.RESET}")
    print(f"{Colors.RED}❌ Tests échoués: {failed_tests}{Colors.RESET}")
    
    # Résultat global
    if failed_tests == 0:
        print(f"\n{Colors.GREEN}{Colors.BOLD}🎉 TOUS LES TESTS ONT RÉUSSI! 🎉{Colors.RESET}\n")
        return 0
    else:
        print(f"\n{Colors.RED}{Colors.BOLD}⚠️  CERTAINS TESTS ONT ÉCHOUÉ! ⚠️{Colors.RESET}\n")
        return 1

if __name__ == "__main__":
    sys.exit(main())