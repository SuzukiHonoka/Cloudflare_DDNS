package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/cloudflare/cloudflare-go"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	_ "strings"
	"time"
)

const (
	// config name
	config = "conf.json"
	// getIPApi api
	getIPApi = "https://ip.starx.ink/"
)

func loadConf() *Credential {
	flagConfigPath := flag.String("conf", ".", "Config file path")
	flag.Parse()
	// parse params
	var confPath string
	if *flagConfigPath == "." {
		// current dir
		wd, _ := os.Getwd()
		confPath = path.Join(wd, config)
	} else {
		// specific dir
		confPath = *flagConfigPath
	}
	// load config
	var conf *Credential
	if _, err := os.Stat(confPath); err == nil {
		// file exist
		fd, _ := os.ReadFile(confPath)
		err = json.Unmarshal(fd, &conf)
		check(err)
	} else {
		// ask now
		conf = &Credential{
			Key:      getInput("Cloudflare APK KEY"),
			TokenKey: getInputBool("Are you using **TOKEN** instead of **GLOBAL** KEY"),
			Email:    getInput("Cloudflare Account Email"),
			Domain:   getInput("Domain"),
			Target:   getInput("Target"),
			IPV6:     getInputBool("prefer IPV6"),
		}
		// save
		confJson, _ := json.Marshal(conf)
		err = os.WriteFile(confPath, confJson, os.ModePerm)
		check(err)
	}
	return conf
}

func update(conf *Credential) {
	log.Println("Started at", time.Now().Format(time.RFC1123))
	// mix full name
	fullName := conf.Target + "." + conf.Domain
	// Get Target A record
	currentRecord, err := net.LookupIP(fullName)
	check(err)
	// if target not bounded with one host
	if len(currentRecord) == 0 {
		panic("dns record not found, you should add one first")
	}
	// get outgoing IP
	var dialer net.Dialer
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if conf.IPV6 {
			// tcp6 only but seems cloudflare api doesn't support it yet
			return dialer.DialContext(ctx, "tcp6", addr)
		}
		return dialer.DialContext(ctx, "tcp4", addr)
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30, // sometimes it stuck by jam
	}
	resp, err := client.Get(getIPApi)
	check(err)
	if resp.StatusCode != 200 {
		panic("IP API response code error")
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	check(err)
	outgoingIP := string(body)
	// check if needed to update
	if currentRecord[0].String() == outgoingIP {
		log.Println("Record equivalent, noting to do.")
		return
	}
	// update process
	log.Printf("DNS record does not match the outgoing IP\nRecord: %s Current: %s\n", currentRecord[0].String(), outgoingIP)
	// create a api instance
	var api *cloudflare.API
	if conf.TokenKey {
		api, err = cloudflare.NewWithAPIToken(conf.Key)
	} else {
		// global key
		api, err = cloudflare.New(conf.Key, conf.Email)
	}
	check(err)
	// get domain zoneId
	id, err := api.ZoneIDByName(conf.Domain)
	check(err)
	// create a context for background
	ctx := context.Background()
	// find target rid
	dnsRecords, _ := api.DNSRecords(ctx, id, cloudflare.DNSRecord{})
	var rid string
	for _, el := range dnsRecords {
		if el.Name == fullName {
			rid = el.ID
			break
		}
	}
	if len(rid) == 0 {
		panic("rid not found")
		// TODO: create record if now found
	}
	// final update
	err = api.UpdateDNSRecord(ctx, id, rid, cloudflare.DNSRecord{Content: outgoingIP})
	check(err)
	log.Println("Record updated")
}

func main() {
	conf := loadConf()
	update(conf)
}
