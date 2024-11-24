package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

//go:embed AESDecrypt.js
var aesDecryptScript embed.FS

//go:embed secret.html
var secretHtml embed.FS

func printUsage() {
	fmt.Println("Usage: ./hugocrypter [command]")
	fmt.Println("Available commands:")
	fmt.Println("  pre    - Execute pre-processing tasks")
	fmt.Println("  post   - Execute post-processing tasks")
	fmt.Println("--help   - Show this help message")
	fmt.Println(" -f   	  - Force overwrite in pre-build process step")
	fmt.Println("If no command specified, full process will be executed: pre-process -> hugo (build) -> post-process")
}

func prebuild_process(force bool) {
	log.Println("pre-build processing....")
	aerDecryptJsFilePath := "static/js/AESDecrypt.js"
	if _, err := os.Stat(aerDecryptJsFilePath); force || errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path.Dir(aerDecryptJsFilePath), os.ModePerm)
		if err != nil {
			log.Println(" ", err)
			return
		}
		err = copyFile(filepath.Base(aerDecryptJsFilePath), aerDecryptJsFilePath, aesDecryptScript)
		if err != nil {
			log.Fatalf("data.CopyFile: AESDecrypt.js gets error %v", err)
		}
		log.Printf("Successfully created %s \n", aerDecryptJsFilePath)
	}

	secretShortcodeFilePath := "layouts/shortcodes/secret.html"
	if _, err := os.Stat(secretShortcodeFilePath); force || errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path.Dir(secretShortcodeFilePath), os.ModePerm)
		if err != nil {
			log.Println("layouts/shortcodes create fail:", err)
			return
		}
		err = copyFile(filepath.Base(secretShortcodeFilePath), secretShortcodeFilePath, secretHtml)
		if err != nil {
			log.Fatal("data.CopyFile: secret.html gets error", err)
		}
		log.Printf("Successfully created %s \n", secretShortcodeFilePath)
	}
}

func postbuild_process() {
	log.Println("post-build processing....")
	// solve html files in public folder
	err := walkHTMLFiles()
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func full_process(force bool) {
	prebuild_process(force)
	output, err := exec.Command("hugo").Output()
	fmt.Println(string(output))
	if err != nil {
		log.Fatalln("cmd.Output() gets error", err)
	}

	postbuild_process()
}

func main() {
	force := false
	for _, arg := range os.Args {
		if arg == "-f" {
			force = true
		}
	}
	if force || len(os.Args) < 2 {
		full_process(force)
		return
	}

	command := os.Args[1]
	switch command {
	case "pre":
		prebuild_process(force)
	case "post":
		postbuild_process()
	case "--help":
		printUsage()
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}
