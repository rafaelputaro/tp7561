package internal

import (
	"bufio"
	"io"
	"os"
	"tp/common"
)

// Permite leer un archivo de a bloques
type FileReader struct {
	fd           *os.File
	reader       *bufio.Reader
	currentBlock []byte
	blockNumber  int
	eof          bool
}

// Retorna un nuevo file reader listo para ser utilizado el cuál permite leer un archivo por bloques
func NewFileReader(filePath string) (*FileReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		common.Log.Errorf(MSG_FILE_COULD_NOT_BE_OPONED, err)
		return nil, err
	}
	r := bufio.NewReader(file)
	toReturn := FileReader{
		fd:           file,
		reader:       r,
		currentBlock: []byte{},
		blockNumber:  INVALID_BLOCK_NUMBER,
		eof:          false,
	}
	return &toReturn, nil
}

// Cierra el archivo asociado
func (file *FileReader) Close() {
	if err := file.fd.Close(); err != nil {
		common.Log.Errorf(MSG_ERROR_CLOSING_FILE, err)
	}
}

// Lee el siguiente bloque del archivo retornando bajo el siguiente formato <bloque actual>,<numero de bloque><eof><error>
// En caso de error retorna: <nil><-1><false><err>
// En caso de eof:  <nil><-1><true><nil>
func (file *FileReader) Next() ([]byte, int, bool, error) {
	if file.eof {
		return nil, INVALID_BLOCK_NUMBER, true, nil
	}
	b := make([]byte, BLOCK_SIZE)
	n, err := file.reader.Read(b)
	// si es eof retorna <nil><-1><true><nil>
	if err == io.EOF {
		file.eof = true
		return nil, INVALID_BLOCK_NUMBER, true, nil
	}
	// si hay error retorna <nil><-1><false><err>
	if err != nil {
		common.Log.Errorf(MSG_ERROR_READING_FILE, err)
		return nil, file.blockNumber, file.eof, err
	}
	file.currentBlock = b[0:n]
	file.blockNumber++
	return file.currentBlock, file.blockNumber, file.eof, nil
}

// Retorna el último bloque leído junto al número bloque y si se ha alcanzado el final del archivo.
func (file *FileReader) CurrentBlock() ([]byte, int, bool) {
	return file.currentBlock, file.blockNumber, file.eof
}
