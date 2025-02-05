docker stop sqlite-server
docker rm sqlite-server
echo "🛑 Arrêt du serveur SQLite..."
echo "🚀 Démarrage du serveur SQLite..."
docker build -t sqlite-server .
echo "sleep 4"
sleep 4
docker run -d -p 8181:8181 --name sqlite-server sqlite-server
echo "✅ Le serveur est démarré et accessible sur le port 8181"
echo "verify if server is running"
docker ps 