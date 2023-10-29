Si quieres implementar el cliente y el servidor por separado
tienes que generar los .pb.go con el order.proto en ambas maquinas o
para simularlo en dos directorios distintos.

Asegurate que el go.mod este bien definido en ambos directorios
Por el momento se puede trabajar en un solo directorio asi que no 
deberia haber problema si no quieres complicarte con esto ahora

Comando para generar los stubs:
protoc --go_out=. --go-grpc_out=. proto/order.proto


INSTRUCCIONES MONGODB

1.
sudo apt-get install gnupg curl

2. importar gpg public key
curl -fsSL https://pgp.mongodb.com/server-7.0.asc | \
   sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg \
   --dearmor

3. Create the list file /etc/apt/sources.list.d/mongodb-org-7.0.list for your version of Ubuntu.
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list

4. 
sudo apt-get update

5.
sudo apt-get install -y mongodb-org

6. pausar auto updates
echo "mongodb-org hold" | sudo dpkg --set-selections
echo "mongodb-org-database hold" | sudo dpkg --set-selections
echo "mongodb-org-server hold" | sudo dpkg --set-selections
echo "mongodb-mongosh hold" | sudo dpkg --set-selections
echo "mongodb-org-mongos hold" | sudo dpkg --set-selections
echo "mongodb-org-tools hold" | sudo dpkg --set-selections

7.
sudo systemctl start mongod


db.createUser({
	user:"admin",
	pwd:"admin",
	roles:["root"]
})

use bookstore
db.createUser({
	user:"admin",
	pwd:"admin",
	roles:["dbOwner"]
})

