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
    dict_data = {"id": config["id"],
        "type": "depository",
        "ip": config["ip"],
        "port": config["reg_port"],
        "tag": "latest",
        "is_contain_fast_netspeed": config["performance"],
    }
    json_data = json.dumps(dict_data)
    s_addr = "http://" + s_conf['ip']
    interface = s_addr+s_conf['handler']
    print(time.ctime(time.time()), "The filestore worker node's info is sent to the scheduler's HTTP interface:", interface)
    
    ret = requests.post(s_addr+s_conf['handler'], json_data)
    # ret = 1
    if ret:
        print(time.ctime(time.time()), "Succeed in filestore worker node online with info:", json_data)