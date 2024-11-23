package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Generate a random number less than `modulus`
func randomNumber(modulus *big.Int) *big.Int {
	n, err := rand.Int(rand.Reader, modulus)
	if err != nil {
		panic(err)
	}
	return n
}

func main() {
	// Modulus and secret
	modulus := big.NewInt(23)                                   // Example modulus (a small prime number)
	secret := big.NewInt(5)                                     // Prover's secret (e.g., x)
	squared := new(big.Int).Exp(secret, big.NewInt(2), modulus) // x^2 % modulus

	fmt.Println("=== Fiat-Shamir Zero Knowledge Proof ===")
	fmt.Printf("Modulus: %d\n", modulus)
	fmt.Printf("Prover's Secret (hidden): %d\n", secret)
	fmt.Printf("Value to Prove Knowledge Of (squared): %d\n", squared)

	// Step 1: Prover generates a random value and computes its square
	random := randomNumber(modulus)
	randomSquare := new(big.Int).Exp(random, big.NewInt(2), modulus)
	fmt.Printf("Prover sends random squared (commitment): %d\n", randomSquare)

	// Step 2: Verifier generates a random challenge (0 or 1)
	challenge := randomNumber(big.NewInt(2)) // Randomly choose 0 or 1
	fmt.Printf("Verifier sends challenge: %d\n", challenge)

	// Step 3: Prover computes the response
	response := new(big.Int)
	if challenge.Cmp(big.NewInt(0)) == 0 {
		// Response is the random value if challenge == 0
		response.Set(random)
	} else {
		// Response is (random * secret) % modulus if challenge == 1
		response.Mul(random, secret)
		response.Mod(response, modulus)
	}
	fmt.Printf("Prover sends response: %d\n", response)

	// Step 4: Verifier checks the proof
	left := new(big.Int).Exp(response, big.NewInt(2), modulus) // response^2 % modulus
	right := new(big.Int)
	if challenge.Cmp(big.NewInt(0)) == 0 {
		// right = randomSquare
		right.Set(randomSquare)
	} else {
		// right = (randomSquare * squared) % modulus
		right.Mul(randomSquare, squared)
		right.Mod(right, modulus)
	}

	fmt.Printf("Verifier checks: Left (%d) == Right (%d)? ", left, right)
	if left.Cmp(right) == 0 {
		fmt.Println("Proof Verified!")
	} else {
		fmt.Println("Proof Failed!")
	}
}
