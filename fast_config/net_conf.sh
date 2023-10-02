# 网卡接口
interface="enp6s18"
# interface="docker0"

# 限制
bw="1000mbit"
l="1ms"

# 清除现有的队列规则
tc qdisc del dev "$interface" root

# 新增队列规则
tc qdisc add dev "$interface" root handle 1: htb default 10
tc class add dev "$interface" parent 1: classid 1:10 htb rate "$bw" burst 15k
tc qdisc add dev "$interface" parent 1:10 handle 10: netem delay "$l"

echo "已成功限制 $interface 网卡接口的出口带宽为 $bw, 延迟为 $l"

tc -s qdisc show dev "$interface"