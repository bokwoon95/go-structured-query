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
	versionLatest = "devlab-mysql"
	version_5_7   = "devlab-mysql-5-7"
)

var (
	resetFlag   = flag.Bool("reset", false, "reset database")
	allFlag     = flag.Bool("all", false, "reset database then load triggers, views, functions, data")
	versionFlag = flag.String("version", "", "which mysql version to use. recognized values: latest, 5.7")
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
	case "5.7":
		containerName = version_5_7
	default:
		fmt.Printf("unrecognized version %s", *versionFlag)
		return
	}
	switch {
	case *resetFlag:
		cmd := exec.Command("docker", "exec", containerName, "mysql", "--user=root", "--password=root", "--database=devlab", "--execute='source /testdata/mysql/init.sql'")
		fmt.Println(cmd)
	case *allFlag:
		var allfiles []string
		switch containerName {
		case versionLatest:
			for _, dir := range []string{"views" /* "triggers", "functions", "data" */} {
				files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
				if err != nil {
					log.Fatalln(err)
				}
				allfiles = append(allfiles, files...)
			}
			cmd := exec.Command("docker", "exec", containerName, "mysql", "--user=root", "--password=root", "--database=devlab")
			mysqlCmd := fmt.Sprint(cmd)
			fmt.Println(mysqlCmd + " --execute='source /testdata/mysql/init.sql';")
			for _, file := range allfiles {
				fmt.Println(mysqlCmd + " --execute='source /testdata/mysql/" + file + "';")
			}
		case version_5_7:
			cmd := exec.Command("docker", "exec", containerName, "mysql", "--user=root", "--password=root", "--database=devlab", "--execute='source /testdata/mysql/init.sql'")
			fmt.Println(cmd)
		}
	}
	cmd := exec.Command("docker", "exec", containerName, "mysql", "--user=root", "--password=root", "--database=devlab", "--execute='source /testdata/mysql/data.sql'")
	fmt.Println(cmd)
}
