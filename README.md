# tp7561

## Plan De Desarrollo

Crear un modulo Peer el cual presentará al exterior las funciones grpc para el manejo de la DHT y las operaciones sobre los archivos del IPFS. 
Los valores a almacenar en la DHT van a ser los nombres de los archivos dados como strings.
Las claves se encriptan mediante SHA1.
Tanto las funcionalidad de la DHT y del IPFS se implementan en paquetes separados del módulo Peer.

![Diagrama De Clases Peer](./docs/DiagramaDeClasesPeer.png)

NOTA: El IPFS actualmente no tiene funciones definidas ya que aguardo requerimientos al respecto, sólo esta en el diamgrama para ilustrar su lugar en la arquitectura.


## Preguntas para el ayudante:

Sobre DHT:

1) ¿Esta bien que al inicio todos los pares hagan Ping al boostrap node y que lo agreguen a la tabla? Por lo tanto el bootstrap va a agregar a cada uno de ellos a la tabla mientras entren.

2) ¿Cuando se recibe un "store" se debe almacenar siempre la clave localmente?

3) Use k=20 contactos por prefijo, que según lo que he leído es lo que se usa normalmente ¿Esta bien?

4) ¿Hace falta usar cache de claves?

5) ¿Tengo que hacer control de colisión de claves de alguna manera?

Sobre IPFS:

1) Requerimientos de la implementación.


## Notas para ejecución:

Se dispone del archivo config.ini el cuál permite configurar la cantidad de pares a ejecutar entre otras cosas:

```
[DEFAULT]
# --------------- PAIRS ---------------
NUMBER_OF_PAIRS = 2
ENTRIES_PER_K_BUCKET = 20
```
Iniciar o crear en entorno virtual:
```
python3 -m venv myenv

source myenv/bin/activate
```
Configurar entorno para protobufer:
```
export PATH=$PATH:$(go env GOPATH)/bin
export PATH=$PATH:$HOME/go/bin
export PATH=$PATH:/usr/local/go/bin
```
Si hace falta instalar Jinja:
```
pip install Jinja2
```
Ejecutar:
```
make-docker-compose-up
```
Ver log:
```
make-docker-compose-logs
```
Detener contendores:
```
make-docker-compose-down
```