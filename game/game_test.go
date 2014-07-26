package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldReturnAllSevenLetterLongWords(t *testing.T) {
	// Given
	words := []string{"Ben", "Dale", "cars", "glasses", "cards", "drivers", "called"}

	// When
	extractedWords := extractStringsOfLength(7, words)

	// Then
	assert.Equal(t, []string{"glasses", "drivers"}, extractedWords, "Error. Should have only extracted seven letter words.")
}

func TestShouldReturnAnyFiveWordsAtRandomFromGivenListOfWords(t *testing.T) {
	// Given
	words := []string{"Ben", "Dale", "cars", "glasses", "cards", "drivers", "called"}

	// When
	result := extractSubsetOfStringsAtRandom(5, words)

	// Then
	assert.Equal(t, 5, len(result), "Error. Should have returned five words.")
}

func TestShouldReturnTrueIfWordExistsInGivenSlice(t *testing.T) {
	// Given
	slice := []string{"hello", "world"}

	// When
	result := stringSliceContains("hello", slice)

	// Then
	assert.Equal(t, true, result, "Error. Should have returned true.")
}

func TestShouldReturnFalseIfWordDoesNotExistInGivenSlice(t *testing.T) {
	// Given
	slice := []string{"Hello", "World"}

	// When
	result := stringSliceContains("cards", slice)

	// Then
	assert.Equal(t, false, result, "Error. Should have returned false.")
}

func TestShouldReturnNumberOfLettersInCorrectPlace(t *testing.T) {
	// Given
	attempt := "hello"
	password := "heyho"

	// When
	result := calculateNumberOfCorrectLetters(attempt, password)

	// Then
	assert.Equal(t, 3, result, "Error. Three of the letters are correct. Expected three.")
}

func TestShouldReturnZeroWhenNoLettersAreInCorrectPlace(t *testing.T) {
	// Given
	attempt := "hello"
	password := "crisp"

	// When
	result := calculateNumberOfCorrectLetters(attempt, password)

	// Then
	assert.Equal(t, 0, result, "Error. No letters are correct. Expected zero.")
}

func TestShouldConvertAllLowercaseLettersInCollectionOfWordsToUpperCase(t *testing.T) {
	// Given
	words := []string{"ben", "Dale"}

	// When
	convertStringsToUpperCase(words)

	// Then
	assert.Equal(t, []string{"BEN", "DALE"}, words, "Error. Should have converted all words to uppercase.")
}

func TestShouldReturnFalseIfUserAttemptIsNotOneOfProvidedOptions(t *testing.T) {
	// Given
	attempt := "LAVEL"
	possiblePasswords := []string{"HELLO", "RUBIX", "CRISP", "CREEP", "BIRDS", "SHEEP"}

	// When
	result := isValidAttempt(attempt, possiblePasswords)

	// Then
	assert.Equal(t, false, result, "Error. Should have returened false as attempt is not one of the options.")
}

func TestShouldReturnTrueIfUserAttemptIsOneOfProvidedOptions(t *testing.T) {
	// Given
	attempt := "CAPCOM"
	possiblePasswords := []string{"MARVEL", "CAPCOM"}

	// When
	result := isValidAttempt(attempt, possiblePasswords)

	// Then
	assert.Equal(t, true, result, "Error. Should have returned true as attempt is one of the options.")
}

func TestShouldBuildRoundWithFourAttemptsWithFivePasswordsWhichAreSevenCharactersLongWithOneCorrectPassword(t *testing.T) {
	// Given
	wordList := []string{"carpool", "ballboy", "primark", "footman", "vanhire", "yellow", "red", "green", "blue"}
	attemptsPerRound := 4
	lengthOfPasswords := 7
	numberOfPasswords := 5

	// When
	round := buildRound(attemptsPerRound, lengthOfPasswords, numberOfPasswords, wordList)

	// Then
	assert.Equal(t, 4, round.attemptsLeft, "Error. Should have returned round with 4 attempts left.")
	assert.NotEmpty(t, round.correctWord, "Error. Should have returned round with chosen correct password.")
	assert.Equal(t, 5, len(round.possiblePasswords), "Error. Should have returned round with five passwords.")
	for i := 0; i < len(round.possiblePasswords); i++ {
		assert.Equal(t, 7, len(round.possiblePasswords[i]), "Error. Possible password is not correct length.")
	}
}
