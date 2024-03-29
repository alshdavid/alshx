package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/alshdavid/alshx/tools/alshx/platform/archive"
	"github.com/alshdavid/alshx/tools/alshx/platform/download"
	"github.com/alshdavid/alshx/tools/alshx/platform/flags"
	"github.com/alshdavid/alshx/tools/alshx/platform/github"
	"github.com/alshdavid/alshx/tools/alshx/platform/logging"
	"github.com/alshdavid/alshx/tools/alshx/platform/meta"
)

func main() {
	var err error
	var args = flags.NewArgs()
	var logger = logging.NewLogger(args.Verbose)

	if !args.Update && !args.Install {
		fmt.Println("Alshx Command Line Utilities Version:", meta.Version, meta.ReleaseDate)
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("alshx update")
		fmt.Println("alshx update --path /usr/local/alshx")
		os.Exit(0)
	}

	var alshxPath = args.Path
	if alshxPath == "" {
		e, _ := os.Executable()
		alshxPath = filepath.Join(path.Dir(e))
	}

	fmt.Println("Alshx Command Line Utilities")
	fmt.Println("Folder to install to:", alshxPath)

	if args.Dry {
		logger.Log("25% Downloading new release")
		logger.Log("50% Removing existing binaries")
		logger.Log("75% Unpacking binaries")
		logger.Log("---")
		if args.Install {
			logger.Log("Alshx utilities installed to:", alshxPath)
		}
		if args.Update {
			logger.Log("Alshx utilities updated at:", alshxPath)
		}
		os.Exit(0)
	}

	tempDir, err := ioutil.TempDir("", "alshx")
	if err != nil {
		logger.Error("ERROR: Unable to create temp directory")
		os.Exit(1)
	}

	var archivePath = filepath.Join(tempDir, github.GetArchiveName())

	defer os.RemoveAll(tempDir)

	logger.Log("---")

	logger.Log("25% Downloading new release")
	err = download.File(github.GetReleaseUrl(), archivePath)
	if err != nil {
		logger.Error("ERROR: Unable to download archive")
		os.Exit(1)
	}

	logger.Log("50% Removing existing binaries")
	binaries, _ := ioutil.ReadDir(alshxPath)

	for _, file := range binaries {
		err = os.RemoveAll(filepath.Join(alshxPath, file.Name()))
		if err != nil {
			logger.Error("ERROR: Unable to remove directory:", alshxPath)
			os.Exit(1)
		}
	}

	logger.Log("75% Unpacking binaries")
	err = os.MkdirAll(alshxPath, 0755)
	if err != nil {
		logger.Error("ERROR: Unable to make directory:", alshxPath)
		os.Exit(1)
	}

	err = archive.Unzip(archivePath, alshxPath)
	if err != nil {
		logger.Error("ERROR: Unable to download archive")
		os.Exit(1)
	}

	logger.Log("---")
	if args.Install {
		logger.Log("Alshx utilities installed to:", alshxPath)
	}
	if args.Update {
		logger.Log("Alshx utilities updated at:", alshxPath)
	}
}
