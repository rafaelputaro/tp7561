package ipfs

import "tp/peer/dht"

func AddFile(node dht.Node, fileName string) error {
	/**
		@TODO
		a) Leer archivo el cual se encuentra en el volumen local
		b) Fraccionarlo en bloques
		c) Crear el primer bloque, guardarlo localmente y agregar key localmente. (AddBlock)
		d) Mientras haya bloques, enviarlos a vecinos con <key-block><fileName-block><data block>
	**/

	return nil
}

func GetFile(node dht.Node, fileName string) (string, error) {
	/**
		@TODO
		a) En base al nombre del archivo calcular key
		b) Con la key buscar un nodo que la tenga
		c) Pedirle bloque al nodo
		d) Guardar el bloque localmente
		e) El bloque pedido tendrá la key del siguiente bloque......
		f) Una vez que tengo todos los bloques reconstruyo el archivo
	**/
	return "", nil
}

// 256Kb
func AddBlock(node dht.Node, key []byte, fileName string, data []byte) error {
	/**
		@TODO
		a) Add Key al nodo con el fileName
		b) Guardar el archivo localmente
	**/
	return nil
}

// 256Kb
func GetBlock(node dht.Node, key []byte) ([]byte, error) {
	/**
		@TODO
		a) Desde node obtener el fileName con la key
		b) Leer archivo del volumen
	**/
	return nil, nil
}
