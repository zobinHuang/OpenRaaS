1. 用 ssh 连接 gpu1 服务器的 broscloud 账号：

```
ssh broscloud@kb109.dynv6.net -p 10901
% 密码: kb109_xusir  
```

可以在电子科大内任何网络连接（流量不行）

2. 单独开一个终端运行前端，执行：

```sh
cd /home/broscloud/Code/OpenRaaS/web/ant-client-page
npm install
npm start
```

3. 单独开一个终端运行后端，执行：

```sh
cd /home/broscloud/Code/OpenRaaS/backstage
sudo docker-compose up
```

4. 单独开一个终端运行 provider，执行：

```sh
cd /home/broscloud/Code/OpenRaaS/provider/serverd
sudo ./serverd
```

注意，上述过程只能在 109 内调试全部功能，出了 109 不能进入游戏（但是可以调试 log 等地方）