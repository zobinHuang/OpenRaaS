import requests
import json
import time
from config_reader import read_config

if __name__ == "__main__":
    # 1. read conf
    config = read_config("config.yaml")
    s_conf = read_config("scheduler_config.yaml")

    # 2. send notificatioin to scheduler
    for i in range(config['app_num']):
        print(time.ctime(time.time()), "Start to update the No.", i+1, "(of", config['app_num'], ") APP.")
        tag = "app"+str(i+1)+"_"
        dict_data = {
            "application_id": config[tag+"id"],
            "application_name": config[tag+"name"],
            "application_path": config[tag+"path"],
            "applictaion_file": config[tag+"file"],
            "hwkey": config[tag+"hwkey"],
            "operating_system": config[tag+"os"],
            "create_user": config[tag+"create_user"],
            "description": config[tag+"description"],
            "new_file_store_id": config["fs_id"],
            "is_provider_req_gpu": config[tag+"provider_performance"],
            "is_filestore_req_fast_netspeed": config[tag+"filestore_performance"],
            "is_depository_req_fast_netspeed": config[tag+"depository_performance"],
        }
        json_data = json.dumps(dict_data)
        s_addr = "http://" + s_conf["ip"] + ":" + str(s_conf["port"])
        interface = s_addr+s_conf['handler']
        print(time.ctime(time.time()), "APP info is sent to the scheduler's HTTP interface:", interface)
        print(time.ctime(time.time()), "APP online with info:", json_data)
        
        ret = requests.post(interface, json_data)
        print(time.ctime(time.time()), "Get answer from the scheduler:", ret)