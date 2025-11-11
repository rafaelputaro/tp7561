import sys
from jinja2 import Environment, FileSystemLoader 
import configparser

PEER_PORT_GRPC = 50051
PEER_PORT_TCP = 8080
CLIENT_PORT_TCP = 8080

def main():
    config = configparser.ConfigParser()
    config.read("config.ini")
    default_config = config["DEFAULT"]
    env = Environment(loader=FileSystemLoader("templates"))
    template = env.get_template("docker-compose-dev.yaml.jinja")
    output = template.render(
        number_of_pairs=int(default_config["NUMBER_OF_PAIRS"]),
        size_of_peer_groups=int(default_config["SIZE_OF_PEER_GROUPS"]),
        metrics_host=default_config["METRICS_HOST"],
        metrics_base_port=int(default_config["METRICS_BASE_PORT"]),
        prometheus_port=default_config["PROMETHEUS_PORT"],
        number_of_clients=int((default_config["NUMBER_OF_CLIENTS"])),
        entries_per_k_bucket=int(default_config["ENTRIES_PER_K_BUCKET"]),
        login_format_for_keys=default_config["LOGIN_FORMAT_FOR_KEYS"],
        host_folder=default_config["HOST_FOLDER"],
        client_folder=default_config["CLIENT_FOLDER"],
        app_folder=default_config["APP_FOLDER"],        
        input_data_folder=default_config["INPUT_DATA_FOLDER"],
        store_folder=default_config["STORE_FOLDER"],       
        search_workers=int((default_config["SEARCH_WORKERS"])),
        peer_port_grpc = PEER_PORT_GRPC,
        peer_port_tcp = PEER_PORT_TCP,
        clien_port_tpc = CLIENT_PORT_TCP,
    )
    with open("docker-compose-dev.yaml", "w") as f:
        f.write(output)

if __name__ == '__main__':
    if len(sys.argv) != 1:
        print("Usage: python3 compose-generator.py")
        sys.exit(1)

    main()