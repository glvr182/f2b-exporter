# F2B-Exporter

[![GoDoc](https://godoc.org/github.com/glvr182/f2b-exporter?status.svg)](https://godoc.org/github.com/glvr182/f2b-exporter)
[![GitHub tag](https://img.shields.io/github/tag/glvr182/f2b-exporter.svg)]()
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/glvr182/f2b-exporter)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/glvr182/f2b-exporter)

This is a simple Fail2Ban prometheus exporter

## Installation

### From source
You can clone this repository from git `https://github.com/glvr182/f2b-exporter.git`.  
Then all you have to do is run `go build` and you're done!

### Docker
Using the following command you can run this program with docker.  
NOTE: When running the docker image you might want to mount the certificates since some remotes use TLS.
```
docker run \
-d \
--name f2b-exporter \
-v /var/lib/fail2ban/fail2ban.sqlite3:/var/lib/fail2ban/fail2ban.sqlite3 \
-v /etc/ssl/certs:/etc/ssl/certs \
glvr182/f2b-exporter
```
Depending on your deployment you might want to expose the configured port (default 8080) like this:
```
docker run \
-d \
-p 8080:8080 \
--name f2b-exporter \
-v /var/lib/fail2ban/fail2ban.sqlite3:/var/lib/fail2ban/fail2ban.sqlite3 \
-v /etc/ssl/certs:/etc/ssl/certs \
glvr182/f2b-exporter
```

## Configuration
The exporter has a few settings that you can tweak using the cli or using env variables:
| cli        | env          | default                            |
|------------|--------------|------------------------------------|
| --port     | F2B_PORT     | 8080                               |
| --database | F2B_DATABASE | /var/lib/fail2ban/fail2ban.sqlite3 |
| --remote   | F2B_REMOTE   | freeGeoIP                          |

To add the exporter to prometheus a simple config like this would do the trick:  
NOTE: this is from prometheus, NOT for this exporter.
```
global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'f2b-exporter'
    scrape_interval: 5s
    static_configs:
    - targets: ['localhost:8080']
```
Or using docker and a dedicated network:
```
global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'f2b-exporter'
    scrape_interval: 5s
    static_configs:
    - targets: ['f2b-exporter:8080']
```