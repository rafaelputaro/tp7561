package ipfs

import (
	"strconv"
	"tp/common"
	"tp/peer/dht"
	"tp/peer/helpers"
	"tp/peer/ipfs/internal"
)

type Ipfs struct {
	Node *dht.Node
}

// Retorna una instancia de Ipfs lista para ser utilizada
func NewIpfs(node *dht.Node) *Ipfs {
	toReturn := &Ipfs{
		Node: node,
	}
	return toReturn
}

func (ipfs *Ipfs) AddFile(fileName string) error {
	/**
		@TODO
		a) Leer archivo el cual se encuentra en el volumen local
		b) Fraccionarlo en bloques
		NO c) Crear el primer bloque, guardarlo localmente y agregar key localmente. (AddBlock)
		d) Mientras haya bloques, enviarlos a vecinos con <key-block><fileName-block><data block>
	**/
	reader, err := internal.NewFileReader(fileName)
	if err != nil {
		return err
	}
	defer reader.Close()
	// lectura primer bloque
	blkData, blkNum, _, err := reader.Next()
	// si hay error finalizar retornando error
	if err != nil {
		return err
	}
	blkName := generateBlockName(fileName, blkNum)
	blkKey := helpers.GetKey(blkName)
	end := false
	for !end {
		// leer un nuevo bloque
		nextBlkData, nextBlkNum, nextEof, nextErr := reader.Next()
		nextName := helpers.NULL_KEY_SOURCE_DATA
		// si no hay más bloques el nombre del siguiente es nulo
		if nextEof || nextErr != nil {
			end = true
		} else {
			nextName = generateBlockName(fileName, nextBlkNum)
		}
		// key del siguiente bloque
		nextKey := helpers.GetKey(nextName)
		// generar bloque
		blockToStore := generateBlockToStore(blkData, blkKey, nextKey)

		// @TODO enviar efectivamente a los vecinos
		common.Log.Infof("%v", blockToStore)
		common.Log.Infof("%v", blkKey)
		// todo el biri biri de enviar a vecinos

		// actualizar para el siguiente ciclo
		blkData = nextBlkData
		blkKey = nextKey
		blkName = nextName
	}
	return err
}

func (ipfs *Ipfs) GetFile(node dht.Node, fileName string) (string, error) {
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
func (ipfs *Ipfs) AddBlock(node dht.Node, key []byte, fileName string, data []byte) error {
	/**
		@TODO
		a) Add Key al nodo con el fileName
		b) Guardar el archivo localmente
	**/
	return nil
}

// 256Kb
func (ipfs *Ipfs) GetBlock(node dht.Node, key []byte) ([]byte, error) {
	/**
		@TODO
		a) Desde node obtener el fileName con la key
		b) Leer archivo del volumen
	**/
	return nil, nil
}

// Construye un bloque de la siguiente manera <blockKey><nextBlockKey><data>
func generateBlockToStore(data []byte, key []byte, nextKey []byte) []byte {
	block := []byte{}
	block = append(block, key...)
	block = append(block, nextKey...)
	block = append(block, data...)
	return block
}

// Retorna el nombre del bloque en base al nombre del archivo y el número de bloque
func generateBlockName(fileName string, blockNumber int) string {
	toReturn := fileName
	if blockNumber != 0 {
		toReturn += ".part" + strconv.Itoa(blockNumber)
	}
	return toReturn
}
