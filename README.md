# Highlander

A Mutex for Codeship Pro builds, used to protect deployment steps across multiple builds.

`codeship.env`: (and `jet encrypt codeship.env codeship.env.encrypted`)
```
CODESHIP_ORGANIZATION=foo
CODESHIP_USERNAME=bar
CODESHIP_PASSWORD=baz
```

`codeship-services.xml`:
```yaml
highlander:
  image: ghcr.io/tailsdotcom/highlander:latest
  encrypted_env_file: codeship.env.encrypted
```

`codeship-steps.xml`:
```yaml
# CI Steps
- service: highlander
  command: highlander
  tag: (^staging$|^prod$)
# CD Steps
```
