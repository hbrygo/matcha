#!/bin/bash

CONTAINER_NAME="sqlite-server"

echo "🛑 Arrêt du serveur SQLite..."

# Vérifier si le conteneur existe et tourne
if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    # Arrêter le conteneur
    echo "⏹️  Arrêt du conteneur..."
    docker stop $CONTAINER_NAME
    
    # Supprimer le conteneur
    echo "🗑️  Suppression du conteneur..."
    docker rm $CONTAINER_NAME
    
    echo "✅ Le serveur a été arrêté avec succès"
else
    echo "ℹ️  Aucun serveur en cours d'exécution"
fi