### Nginx Config

```bash
sudo docker run --name nginx --restart always -p 80:80 -v /home/broscloud/Code/nginx/logs:/var/log/nginx -v /home/broscloud/Code/nginx/nginx.conf:/etc/nginx/nginx.conf -v /home/broscloud/Code/nginx/conf.d:/etc/nginx/conf.d -d nginx
```