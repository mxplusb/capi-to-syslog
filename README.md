## CAPI to Syslog

[![Build Status](https://travis-ci.org/mxplusb/capi-to-syslog.svg)](https://travis-ci.org/mxplusb/capi-to-syslog)

Forwards configured events to the downstream syslog via an application instance.

### Configure

Edit `manifest.yml`, setting appropriate values for `env` variables.

### Deploy

```bash
cf push
```

### (Optional) Build binary

```bash
gox
for i in $(ls capi-to-syslog_*); do
  chmod a+x $i
  tar -cvzf $i.tar.gz $i
  rm $i
done
```
