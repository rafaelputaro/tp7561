package config_fm

import "tp/peer/helpers"

const INVALID_BLOCK_NUMBER = -1
const BLOCK_SIZE = 256 * 1024 // tamaño de los bloques en bytes
const HEADER_BLOCK_FILE_SIZE = 2 * helpers.LENGTH_KEY_IN_BYTES
const MAX_BLOCK_FILE_SIZE = HEADER_BLOCK_FILE_SIZE + BLOCK_SIZE
const DOWLOAD_SUB_DIRECTORY = "down"
const RESTORE_SUB_DIRECTORY = "restore"
const UPLOAD_SUB_DIRECTORY = "upload"
