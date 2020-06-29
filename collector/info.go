package collector

import (
	"github.com/shirou/gopsutil/host"
	"log"
	"net"
)

type Host struct {
	HostName string
	IP  string
}

func GetHost() *Host{
	hostinfo := Host{}
	info , err := host.Info()
	if err != nil {
		log.Printf("%v\n",err)
		return nil
	}
	addrs ,err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("%v\n",err)
		return nil
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//多网卡  多ip那么就ip拼接
				if hostinfo.IP == "" {
					hostinfo.IP = ipnet.IP.String()
				}else {
					hostinfo.IP += ","+ipnet.IP.String()
				}
			}
		}
	}
	hostinfo.HostName = info.Hostname
	return  &hostinfo
}
