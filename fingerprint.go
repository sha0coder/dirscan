package main

import "strings"

const (
	RespRaw = iota
	RespHTML
	RespEXE
	RespELF
	RespPHP
	RespDIR
	RespIMG
	RespZIP
	RespOther
)

// Fingerprint object for fingerprinting responses
func Fingerprint(data string) (int, string) {
	hdr := strings.ToUpper(data[:8])

	if strings.Contains(data, "<h1>Index of /") {
		return RespDIR, "<dir>"
	}
	if strings.Contains(data, "<!DOCTYPE") || strings.Contains(data, "<HTML") || strings.Contains(data, "<html") {
		return RespHTML, "<html>"
	}
	if strings.Contains(hdr, "MZ") {
		return RespEXE, "<exe>"
	}
	if strings.Contains(hdr, "ELF") {
		return RespELF, "<elf>"
	}
	if strings.Contains(hdr, "PKZ") {
		return RespZIP, "<zip>"
	}
	if strings.Contains(hdr, "<?PHP") {
		return RespPHP, "<php>"
	}

	var imgs = []string{"GIF", "PNG", "JPG", "JPEG", "TIF", "ICO", "BMP"}
	for _, img := range imgs {
		if strings.Contains(hdr, img) {
			return RespIMG, "<img>"
		}
	}

	return RespOther, "??"
}
