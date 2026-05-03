package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DaanyaalSobani/go-snake-ssh/game"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const hostKeyPath = "host_key"

func ServeSSH(addr string, hostKey gossh.Signer) error {
	server := &ssh.Server{
		Addr:    addr,
		Handler: handleSession,
	}
	server.AddHostKey(hostKey)
	log.Printf("snake server on %s", addr)
	return server.ListenAndServe()
}

func handleSession(s ssh.Session) {
	pty, _, isPty := s.Pty()
	if !isPty {
		s.Write([]byte("PTY required. Use: ssh -t ...\r\n"))
		return
	}

	// size the game to the player's terminal:
	//   - logical cells render 2 chars wide (CellWidth = 2), so width is half
	//   - subtract 1 row for the HUD line
	cfg := game.DefaultConfig()
	cfg.Width = pty.Window.Width / 2
	cfg.Height = pty.Window.Height - 1

	if cfg.Width < 12 || cfg.Height < 8 {
		fmt.Fprintf(s,
			"Terminal too small (%dx%d). Resize to at least 24x9 and reconnect.\r\n",
			pty.Window.Width, pty.Window.Height)
		return
	}

	// welcome banner — stays on screen until the player hits a key
	fmt.Fprint(s, "\r\n")
	fmt.Fprintf(s, "  Welcome to snake-ssh, %s!\r\n", s.User())
	fmt.Fprint(s, "\r\n  WASD to move   Ctrl-C to quit\r\n")
	fmt.Fprint(s, "\r\n  Press any key to start...\r\n")

	// gate the game on a single keystroke so the banner is actually readable
	if _, err := s.Read(make([]byte, 1)); err != nil {
		return
	}

	log.Printf("session start: user=%q remote=%s size=%dx%d",
		s.User(), s.RemoteAddr(), cfg.Width, cfg.Height)

	if err := game.Run(s, s, cfg); err != nil {
		log.Printf("game ended for %q: %v", s.User(), err)
	}

	log.Printf("session end: user=%q", s.User())
}

func loadHostKey(path string) (gossh.Signer, error) {
	pem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return gossh.ParsePrivateKey(pem)
}

func main() {
	hostKey, err := loadHostKey(hostKeyPath)
	if err != nil {
		log.Fatalf("load host key from %s: %v\n\n"+
			"Generate one with:\n"+
			"  ssh-keygen -t ed25519 -f %s -N \"\"\n",
			hostKeyPath, err, hostKeyPath)
	}
	log.Fatal(ServeSSH(":2222", hostKey))
}
