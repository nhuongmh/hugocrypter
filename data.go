package main

import (
	"embed"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func copyFile(sourcePath, destinationPath string, content embed.FS) error {
	// Read content from embedded files
	file, err := content.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create the target file and write the contents
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, file)
	if err != nil {
		return err
	}
	return nil
}

func walkHTMLFiles() error {
	err := filepath.WalkDir("public", func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".html" {
			getData(path)
		}

		return nil
	})

	return err
}

func getData(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error reading file:", path)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(file)))
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return
	}

	doc.Find("body").AppendHtml(`<script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>`)
	doc.Find("body").AppendHtml(`<script src="https://cdnjs.cloudflare.com/ajax/libs/crypto-js/3.1.9-1/crypto-js.js"></script>`)
	doc.Find("body").AppendHtml(`<script src="/js/AESDecrypt.js"></script>`)
	secretElements := doc.Find("div#secret")
	if secretElements.Length() == 0 {
		return
	}
	log.Printf("Processing: %v \n", path)
	passwordAttr, _ := secretElements.Attr("password")
	innerText := secretElements.Text()

	secretElements.RemoveAttr("password")
	encryptedPassword := GetEncryptedPassword(passwordAttr)
	log.Printf("  Encrypted password: %v \n", encryptedPassword)
	encryptedContent, err := AESEncrypt(innerText, encryptedPassword)
	if err != nil {
		log.Fatal("crypto.AESEncrypt(innerText, encryptedPassword) gets err", err)
	}
	secretElements.SetText(encryptedContent)
	newHtml, err := doc.Html()
	if err != nil {
		log.Fatal("doc.Html() gets err: ", err)
	}
	err = os.WriteFile(path, []byte(newHtml), 0644)
	if err != nil {
		log.Fatal("os.WriteFile(path, []byte(newHtml), 0644) gets err: ", err)
	}
}
