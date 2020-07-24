package collector

import (
	"bufio"
	"os"
	"strings"
	"log"
)

var (
	Scrapers = map[Scraper]bool{
		CpuLoad{}: true,
		&Rest{}: true,
	}
)



type Domain struct {
	Domaininfo string
	Ip string
}

func init() {
	domains = ReadDomain()
}


func ReadDomain() []Domain{
	domains := []Domain{}
	file,err := os.Open("domain.txt")
	if err != nil {
		log.Printf("文件不存在 %s\n",err)
		return nil
	}
	reader := bufio.NewReader(file)
	for {
		line,_,err := reader.ReadLine()
		if err != nil {
			break
		}
		s := strings.Fields(string(line))
		if len(s) < 2 {
			log.Printf("此域名无效  %s\n",s[0])
			continue
		}
		domain :=  Domain{
			Domaininfo: s[0],
			Ip:         s[1],
		}
		domains = append(domains,domain)
	}
	return domains
}

