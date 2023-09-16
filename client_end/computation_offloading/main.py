import requests
import threading
import json
from datetime import datetime
import time
from cheroot import wsgi
from wsgidav.wsgidav_app import WsgiDAVApp
import yaml


def start_fs_server(conf):
    print(time.ctime(time.time()), "Starting file-sharing setup...",)
    app = WsgiDAVApp(conf)

    server_args = {
        "bind_addr": (conf["host"], conf["port"]),
        "wsgi_app": app,
    }
    server = wsgi.Server(**server_args)

    print(time.ctime(time.time()), "Filestore microservice is setup successfully with config:", json.dumps(conf))
    
    try:
        server.start()
    except KeyboardInterrupt:
        print("Received Ctrl-C: stopping...")
    finally:
        server.stop()

def init_all(config):
    # 1. start webdav server
    fs_conf = {}
    fs_conf["host"] = "0.0.0.0"
    fs_conf["port"] = config["fs_port"]
    fs_conf["provider_mapping"] = {
        config["fs_directory"]: "./workspace",
    }
    fs_conf["http_authenticator"] = {
        "trusted_auth_header": None,
        "domain_controller": None,
        "accept_basic": True,  # Pass false to prevent sending clear text passwords
        "accept_digest": True,
        "default_to_digest": True,
    }
    fs_conf["simple_dc"] = {
        "user_mapping": {
            "*": {
                config["fs_user"]: {
                    "password": str(config["fs_pwd"]),
                    "roles": ["editor", "admin"],
                }
            },
            "/pub": True
        },
    }
    fs_conf["verbose"] = 1

    t = threading.Thread(target=start_fs_server, args=(fs_conf,))
    t.start()


    # 2. filestore online
    dict_data = {
        "id": config["id"],
        "ip": config["ip"],
        "port": str(config["fs_port"]),
        "protocol": config["fs_protocol"],
        "directory": config["fs_directory"],
        "username": config["fs_user"],
        "password": str(config["fs_pwd"]),
        "is_contain_fast_netspeed": config["performance"],
        "mem": config["mem"],
    }
    json_data = json.dumps(dict_data)
    s_addr = "http://" + config["scheduler_ip"] + ":" + str(config["scheduler_port"])
    interface = s_addr + config["handler_filestore_online"]
    headers = {
        "type": "filestore",
    }
    
    requests.post(interface, params = headers, data = json_data)
    
    # 3. app online
    
    for i in range(config['app_num']):
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
            "new_file_store_id": config["id"],
            "is_provider_req_gpu": config[tag+"provider_performance"],
            "is_filestore_req_fast_netspeed": config[tag+"filestore_performance"],
            "is_depository_req_fast_netspeed": config[tag+"depository_performance"],
            "image_name": config[tag+"image_name"],
        }
        json_data = json.dumps(dict_data)
        interface = s_addr + config["handler_app_online"]
        
        requests.post(interface, json_data)

if __name__ == "__main__":
    # 1. read conf
    with open("config.yaml", 'r') as ymlfile:
        config = yaml.safe_load(ymlfile)

    # 2. init working environment & register
    init_all(config)
    
    # 3. get microservice schedule from the scheduler node
    s_addr = "http://" + config["scheduler_ip"] + ":" + str(config["scheduler_port"])
    interface = s_addr + config["handler_computation_offloading"]
    headers = {
        "app_id": config["app1_id"],
    }
    
    response = requests.get(interface, params = headers)
    info = json.loads(response.json()["info"])
    # print(f"Get json from {interface}: {response.status_code}, {info}")
    
    depository = "DockerHub"
    
    # 4. start container on provider
    provider = info["provider_core"]
    p_addr = "http://" + provider["ip"] + ":3080"
    interface = p_addr + "/api/daemon/createinstance"
    
    print(f"Selected provider: {provider['ip']}, depository: {depository}")
    
    tag = "app1_"
    dict_data = {
        "run_in_linux": True,
        "application_name": config[tag+"name"],
        "application_path": config[tag+"path"],
        "application_file": config[tag+"file"],
        "hwkey": config[tag+"hwkey"],
        "image_name": config[tag+"image_name"],
        "filestore_list": info["filestore_list"],
        # "app_option": "detector train ./data ./yolov4-tiny.cfg ; mv ../yolov4-tiny_last.weights ./backup/yolov4-tiny_last.weights"
    }
    json_data = json.dumps(dict_data)
    print(f"Post to {interface} with config: {json_data}")
    headers = {'Content-Type': 'application/json'} 
    response = requests.post(interface, data=json_data, headers=headers)
    vmid = response.json()["vmid"]
    print(response, response.text)
    
    # 5. output
    
    interface = p_addr + "/api/daemon/checkinstancebyvmid"
    headers = {'vmid': vmid} 
    
    backup_folder = "workspace/backup"
    
    while True:
        time.sleep(5)
        response = requests.get(interface, params = headers)
        
        try:
            if response.status_code == 200:
                if 'log' in response.json():
                    print(f"{datetime.now()} ---- 正在进行计算卸载，最近几条日志内容如下: \n {response.json()['log']}")
            else:
                print("训练已完毕，请在 'workspace/backup' 文件夹中查看下载好的神经网络梯度文件")
                break
        except:
            print("训练已完毕，请在 'workspace/backup' 文件夹中查看下载好的神经网络梯度文件")
            break
        
        
    