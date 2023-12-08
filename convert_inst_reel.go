package main

import "strings"

func convertInstReel(srcPath string) (res string, changed bool) {
	if strings.Contains(srcPath, "www.instagram.com/reel") {
		return strings.Replace(srcPath, "www.instagram.com/reel", "ddinstagram.com/reel", -1), true
	}
	return srcPath, false
}
