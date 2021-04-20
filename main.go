package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

const (
	start              = "Start game"
	end                = "End game"
	apiUrl             = "https://opentdb.com"
	minQuestionsNumber = 1
	maxQuestionsNumber = 10
	questionType       = "multiple"
	encodeType         = "url3986"
	selectSize         = 10
)

var difficultyList = []string{"Easy", "Medium", "Hard", "Random"}

type Categories struct {
	TriviaCategories []CategoryItem `json:"trivia_categories"`
}

type CategoryItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Questions struct {
	ResponseCode int            `json:"response_code"`
	Results      []QuestionItem `json:"results"`
}

type QuestionItem struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

var clear map[string]func() 									// create a map for storing clear functions

func init() {
	clear = make(map[string]func()) 							// Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") 							// Linux example, its tested
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") 				// Windows example, its tested
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func clearScreen() {
	value, ok := clear[runtime.GOOS] 							// runtime.GOOS -> linux, windows, darwin etc.
	if ok {                         							// if we defined a clear func for that platform:
		value() 					 							// we execute it
	} else { 						 							// unsupported platform
		panic("OS is unsupported! Can't clear terminal screen")
	}
}

// unescape from rfc 3986 encoding
func unescapeString(escapedString string) string {
	unescapedAnswer, err := url.QueryUnescape(escapedString)
	if err != nil {
		log.Fatal(err)
	}
	return unescapedAnswer
}

func getCategoriesList() Categories {
	u, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	u.Path = path.Join(u.Path, "api_category.php")

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var categoriesList Categories
	if err = json.Unmarshal(body, &categoriesList); err != nil {
		log.Fatal(err)
	}

	return categoriesList
}

func getQuestionsList(difficultySelect string, questionsNumber string, categorySelect string, categoriesJson Categories) Questions {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.URL.Path = path.Join(req.URL.Path, "/api.php")
	q := req.URL.Query()

	if difficultySelect != "Random" {
		q.Add("difficulty", strings.ToLower(difficultySelect))
	}

	if categorySelect != "Random" {
		var categoryIndex int
		for _, v := range categoriesJson.TriviaCategories {
			if v.Name == categorySelect {
				categoryIndex = v.Id
			}
		}
		q.Add("category", strconv.Itoa(categoryIndex))
	}
	q.Add("amount", questionsNumber)
	q.Add("type", questionType)
	q.Add("encode", encodeType)

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	var questions Questions
	if err = json.Unmarshal(body, &questions); err != nil {
		log.Fatal(err)
	}
	return questions
}

func promptSelect(label string, itemList []string) string {
	clearScreen()
	promptSelect := promptui.Select{
		Label: label,
		Items: itemList,
		Size:  selectSize,
	}

	_, userAnswer, err := promptSelect.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	return userAnswer
}

func promptEnter(label string, validator promptui.ValidateFunc) string {
	clearScreen()
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validator,
	}

	questionsNumber, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	return questionsNumber
}

func main() {
	for {
		mainSelect := promptSelect("Welcome to Trivia app", []string{start, end})

		if mainSelect == start {
			// Select difficulty
			difficultySelect := promptSelect("Select difficulty", difficultyList)

			// Get categories list from API and select category
			var categoriesJson = getCategoriesList()
			var categoriesList []string
			categoriesList = append(categoriesList, "Random")
			for _, v := range categoriesJson.TriviaCategories {
				categoriesList = append(categoriesList, v.Name)
			}

			categorySelect := promptSelect("Select category: ", categoriesList)

			// Select number of question
			validate := func(input string) error {
				inputInt, err := strconv.ParseInt(input, 10, 64)
				if err != nil {
					return errors.New("invalid number")
				}
				if inputInt < minQuestionsNumber || inputInt > maxQuestionsNumber {
					return errors.New("number should be between " + strconv.Itoa(minQuestionsNumber) + " and " + strconv.Itoa(maxQuestionsNumber))
				}
				return nil
			}

			questionsNumber := promptEnter("Enter number of questions", validate)

			// Get questions list
			var questions = getQuestionsList(difficultySelect, questionsNumber, categorySelect, categoriesJson)
			// Start of game
			var successAnswersCount int

			for _, v := range questions.Results {
				// Create slice with answer options
				var answerOptions []string
				for _, v := range v.IncorrectAnswers {
					answerOptions = append(answerOptions, unescapeString(v))
				}
				answerOptions = append(answerOptions, unescapeString(v.CorrectAnswer))

				// Shuffling slice with answer options
				rand.Shuffle(len(answerOptions), func(i, j int) {
					answerOptions[i], answerOptions[j] = answerOptions[j], answerOptions[i]
				})

				userAnswer := promptSelect(unescapeString(v.Question), answerOptions)

				clearScreen()
				fmt.Println("Your answer: ", userAnswer)
				if userAnswer != unescapeString(v.CorrectAnswer) {
					fmt.Printf("You are wrong, correct answer is: %s\n\n", unescapeString(v.CorrectAnswer))
					fmt.Print("Press 'Enter' to continue...")
					if _, err := bufio.NewReader(os.Stdin).ReadBytes('\n'); err != nil {
						log.Fatal(err)
					}
				} else {
					fmt.Printf("You are right!\n\n")
					fmt.Print("Press 'Enter' to continue...")
					if _, err := bufio.NewReader(os.Stdin).ReadBytes('\n'); err != nil {
						log.Fatal(err)
					}
					successAnswersCount += 1
				}
			}

			clearScreen()
			if floatNumber, err := strconv.ParseFloat(questionsNumber, 64); err == nil {
				fmt.Printf("Success rate is %0.2f %%\n", float64(successAnswersCount)/floatNumber*100)
				fmt.Printf("%d from %d correct answers\n", successAnswersCount, int(floatNumber))
				fmt.Print("Press 'Enter' to continue...")
				if _, err = bufio.NewReader(os.Stdin).ReadBytes('\n'); err != nil {
					log.Fatal(err)
				}
			}
		} else if mainSelect == end {
			os.Exit(0)
		}
	}
}
