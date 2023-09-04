import yaml

def read_config(path):
    """
    Return python dict from .yml file.
    
    Usage: config = read_config('config.yml')

    Args:
        path (str): path to the .yml config.

    Returns (dict): configuration object.
    """
    with open(path, 'r') as ymlfile:
        cfg = yaml.safe_load(ymlfile)
    return cfg