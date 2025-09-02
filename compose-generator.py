import sys
from jinja2 import Environment, FileSystemLoader # pyright: ignore[reportMissingImports]
import configparser

def main():
    config = configparser.ConfigParser()
    config.read("config.ini")
    default_config = config["DEFAULT"]
    env = Environment(loader=FileSystemLoader("templates"))
    template = env.get_template("docker-compose-dev.yaml.jinja")
    output = template.render(
        number_of_pairs=int(default_config["NUMBER_OF_PAIRS"]),
        entries_per_k_bucket=int(default_config["ENTRIES_PER_K_BUCKET"]),
    )
    with open("docker-compose-dev.yaml", "w") as f:
        f.write(output)

if __name__ == '__main__':
    if len(sys.argv) != 1:
        print("Usage: python3 compose-generator.py")
        sys.exit(1)

    main()