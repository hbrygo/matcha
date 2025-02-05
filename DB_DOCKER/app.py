from flask import Flask, jsonify, request
from flask_sqlalchemy import SQLAlchemy
import os
import sqlite3

app = Flask(__name__)
# Configuration de SQLite avec le nouveau chemin
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///data/users.db'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
db = SQLAlchemy(app)

# SQL pour créer la table USERS
CREATE_TABLE_SQL = """
CREATE TABLE IF NOT EXISTS USERS (
    uid INTEGER PRIMARY KEY AUTOINCREMENT,
    nom VARCHAR(80) NOT NULL,
    prenom VARCHAR(80) NOT NULL,
    email VARCHAR(120) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
)
"""

# SQL pour insérer un utilisateur test
INSERT_TEST_USER_SQL = """
INSERT INTO USERS (nom, prenom, email, password)
SELECT 'Test', 'User', 'test@example.com', 'testpassword'
WHERE NOT EXISTS (SELECT 1 FROM USERS WHERE email = 'test@example.com')
"""

def init_db():
    # Assure que le répertoire data existe
    os.makedirs('data', exist_ok=True)
    
    # Connexion à la base de données
    conn = sqlite3.connect('data/users.db')
    cursor = conn.cursor()
    
    try:
        # Création de la table
        cursor.execute(CREATE_TABLE_SQL)
        
        # Insertion de l'utilisateur test
        cursor.execute(INSERT_TEST_USER_SQL)
        
        # Validation des changements
        conn.commit()
        print("Base de données initialisée avec succès")
        
    except Exception as e:
        print(f"Erreur lors de l'initialisation de la base de données: {e}")
        conn.rollback()
    finally:
        cursor.close()
        conn.close()

# Création de la base de données et des tables si elles n'existent pas
with app.app_context():
    init_db()

@app.route('/check-db', methods=['GET'])
def check_db():
    try:
        # Vérifie si le fichier de base de données existe
        if os.path.exists('data/users.db'):
            # Connexion à la base et comptage des utilisateurs
            conn = sqlite3.connect('data/users.db')
            cursor = conn.cursor()
            cursor.execute("SELECT COUNT(*) FROM USERS")
            users_count = cursor.fetchone()[0]
            cursor.close()
            conn.close()
            
            return jsonify({
                "message": "db trouvé",
                "details": f"Table USERS existe avec {users_count} utilisateurs"
            }), 200
        else:
            return jsonify({"message": "db non trouvée"}), 404
    except Exception as e:
        return jsonify({"message": f"Erreur: {str(e)}"}), 500


@app.route('/create_user', methods=['POST'])
def create_user():
    try:
        data = request.get_json()
        
        # Vérifie si toutes les données requises sont présentes
        required_fields = ['nom', 'prenom', 'email', 'password']
        if not all(key in data for key in required_fields):
            return jsonify({
                "message": "Données manquantes",
                "required_fields": required_fields
            }), 400

        conn = sqlite3.connect('data/users.db')
        cursor = conn.cursor()
        
        # Insertion du nouvel utilisateur
        cursor.execute("""
            INSERT INTO USERS (nom, prenom, email, password) 
            VALUES (?, ?, ?, ?)
        """, (data['nom'], data['prenom'], data['email'], data['password']))
        
        conn.commit()
        cursor.close()
        conn.close()
        
        return jsonify({"message": "Utilisateur créé avec succès"}), 200
        
    except sqlite3.IntegrityError:
        return jsonify({"message": "Email déjà utilisé"}), 409
    except Exception as e:
        return jsonify({"message": f"Erreur: {str(e)}"}), 500

@app.route('/get_user', methods=['POST'])
def get_user():
    try:
        data = request.get_json()
        
        conn = sqlite3.connect('data/users.db')
        cursor = conn.cursor()
        
        if 'uid' in data:
            # Recherche par UID
            cursor.execute("""
                SELECT uid, nom, prenom, email 
                FROM USERS 
                WHERE uid = ?
            """, (data['uid'],))
        elif 'email' in data and 'password' in data:
            # Recherche par email et password
            cursor.execute("""
                SELECT uid, nom, prenom, email 
                FROM USERS 
                WHERE email = ? AND password = ?
            """, (data['email'], data['password']))
        else:
            return jsonify({
                "message": "Données manquantes",
                "details": "Fournir soit uid, soit email et password"
            }), 400
            
        user = cursor.fetchone()
        cursor.close()
        conn.close()
        
        if user:
            return jsonify({
                "user": {
                    "uid": user[0],
                    "nom": user[1],
                    "prenom": user[2],
                    "email": user[3]
                }
            }), 200
        else:
            return jsonify({"message": "Utilisateur non trouvé"}), 404
            
    except Exception as e:
        return jsonify({"message": f"Erreur: {str(e)}"}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8181)