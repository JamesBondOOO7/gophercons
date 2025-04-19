### Why Two Separate Multiplexers Could be Better:
1. Isolation of Concerns
* The `/devices` endpoint provides business logic.
* The `/metrics` endpoint is strictly infrastructure-related.
* Keeping them separate means infrastructure issues (monitoring, metrics scraping) won't accidentally interfere with your main API logic.

2. Security and Firewalling
* You can expose the metrics port (8081) only internally to monitoring systems (like Prometheus) and not publicly accessible.
* The API (8080) can have its security rules independently, offering clearer boundaries.

3. Simplified Logging and Debugging
* Logs and errors from metrics scraping and from business logic endpoints are cleanly separated.
* Troubleshooting performance issues becomes simpler, as you see clearly which port or endpoint is having an issue.

### OTHER IMP INSIGHTS
* `select {}` is a clean, idiomatic way to signal clearly that the program should run indefinitely.
  * If used in `main()`, then it only stops when main goroutine is interrupted
  * If the `main goroutine` exits, every other goroutine instantly stops, regardless of what they’re doing.
  * You can also gracefully handle such cases in the `main()` function itself.
* `init()` in Golang
  * Must Read : https://www.digitalocean.com/community/tutorials/understanding-init-in-go
* Docker mappings
  * Eg: `ports:<host_port>:<container_port>`
  * Any external service knows your machine's port
  * But your machine knows the container's port
* Prometheus `scrape_interval`
  * Data aggregation time limit
  * Eg: 10 sec -> data points aggregated at 10 sec interval

### Commands
Add a device
```
curl -d '{"id": 3, "mac": "96-40-D1-32-D7-1A", "firmware": "3.03.00"}' localhost:8080/devices
```
Get the list of device
```
curl localhost:8080/devices 
```
Get/Hit the metrics endpoint (which Prometheus is polling)
```
curl localhost:8081/metrics
```