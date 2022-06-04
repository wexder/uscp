# USCP Unnamed simple configuration parser

This is extremely simple tool for loading configurations.

Currently it can load configuration from `json` or `yaml`, but can be easily extended.

## Usage

```go
type Conf struct {
	Test string `uscp:"required"`
	Obj  struct {
		K time.Duration
	}
}

func main() {
	u := uscp.New()
	u.SetConfigName("conf")
	u.AddConfigPath("./test")
	u.AddFileType("yaml")
	err := u.ReadInConfiguration()
	if err != nil {
		panic(err)
	}
	c := Conf{}
	err = u.Unmarshal(&c)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", c)
}
```

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
