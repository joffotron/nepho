package main

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/joffotron/nepho/cloudformation"
	"github.com/joffotron/nepho/preprocess"
)

func createWithFile(stackName, fileName, paramsFile string) {
	fmt.Printf("Creating %s stack from config: %s (%s)...\n\n", stackName, fileName, paramsFile)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	var yamlData map[string]interface{}
	err = yaml.Unmarshal(file, &yamlData)
	if err != nil {
		panic(err)
	}

	translated := preprocess.Translate(yamlData)
	yamlOut, err := yaml.Marshal(translated)
	fmt.Println(string(yamlOut))

	cfn := cloudformation.New(stackName)
	err = cfn.Create(string(yamlOut))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(yamlOut))
}

func createWithPath(stackName, path, paramsFile string) {
	fmt.Printf("Creating %s stack from files in: %s (%s)...\n\n", stackName, path, paramsFile)
	panic("Not implemented, hah :p")
}
