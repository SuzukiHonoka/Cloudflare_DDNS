package wrapper

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"net"
)

type Wrapper struct {
	Name     string
	ZoneID   string
	RecordID string
	API      *cloudflare.API
}

func New(credential *Credential) (*Wrapper, error) {
	var err error
	var api *cloudflare.API
	if credential.IsGlobal {
		api, err = cloudflare.New(credential.Key, credential.Email)
	} else {
		api, err = cloudflare.NewWithAPIToken(credential.Key)
	}
	return &Wrapper{
		API: api,
	}, err
}

func (w *Wrapper) Init(t *Target) (err error) {
	w.Name = t.String()
	w.ZoneID, err = w.API.ZoneIDByName(t.Domain)
	if err != nil {
		return err
	}
	err = w.setRecordID()
	return err
}

func (w *Wrapper) setRecordID() error {
	ds, err := w.API.DNSRecords(context.Background(), w.ZoneID, cloudflare.DNSRecord{
		Name: w.Name,
	})
	if err != nil {
		return err
	}
	count := len(ds)
	if count == 0 {
		return fmt.Errorf("record: %s does not exist", w.Name)
	} else if count > 1 {
		return fmt.Errorf("record: %s has too many records, count: %d", w.Name, count)
	}
	w.RecordID = ds[0].ID
	return nil
}

func (w *Wrapper) CreateRecord(ip net.IP) (string, error) {
	recordType := "A"
	if ip.To4() == nil {
		recordType = "AAAA"
	}
	resp, err := w.API.CreateDNSRecord(context.Background(), w.ZoneID, cloudflare.DNSRecord{
		Type:    recordType,
		Name:    w.Name,
		Content: ip.String(),
		TTL:     60,
	})
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("create [%s] record: [%s] with content: [%s] failed", recordType, w.Name, ip)
	}
	return resp.Result.ID, nil
}

func (w *Wrapper) GetRecordContent() (string, error) {
	ds, err := w.API.DNSRecord(context.Background(), w.ZoneID, w.RecordID)
	if err != nil {
		return "", err
	}
	return ds.Content, nil
}

func (w *Wrapper) UpdateRecord(ip net.IP) error {
	return w.API.UpdateDNSRecord(context.Background(), w.ZoneID, w.RecordID, cloudflare.DNSRecord{
		Content: ip.String(),
	})
}
