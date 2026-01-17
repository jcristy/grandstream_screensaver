package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Screensaver struct {
	XMLName xml.Name `xml:"screensaver"` // Root element name "screensaver"
	// Images is a slice of an anonymous struct.
	// Each element of this slice will be marshalled as an <image> element.
	Images []ImagePathHolder `xml:"image"` // This tag applies to each element in the slice to become <image>
}

type ImagePathHolder struct {
	Path string `xml:"path,attr"` // 'path' attribute for the <image> tag
}

// screensaverHandler handles requests for /screensaver.xml
func screensaverHandler(w http.ResponseWriter, r *http.Request) {
	var serverHost string
	if envHost := os.Getenv("SERVER_EXTERNAL_HOST"); envHost != "" {
		serverHost = envHost
	} else if r.Host != "" {
		serverHost = r.Host
	} else {
		serverHost = "localhost:8080" // Fallback default
	}

	s := Screensaver{
		XMLName: xml.Name{
			Space: "",
			Local: "",
		},
		Images: []ImagePathHolder{},
	}
	images, err := getImageFiles("images")
	for _, image := range images {
		s.Images = append(s.Images, ImagePathHolder{Path: fmt.Sprintf("http://%s/%s", serverHost, image)})
	}

	xmlOutput, err := xml.Marshal(s)
	if err != nil {
		log.Fatalf("Error marshalling to XML: %v", err)
	}

	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write([]byte(xmlOutput))
	if err != nil {
		log.Printf("Error writing screensaver.xml response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Handler for screensaver.xml
	http.HandleFunc("/screensaver.xml", screensaverHandler)

	// Handler for files in the images/ directory
	// http.StripPrefix removes the leading slash before looking up the file
	http.Handle("/", http.FileServer(http.Dir("./images")))

	log.Printf("Server starting on port %s", port)
	log.Printf("Access screensaver.xml at http://localhost:%s/screensaver.xml", port)
	log.Printf("Access image files at http://localhost:%s/<image_name.jpg>", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getImageFiles reads the specified directory and returns a slice of image filenames.
// It filters by common image extensions.
func getImageFiles(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory '%s': %w", dirPath, err)
	}

	var imageFiles []string
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
	}

	for _, entry := range entries {
		if !entry.IsDir() { // Only process files, not subdirectories
			filename := entry.Name()
			ext := strings.ToLower(filepath.Ext(filename)) // Get file extension
			if imageExtensions[ext] {
				imageFiles = append(imageFiles, filename)
			}
		}
	}
	return imageFiles, nil
}
