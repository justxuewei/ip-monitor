# IP Monitor

IP Monitor is an application running on Linux to monitor your IPs. If your IPs
are changed, the application will push a message to let you known via
[ServerChan](https://sct.ftqq.com/).

The application runs with crontab to make all functionalities work. Here's a
typical example that checks IPs in every 5 minutes.

```bash
$ go build -o ipmonitor cmd/main.go
$ sudo mv ipmonitor /usr/local/bin
$ crontab -e
# insert a new line:
# */5 * * * * ipmonitor --key {sendkey} --name {server name (optional)}
```
