## CAPI to Syslog

Forwards configured events to the downstream syslog via an application instance.

To build: `gox; for i in $(ls capi-to-syslog_*); do chmod a+x $i; tar -cvzf $i.tar.gz $i; rm $i; done`

To set up the environment:
```
cf set-env capi-to-syslog GOPACKAGENAME capi-to-syslog # only if CF doesn't have internet access.
cf set-env capi-to-syslog CAPI_CLIENT_ID MyUniqueUaaClient
cf set-env capi-to-syslog CAPI_CLIENT_SECRET MyClientSecret123!
cf set-env capi-to-syslog CAPI_SYSTEM_URI system.internal.domain
cf set-env capi-to-syslog CAPI_EVENTS "type:audit.app.ssh-authorized,type:audit.app.ssh-unauthorized,type:audit.service_key.create,type:audit.service_key.delete,type:audit.space.create,type:audit.app.create,type:audit.app.update"
```

To use with the Go buildpack:
```bash
cf push capi-to-syslog -u none --no-route
```

To use with the Binary Buildpack:
```bash
cf push capi-to-syslog -c "chmod a+x capi-to-syslog && ./capi-to-syslog" --no-route -u none
```

