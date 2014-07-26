package game

import (
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

/**
 * Starts the game. Builds up a new player struct and round struct.
 * Sets the players score to 0.
 * Builds a round with a number of attempts, a set length of password and passes in
 * the word list.
 * At the end of each round, the player's score is printed to the console.
 * Connection is also passed in so data can be written back to the connecting client.
 */
func StartGame(filename string, connection net.Conn, attemptsPerRound int, lengthOfPasswords int, numberOfPasswords int) {
	defer connection.Close()
	wordList := loadStringsFromDisctionaryFile(filename)

	player := player{0} // Setup new player

	for {
		round := buildRound(attemptsPerRound, lengthOfPasswords, numberOfPasswords, wordList)
		playRound(&player, round, connection)
		printPlayerScore(player, connection)
		connection.Write([]byte("PLAY AGAIN? (Y/N): "))
		input := getUserInput(connection)
		if input != "Y" {
			connection.Write([]byte("\nTHANKS FOR PLAYING!"))
			printPlayerScore(player, connection)
			break
		}
	}

}

/**
 * Builds up a Round struct.
 * Takes the number of attempts the round should have,
 * the length of the words provided and the list of words
 * to take them from.
 *
 * First extracts the words of a given length from the provided list of words.
 * Then extracts a number of those words at random. Next, it converts those 10 words
 * to upper case and finally extracts one word at random to be the
 * right answer.
 */
func buildRound(attemptsPerRound int, lengthOfPasswords int, numberOfPasswords int, wordList []string) round {
	certainLengthWords := extractStringsOfLength(lengthOfPasswords, wordList)
	possiblePasswords := extractSubsetOfStringsAtRandom(numberOfPasswords, certainLengthWords)
	convertStringsToUpperCase(possiblePasswords)
	correctWord := extractSubsetOfStringsAtRandom(1, possiblePasswords)[0]
	return round{attemptsPerRound, possiblePasswords, correctWord}
}

/**
 * Plays the round. Requires a player, a round and a connection.
 * Starts by printing the rounds header. Then prints the possible passwords to the connected client.
 * For as long as the player has lives, waits for the user input, evaluates the user's input and responds
 * appropriately.
 *
 * If the user correctly guesses password, increment the player's score by 1, print ACCESS GRANTED and
 * break out of loop.
 *
 * If the user incorrectly guesses password BUT their attempt was one of the provided possible passwords,
 * print ENTRY DENIED and calculate and print the number of letters in the correct place and continue.
 *
 * If the user incorrectly guesses password AND their attempt was NOT one of the provided possible passwords,
 * print ENTRY DENIED and print 0/n letters in correct place.
 */
func playRound(player *player, round round, connection net.Conn) {
	connection.Write([]byte("\n----------------------------------------\n"))
	connection.Write([]byte("\nROBCO INDUSTRIES (TM) TERMALINK PROTOCOL\n"))
	connection.Write([]byte("\nENTER PASSWORD NOW\n\n"))
	printPossiblePasswords(round.possiblePasswords, connection)
	connection.Write([]byte("\n----------------------------------------\n"))
	for i := round.attemptsLeft; i > 0; i-- {
		connection.Write([]byte("\n" + strconv.Itoa(i) + " ATTEMPT(S) LEFT")) // Print number of attempts left
		connection.Write([]byte("\nENTER PASSWORD: "))
		input := getUserInput(connection) // Get user's attempt
		if isValidAttempt(input, round.possiblePasswords) && input == round.correctWord {
			// User guess correctly.
			player.score++
			connection.Write([]byte("ACCESS GRANTED.\n"))
			break
		} else if isValidAttempt(input, round.possiblePasswords) {
			// User guessed incorrectly but was valid attempt
			connection.Write([]byte("ENTRY DENIED. "))
			connection.Write([]byte(strconv.Itoa(calculateNumberOfCorrectLetters(input, round.correctWord)) + "/" + strconv.Itoa(len(round.correctWord)) + " CORRECT.\n"))
		} else {
			// User did not provide a valid attempt
			connection.Write([]byte("ENTRY DENIED. "))
			connection.Write([]byte("0/" + strconv.Itoa(len(round.correctWord)) + " CORRECT.\n"))
		}
	}
}

// Prints the player's score
func printPlayerScore(player player, connection net.Conn) {
	connection.Write([]byte("\nSCORE:" + strconv.Itoa(player.score) + "\n\n"))
}

// Returns user attempt as a string
func getUserInput(connection net.Conn) string {
	var input [512]byte
	n, err := connection.Read(input[0:])
	if err != nil {
		connection.Close()
	}
	return strings.TrimSpace(strings.ToUpper(string(input[0:n])))
}

// Returns true if attempt is in list of provided passwords
func isValidAttempt(attempt string, possiblePasswords []string) bool {
	if stringSliceContains(attempt, possiblePasswords) {
		return true
	}
	return false
}

// Prints the all of the possible passwords within a given round
func printPossiblePasswords(possiblePasswords []string, connection net.Conn) {
	for i := 0; i < len(possiblePasswords); i++ {
		charsToPrepend, charsToAppend := generateRandomCharacterStrings()
		connection.Write([]byte(charsToPrepend + possiblePasswords[i] + charsToAppend + "\n"))
	}
}

// Each password, when displayed, is surrounded by random characters.
// Each word is surrounded by 10 characters
func generateRandomCharacterStrings() (string, string) {
	randomCharacters := ";()[]*&^$.-=<>+#_!?@'/|"
	rand.Seed(time.Now().UnixNano())
	noOfCharsToPrepend := rand.Intn(10)
	noOfCharsToAppend := 10 - noOfCharsToPrepend

	charsToPrepend := ""
	for i := 0; i < noOfCharsToPrepend; i++ {
		charsToPrepend += string(randomCharacters[rand.Intn(len(randomCharacters))])
	}

	charsToAppend := ""
	for i := 0; i < noOfCharsToAppend; i++ {
		charsToAppend += string(randomCharacters[rand.Intn(len(randomCharacters))])
	}

	return charsToAppend, charsToPrepend

}

// Calculates the number of correct letters the user has. Iterates over each letter in word.
func calculateNumberOfCorrectLetters(attempt string, correctPassword string) int {
	correctLetterCount := 0

	for i := 0; i < len(attempt); i++ {
		if attempt[i] == correctPassword[i] {
			correctLetterCount++
		}
	}

	return correctLetterCount
}

// Extracts all strings from given string slice that are a certain length
func extractStringsOfLength(length int, strings []string) []string {
	extractedStrings := []string{}
	for i := 0; i < len(strings); i++ {
		if len(strings[i]) == length {
			extractedStrings = append(extractedStrings, strings[i])
		}
	}
	return extractedStrings
}

// Extracts a given number of strings from provided list of strings
func extractSubsetOfStringsAtRandom(amount int, strings []string) []string {
	extractedStrings := []string{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < amount; i++ {
		possibleString := strings[rand.Intn(len(strings))]
		if stringSliceContains(possibleString, extractedStrings) {
			i--
			continue
		} else {
			extractedStrings = append(extractedStrings, possibleString)
		}
	}
	return extractedStrings
}

// returns true if string is already in slice of strings
func stringSliceContains(word string, words []string) bool {
	for i := 0; i < len(words); i++ {
		if word == words[i] {
			return true
		}
	}
	return false
}

// Iterates over a string slice and converts all strings to uppercase
func convertStringsToUpperCase(slice []string) {
	for i := 0; i < len(slice); i++ {
		slice[i] = strings.ToUpper(slice[i])
	}
}

// Loads words from dictionary file and returns as string slice
func loadStringsFromDisctionaryFile(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	words := strings.Split(string(data), "\n")

	// Trim all whitespace from words
	for i := 0; i < len(words); i++ {
		words[i] = strings.TrimSpace(words[i])
	}

	return words

}

type player struct {
	score int
}

type round struct {
	attemptsLeft      int
	possiblePasswords []string
	correctWord       string
}
