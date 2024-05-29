
package main

import (
	"strings"
)

// SeparateArtistAndTitle takes a filename in the format "Artist-Title.mp3"
// and returns the artist and title as separate strings.
func SeparateArtistAndTitle(filename string) (string, string) {
	// Remove the file extension
	filenameWithoutExtension := strings.TrimSuffix(filename, ".mp3")
	// Split the string by the delimiter '-'
  filenameWithoutExtension = replaceUnderscoreWithSpace(filenameWithoutExtension)
	parts := strings.Split(filenameWithoutExtension, "-")

	// Return the artist and title separately
	if len(parts) == 2 {
		artist := strings.TrimSpace(parts[0])
		title := strings.TrimSpace(parts[1])
		return artist, title
	}
  
	return "Unknown", filenameWithoutExtension
}
func replaceUnderscoreWithSpace(input string) string {
	return strings.ReplaceAll(input, "_", " ")
}
// func main() {
// 	// Example usage
// 	filename := "Sabaton-LastStand.mp3"
// 	artist, title := SeparateArtistAndTitle(filename)
// 	fmt.Printf("Artist: %s, Title: %s\n", artist, title)
// }
