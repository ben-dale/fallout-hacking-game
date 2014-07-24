package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldReturnAllSevenLetterLongWords(t *testing.T) {
	// Given
	words := []string{"Ben", "Dale", "cars", "glasses", "cards", "drivers", "called"}

	// When
	extractedWords := extractWordsOfLength(7, words)

	// Then
	assert.Equal(t, []string{"glasses", "drivers"}, extractedWords, "Error. Should have only extracted seven letter words.")
}

func TestShouldReturnAnyFiveWordsAtRandomFromGivenListOfWords(t *testing.T) {
	// Given
	words := []string{"Ben", "Dale", "cars", "glasses", "cards", "drivers", "called"}

	// When
	result := extractSubsetOfWordsAtRandom(5, words)

	// Then
	assert.Equal(t, 5, len(result), "Error. Should have returned five words.")
}

func TestShouldReturnTrueIfWordExistsInGivenSlice(t *testing.T) {
	// Given
	slice := []string{"hello", "world"}

	// When
	result := wordListContains("hello", slice)

	// Then
	assert.Equal(t, true, result, "Error. Should have returned true.")
}

func TestShouldReturnFalseIfWordDoesNotExistInGivenSlice(t *testing.T) {
	// Given
	slice := []string{"Hello", "World"}

	// When
	result := wordListContains("cards", slice)

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
	convertWordsToUpperCase(words)

	// Then
	assert.Equal(t, []string{"BEN", "DALE"}, words, "Error. Should have converted all words to uppercase.")
}

func TestShouldReturnFalseIfUserAttemptIsNotOneOfProvidedOptions(t *testing.T) {
	// Given
	attempt := "LAVEL"
	options := []string{"HELLO", "RUBIX", "CRISP", "CREEP", "BIRDS", "SHEEP"}

	// When
	result := isValidAttempt(attempt, options)

	// Then
	assert.Equal(t, false, result, "Error. Should have returened false as attempt is not in options.")
}
