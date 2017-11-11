package nepho

import "fmt"

func CreateWithFile(stackName, file, paramsFile string) {
	fmt.Printf("Creating %s from config: %s (%s)", stackName, file, paramsFile)
}

func CreateWithPath(stackName, path, paramsFile string) {
	fmt.Printf("Creating %s from config: %s (%s)", stackName, path, paramsFile)
}
