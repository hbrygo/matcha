# Variables
IMAGE_NAME = sqlite-server
CONTAINER_NAME = sqlite-server
PORT = 8181

# Règle par défaut - vérifie si l'image existe et lance le conteneur
all:
	@if ! docker images $(IMAGE_NAME) | grep -q $(IMAGE_NAME); then \
		echo "Building Docker image..." ; \
		docker build -t $(IMAGE_NAME) . ; \
	else \
		echo "Docker image already exists" ; \
	fi
	@if ! docker ps -a | grep -q $(CONTAINER_NAME); then \
		echo "Starting container..." ; \
		docker run -d -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME) ; \
	else \
		echo "Container already exists, starting if not running..." ; \
		docker start $(CONTAINER_NAME) 2>/dev/null || true ; \
	fi

# Arrête le conteneur (équivalent à docker stop)
down:
	@echo "Stopping container..."
	@docker stop $(CONTAINER_NAME) || true

# Arrête le conteneur sans supprimer l'image et la db
clean:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

# Supprime tout (conteneur, image et db)
fclean: clean
	docker rmi $(IMAGE_NAME) || true
	docker rm $(CONTAINER_NAME) || true
	rm -rf data/users.db || true

# Reconstruit tout depuis zéro
re: fclean all

# Empêche les conflits avec les fichiers du même nom
.PHONY: all clean fclean re down