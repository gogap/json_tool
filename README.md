# json_tool
a tool for json config file, we could get and set key-values in folders

`test.json`
```json
{
	"a":"",
	"b":"2"
}


```
#### set value

```bash
NAME:
   set - set json value

USAGE:
   command set [command options] [arguments...]

OPTIONS:
   --file, -f [--file option --file option]	json files or folders
   --ext, -e ".json"				file ext, default is json
   --key, -k 					the key to be set
   --value, -v 					the key to be set
   --type, -t "string"				the value type, default is string, (string|int|float|object)
```


```bash
$ json_tool get -f "./" -e ".json" -k a -v 1 -t string
```

the `test.json` will be modify as follow:

```json
{
	"a":"1",
	"b":"2"
}

```bash
$ json_tool get -f "./" -e ".json" -k a -v {\"c\":1} -t object
```

the `test.json` will be modify as follow:

```json
{
	"a":{"c":1},
	"b":"2"
}


#### get value

```bash
NAME:
   get - get json value

USAGE:
   command get [command options] [arguments...]

OPTIONS:
   --file, -f [--file option --file option]	json files or folders
   --ext, -e ".json"				file ext, default is json
   --key, -k
```


```bash
$ json_tool get -f "./" -e ".json" -k a
```

```bash
./a.json:
{
   "a": {
     "b": 1
   }
 }

 ./b.json:
{
   "a": {
     "c": 1
   }
 }

 ...
```