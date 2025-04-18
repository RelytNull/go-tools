package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"bytes"
)

func pad(plaintext []byte) []byte {
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

func unpad(plaintext []byte) ([]byte, error) {
	length := len(plaintext)
	if length == 0 {
		return nil, fmt.Errorf("inpit is empty")
	}
	padding := plaintext[length-1]
	if int(padding) > length {
		return nil, fmt.Errorf("padding size is invalid")
	}
	return plaintext[:length-int(padding)], nil
}


func encryptFile(key []byte, filepath string) error {
	plaintext, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Pad plaintext
	plaintext = pad(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Use CBC with random initialization vector (IV)
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	cfb := cipher.NewCBCEncrypter(block, iv)

	// Create a byte slice to hold the ciphertext
	ciphertext := make([]byte, len(plaintext))
	cfb.CryptBlocks(ciphertext, plaintext) // Modify ciphertext in place
	ciphertext = append(iv, ciphertext...) // Prepend IV to the ciphertext

	//ciphertext := cfb.CryptBlocks(make([]byte, aes.BlockSize), plaintext)
	//ciphertext = append(iv, ciphertext...)

	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return ioutil.WriteFile(filepath+".enc", []byte(encodedCiphertext), 0644)

}

func decryptFile(key []byte, filepath string) error {
	ciphertext, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := decodedCiphertext[:aes.BlockSize]
	ciphertext = decodedCiphertext[aes.BlockSize:]

	cfb := cipher.NewCBCDecrypter(block, iv)

	// Create a byte slice to hold the plaintext
	plaintext :=make([]byte, len(ciphertext))
	cfb.CryptBlocks(plaintext, ciphertext) // Modify plaintext in place

	// Unpad the plaintext
	plaintext, err = unpad(plaintext)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath+".dec", plaintext, 0644)

	//plaintext := cfb.CryptBlocks(make([]byte, aes.BlockSize), ciphertext)

	//return ioutil.WriteFile(filepath+".dec", plaintext, 0644)

}

func main() {
	
	var (
		filepath, action string
	 	password 	 []byte
	    )
	    
	fmt.Println("Enter your secret password: ")
	fmt.Scanln(&password)

	// replace with your desired password
	//password := "your_secret_password"
	
	// Ensure the password is of valid length(16,24, or 32 foe AES)
	key := make([]byte, 32) // Create a 32-byte key
	copy(key, password)	// Copy the password into the key (truncates if needed)

	fmt.Println("Enter the file to encrypt or decrypt:")
	// Choose file to encrypt/decrypt
	fmt.Scanln(&filepath)

	// Choose between encryption/decryption

	fmt.Println("Choose and action: encrypt or decrypt?")
	fmt.Scanln(&action)

	if action == "encrypt" {
		err := encryptFile(key, filepath)
		if err != nil {
			fmt.Println("Error encrypting file:", err)
			return
		}
		fmt.Println("File encrypted successfully!")
	} else if action == "decrypt" {
		err := decryptFile(key, filepath)
		if err != nil {
			fmt.Println("Error decrypting file:", err)
			return
		}
		fmt.Println("File decrypted successfully!")
	} else {
		fmt.Println("Invalid action, Please choose 'encrypt' or 'decrypt'")
	}
}
