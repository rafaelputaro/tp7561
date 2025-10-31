# tp7561

## Introducción DHS
@TODO desarrollar el tema de DHS

### Kademlia - a Distributed Hash Table implementation | Paper Dissection and Deep-dive

[![](https://markdown-videos-api.jorgenkh.no/youtube/_kCHOpINA5g)](https://www.youtube.com/watch?v=_kCHOpINA5g)

@TODO explicar Kademlia en base al vídeo y el siguiente: https://xlattice.sourceforge.net/components/protocol/kademlia/specs.html


## Introducción IPFS

[![](https://markdown-videos-api.jorgenkh.no/youtube/-ZC1-M3biyo)](https://youtu.be/-ZC1-M3biyo)

@TODO explicar IPFS en base al vídeo y a los siguientes links: https://www.conquerblocks.com/post/ipfs

https://about.ipfs.io/


## Implementación del sistema

Se tiene un modulo Peer el cual presentará al exterior las funciones grpc para el manejo de la DHT y las operaciones sobre los archivos del IPFS. 
Los valores a almacenar en la DHT van a ser los nombres de los archivos dados como strings.
Las claves se encriptan mediante SHA256.

![Diagrama De Clases Peer](./docs/DiagramaDeClasesPeer.png)


Idea sobre la estructura de los archivos:

{
    <contenido_bloque>,
    <id_siguiente_bloque>
}


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
Instalar lorem-ipsum generator p/crear set's de datos /tmp:
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

NOTA: make docker-compose-logs | grep -i <exp> | grep -i <exp>


Guía Prometheus:

https://last9.io/blog/docker-monitoring-with-prometheus-a-step-by-step-guide/

https://youtu.be/WUBjlJzI2a0

https://youtu.be/sNk9NkgTOLs

https://www.youtube.com/watch?v=sNk9NkgTOLs&t=371s

La IA me sugirió: Usar una herramienta de visualización de grafos como Grafana con panel de grafos para ver el derrotero de los mensajes en los nodos.

http://localhost:9090/

http://localhost:3000/login



services:
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    restart: unless-stopped
    ports:
      - "3000:3000"
    volumes:
      - grafana-data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=MiPassSegura2024
      - GF_USERS_ALLOW_SIGN_UP=false