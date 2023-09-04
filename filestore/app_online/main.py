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
            "is_depositary_req_fast_netspeed": config[tag+"depository_performance"],
        }
        json_data = json.dumps(dict_data)
        s_addr = "http://" + s_conf["ip"]
        interface = s_addr+s_conf['handler']
        print(time.ctime(time.time()), "The APP info is sent to the scheduler's HTTP interface:", interface)
        
        # ret = requests.post(s_addr+s_conf['handler'], json_data)
        ret = 1
        if ret:
            print(time.ctime(time.time()), "Succeed in APP online with info:", json_data)