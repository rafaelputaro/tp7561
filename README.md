# tp7561

## Plan De Desarrollo

Crear un modulo Peer el cual presentará al exterior las funciones grpc para el manejo de la DHT y las operaciones sobre los archivos del IPFS. 
Los valores a almacenar en la DHT van a ser los nombres de los archivos dados como strings (posteriormente puede sumarse la durabilidad u otras características).
Las claves se encriptan mediante SHA256.
Tanto las funcionalidad de la DHT y del IPFS se implementan en paquetes separados del módulo Peer.

![Diagrama De Clases Peer](./docs/DiagramaDeClasesPeer.png)

NOTA: El IPFS actualmente no tiene funciones definidas ya que aguardo requerimientos al respecto, sólo esta en el diamgrama para ilustrar su lugar en la arquitectura.



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
Instalar lorem-ipsum generator:
```
pip install lorem-text
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