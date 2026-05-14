# IP Monitor

IP Monitor is an application running on Linux to monitor your IP changes. It is
a DDNS-like way. If your IPs are changed, the application will push a message to
let you know via a webhook.

The application runs with crontab to make all functionalities work. Here's a
typical example that checks IPs in every 5 minutes.

```bash
$ make build
$ ipmonitor help
$ ipmonitor version
$ sudo mv ipmonitor /usr/local/bin
$ crontab -e
# insert a new line:
# */5 * * * * ipmonitor --webhook-url='https://example.com/api/messages?token=TOKEN&message={message}' --name={server name (optional)} --heartbeat=true --devices={links (optional, e.g. "lo,enp0")}
```
