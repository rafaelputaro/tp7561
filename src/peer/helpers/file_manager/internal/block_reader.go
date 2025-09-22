package internal

import (
	"errors"
	"io"
	"os"
	"tp/common"
	"tp/peer/helpers"
)

const MSG_ERROR_READING_HEADER = "error reading header"

// <blockKey><nextBlockKey><data>
func ReadBlock(path string) ([]byte, []byte, []byte, error) {
	// abrir archivo
	f, err := os.Open(path)
	if err != nil {
		common.Log.Debugf(MSG_FILE_COULD_NOT_BE_OPONED, err)
		return nil, nil, nil, err
	}
	defer f.Close()
	fileContent := make([]byte, MAX_BLOCK_FILE_SIZE)
	// leer bloque completo
	nBytes, err := f.Read(fileContent)
	if err != nil {
		if err != io.EOF {
			common.Log.Debugf(MSG_ERROR_READING_FILE, err)
			return nil, nil, nil, err
		}
	}
	// parseo de contenido del archivo
	if nBytes > HEADER_BLOCK_FILE_SIZE {
		return nil, nil, nil, errors.New(MSG_ERROR_READING_HEADER)
	}
	// parseo header
	blockKey := fileContent[:helpers.LENGTH_KEY_IN_BITS-1]
	nextBlockKey := fileContent[helpers.LENGTH_IN_BYTES : HEADER_BLOCK_FILE_SIZE-1]
	// paseo data
	data := []byte{}
	if nBytes-HEADER_BLOCK_FILE_SIZE > 0 {
		data = fileContent[HEADER_BLOCK_FILE_SIZE : nBytes-1]
	}
	return blockKey, nextBlockKey, data, nil
}
