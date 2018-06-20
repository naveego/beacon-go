# Generator for Client for Beacon API

> see https://aka.ms/autorest

To regenerate client, run `autorest` in this directory, or `mage GenerateBeaconClient` in the root of the repo.

If you don't have autorest, run `npm install -g autorest` first.

## Generation Config:

```yaml
input-file: http://localhost:9005/api/swagger.json # full Unicode support

go:
  output-folder: .  
  add-credentials: true
  namespace: beacon

directive:
  from: swagger-document # do it globally (in case there are multiple input OpenAPI definitions)
  where: $.paths..pattern
  transform: return ".*"
  reason: Client side regex pattern interpretation is buggy.
```
