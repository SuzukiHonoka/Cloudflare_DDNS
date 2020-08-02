package main

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path"
)

const (
	config = "conf.json"
	getIP  = "http://ip.03k.org/"
	fatal  = 0
	normal = 1
	warn   = 2
)

type API struct {
	Key    string
	Email  string
	Domain string
	Target string
}

var (
	cwd      string
	confPath string
	conf     API
	ogIP     string
)

func dealE(err error, mode int) {
	if err != nil {
		switch mode {
		case fatal:
			color.HiRed("Fatal Error:", err)
			os.Exit(1)
		case normal:
			color.HiRed("Error:", err)
		case warn:
			color.HiYellow("Warn:", err)
		}
	}
}

func ask(qua string) string {
	color.Cyan("Please Enter", qua, ":")
	reader := bufio.NewReader(os.Stdin)
	rp, err := reader.ReadString('\n')
	dealE(err, fatal)
	return strings.Trim(strings.Trim(rp, "\r\n"), "\n")
}

func main() {
	start := time.Now()
	color.HiRed("Started at " + time.Now().String())
	// FILL API
	cwd, _ = os.Getwd()
	confPath = path.Join(cwd, config)
	if _, err := os.Stat(confPath); err == nil {
		fd, _ := ioutil.ReadFile(confPath)
		err := json.Unmarshal(fd, &conf)
		dealE(err, fatal)
		color.Green("READ CONFIG FILE SUCCEED.")
	} else {
		conf = API{
			Key:    ask("APK KEY"),
			Email:  ask("Email"),
			Domain: ask("Domain"),
			Target: ask("Target"),
		}
		confJson, _ := json.Marshal(conf)
		err = ioutil.WriteFile(confPath, confJson, os.ModePerm)
	}
	// mix full name
	fullName := conf.Target + "." + conf.Domain
	// Get Target A record
	targetA, err := net.LookupIP(fullName)
	dealE(err, fatal)
	color.Blue("DNS LOOKUP SUCCEED.")
	// if target not bounded with one host
	first := targetA[0]
	for _, ip := range targetA {
		if ip.String() != first.String() {
			color.HiRed("Fatal Error: Target has more than one host record.")
			os.Exit(1)
		}
	}
	// Get outgoing IP
	resp, err := http.Get(getIP)
	dealE(err, fatal)
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		dealE(err, fatal)
		ogIP = string(body)
	} else {
		color.HiRed("Fatal Error: IP API response error.")
		os.Exit(1)
	}
	color.Yellow("GET OUTGOING IP SUCCEED.")
	// create a api instance
	api, err := cloudflare.New(conf.Key, conf.Email)
	dealE(err, fatal)
	color.Green("GET API INSTANCE SUCCEED!!")
	// get domain zoneid
	id, err := api.ZoneIDByName(conf.Domain)
	dealE(err, fatal)
	color.Green("Get ZONEID SUCCEED!!")
	// get target rid
	rlist, _ := api.DNSRecords(id, cloudflare.DNSRecord{})
	var rid string
	var cIP string
	for _, el := range rlist {
		if el.Name == fullName {
			rid = el.ID
			cIP = el.Content
			color.Magenta("GET RID SUCCEED!")
			break
		}
	}
	if len(rid) == 0 {
		// Create record TODO
	}

	if cIP != ogIP {
		color.HiYellow("CURRENT DNS RECORD DOES NOT MATCH THE OUTGOING IP: " + cIP)
		err := api.UpdateDNSRecord(id, rid, cloudflare.DNSRecord{Content: ogIP})
		if err == nil {
			color.Green("UPDATE SUCCEED!")
		} else {
			dealE(err, fatal)
		}

	} else {
		color.HiWhite("CURRENT DNS RECORD MATCH THE OUTGOING IP: " + ogIP)
	}
	color.HiRed("PROCESS ENDED. TOTAL COST: " + time.Now().Sub(start).String() + "ms")
}
