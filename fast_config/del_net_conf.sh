# 网卡接口
interface="enp6s18"

# 清除现有的队列规则
tc qdisc del dev "$interface" root

tc -s qdisc show dev "$interface"