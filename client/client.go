package client

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"s2p_minichat/crypto"
	"s2p_minichat/utils"
)

func StartClient() {
	reader := bufio.NewReader(os.Stdin)

	// === CONNECT ===
	fmt.Print("Server IP: ")
	serverIP, _ := reader.ReadString('\n')
	serverIP = strings.TrimSpace(serverIP)

	conn, err := net.Dial("tcp", serverIP+":9000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[System] Connected to server!")

	// === KEY EXCHANGE ===
	pub, priv, _ := crypto.GenerateKeyPair()
	var serverPub [32]byte
	conn.Read(serverPub[:])
	conn.Write(pub[:])

	shared, _ := crypto.ComputeSharedSecret(priv, serverPub)
	fmt.Println("[System] Secure channel established.")

	// === AUTH ===
	fmt.Print("Choose SIGNUP or LOGIN: ")
	action, _ := reader.ReadString('\n')
	action = strings.TrimSpace(action)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	payload := fmt.Sprintf("%s|%s|%s", action, username, password)
	encPayload, _ := crypto.Encrypt(shared, []byte(payload))
	ClientSendMessage(conn, encPayload)

	resp, _ := ClientReadMessage(conn)
	decResp, _ := crypto.Decrypt(shared, resp)

	fmt.Println("[Server]:", string(decResp))

	if !strings.HasPrefix(string(decResp), "SUCCESS") {
		fmt.Println("[System] Authentication failed.")
		return
	}

	fmt.Println("[System] Chat started. Type 'exit' to quit.")

	// === CHAT ===
	go clientReceive(conn, shared)
	clientSend(conn, shared)
}

// ================= UTIL =================

func ClientReadMessage(conn net.Conn) ([]byte, error) {
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

func ClientSendMessage(conn net.Conn, data []byte) error {
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))

	conn.Write(lenBuf)
	_, err := conn.Write(data)
	return err
}

// ================= CHAT =================

func clientReceive(conn net.Conn, key []byte) {
	for {
		msg, err := ClientReadMessage(conn)
		if err != nil {
			fmt.Println("\n[System] Server disconnected.")
			os.Exit(0)
		}

		dec, _ := crypto.Decrypt(key, msg)

		// ✅ CLEAN UI
		fmt.Printf("\r\n[Server]: %s\n", string(dec))
		fmt.Print("You: ")
	}
}

func clientSend(conn net.Conn, key []byte) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("You: ")
		scanner.Scan()
		text := scanner.Text()

		if text == "exit" {
			conn.Close()
			os.Exit(0)
		}

		// ✅ VALIDATION
		if err := utils.ValidateMessage(text); err != nil {
			fmt.Println("[Error]:", err)
			continue
		}

		// ✅ SANITIZE
		text = utils.SanitizeMessage(text)

		enc, _ := crypto.Encrypt(key, []byte(text))
		ClientSendMessage(conn, enc)
	}
}