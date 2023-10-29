Grupo 12:

Integrantes:

Inti Vidal 201973100-9
Tanay Gonzalez 201973005-3

Consideraciones:

** Esta implementacion contempla una version local del programa sobre la tarea, vale decir, que es funcional mediante 5 terminales desde la misma maquina,
la unica discrepancia se da en que no pudimos testear ni arreglar el codigo para su funcionamiento de forma distribuida, debido a los problemas tecnicos que
presenta el labcomp el dia de la entraga (lo cual se atendio en el correo enviado), esto es, las conexiones que necesita rabbitmq y mongodb para poder acceder de una maquina a otra, aun asi, la funcionalidad del codigo se mantiene cumplida, y la implementacion distribuida bastaria realizar pequeños cambios.

Con esto la conexion mediante gRPC y rabbitmq de forma local funcionan de manera correcta y el codigo es usable en su estado.

** las ejecuciones de los archivos es tal cual se plantea en las instrucciones de la tarea.

Instrucciones para la instalacion:

Se deben instalar las tecnologias correspondientes y obvias como golang y demas.
De forma especifica las teconologias usadas son:

Protoc:
Seguir las instrucciones de instalacion: https://grpc.io/docs/protoc-installation/

gRPC:
Seguir las instrucciones de instalacion: https://grpc.io/docs/languages/go/quickstart/

Es necesario para poder generar los stubs y los archivos que los definen y poder usarlos en las maquinas 1 y 2 la primera ves, en este caso ambos archivos
se comparten en diferentes directorios para simular la conexion remota. estos archivos ya estan creados y no es necesario hacer su generacion de nuevo.

MongoDB:

Seguir las instrucciones de instalacion: https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-ubuntu/
(Esto en el caso de una forma distribuida en las maquinas 2 y 3)

Luego la implementacion siguiente va en la segunda maquina virtual(debido a la eventualidad, es solo necesario generarlo en la maquina local)

1)
sudo systemctl start mongod
(Inicia el servicio de mongodb)

2)
db.createUser({
	user:"admin",
	pwd:"admin",
	roles:["root"]
})
(Declara al usuario maestro para poder manipular la bd)

O si no funciona por algun motivo usar:
db.createUser({
	user:"admin",
	pwd:"admin",
	roles:["dbOwner"]
})

3)
use bookstore
(Crea la base de datos en donde estaran almacenadas las colecciones)

Rabbitmq:

Seguir las largas instrucciones de instalacion: https://www.rabbitmq.com/install-debian.html

luego rabbitmq deberia ser instalado en las maquinas 2 y 3, para poder hacer funcionar las conexiones entre las maquinas (de nuevo, debido al problema, es solo necesario instalarlo de forma local)

1)
systemctl start rabbitmq-server
(inicia el servicio de rabbitmq en la maquina)

2) (Opcional)
sudo systemctl stop rabbitmq-server
(si quieres detener el servicio de rabbitmq)

3) (Opcional)
sudo systemctl status rabbitmq-server
(para ver el estatus del servidor)
