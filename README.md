# tp7561

Plan de Implementación: 

DHT: Implementar chord con un anillo de tamaño "n" en principio con búsqueda secuencial. Luego ir agregando fingertable. alta/baja de nodos y replicación. Implementado con gRPC, protobuf en Golang.

-----------------------------------------------------------------------------

DHT:

https://es.wikipedia.org/wiki/Tabla_de_hash_distribuida


IPFS (Interplanetary File System).
Búsqueda por contenido en nodos que se dan de alta o baja dinámica. CID resultado de aplicarle una función de hash al contenido. Una tabla del sistema mapea el CID a ciertas ubicaciones prefiriendo la más cercana en una búsqueda. 




https://hazelcast.com/foundations/distributed-computing/distributed-hash-table/

https://stackoverflow.com/questions/144360/simple-basic-explanation-of-a-distributed-hash-table-dht

https://www.geeksforgeeks.org/system-design/distributed-hash-tables-with-kademlia/

https://github.com/savoirfairelinux/opendht

https://www.google.com/search?q=arquitectura+para+distributed+hash+table&oq=arquitectura+para+distributed+hash+table&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIHCAEQIRigATIHCAIQIRigATIHCAMQIRifBTIHCAQQIRifBTIHCAUQIRifBTIHCAYQIRifBTIHCAcQIRifBTIHCAgQIRifBTIHCAkQIRifBdIBCTE0NTM1ajBqN6gCALACAA&sourceid=chrome&ie=UTF-8



GRPC:

https://earthly.dev/blog/golang-grpc-example/

VENV:

python3 -m venv myenv

source myenv/bin/activate

pip install Jinja2

