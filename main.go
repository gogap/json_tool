package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "json_tool"
	app.Usage = "a tool for batch set json values"

	app.Commands = []cli.Command{
		{
			Name:   "set",
			Usage:  "set json value",
			Action: set,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "file, f",
					Usage: "json files or folders",
					Value: new(cli.StringSlice),
				}, cli.StringFlag{
					Name:  "ext, e",
					Usage: "file ext, default is json",
					Value: ".json",
				},
				cli.StringFlag{
					Name:  "key, k",
					Usage: "the key to be set",
				}, cli.StringFlag{
					Name:  "value, v",
					Usage: "the key to be set",
				}, cli.StringFlag{
					Name:  "type, t",
					Usage: "the value type, default is string, (string|int|float|object)",
					Value: "string",
				}, cli.BoolFlag{
					Name:  "create, c",
					Usage: "if create = true, the key will be create if env file not exist the key",
				},
			},
		},
		{
			Name:   "get",
			Usage:  "get json value",
			Action: get,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "file, f",
					Usage: "json files or folders",
					Value: new(cli.StringSlice),
				}, cli.StringFlag{
					Name:  "ext, e",
					Usage: "file ext, default is json",
					Value: ".json",
				},
				cli.StringFlag{
					Name:  "key, k",
					Usage: "the key to be get",
				},
			},
		},
	}
	app.Run(os.Args)
}

func expandFiles(fileAndFolders []string, fileExt string) (files []string, err error) {
	allFiles := []string{}
	walkFunc := func(path string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return
		}

		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(path); ext == fileExt {
			absPath := path
			if !filepath.IsAbs(path) {
				if absPath, e = filepath.Abs(path); e != nil {
					return e
				}
			}
			allFiles = append(allFiles, absPath)
		}

		return nil
	}

	for _, fileOrFolder := range fileAndFolders {
		filepath.Walk(fileOrFolder, walkFunc)
	}

	files = allFiles
	return
}

func set(c *cli.Context) {
	files := c.StringSlice("file")
	ext := c.String("ext")
	key := c.String("key")
	value := c.String("value")
	Type := c.String("type")
	create := c.Bool("create")

	if files == nil || len(files) == 0 {
		fmt.Println("file list is empty")
		return
	}

	if key == "" {
		fmt.Println("key is empty")
		return
	}

	jsonFiles, e := expandFiles(files, ext)
	if e != nil {
		fmt.Println(e)
		return
	}

	for _, jsonFile := range jsonFiles {
		fmt.Printf("%s:\n", jsonFile)
		v := map[string]interface{}{}
		data, e := ioutil.ReadFile(jsonFile)
		if e != nil {
			fmt.Println(e)
			return
		}

		e = json.Unmarshal(data, &v)
		if e != nil {
			fmt.Println(e)
			return
		}

		if _, exist := v[key]; !exist && !create {
			continue
		}

		switch Type {
		case "string":
			v[key] = value
		case "int":
			{
				if iV, e := strconv.ParseInt(value, 10, 64); e != nil {
					fmt.Println(e)
					return
				} else {
					v[key] = iV
				}
			}
		case "float":
			{
				if fV, e := strconv.ParseFloat(value, 10); e != nil {
					fmt.Println(e)
					return
				} else {
					v[key] = fV
				}
			}
		case "object":
			{
				obj := map[string]interface{}{}
				if e := json.Unmarshal([]byte(value), &obj); e != nil {
					fmt.Println(e)
					return
				} else {
					v[key] = obj
				}
			}
		}

		if fi, e := os.Stat(jsonFile); e != nil {
			fmt.Println(e)
			return
		} else {
			if newData, e := json.MarshalIndent(v, " ", "  "); e != nil {
				fmt.Println(e)
				return
			} else {
				if e := ioutil.WriteFile(jsonFile, newData, fi.Mode()); e != nil {
					fmt.Println(e)
					return
				}
			}
		}
	}
}

func get(c *cli.Context) {
	files := c.StringSlice("file")
	ext := c.String("ext")
	key := c.String("key")

	if files == nil || len(files) == 0 {
		fmt.Println("file list is empty")
		return
	}

	if key == "" {
		fmt.Println("key is empty")
		return
	}

	jsonFiles, e := expandFiles(files, ext)
	if e != nil {
		fmt.Println(e)
		return
	}

	for _, jsonFile := range jsonFiles {
		fmt.Printf("%s:\n", jsonFile)
		v := map[string]interface{}{}
		data, e := ioutil.ReadFile(jsonFile)
		if e != nil {
			fmt.Println(e)
			return
		}

		e = json.Unmarshal(data, &v)
		if e != nil {
			fmt.Println(e)
			return
		}

		if val, exist := v[key]; exist {
			newKV := map[string]interface{}{key: val}
			if b, e := json.MarshalIndent(newKV, " ", "  "); e != nil {
				fmt.Println(e)
			} else {
				fmt.Println(string(b))
				return
			}
		} else {
			fmt.Println("the key of", key, "not exist")
		}

		fmt.Println("")
	}
}
