import sys
from jinja2 import Environment, FileSystemLoader 
import configparser

def main():
    config = configparser.ConfigParser()
    config.read("config.ini")
    default_config = config["DEFAULT"]
    env = Environment(loader=FileSystemLoader("templates"))
    template = env.get_template("prometheus.yaml.jinja")
    output = template.render(
        number_of_pairs=int(default_config["NUMBER_OF_PAIRS"]),
        metrics_host=default_config["METRICS_HOST"],
        metrics_base_port=int(default_config["METRICS_BASE_PORT"]),
        prometheus_port=default_config["PROMETHEUS_PORT"],
        number_of_clients=int((default_config["NUMBER_OF_CLIENTS"])),
    )
    with open("./src/prometheus/prometheus.yml", "w") as f:
        f.write(output)

if __name__ == '__main__':
    if len(sys.argv) != 1:
        print("Usage: python3 compose-generator.py")
        sys.exit(1)

    main()