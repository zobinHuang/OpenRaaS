## 1. Frontend

```
cd web/ant-client-page
npm install
npm start
```

## 2. Nginx

Modify Nginx configuration file at <code>nginx/conf.d/static.conf</code>:

```nginx
upstream backstage{
    server [Ip and port of Global Master];
}

upstream web{
    server [Ip and port of Web frontend server];
}

server{
    listen 80;
    server_name [Server Domain];

    ...
```

Run OpenRaaS Nginx proxy:

```bash
sudo docker run --name nginx \
    --restart always \
    -p 80:80 \
    -v /home/broscloud/Code/OpenRaaS/nginx/logs:/var/log/nginx \
    -v /home/broscloud/Code/OpenRaaS/nginx/nginx.conf:/etc/nginx/nginx.conf \
    -v /home/broscloud/Code/OpenRaaS/nginx/conf.d:/etc/nginx/conf.d \
    -d nginx
```

## 3. Backend

```bash
cd backstage
sudo docker-compose up
```

## 4. Provider

注意，要先去修改 `/provider/run-streamer.sh` 和 `/provider/stop-streamer.sh` 里的相对目录 `cd ../../provider/`!
要站在当前目录是 `/winecontainer/serverd/` 的角度，去找到上述两个 sh 文件所在目录的相对路径。

如果有 `serverd` 文件：

```
cd winecontainer/serverd
sudo ./serverd
```

如果需要重新变异 `serverd` 文件：
1. 如果有，就先删掉 `go.mod` 和 `go.sum`
2. 运行 `go mod init` 和 `go mod tidy`
3. 运行 `go build`
4. 运行 `sudo ./serverd` 二进制文件