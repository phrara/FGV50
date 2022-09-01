import re
import sys

def getmac(ip):
    obj = re.compile(rf"{ip}\s*(?P<mac>.*?)\s")
    with open("./tools/arp.txt", mode="r", encoding="gbk") as f:
        ret = obj.search(f.read())
        print(ret.group("mac"))
    return ret.group("mac")


def getdev(key):
    obj = re.compile(rf"{key}..........(?P<res>.*)", )
    with open("./python/oui.txt", mode="r", encoding="utf-8") as f:
        ret = obj.search(f.read())
    print(ret.group("res"))
    return ret.group("res")


ip=sys.argv[1]
# ip = "192.168.133.128"
mac = getmac(ip)
iden = mac[:8].upper()
getdev(iden)
