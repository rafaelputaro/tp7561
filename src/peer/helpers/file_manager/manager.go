package file_manager

import (
	"errors"
	"os"
	"strconv"
	"tp/common"
	"tp/peer/helpers"
	"tp/peer/helpers/file_manager/internal"
)

const MSG_ERROR_FILE_EXIST = "the file exists"
const MSG_ERROR_CREATING_FILE = "Error creating file: %v"
const MSG_ERROR_WRITING_FILE = "Error writing file: %v"
const MSG_FILE_ADDED = "File %v added with a total of %v blocks"

type ProcessBlockCallBack func(key []byte, fileName string, data []byte) error

// Leer un archivo por bloques y enviando los mismos como bloques con el formato
// <blockKey><nextBlockKey><data> los cuales son procesados con el callback parámetro
func AddFile(fileName string, processBlock ProcessBlockCallBack) error {
	reader, err := internal.NewFileReader(generateInputFilePath(fileName))
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
			blkNum = nextBlkNum
		}
		// key del siguiente bloque
		nextKey := helpers.GetKey(nextName)
		// generar bloque
		blockToStore := generateBlockToStore(blkData, blkKey, nextKey)
		// enviar bloque a vecinos
		processBlock(blkKey, blkName, blockToStore)
		// actualizar para el siguiente ciclo
		blkKey = nextKey
		blkName = nextName
		blkData = nextBlkData
	}
	common.Log.Infof(MSG_FILE_ADDED, fileName, blkNum+1)
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

// Escribe un bloque en un archivo localmente. Retorna error si el archivo ya existe
// o si se presenta algún error de acceso a disco
func StoreBlock(fileName string, data []byte) error {
	filepath := generateIpfsFilePath(fileName)
	// chequear si el archivo ya existe
	if fileExists(filepath) {
		common.Log.Debugf(MSG_ERROR_FILE_EXIST)
		return errors.New(MSG_ERROR_FILE_EXIST)
	}
	file, err := os.Create(filepath)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_CREATING_FILE, err)
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_WRITING_FILE, err)
		return err
	}
	return nil
}

// 256Kb
func GetBlock(key []byte) ([]byte, error) {
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

// Retorna el path completo del archivo de acuerdo a la configuración
func generateInputFilePath(fileName string) string {
	return LocalStorageConfig.InputDataFolder + "/" + fileName
}

// Retorna el path completo del archivo de acuerdo a la configuración
func generateIpfsFilePath(fileName string) string {
	return LocalStorageConfig.StoreIpfsFolder + "/" + fileName
}

// Retorna verdadero si el archivo existe
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return false
	} else {
		return !errors.Is(err, os.ErrNotExist)
	}
}
