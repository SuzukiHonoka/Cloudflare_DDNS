# Cloudflare DDNS
A DDNS tool helps you to update your specific DNS record with your realtime public IP address.

## How does it work?
Firstly, it checks your public IP address by fetching specific API which your already written in your configuration file.
After your public IP has been fetched, it'll make a dns lookup operation for compare if the public IP address matches
your current DNS record. It will update your DNS record according to the result of compare.

## Requirement
You'll have to got your domain added in your Cloudflare Dashboard, then get your Global/Token API key from here for auth process. 
If the specific DNS record you provided does not exist, it will create a new one with your current public IP address,
you may see "record not found" in this case, it's because DNS record has not been synced crossing DNS servers, so just be chill.

## Get Started!
The program default config path were set to the current working dir (use `-c` to specific config path).  
You'll have to pass a config file to the tool before you use it.

## Configuration example
```json
{
  "api": {
    "url": "https://example.ip-api.com",
    "ipv6": false
  },
  "dns": "8.8.8.8",
  "credentials": {
    "global": false,
    "key": "your cloudflare api key",
    "email": "your cloudflare account email, only required if you're using global key"
  },
  "target": {
    "domain": "example.com",
    "sub": "sub.level.name"
  }
}
```

## Tip
You may set up a cronjob for this tool if your public IP address changes frequently.