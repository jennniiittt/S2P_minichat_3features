# 🔐 S2P_minichat_3features  
## Secure Mini Chat Application (S2P_minichat)

---

## 📌 Overview
This project is a **secure client-server chat application** developed in Go, with a strong emphasis on **secure coding practices**. It demonstrates key security concepts including authentication, input validation, cryptography, and defense in depth.

---

## 🎯 Objectives
- Implement a functional chat system  
- Apply secure authentication mechanisms  
- Prevent common vulnerabilities (OWASP Top 10)  
- Protect data using cryptographic techniques  
- Demonstrate layered security (Defense in Depth)  

---

## ⚙️ Features

### 🔑 Authentication
- User registration and login  
- Passwords securely hashed using **bcrypt**  
- Prevents plain-text password storage  

### 🧪 Input Validation
- Username and password constraints  
- Message length validation  
- Sanitization of user inputs  

### 🔐 Cryptography
- **Curve25519** for secure key exchange  
- **SHA-256** for shared secret derivation  
- **AES-GCM** for message encryption  

### 📁 Secure Storage
- User credentials stored in `users.json`  
- Passwords stored as hashed values only  

### 🛡️ Defense in Depth
Multiple layers of security:
- Input validation  
- Authentication  
- Encryption  
- Secure storage  

---

## 🧱 System Architecture

Client → Server → Storage (JSON)  
        ↓  
   Security Layers:  
   - Validation  
   - Authentication  
   - Cryptography  

---

## ⚠️ Security Considerations

### OWASP A02 – Cryptographic Failures
- Passwords are hashed using bcrypt  
- Messages are encrypted using AES-GCM  

### OWASP A07 – Identification & Authentication Failures
- Strong password enforcement  
- Secure login verification  
- Duplicate user prevention  

---

## ▶️ How to Run
go run .


### 1. Start Server
type "server" when prompted.

### 2. Start Client
- Run the same program  
- Type "client" for client mode
- Connect to the server

### 3. Usage
- Register a new account  
- Login with credentials  
- Start sending messages  

---

## 🧪 Demonstration Scenarios
- Weak password → rejected  
- Empty/invalid message → blocked  
- Successful login with hashed password  
- Encrypted message communication  

---

## 📜 Logging
- Basic error and authentication feedback implemented  
- Can be extended to include:
  - Login attempt tracking  
  - Intrusion detection  

---

## 🚀 Future Improvements
- Add structured logging system  
- Implement HTTPS / TLS  
- Add session/token-based authentication  
- Replace JSON with a secure database  

---

## 📚 Learning Outcomes
- Applied secure coding principles in a real system  
- Understood OWASP Top 10 vulnerabilities  
- Implemented encryption and hashing techniques  
- Designed layered security architecture  

---

## 👨‍💻 Technologies Used
- Go (Golang)  
- bcrypt (password hashing)  
- AES-GCM (encryption)  
- Curve25519 (key exchange)  
- JSON (data storage)  

---

## 🏁 Conclusion
This project demonstrates how a simple chat application can be secured by integrating authentication, validation, cryptography, and defense-in-depth strategies aligned with modern security standards.