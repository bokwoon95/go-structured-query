package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	versionLatest = "devlab-postgres"
	version_9_5   = "devlab-postgres-9-5"
)

var (
	resetFlag   = flag.Bool("reset", false, "reset database")
	allFlag     = flag.Bool("all", false, "reset database then load triggers, views, functions, data")
	versionFlag = flag.String("version", "", "which postgres version to use. recognized values: latest, 9.5")
)

var (
	// sourcefile is the path to this file
	_, sourcefile, _, _ = runtime.Caller(0)
	// rootdir is the directory containing this source file
	rootdir = filepath.Dir(sourcefile) + string(os.PathSeparator)
)

func main() {
	flag.Parse()
	if !*resetFlag && !*allFlag {
		flag.PrintDefaults()
		return
	}
	var containerName string
	switch *versionFlag {
	case "":
		fallthrough
	case "latest":
		containerName = versionLatest
	case "9.5":
		containerName = version_9_5
	default:
		fmt.Printf("unrecognized version %s", *versionFlag)
		return
	}
	switch {
	case *resetFlag:
		cmd := exec.Command("docker", "exec", containerName, "psql", "--variable=ON_ERROR_STOP=1", "--username=postgres", "--dbname=devlab", "--file='/testdata/postgres/init.sql'")
		fmt.Println(cmd)
	case *allFlag:
		var allfiles []string
		for _, dir := range []string{"views", "triggers", "functions", "data"} {
			files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
			if err != nil {
				log.Fatalln(err)
			}
			allfiles = append(allfiles, files...)
		}
		switch containerName {
		case versionLatest:
			cmdArgs := []string{"exec", containerName, "psql", "--variable=ON_ERROR_STOP=1", "--username=postgres", "--dbname=devlab", "--file='/testdata/postgres/init.sql'"}
			for _, file := range allfiles {
				cmdArgs = append(cmdArgs, "--file=/testdata/postgres/"+file)
			}
			cmd := exec.Command("docker", cmdArgs...)
			fmt.Println(cmd)
		case version_9_5:
			cmd := exec.Command("docker", "exec", containerName, "psql", "--variable=ON_ERROR_STOP=1", "--username=postgres", "--dbname=devlab")
			psqlCmd := fmt.Sprint(cmd)
			fmt.Println(psqlCmd + " --file=/testdata/postgres/init.sql;")
			for _, file := range allfiles {
				fmt.Println(psqlCmd + " --file=/testdata/postgres/" + file + ";")
			}
		}
	}
}
