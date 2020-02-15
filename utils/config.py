import yaml
import os


def load_config(filename='config.yml'):
    if os.path.isfile(filename):
        try:
            with open(filename, 'r') as f:
                config = yaml.safe_load(f)
        except (yaml.YAMLError, IOError):
            raise Exception(f'{filename} either not readable or contains malformed yaml')
    else:
        raise Exception(f'{filename} does not exists')

    return config
