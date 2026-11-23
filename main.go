// Launcher kiosk antrean Erzap (Windows).
// Setara dengan erzap-kiosk-antrean.bat: baca link dari link.md di sebelah .exe,
// kalau belum ada minta paste link lalu simpan ke link.md, lalu buka kiosk.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func browserPath() string {
	pf := os.Getenv("ProgramFiles")
	pf86 := os.Getenv("ProgramFiles(x86)")
	local := os.Getenv("LocalAppData")
	candidates := []string{
		filepath.Join(pf, `Google\Chrome\Application\chrome.exe`),
		filepath.Join(pf86, `Google\Chrome\Application\chrome.exe`),
		filepath.Join(local, `Google\Chrome\Application\chrome.exe`),
		filepath.Join(pf86, `Microsoft\Edge\Application\msedge.exe`),
		filepath.Join(pf, `Microsoft\Edge\Application\msedge.exe`),
	}
	for _, c := range candidates {
		if !filepath.IsAbs(c) { // variabel lingkungan kosong
			continue
		}
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

func isURL(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), "http")
}

// readLink mengambil baris pertama di link.md yang diawali http.
func readLink(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if isURL(line) {
			return line
		}
	}
	return ""
}

func pause(msg string) {
	fmt.Println(msg)
	fmt.Println("Tekan Enter untuk menutup.")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func hideConsole() {
	syscall.NewLazyDLL("kernel32.dll").NewProc("FreeConsole").Call()
}

func main() {
	browser := browserPath()
	if browser == "" {
		pause("Chrome atau Microsoft Edge tidak ditemukan di PC ini.")
		os.Exit(1)
	}

	exe, err := os.Executable()
	if err != nil {
		pause("Gagal menentukan lokasi program: " + err.Error())
		os.Exit(1)
	}
	linkFile := filepath.Join(filepath.Dir(exe), "link.md")

	url := readLink(linkFile)
	if url == "" {
		in := bufio.NewReader(os.Stdin)
		for {
			fmt.Println("==============================================")
			fmt.Println("  Link kiosk antrean belum diatur.")
			fmt.Println("  Paste link web kiosk di bawah ini, lalu Enter.")
			fmt.Println("  Tips: klik kanan di jendela ini untuk paste.")
			fmt.Println("==============================================")
			fmt.Println()
			fmt.Print("Link: ")
			line, _ := in.ReadString('\n')
			line = strings.TrimSpace(line)
			if isURL(line) {
				url = line
				break
			}
			fmt.Println()
			fmt.Println("Link harus diawali http atau https. Coba lagi.")
			fmt.Println()
		}
		if err := os.WriteFile(linkFile, []byte(url+"\r\n"), 0o644); err != nil {
			fmt.Println("Peringatan: gagal menyimpan link.md: " + err.Error())
			fmt.Println("Kiosk tetap dibuka, tapi link akan ditanya lagi lain kali.")
			fmt.Println()
		}
	}

	profile := filepath.Join(os.Getenv("LocalAppData"), "ErzapKioskAntrean")
	cmd := exec.Command(browser,
		"--app="+url,
		"--start-fullscreen",
		"--kiosk-printing",
		"--user-data-dir="+profile,
	)
	if err := cmd.Start(); err != nil {
		pause("Gagal membuka browser: " + err.Error())
		os.Exit(1)
	}
	hideConsole()
}
