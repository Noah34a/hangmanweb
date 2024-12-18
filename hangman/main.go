package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates *template.Template
var words = map[string][]string{
	"facile":    {"banane", "pomme", "chaise", "porte"},
	"difficile": {"procrastination", "anachronisme", "endogene", "cacophonie"},
	"pays":      {"bresil", "france", "japon", "canada"},
	"marque": {"azzaro","zara","rolex","celine","lanvin","omega",},
}

type GameData struct {
	Category     string
	Word         string
	MaskedWord   string
	AttemptsLeft int
	Message      string
	Guesses      string
	Image        string
}

var gameState = GameData{}
func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/game", gameHandler)

	log.Println("Serveur démarré sur http://localhost:3030")
	log.Fatal(http.ListenAndServe(":3030", nil))
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		category := r.FormValue("category")
		if category == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		gameState = GameData{
			Category:     category,
			Word:         randomWord(category),
			MaskedWord:   maskWord(randomWord(category)),
			AttemptsLeft: 9,
			Guesses:      "",
			Message:      "",
			Image:        "/static/image/pendu9.png",
		}
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "index.html", nil)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		guess := r.FormValue("letter")
		if len(guess) != 1 || strings.Contains(gameState.Guesses, guess) {
			gameState.Message = "Lettre invalide ou déjà utilisée."
		} else {
			gameState.Guesses += guess
			if strings.Contains(gameState.Word, guess) {
				gameState.MaskedWord = updateMaskedWord(gameState.Word, gameState.Guesses)
			} else {
				gameState.AttemptsLeft--
			}
		}
		gameState.Image = "static/image/pendu" + string('0'+rune(9-gameState.AttemptsLeft)) + ".png"
		if gameState.AttemptsLeft == 0 {
			gameState.Message = "Désolé, vous avez perdu ! Le mot était : " + gameState.Word
		} else if !strings.Contains(gameState.MaskedWord, "_") {
			gameState.Message = "Félicitations, vous avez gagné ! 🎉"
		}
	}
	templates.ExecuteTemplate(w, "game.html", gameState)
}
func randomWord(category string) string {
	wordsInCategory := words[category]
	return wordsInCategory[0] 
}
func maskWord(word string) string {
	return strings.Repeat("_ ", len(word))
}
func updateMaskedWord(word, guesses string) string {
	masked := ""
	for _, letter := range word {
		if strings.ContainsRune(guesses, letter) {
			masked += string(letter) + " "
		} else {
			masked += "_ "
		}
	}
	return masked
}
