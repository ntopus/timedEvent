# timedQueue Service
A Service to support timed queue lib. Users should start this service before using timed queue lib, and both should share the same queue and database.

## Run development

Add `timedEvent.db.ivanmeca.com.br` to `/etc/hosts`.

Run command:
```bash
make run-dev
```