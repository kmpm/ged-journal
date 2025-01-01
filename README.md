# Go Elite Dangerous Journal Parser




## Usage
### Windows
```powershell

.\ged-journal.exe -l debug | Tee-Object log.jsonl

```




## Nats
```powershell
docker network create nats
docker run -d --name nats --network nats --rm -p 4222:4222 -p 8222:8222 nats --http_port 8222
```