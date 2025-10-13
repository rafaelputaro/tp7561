from lorem_text import lorem
import configparser
import os
import random

#BASE_FOLDER = "/tmp/data" 
BASE_PEER_FOLDER="peer-"
BASE_CLIENT_FOLDER="client-"
BASE_FILE_NAME_PEER = "file-"
BASE_FILE_NAME_CLIENT = "filec-"
FILE_EXT = "txt"
MSG_GENERATING_DATASET_FOR_PEERS = "Generating dataset for peers\n"
MSG_GENERATING_DATASET_FOR_CLIENTS = "Generating dataset for clients\n"
MSG_GENERATING_DATASET = "Generating dataset for number %d\n"

def main():
    # leer configuración
    config = configparser.ConfigParser()
    config.read("config.ini")
    default_config = config["DEFAULT"]
    number_of_files_clients = int(default_config["NUMBER_OF_FILES_CLIENTS"])
    number_of_files_peers = int(default_config["NUMBER_OF_FILES_PEERS"])
    number_of_clients = int(default_config["NUMBER_OF_CLIENTS"])
    number_of_pairs = int(default_config["NUMBER_OF_PAIRS"])
    min_paragraphs = int(default_config["MIN_PARAGRAPHS"])
    max_paragraphs = int(default_config["MAX_PARAGRAPHS"])
    host_folder = default_config["HOST_FOLDER"]
    input_data_folder = default_config["INPUT_DATA_FOLDER"]
    base_folder = host_folder+input_data_folder
    # crear carpeta datos
    os.makedirs(base_folder, exist_ok=True)
    # generar datos para cada par
    print(MSG_GENERATING_DATASET_FOR_PEERS)
    for pair_number in range(1, number_of_pairs+1):
        generate_data_for_pair(base_folder, pair_number, number_of_files_peers, min_paragraphs, max_paragraphs)
    # generar datos para cada cliente
    print(MSG_GENERATING_DATASET_FOR_CLIENTS)
    for client_number in range(1, number_of_clients+1):
        generate_data_for_client(base_folder, client_number, number_of_files_clients, min_paragraphs, max_paragraphs)

# generar archivos para un par específico en la carpeta del par dado    
def generate_data_for_pair(base_folder, pair_number, number_of_files, min_paragraphs, max_paragraphs):
    generate_data_for_entity(base_folder, pair_number, BASE_FILE_NAME_PEER, BASE_PEER_FOLDER, number_of_files, min_paragraphs, max_paragraphs)

# generar archivos para un par específico en la carpeta del par dado    
def generate_data_for_client(base_folder, pair_number, number_of_files, min_paragraphs, max_paragraphs):
    generate_data_for_entity(base_folder, pair_number, BASE_FILE_NAME_CLIENT, BASE_CLIENT_FOLDER, number_of_files, min_paragraphs, max_paragraphs)

# generar archivos para una entidad específica en la carpeta de la entidad dada
def generate_data_for_entity(base_folder, entity_number, base_file_name, base_entity_folder, number_of_files, min_paragraphs, max_paragraphs):
    folder = f"{base_folder}/{base_entity_folder}{entity_number}"
    if os.path.isdir(folder):
        return
    print(MSG_GENERATING_DATASET%(entity_number))
    os.makedirs(folder, exist_ok=True)
    for file_num in range(1, number_of_files+1):
        file_path = f"{folder}/{base_file_name}{entity_number}-{file_num}.{FILE_EXT}"        
        with open(file_path, "w") as out:
            begin = f"[begin, {file_path}]\n"
            end = f"\n[end, {file_path}]"
            content = lorem.paragraphs(random.randint(min_paragraphs, max_paragraphs))
            out.write(f"{begin}{content}{end}")

if __name__ == '__main__':
    main()