import requests
import threading
import json
import time
from config_reader import read_config


if __name__ == "__main__":
    # 1. read conf
    config = read_config("config.yaml")
    s_conf = read_config("scheduler_config.yaml")

    # 2. send notificatioin to scheduler
    dict_data = {}
    json_data = json.dumps(dict_data)
    s_addr = "http://" + s_conf['ip']
    
    # ret = requests.post(s_addr+s_conf['handler'], json_data)