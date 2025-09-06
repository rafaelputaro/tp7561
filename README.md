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

https://github.com/7574-sistemas-distribuidos/grpc-example/blob/master/simple-rpc/client/client.go

VENV:

python3 -m venv myenv

source myenv/bin/activate

pip install Jinja2


export PATH=$PATH:$(go env GOPATH)/bin
export PATH=$PATH:$HOME/go/bin
export PATH=$PATH:/usr/local/go/bin



Anillo:



Kademlia:

https://es.wikipedia.org/wiki/Kademlia

https://es.wikipedia.org/wiki/Kad

https://youtu.be/j5cLOODKccI

https://www.wikiwand.com/es/articles/Red_Kad

Vídeo piola:

https://youtu.be/_kCHOpINA5g




grpcurl -plaintext -d  \
  '{}' \
  127.0.0.1:50051 protopb.Operations/Ping
{
 
}

grpc_cli ls localhost:50051 -l

$ grpcurl -plaintext -d  \
  '{ "description": "christmas eve bike class" }' \
  localhost:50051 api.v1.Activity_Log/Insert
{
  "id": 1
}


grpcurl -plaintext -d '{}' \
   -proto ./src/peer/helpers/rpc_ops/protobuf/peer.proto \
    localhost:50051 protopb/Options/Ping
{
  
}