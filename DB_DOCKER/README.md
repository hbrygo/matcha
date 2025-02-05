########################################
Makefile 
########################################
pour le makefile , aller dans le dossier faire make. 

Commandes makefile 
make ==> crée l'image docker si elle ne existe pas et / ou lance le docker 
make clean ==> arrête le docker + supprime l'image docker 
make fclean ==> supprime l'image + la db
make down ==> arrete le docker 

#########################################
la db 
########################################
la db est persistante donc meme si tu arrete le docker les données serons à nouveau la 
juste ne supprime pas l'image docker ou l'espace alloué pr la db (sauf si c'est ton intention de la supprimer)



############################################################################
Pour créer un nouvel utilisateur via l'API, voici les spécifications :
############################################################################

Endpoint : POST http://localhost:8181/create_user

Format des données (JSON) :
{
    "nom": "string",        // Le nom de famille de l'utilisateur
    "prenom": "string",     // Le prénom de l'utilisateur
    "email": "string",      // L'email unique de l'utilisateur
    "password": "string"    // Le mot de passe de l'utilisateur
}

Exemple de requête curl :
curl -X POST http://localhost:8181/create_user \
     -H "Content-Type: application/json" \
     -d '{
           "nom": "Dupont",
           "prenom": "Jean",
           "email": "jean.dupont@example.com",
           "password": "motdepasse123"
         }'

Réponses possibles :

200 : Utilisateur créé avec succès
400 : Données manquantes
409 : Email déjà utilisé
500 : Erreur serveur
⚠️ Notes importantes :

Tous les champs sont obligatoires (nom, prenom, email, password)
L'email doit être unique dans la base de données
Les champs ne doivent pas être vides
Le mot de passe doit être sécurisé :
Évitez les mots de passe trop simples
Préférez une combinaison de lettres, chiffres et caractères spéciaux
Longueur minimale recommandée : 8 caractères
Exemple de données valides :
{
    "nom": "Dupont",
    "prenom": "Jean",
    "email": "jean.dupont@example.com",
    "password": "MonMotDePasse123!"
}



#########################################################################################
Pour récupérer les informations d'un utilisateur via l'API, voici les spécifications :
#########################################################################################

Endpoint : POST http://localhost:8181/get_user

Deux méthodes possibles pour récupérer les informations :

Par UID :
{
    "uid": integer    // L'identifiant unique de l'utilisateur
}

Par email et password :
{
    "email": "string",     // L'email de l'utilisateur
    "password": "string"   // Le mot de passe de l'utilisateur
}

Exemples de requêtes curl :

Recherche par UID :
curl -X POST http://localhost:8181/get_user \
     -H "Content-Type: application/json" \
     -d '{
           "uid": 1
         }'

Recherche par email/password :
curl -X POST http://localhost:8181/get_user \
     -H "Content-Type: application/json" \
     -d '{
           "email": "test@example.com",
           "password": "testpassword"
         }'

Réponses possibles :

200 : Utilisateur trouvé, avec les données suivantes :
{
    "user": {
        "uid": 1,
        "nom": "Test",
        "prenom": "User",
        "email": "test@example.com"
    }
}

400 : Données manquantes ou format incorrect
404 : Utilisateur non trouvé
500 : Erreur serveur
⚠️ Notes importantes :

Vous devez fournir SOIT l'uid, SOIT la combinaison email/password
Les informations sensibles (mot de passe) ne sont pas incluses dans la réponse
Les champs fournis ne doivent pas être vides
Utilisez de préférence une connexion sécurisée (HTTPS) pour transmettre les données d'authentification
