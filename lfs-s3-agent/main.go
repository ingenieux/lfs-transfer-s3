package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docopt/docopt.go"
	"github.com/dustin/go-humanize"
	"github.com/ingenieux/lfs-transfer-s3"
	"os"
	"reflect"
	"runtime"
	"strings"
)

const DOC = `
lfs-s3-agent.

Usage:
  lfs-s3-agent [options]
  lfs-s3-agent -h | --help
  lfs-s3-agent --version

Options:
  --loglevel=<level>   Log Level [default: info].
  --part-size=<size>   Size of the Part (set to 0 do disable multipart upload) [default: 100MiB].
  --region=<region>    Region to use [default: us-east-1].
  -h --help            this message
  --version            Show version.
`

func main() {
	log.SetOutput(os.Stderr)

	argsToUse := make([]string, 0)

	if len(os.Args) > 1 {
		argsToUse = strings.Split(os.Args[1], " ")
	}

	version := fmt.Sprintf("0.0.2@%s(%s)",
		reflect.TypeOf(lfstransfers3.TransferAgent{}).PkgPath(),
		runtime.Version())

	args, _ := docopt.Parse(DOC, argsToUse, true, version, true, true)

	partSize := lfstransfers3.DEFAULT_UPLOAD_PART_SIZE

	if partSizeStr, ok := args["--part-size"].(string); ok {
		if partSizeToUse, err := humanize.ParseBigBytes(partSizeStr); nil != err {
			panic(err)
		} else {
			partSize = partSizeToUse.Int64()
		}
	}

	log.SetLevel(log.InfoLevel)

	if logLevelToUse, ok := args["--loglevel"].(string); ok {
		switch strings.ToLower(logLevelToUse) {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		}
	}

	if logFilePath, ok := args["--logfile"].(string); ok {
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

		if nil != err {
			panic(err)
		}

		log.SetOutput(file)
	}

	tempDirPath := os.TempDir()

	region := "us-east-1"

	if regionToUse, ok := args["--region"].(string); ok {
		region = regionToUse
	}

	parameters := &lfstransfers3.Parameters{
		PartSize: partSize,
		TempDir:  tempDirPath,
		Region:   region,
	}

	agent, err := lfstransfers3.NewAgent(parameters)

	if nil != err {
		panic(err)
	}

	rc, err := agent.Execute()

	if nil != err {
		panic(err)
	}

	os.Exit(rc)
}
