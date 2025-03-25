mongosh --port 27017

db.createUser(
   {
     user: "admin",
     pwd: "8BlanchE8", // or cleartext password
     roles: [
       { role: "userAdminAnyDatabase", db: "admin" },
       { role: "readWriteAnyDatabase", db: "admin" }
     ]
   }
 )
mongosh --port 27017 -u "admin" --authenticationDatabase 'admin' -p

 mongosh --port 27017 --authenticationDatabase "admin" -u "admin" -p "8BlanchE8"

mongo admin -u "admin" -p "8BlanchE8" --authenticationDatabase admin

Pull the MongoDB Docker Image

docker pull mongodb/mongodb-community-server:latest

Run the Image as a Container

docker run --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest


pull the latest version of the MongoDB image.

docker compose pull

get Docker to move your already running container to the new release.

docker compose up -d

docker-compose up -d
