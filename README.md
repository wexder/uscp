# USCP Unnamed simple configuration tool

This is extremely simple tool for loading configurations.

Currently it can load configuration from `json` or `yaml`, but can be easily extended.

## Merging
One of important properties of **USCP** is configuration merging.
When you load two configuration files with, we merge them together.
For example
```yaml
test: t
hacker:
  name: Viktor
```
and
```yaml
hacker.name: Test
```
Will result in 
```yaml
test: t
hacker:
  name: Test
```
This then will be serialized into given struct
