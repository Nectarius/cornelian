# cornelian
go run main.go

templ generate ./...

git add .
git commit -m "introduce something ..."
git commit -a --amend -m "remove files"
// add admin panel
// add caching ...

// update redirect url
https://console.cloud.google.com/apis/credentials?inv=1&invt=AbmJwA&project=liibeillius

docker build -t cornelian .

docker run -p 443:443 cornelian
docker run -p 5120:5120 cornelian
docker run -p 5120:5120 --rm -ti cornelian
docker run -d -p 5120:5120 cornelian
docker run -d -p 443:443 cornelian

docker build -t cornelian -f Dockerfile ./cornelian

docker load -i kornelian.tar
docker load -i kornelian_with_mkcert.tar

docker save -o kornelian.tar cornelian:latest

https://stackoverflow.com/questions/61054088/how-do-i-get-ssl-certificates-for-my-golang-application

traefik ? https://docs.traefik.io/




sudo mkcert -cert-file /home/nefarius/certificates/cert.pem \
       -key-file /home/nefarius/certificates/key.pem \
       localhost 127.0.0.1 ::1

sudo mkcert -cert-file /home/nefarius/certificates/kornelian.com.pem -key-file /home/nefarius/certificates/kornelian.com.key kornelian.com

sysctl net.ipv4.ip_unprivileged_port_start=443


Some may say, there is a potential security problem: unprivileged users now may bind to the other privileged ports (444-1024). But you can solve this problem easily with iptables, by blocking other ports:

iptables -I INPUT -p tcp --dport 444:1024 -j DROP
iptables -I INPUT -p udp --dport 444:1024 -j DROP



ssh root@80.190.84.21
