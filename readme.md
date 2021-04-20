# Quiz application, which based on the Trivia API

## Description

Simple application, which communicates with [the Trivia API](https://opentdb.com/api_config.php), for get database of question  
and check user answers

## Prerequisites

Installed go 1.14+ or Docker

## Used external libs

[promptui](https://github.com/manifoldco/promptui) - for provide a user-friendly console interface

## How to run
If you have preinstalled go 1.14:
* `go mod dowload`
* `go run main.go`

In case of using Docker:
* `docker build -t trivia-app:1.0 .`
* `docker run -it trivia-app:1.0`

## Example of work

```
Use the arrow keys to navigate: ↓ ↑ → ←
? Welcome to Trivia app:
▸ Start game
End game  
-----
? Select difficulty:
▸ Easy
Medium
Hard
Random
-----
? Select category:
Random
General Knowledge
Entertainment: Books
Entertainment: Film
Entertainment: Music
Entertainment: Musicals & Theatres
Entertainment: Television
Entertainment: Video Games
Entertainment: Board Games
Entertainment: Math
▸ Science: Mathematics
↓   Science & Nature
-----
✗ Enter number of questions: 3
-----
? In the hexadecimal system, what number comes after 9?: 
  ▸ The Letter A
    10
    16
    The Number 0
-----
Your answer:  The Letter A
You are right!

Press 'Enter' to continue...
-----
? How many books are in Euclid's Elements of Geometry?: 
  ▸ 13
    8
    17
    10
-----
Your answer:  13
You are right!

Press 'Enter' to continue...
-----
? What is the Roman numeral for 500?: 
    C
    L
  ▸ D
    X
-----
Your answer:  L
You are wrong, correct answer is: D
-----
Success rate is 66.67 %
2 from 3 correct answers
Press 'Enter' to continue...

```
