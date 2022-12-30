package entry

import (
	"cf_ddns/internal/config"
	"cf_ddns/internal/dns"
	"cf_ddns/internal/proxy"
	"cf_ddns/utils"
	"log"
	"net"
)

func Handle(conf *config.Config) {
	conf.DNS.SetDefault()
	//fmt.Printf("%+v\n", conf.Target)
	ip, err := conf.API.GetIP()
	utils.CheckErr(err)
	log.Printf("Public IP: %s", ip)
	ok, err := dns.EqualsTo(conf.Target.Name(), ip)
	if err != nil {
		if err.(*net.DNSError).IsNotFound {
			log.Fatalln("record not found, it might be a dns issue or you just created a new record")
		}
		log.Fatalln(err)
	}
	if ok {
		log.Println("Matches")
		return
	}
	api, err := proxy.New(conf.Creds)
	utils.CheckErr(err)
	err = api.Init(conf.Target, ip)
	utils.CheckErr(err)
	err = api.UpdateRecord(ip)
	utils.CheckErr(err)
	log.Println("Record Updated")
}
