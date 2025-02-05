#!/bin/bash

# Nom de l'image et du conteneur
IMAGE_NAME="sqlite-server"
CONTAINER_NAME="sqlite-server"
VOLUME_NAME="sqlite-data"

echo "🚀 Démarrage du serveur SQLite..."

# Vérifier si le volume existe, sinon le créer
if [[ "$(docker volume ls -q -f name=$VOLUME_NAME)" == "" ]]; then
    echo "📁 Création du volume de données..."
    docker volume create $VOLUME_NAME
fi

# Vérifier si l'image existe, sinon la construire
if [[ "$(docker images -q $IMAGE_NAME 2> /dev/null)" == "" ]]; then
    echo "🏗️  Construction de l'image Docker..."
    docker build -t $IMAGE_NAME .
fi

# Vérifier si le conteneur existe déjà
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    # Si le conteneur existe mais n'est pas en cours d'exécution, le supprimer
    if [ ! "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
        echo "🧹 Suppression de l'ancien conteneur..."
        docker rm $CONTAINER_NAME
    else
        echo "⚠️  Le serveur est déjà en cours d'exécution!"
        exit 0
    fi
fi

# Démarrer le nouveau conteneur
echo "▶️  Démarrage du serveur..."
docker run -d \
    --name $CONTAINER_NAME \
    -p 8181:8181 \
    -v $VOLUME_NAME:/app/data \
    $IMAGE_NAME

# Vérifier que le conteneur est bien démarré
if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    echo "✅ Le serveur est démarré et accessible sur le port 8181"
else
    echo "❌ Erreur lors du démarrage du serveur"
    echo "📝 Logs du conteneur:"
    docker logs $CONTAINER_NAME
fi