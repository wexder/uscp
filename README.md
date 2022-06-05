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

## Required
You can mark struct fields with tag `uscp:"required"`.
When unmarshallig the **USCP** will check if value of this field is non [zero value](https://go.dev/ref/spec#The_zero_value)
if value is zero error is returned

## ENV vars
You can bind env vars in two ways.
Either by using the autobind path pattern.
```
type Conf struct {
	Obj  struct {
		M int32
	}
}

func main() {
	u := uscp.New()
	u.SetConfigName("conf")
}
```
In this example you can set he value of `M` with env var `conf_Obj_M`.
Autobinding follows simple pattern of filename followed by struct path you want to set with underscore to connect them all. 

Second path is to use `uscp_env` struct tag.
```go
type Conf struct {
	Test string `uscp_env:"TEST"`
}
```
In this example you can set the value of `Test` with env var `TEST`.


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
