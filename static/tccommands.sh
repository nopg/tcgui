#!/bin/bash
# Initial idea by disk91.com (THX!) with some minor changes by uebi.net
echo "Content-type: text/html"
echo ""
# read
latency=$1
loss=$2
jitter=$3
bandwidth=$4
echo "<html><head><title>Submitted</title></head><body><br>"
# if [ -z "$latency" ] || [ -z "$bw" ] || [ -z "$var" ] || [ -z "$loss" ] ; then
# echo "<font face=arial><b>Result:</b><br><br>Erorr! Go back and try again!<br/><br/>"
# echo "<br><br><br><a href=/>Back</a></font>"
# else
latency2=$(( $latency / 2 ))
loss2=$(( $loss / 2 ))
echo $1
echo $2
echo $3
echo $4
sudo tc qdisc del dev eth1 root
sudo tc qdisc del dev eth2 root
sudo tc qdisc add dev eth1 root handle 1:0 tbf rate ${bandwidth}kbit burst ${bandwidth}K latency 5000ms
sudo tc qdisc add dev eth2 root handle 2:0 tbf rate ${bandwidth}kbit burst ${bandwidth}K latency 5000ms
sudo tc qdisc add dev eth1 parent 1:1 handle 10: netem delay ${latency2}ms ${jitter}ms loss ${loss2}
sudo tc qdisc add dev eth2 parent 2:1 handle 10: netem delay ${latency2}ms ${jitter}ms loss ${loss2}
echo "<font face=arial><b>Result:</b><br><br>"
echo "Latency should now be <b>+${latency}ms</b><br>"
echo "Jitter should now be <b>${var}ms</b><br>"
echo "Bandwidth should now be <b>${bw}kbit</b><br>"
echo "Packet loss should now be <b>${loss}%</b><br>"
# # tc qdisc | tr "\n" "#" | sed -e "s/#/<br\/>/g"
echo "<br><br><br><a href=/>Back</a></font></body></html>"
# fi