package file_manager

import (
	"tp/common"
	"tp/peer/helpers"

	"tp/peer/helpers/file_manager/blocks"
	"tp/peer/helpers/file_manager/source_file"
	"tp/peer/helpers/file_manager/utils"
)

type ProcessBlockCallBack func(key []byte, fileName string, data []byte) error

// Leer un archivo por bloques y enviando los mismos como bloques con el formato
// <blockKey><nextBlockKey><data> los cuales son procesados con el callback parámetro
func AddFile(fileName string, processBlock ProcessBlockCallBack) error {
	reader, err := source_file.NewFileReader(utils.GenerateInputFilePath(fileName))
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
	blkName := blocks.GenerateBlockName(fileName, blkNum)
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
			nextName = blocks.GenerateBlockName(fileName, nextBlkNum)
			blkNum = nextBlkNum
		}
		// key del siguiente bloque
		nextKey := helpers.GetKey(nextName)
		// generar bloque
		blockToStore := blocks.GenerateBlockToStore(blkData, blkKey, nextKey)
		// enviar bloque a vecinos
		processBlock(blkKey, blkName, blockToStore)
		// actualizar para el siguiente ciclo
		blkKey = nextKey
		blkName = nextName
		blkData = nextBlkData
	}
	common.Log.Infof(utils.MSG_FILE_ADDED, fileName, blkNum+1)
	return err
}

func GetFile(fileName string) (string, error) {
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

// Escribe un bloque en un archivo localmente para ser recuperado en una búqueda. Retorna error
// si el archivo ya existe o si se presenta algún error de acceso a disco
func StoreBlock(fileName string, data []byte) error {
	return blocks.StoreBlock(utils.GenerateIpfsFilePath(fileName), data)
}

// Escribe un bloque en un archivo localmente como parte de un archivo a ser recuparado.
// Retorna error si el archivo ya existe o si se presenta algún error de acceso a disco
func StoreBlockOnDownload(fileName string, data []byte) error {
	return blocks.StoreBlock(utils.GenertaIpfsRecoverPath(fileName), data)
}

// Obtiene un block completo con su header y datos.
func GetBlock(fileName string) ([]byte, error) {
	_, data, err := blocks.ReadBlock(utils.GenerateIpfsFilePath(fileName))
	return data, err
}
