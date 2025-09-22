from lorem_text import lorem
import configparser
import os
import random

BASE_FOLDER = "/tmp/data" 
BASE_PEER_FOLDER="peer-"
BASE_FILE_NAME = "file-"
FILE_EXT = "txt"
MSG_GENERATING_DATASET = "Generating dataset for peer number %d\n"

def main():
    # leer configuración
    config = configparser.ConfigParser()
    config.read("config.ini")
    default_config = config["DEFAULT"]
    number_of_files = int(default_config["NUMBER_OF_FILES"])
    number_of_pairs = int(default_config["NUMBER_OF_PAIRS"])
    min_paragraphs = int(default_config["MIN_PARAGRAPHS"])
    max_paragraphs = int(default_config["MAX_PARAGRAPHS"])
    # crear carpeta datos
    os.makedirs(BASE_FOLDER, exist_ok=True)
    # generar datos para cada par
    for pair_number in range(1, number_of_pairs+1):
        generate_data_for_pair(pair_number, number_of_files, min_paragraphs, max_paragraphs)

# generar archivos para un par específico en la carpeta del par dado    
def generate_data_for_pair(pair_number, number_of_files, min_paragraphs, max_paragraphs):
    folder = f"{BASE_FOLDER}/{BASE_PEER_FOLDER}{pair_number}"
    if os.path.isdir(folder):
        return
    print(MSG_GENERATING_DATASET%(pair_number))
    os.makedirs(folder, exist_ok=True)
    for file_num in range(1, number_of_files):
        file_path = f"{folder}/{BASE_FILE_NAME}{pair_number}-{file_num}.{FILE_EXT}"        
        with open(file_path, "w") as out:
            begin = f"[begin, {file_path}]\n"
            end = f"\n[end, {file_path}]"
            content = lorem.paragraphs(random.randint(min_paragraphs, max_paragraphs))
            out.write(f"{begin}{content}{end}")

if __name__ == '__main__':
    main()