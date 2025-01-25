package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	// Генерация закрытого ключа
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Ошибка генерации ключа:", err)
		os.Exit(1)
	}

	// Создание самоподписанного сертификата
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	cert := &x509.Certificate{
		SerialNumber: randSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{"My Org"},
		},
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if err != nil {
		fmt.Println("Ошибка создания сертификата:", err)
		os.Exit(1)
	}

	// Сохранение закрытого ключа
	keyFile, err := os.Create("server-go.key")
	if err != nil {
		fmt.Println("Ошибка при создании файла ключа:", err)
		os.Exit(1)
	}
	defer keyFile.Close()
	if err := pem.Encode(keyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		fmt.Println("Ошибка при кодировании ключа:", err)
	}

	// Сохранение сертификата
	certFile, err := os.Create("server-go.crt")
	if err != nil {
		fmt.Println("Ошибка при создании файла сертификата:", err)
		os.Exit(1)
	}
	defer certFile.Close()
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		fmt.Println("Ошибка при кодировании сертификата:", err)
	}

	// Настройка TLS
	config := &tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{certDER}, PrivateKey: priv}}}
	ln, err := tls.Listen("tcp", "localhost:8080", config)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("TLS-сервер запущен и внимает порту 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении:", err)
		return
	}
	fmt.Println("Полученное сообщение:", message)
	_, err = conn.Write([]byte("Сообщение успешно доставлено\n"))
	if err != nil {
		fmt.Println("Ошибка при отправке:", err)
	}
}

func randSerialNumber() *big.Int {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		panic(err)
	}
	return n
}
