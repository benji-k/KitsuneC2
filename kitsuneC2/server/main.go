package main

import (
	"KitsuneC2/server/db"
	"KitsuneC2/server/listener"
	"time"
)

/*
import (
	"KitsuneC2/lib/cryptography"
	"fmt"
	"log"
)
*/

func main() {

	db.Init()
	l1 := listener.Listener{Type: "tcp", Handler: tcpHandler, Network: "127.0.0.1", Port: 4444}
	l1.Start()
	time.Sleep(time.Minute * 10)
	db.Shutdown()

	/*
		privKey, pubKey, err := cryptography.GenerateRSAKeyPair(2048)
		if err != nil {
			log.Fatal("Could not generate keypair")
		}
		privKeyString := cryptography.RsaPrivateKeyToString(privKey)
		pubKeyString := cryptography.RSAPublicKeyToString(pubKey)

		privKeyUnmarshalled, err := cryptography.StringToRsaPrivateKey(privKeyString)
		if err != nil {
			log.Fatal("Could not unmarshal private Key")
		}
		pubKeyUnmarshalled, err := cryptography.StringToRsaPublicKey(pubKeyString)
		if err != nil {
			log.Fatal("Could not unmarshal pub key")
		}
		plainText := "SomeAESKey"
		cipherText, err := cryptography.EncryptWithRsaPublicKey([]byte(plainText), pubKeyUnmarshalled)
		if err != nil {
			log.Fatal("Could not encrypt plain text with public key")
		}
		fmt.Printf("Ciphertext: %s [%d]", string(cipherText), len(cipherText))
		decryptedCipherText, err := cryptography.DecryptWithRsaPrivateKey(cipherText, privKeyUnmarshalled)
		if err != nil {
			log.Fatal("Could not decrypt cipher text with private key")
		}
		fmt.Printf("Decrypted cipher text: %s [%d]", string(decryptedCipherText), len(decryptedCipherText))
	*/
}
