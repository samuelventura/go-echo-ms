# go-echo-ms

```bash
export ECHO_LOGS=/var/log
export ECHO_ENDPOINT=0.0.0.0:31653
go run .
go install && go-echo-ms
lsof -p PID
curl --unix-socket /tmp/go-echo-ms.state http://localhost/node/echo/
```
