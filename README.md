To build: `gox; for i in $(ls capi-to-syslog_*); do chmod a+x $i; tar -cvzf $i.tar.gz $i; rm $i; done`

```
cf set-env capi-to-syslog CAPI_CLIENT_ID MyUniqueUaaClient
cf set-env capi-to-syslog CAPI_CLIENT_SECRET MyClientSecret123!
cf set-env capi-to-syslog CAPI_SYSTEM_URI system.internal.domain
cf set-env capi-to-syslog CAPI_EVENTS "type:audit.app.ssh-authorized,type:audit.app.ssh-unauthorized,type:audit.service_key.create,type:audit.service_key.delete,type:audit.space.create,type:audit.app.create"
```