# Backlog

- Finish minimal produce / get functionality including volatile memory-backed
  message storage and sell-by-date message eviction housekeeping.
- Make protoc runnable from go generate
- Define default server port in single place
- Configure server URL from env vars (12 factor)
- Turn on TLS
- Add some from of auth (pref JWT-based)
- Redis storage backend instead of volatile.
- Add pub/sub features
- Dockerize as separate microservices: server, redis storage backend, (maybe) 
  auth service.
- Kubernetes deployment to GCP
- Let's encrypt guard on public service end point