package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"regexp"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// `.env` file isn't required, at least if the
		// environment variables are set elsewhere
		fmt.Println(err)
	}

	downloadsDir := os.Getenv("DOWNLOADSDIR")
	args := os.Args
	albumURL := args[1]
	finalizeWithZip := args[len(args)-1] == "-z"

	pageSource := getContents(albumURL)
	modelName, albumName := getAlbumInfo(pageSource)
	imagesFound := crawlImages(pageSource)

	fmt.Println("Found", albumName, "set from", modelName, "!")
	fmt.Println("Found", len(imagesFound), "images in set. Downloading...")

	replaceSpace := regexp.MustCompile(`\s+`)
	albumDir := replaceSpace.ReplaceAllString(strings.TrimSpace(downloadsDir + "/" + modelName + " - " + albumName), " ")

	checkAndCreateDir(downloadsDir)
	checkAndCreateDir(albumDir)
	imagesDownloaded := []string{}

	for i, imageURL := range imagesFound {
		pad := ""

		if i <= 10 {
			pad = "0"
		}
		
		filename := leftPad(strconv.Itoa(i + 1), pad, digitsLen(len(imagesFound))-1) + ".jpg"
		imageOutput := albumDir + "/" + filename
		fmt.Println(imageURL + " -> " + imageOutput)
		imagesDownloaded = append(imagesDownloaded, imageOutput)

		b, _ := saveImage(imageURL, imageOutput)
		fmt.Println("File size:", b)
	}

	if finalizeWithZip {
		err := ZipFiles(albumDir+"/"+albumName+".zip", imagesDownloaded)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Done... Enjoy!")
}
