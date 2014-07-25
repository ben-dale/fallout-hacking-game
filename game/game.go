package game

import (
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const attemptsPerRound int = 4   // Amount of attempts to guess the password the player has
const lengthOfPassword int = 7   // Length of the passwords
const numberOfPasswords int = 10 // Number of passwords to guess from

/**
 * Starts the game. Builds up a new player struct and round struct.
 * Sets the players score to 0.
 * Builds a round with a number of attempts, a set length of password and passes in
 * the word list.
 * At the end of each round, the player's score is printed to the console.
 * Connection is also passed in so data can be written back to the connecting client.
 */
func StartGame(filename string, connection net.Conn) {
	defer connection.Close()
	wordList := loadWordsFromDictionaryFile(filename)

	player := player{0} // Setup new player

	for {
		round := buildRound(attemptsPerRound, lengthOfPassword, wordList)
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
func buildRound(attempts int, wordLength int, wordList []string) round {
	certainLengthWords := extractWordsOfLength(wordLength, wordList)
	possiblePasswords := extractSubsetOfWordsAtRandom(numberOfPasswords, certainLengthWords)
	convertWordsToUpperCase(possiblePasswords)
	correctWord := extractSubsetOfWordsAtRandom(1, possiblePasswords)[0]
	return round{attempts, possiblePasswords, correctWord}
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
	printRoundHeader(connection)
	printPossiblePasswords(round.possiblePasswords, connection)
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
	connection.Write([]byte("\nSCORE:" + strconv.Itoa(player.score) + "\n"))
}

// Prints the round's header
func printRoundHeader(connection net.Conn) {
	connection.Write([]byte("\n----------------------------------------\n"))
	connection.Write([]byte("\nROBCO INDUSTRIES (TM) TERMALINK PROTOCOL\n"))
	connection.Write([]byte("\nENTER PASSWORD NOW\n\n"))
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
	if stringArrayContains(attempt, possiblePasswords) {
		return true
	}
	return false
}

// Iterates over a string slice and converts all passwords to uppercase
func convertWordsToUpperCase(passwords []string) {
	for i := 0; i < len(passwords); i++ {
		passwords[i] = strings.ToUpper(passwords[i])
	}
}

// Prints the all of the possible passwords within a given round
func printPossiblePasswords(possiblePasswords []string, connection net.Conn) {
	for i := 0; i < len(possiblePasswords); i++ {
		connection.Write([]byte(possiblePasswords[i] + "\n"))
	}
	connection.Write([]byte("\n----------------------------------------\n"))
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

// Extracts all words from given word list that are a given length
func extractWordsOfLength(length int, words []string) []string {
	extractedWords := []string{}
	for i := 0; i < len(words); i++ {
		if len(words[i]) == length {
			extractedWords = append(extractedWords, words[i])
		}
	}
	return extractedWords
}

// Extracts a given number of words from provided list of words
func extractSubsetOfWordsAtRandom(amount int, words []string) []string {
	extractedWords := []string{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < amount; i++ {
		possibleWord := words[rand.Intn(len(words))]
		if stringArrayContains(possibleWord, extractedWords) {
			i--
		} else {
			extractedWords = append(extractedWords, possibleWord)
		}
	}
	return extractedWords
}

// returns true if word is already in slice of words
func stringArrayContains(word string, words []string) bool {
	for i := 0; i < len(words); i++ {
		if word == words[i] {
			return true
		}
	}
	return false
}

// Loads words from dictionary file and returns as string slice
func loadWordsFromDictionaryFile(filename string) []string {
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
