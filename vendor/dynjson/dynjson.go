package dynjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"models"
	"os/exec"
	"plugin"
	"time"
)

func Parse(json interface{}) string {
	state := models.State{}

	state.SetRoot(json)
	state.SetDefinitions()

	var definitions string

	for _, definition := range state.Definitions {
		definitions += definition.ToString()
	}

	var definition string

	if state.Root.TypeOf == "object" {
		definition = state.Definition(state.Root).Name()
	} else {
		definition = *state.RecurseKeyName(state.Root)
	}

	return fmt.Sprintf(`package main
	
	import (
		"encoding/json"
	)

	%s
	
	var Root %s
	func JsonDecode(input string) {
		json.Unmarshal([]byte(input), &Root)
	}`, definitions, definition)
}

func ParseString(str string) interface{} {
	var i interface{}

	err := json.Unmarshal([]byte(str), &i)

	if err != nil {
		panic(err)
	}

	content := Parse(i)
	filename := fmt.Sprintf("evaluated-%d", time.Now().UnixNano())

	err = ioutil.WriteFile(filename+".go", []byte(content), 0644)

	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", filename+".so", filename+".go")

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	p, err := plugin.Open(filename + ".so")
	if err != nil {
		panic(err)
	}

	f, err := p.Lookup("JsonDecode")

	if err != nil {
		panic(err)
	}

	f.(func(string))(str)

	v, err := p.Lookup("Root")

	if err != nil {
		panic(err)
	}

	fmt.Println(v)

	return v
}

func ParseFile(filename string) interface{} {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return ParseString(string(content))
}
