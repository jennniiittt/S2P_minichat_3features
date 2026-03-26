package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"s2p_minichat/auth"
	"s2p_minichat/crypto"
	"s2p_minichat/utils"
	"s2p_minichat/config"
)

func StartServer() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[System] Server started on port 9000...")

	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[System] Client connected!")

	// === KEY EXCHANGE ===
	pub, priv, _ := crypto.GenerateKeyPair()
	conn.Write(pub[:])

	var clientPub [32]byte
	conn.Read(clientPub[:])

	shared, _ := crypto.ComputeSharedSecret(priv, clientPub)
	fmt.Println("[System] Secure channel established.")

	// === AUTH ===
	if !authenticateClient(conn, shared) {
		fmt.Println("[System] Authentication failed. Closing.")
		conn.Close()
		return
	}

	fmt.Println("[System] Chat started. Type 'exit' to quit.")

	go serverReceive(conn, shared)
	serverSend(conn, shared)
}

// ================= AUTH =================

func authenticateClient(conn net.Conn, key []byte) bool {
	encMsg, _ := ServerReadMessage(conn)
	msg, _ := crypto.Decrypt(key, encMsg)

	parts := strings.Split(string(msg), "|")
	if len(parts) < 3 {
		ServerSendMessage(conn, mustEncrypt(key, []byte("FAIL|Invalid format")))
		return false
	}

	action := strings.ToUpper(strings.TrimSpace(parts[0]))
	username := strings.TrimSpace(parts[1])
	password := strings.TrimSpace(parts[2])

	if err := utils.ValidateCredentials(username, password); err != nil {
		ServerSendMessage(conn, mustEncrypt(key, []byte("FAIL|"+err.Error())))
		return false
	}

	// ✅ BASIC VALIDATION
	if username == "" || password == "" {
		ServerSendMessage(conn, mustEncrypt(key, []byte("FAIL|Empty credentials")))
		return false
	}

	switch action {

	case "SIGNUP":
		err := auth.Signup(config.UserFilePath, username, password)
		if err != nil {
			ServerSendMessage(conn, mustEncrypt(key, []byte("FAIL|"+err.Error())))
			return false
		}
		ServerSendMessage(conn, mustEncrypt(key, []byte("SUCCESS|Signup successful")))
		return true

	case "LOGIN":
		err := auth.Authenticate(config.UserFilePath, username, password)
		if err != nil {
			ServerSendMessage(conn, mustEncrypt(key, []byte("FAIL|Invalid credentials")))
			return false
		}
		ServerSendMessage(conn, mustEncrypt(key, []byte("SUCCESS|Login successful")))
		return true

	default:
		ServerSendMessage(conn, mustEncrypt(key, []byte("FAIL|Unknown action")))
		return false
	}
}

func mustEncrypt(key, msg []byte) []byte {
	enc, _ := crypto.Encrypt(key, msg)
	return enc
}

// ================= UTIL =================

func ServerReadMessage(conn net.Conn) ([]byte, error) {
	lenBuf := make([]byte, 4)
	_, err := conn.Read(lenBuf)
	if err != nil {
		return nil, err
	}

	msgLen := binary.BigEndian.Uint32(lenBuf)
	msg := make([]byte, msgLen)
	_, err = conn.Read(msg)
	return msg, err
}

func ServerSendMessage(conn net.Conn, data []byte) error {
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))

	conn.Write(lenBuf)
	_, err := conn.Write(data)
	return err
}

// ================= CHAT =================

func serverReceive(conn net.Conn, key []byte) {
	for {
		msg, err := ServerReadMessage(conn)
		if err != nil {
			fmt.Println("\n[System] Client disconnected.")
			os.Exit(0)
		}

		dec, _ := crypto.Decrypt(key, msg)
		text := string(dec)

		// ✅ VALIDATION (server-side protection)
		if err := utils.ValidateMessage(text); err != nil {
			fmt.Printf("\r\n[Blocked Invalid Message]: %v\n", err)
			fmt.Print("Server: ")
			continue
		}

		// ✅ SANITIZE
		text = utils.SanitizeMessage(text)

		// ✅ CLEAN UI
		fmt.Printf("\r\n[Client]: %s\n", text)
		fmt.Print("Server: ")
	}
}

func serverSend(conn net.Conn, key []byte) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Server: ")
		scanner.Scan()
		msg := scanner.Text()

		if msg == "exit" {
			conn.Close()
			os.Exit(0)
		}

		enc, _ := crypto.Encrypt(key, []byte(msg))
		ServerSendMessage(conn, enc)
	}
}