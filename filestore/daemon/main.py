import requests
import threading
import json
import time
from cheroot import wsgi
from wsgidav.wsgidav_app import WsgiDAVApp
from config_reader import read_config


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


if __name__ == "__main__":
    # 1. read conf
    config = read_config("config.yaml")
    s_conf = read_config("scheduler_config.yaml")

    # 2. start webdav server
    fs_conf = {}
    fs_conf["host"] = "0.0.0.0"
    fs_conf["port"] = config["fs_port"]
    fs_conf["provider_mapping"] = {
        config["fs_directory"]: "./storage",
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


    # 3. send notificatioin to scheduler
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
        "type": "filestore",
    }
    print(time.ctime(time.time()), "Filestore worker node's info is sent to the scheduler's HTTP interface:", interface)
    print(time.ctime(time.time()), "Filestore worker node online with info:", json_data)
    
    ret = requests.post(interface, params = headers, data = json_data)
    print(time.ctime(time.time()), "Get answer from the scheduler:", ret)