# Image de base légère
FROM alpine:3.18

# Installation des paquets nécessaires
RUN apk add --no-cache \
    python3 \
    py3-pip \
    sqlite \
    sqlite-dev

# Installation des dépendances Python pour le serveur
RUN pip3 install --no-cache-dir \
    flask \
    flask-sqlalchemy

# Création du répertoire de travail
WORKDIR /app

# Création du répertoire pour la base de données
RUN mkdir -p /app/data

# Copie du fichier de l'application
COPY app.py .

# Exposition du port 8181
EXPOSE 8181

# Lancement de l'application Flask au démarrage
CMD ["python3", "app.py"]