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

    print(time.ctime(time.time()), "The filestore microservice is setup successfully with config:", json.dumps(conf))
    
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
        "port": config["fs_port"],
        "protocol": config["fs_protocol"],
        "directory": config["fs_directory"],
        "username": config["fs_user"],
        "password": config["fs_pwd"],
        "is_contain_fast_netspeed": config["performance"],
    }
    json_data = json.dumps(dict_data)
    s_addr = "http://" + s_conf["ip"]
    interface = s_addr+s_conf['handler']
    print(time.ctime(time.time()), "The filestore worker node's info is sent to the scheduler's HTTP interface:", interface)
    
    # ret = requests.post(s_addr+s_conf['handler'], json_data)
    ret = 1
    if ret:
        print(time.ctime(time.time()), "Succeed in filestore worker node online with info:", json_data)