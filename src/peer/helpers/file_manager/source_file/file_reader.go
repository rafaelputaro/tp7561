package source_file

import (
	"tp/common/files_common"
	"tp/peer/helpers/file_manager/config_fm"
)

// Retorna un nuevo file reader listo para ser utilizado el cuál permite leer un archivo por bloques
func NewFileReader(filePath string) (*files_common.FileReader, error) {
	return files_common.NewFileReader(filePath, config_fm.BLOCK_SIZE)
}
