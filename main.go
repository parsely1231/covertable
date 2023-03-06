package main

import (
	"errors"
	"log"
	"os"

	"golang.org/x/tools/cover"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatalf("failed to parse args: %v", err)
	}

	profiles, err := cover.ParseProfiles(args.inputProfilePath)
	if err != nil {
		log.Fatalf("failed to parse profiles: %v", err)
	}

	profileTree := makeProfileTree(args.profileRootPath)
	for _, profile := range profiles {
		err = profileTree.AddProfile(*profile)
		if err != nil {
			log.Fatal(err)
		}
	}

	csvExporter := makeCSVExporter(profileTree)
	err = csvExporter.Save(args.savePath)
	if err != nil {
		log.Fatal(err)
	}
}

type parsedArgs struct {
	profileRootPath  string
	inputProfilePath string
	savePath         string
}

func parseArgs() (*parsedArgs, error) {
	if len(os.Args) < 3 {
		return nil, errors.New("incomplete command: need $0 profileRootPath inputProfilePath savePath")
	}

	args := &parsedArgs{
		profileRootPath:  os.Args[1],
		inputProfilePath: os.Args[2],
		savePath:         os.Args[3],
	}

	return args, nil
}
