# Cloudflare DDNS

#### A DDNS TOOL MADE BY STARX

## What is it?

It's a tool for update **DNS RECORD** of the target domain when the **RECORD IS MISMATCH WITH YOUR OUTGOING IP.**

## How does it work?

It uses the cloudflare official API to process the requests.

## REQUIREMENT

You are required to get your cloudflare account's global **API KEY** for auth process.  
And you should setup a DNS record **(it can be any value but A record only)** before you using this program.  
Here's a list you need to **know**.

- API KEY
- ACCOUNT EMAIL
- DOMAIN
- TARGET

## GET STARTED

The program default config path were set to the current working dir.  
If you want to specify a path then you need to place it to arg1 to start the program.  
If not, the program'll ask you for the detail then save it to the default path or the path of **ARG1**.

## SOME TIPS

You may setup a cron job for this program.