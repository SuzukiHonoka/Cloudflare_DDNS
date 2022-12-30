package proxy

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"net"
)

type Wrapper struct {
	Name   string
	ZoneID string
	RID    string
	API    *cloudflare.API
}

func New(creds *Credentials) (*Wrapper, error) {
	var err error
	var api *cloudflare.API
	if creds.IsGlobal {
		api, err = cloudflare.New(creds.Key, creds.Email)
	} else {
		api, err = cloudflare.NewWithAPIToken(creds.Key)
	}
	return &Wrapper{
		API: api,
	}, err
}

func (w *Wrapper) Init(t *Target, fallback string) (err error) {
	w.Name = t.Name()
	w.ZoneID, err = w.API.ZoneIDByName(t.Domain)
	if err != nil {
		return err
	}
	if err = w.setRID(); err != nil && fallback != "" {
		log.Println("creating record")
		w.RID, err = w.CreateRecord(fallback)
	}
	return err
}

func (w *Wrapper) setRID() error {
	ds, err := w.API.DNSRecords(context.Background(), w.ZoneID, cloudflare.DNSRecord{
		Name: w.Name,
	})
	if err != nil {
		return err
	}
	count := len(ds)
	switch {
	case count == 0:
		return fmt.Errorf("record: %s does not exist", w.Name)
	case count > 1:
		return fmt.Errorf("record: %s has too many records, count: %d", w.Name, count)
	}
	w.RID = ds[0].ID
	return nil
}

func (w *Wrapper) CreateRecord(content string) (string, error) {
	var RT string
	ip := net.ParseIP(content)
	switch {
	case ip == nil:
		return "", fmt.Errorf("content: [%s] it not a valid", content)
	case ip.To4() == nil:
		RT = "AAAA"
		break
	default:
		RT = "A"
	}
	resp, err := w.API.CreateDNSRecord(context.Background(), w.ZoneID, cloudflare.DNSRecord{
		Type:    RT,
		Name:    w.Name,
		Content: content,
		TTL:     60,
	})
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("create [%s] record: [%s] with content: [%s] failed", RT, w.Name, content)
	}
	return resp.Result.ID, nil
}

func (w *Wrapper) GetRecordContent() (string, error) {
	ds, err := w.API.DNSRecord(context.Background(), w.ZoneID, w.RID)
	if err != nil {
		return "", err
	}
	return ds.Content, nil
}

func (w *Wrapper) UpdateRecord(content string) error {
	return w.API.UpdateDNSRecord(context.Background(), w.ZoneID, w.RID, cloudflare.DNSRecord{
		Content: content,
	})
}
