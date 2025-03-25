mongosh --port 27017



Pull the MongoDB Docker Image

docker pull mongodb/mongodb-community-server:latest

Run the Image as a Container

docker run --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest


pull the latest version of the MongoDB image.

docker compose pull

get Docker to move your already running container to the new release.

docker compose up -d

docker-compose up -d
