# Emails Indexer

Este proyecto consta de dos aplicativos que se deben compilar, se hizo de esta forma para poder reutilizar el modulo `indexer`.

## indexer.go

Esta aplicación recibe como argumento el directorio que se desea analizar, luego recorre todos los subdirectorios en busca de archivos que tengan la estructura de un email. 

Analiza cada archivo para obtener los datos del email, como el remitente, los receptores, el contenido, etc.

Se ejecuta de la siguiente forma

`./indexer -path=enron_mail_20110402 -host=http://localhost:4080 -username="admin" -password=Complexpass#123`

### flags
- `path`: Directorio raiz a analizar
- `host`:  Dirección de la base de datos
- `username`: Usuario de la base de datos
- `password`: Contraseña del usuario de la base de datos

#### TODO

- Analizar los emails cocurrentemente para mejorar performance de la aplicación

## server.go

Esta aplicación levanta un servidor que disponibiliza una intefaz para hacer consultas a la base de datos.

Cuenta con dos rutas implementadas. 

- `"/"`: Sirve la aplicación hecha en Vue.js con Tailwind
- `"/api/search"`: Acepta como parametros la query de busqueda `q` y la página a consultar `page`, hace la consulta a zinc search y devuelve los resultados al front. 

Se ejecuta de la siguiente forma:

`./server -host=http://localhost:4080 -username="admin" -password=Complexpass#123`

### flags

- `host`: Dirección de la base de datos
- `username`: Usuario de la base de datos
- `password`: Contraseña del usuario de la base de datos