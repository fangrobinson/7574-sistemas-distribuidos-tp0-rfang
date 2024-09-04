from copy import deepcopy
import sys
import yaml
from pathlib import Path

class InvalidArgument(Exception):
    ...

BASE_YAML = {
    "name": "tp0",
    "services": {
        "server": {
            "container_name": "server",
            "image": "server:latest",
            "entrypoint": "python3 /main.py",
            "environment": [
                "PYTHONUNBUFFERED=1",
            ],
            "networks": [
                "testing_net"
            ],
            "volumes": [
                "./server/config.ini:/config.ini"
            ]
        },
    },
    "networks": {
        "testing_net": {
            "ipam": {
                "driver": "default",
                "config": [
                    {
                        "subnet": "172.25.125.0/24"
                    }
                ]
            }
        }
    }
}

def validate_and_parse_arguments(path_to_yaml: str, client_n: str):
    try:
        client_n = int(client_n)
        if client_n < 1:
            raise ValueError
    except ValueError:
        raise InvalidArgument("Second argument must be a positive integer.")
    return path_to_yaml, client_n

def main():
    if len(sys.argv) != 3:
        print("Not enough arguments provided.")
    path_to_yaml, client_n = validate_and_parse_arguments(*sys.argv[1:3])
    Path(path_to_yaml).parent.mkdir(parents=True, exist_ok=True)
    yaml_config = deepcopy(BASE_YAML)
    services: dict = yaml_config["services"]
    new_services = {
        f"client{i}": {
            "container_name": f"client{i}",
            "image": "client:latest",
            "entrypoint": "/client",
            "environment": [
                f"CLI_ID={i}",
            ],
            "networks": [
                "testing_net"
            ],
            "depends_on": [
                "server"
            ],
            "volumes": [
                "./client/config.yaml:/config.yaml",
                "./.data/dataset/:/data/"
            ]
        }
        for i in range(1, client_n+1)
    }
    services.update(new_services)
    with open(path_to_yaml, "w") as f:
        yaml.dump(yaml_config, f)

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"[ERROR] {e}")