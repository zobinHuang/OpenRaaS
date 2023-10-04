import requests
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
        "port": str(config["reg_port"]),
        "tag": "latest",
        "is_contain_fast_netspeed": config["performance"],
        "mem": config["mem"],
        "bandwidth": config["bw"],
        "latency": config["latency"]
    }
    json_data = json.dumps(dict_data)
    trimmed_json_data = json_data[:-1]
    trimmed_json_data += ', "inst_history": ' + config["inst_history"] + '}'
    json_data = trimmed_json_data
     
    s_addr = "http://" + s_conf["ip"] + ":" + str(s_conf["port"])
    interface = s_addr + s_conf["handler"]
    headers = {
        "type": "depository",
    }
    print(time.ctime(time.time()), "Depository worker node's info is sent to the scheduler's HTTP interface:", interface)
    print(time.ctime(time.time()), "Depository worker node online with info:", json_data)
    
    ret = requests.post(interface, params = headers, data = json_data)
    print(time.ctime(time.time()), "Get answer from the scheduler:", ret)