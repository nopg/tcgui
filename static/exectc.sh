#!/bin/bash
# Taken from: http://www.uebi.net/howtos/rpiwanem.htm
latency=$1
loss=$2
jitter=$3
bandwidth=$4
echo "<html><head><title>Submitted</title></head><body><br>"
latency2=$(( $latency / 2 ))
loss2=$(( $loss / 2 ))
sudo tc qdisc del dev eth1 root
sudo tc qdisc del dev eth2 root
sudo tc qdisc add dev eth1 root handle 1:0 tbf rate ${bandwidth}kbit burst ${bandwidth}K latency 5000ms
sudo tc qdisc add dev eth2 root handle 2:0 tbf rate ${bandwidth}kbit burst ${bandwidth}K latency 5000ms
sudo tc qdisc add dev eth1 parent 1:1 handle 10: netem delay ${latency2}ms ${jitter}ms loss ${loss2}
sudo tc qdisc add dev eth2 parent 2:1 handle 10: netem delay ${latency2}ms ${jitter}ms loss ${loss2}
echo "<font face=arial><b>Result:</b><br><br>"
echo "Latency should now be <b>+${latency}ms</b><br>"
echo "Jitter should now be <b>${jitter}ms</b><br>"
echo "Bandwidth should now be <b>${bandwidth}kbit</b><br>"
echo "Packet loss should now be <b>${loss}%</b><br>"
echo "<br><br><br><a href=/>Back</a></font></body></html>"