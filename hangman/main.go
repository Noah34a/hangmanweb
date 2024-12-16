package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// Structure des données pour le jeu
type GameData struct {
	WordToGuess    string   // Mot à deviner
	Guesses        []rune   // Lettres devinées
	Attempts       int      // Tentatives restantes
	ImagePendu     string   // Chemin vers l'image actuelle
	GameOver       bool     // État de fin de jeu
	Message        string   // Message affiché (victoire, défaite)
	WordRevealed   string   // Mot complet pour défaite
}

var templates *template.Template
var words = map[string][]string{
	"facile":    {"banc","bureau","cabinet","carreau","chaise","classe","maison","coin","couloir","dossier","video","ecole","ecriture","entree","escalier","interieur"},
	"difficile": {"obsolescence","endogene","procrastination","exsangue","concomitance","peregrination","vicissitude","sagacite","ineffable","anachorete","cacophonie","hermeneutique"},
	"pays":      {"bresil","colombie","equateur","guyane","jordanie","lituanie","malawi","nepal","portugal","pakistan","somalie","croatie","france"},
	"marque":    {"rolex","balanciaga","prada","gucci","celio","jules","asics","chanel","casio","dior","armani","decathlon","azzaro"},
}

var game GameData

func main() {
	// Charger les templates HTML
	templates = template.Must(template.ParseGlob("templates/*.html"))

	// Servir les fichiers statiques (images et CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Définir les routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/guess", guessHandler)

	// Lancer le serveur
	fmt.Println("Serveur démarré sur http://localhost:3030")
	http.ListenAndServe(":3030", nil)
}

// Handler pour la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

// Handler pour démarrer une partie
func startHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	level := r.FormValue("level")
	wordsList := words[level]

	// Initialisation du jeu
	game.WordToGuess = strings.ToUpper(wordsList[0])
	game.Guesses = make([]rune, len(game.WordToGuess))
	for i := range game.Guesses {
		game.Guesses[i] = '_'
	}
	game.Attempts = 9
	game.GameOver = false
	game.ImagePendu = "/static/images/pendu0.png"
	game.Message = ""

	http.Redirect(w, r, "/game", http.StatusFound)
}

// Handler pour afficher le jeu
func gameHandler(w http.ResponseWriter, r *http.Request) {
	data := game
	data.ImagePendu = fmt.Sprintf("/static/images/pendu%d.png", 9-game.Attempts)
	templates.ExecuteTemplate(w, "game.html", data)
}

// Handler pour gérer les entrées du joueur
func guessHandler(w http.ResponseWriter, r *http.Request) {
	if game.GameOver {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	r.ParseForm()
	letter := strings.ToUpper(r.FormValue("letter"))
	if len(letter) != 1 {
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	correctGuess := false
	for i, char := range game.WordToGuess {
		if rune(letter[0]) == char {
			game.Guesses[i] = char
			correctGuess = true
		}
	}

	if !correctGuess {
		game.Attempts--
	}

	// Vérifier les conditions de fin
	if strings.Index(string(game.Guesses), "_") == -1 {
		game.GameOver = true
		game.Message = "Félicitations, vous avez gagné ! 🎉"
	} else if game.Attempts <= 0 {
		game.GameOver = true
		game.Message = "Désolé, vous avez perdu ! 😢"
		game.WordRevealed = game.WordToGuess
	}

	http.Redirect(w, r, "/game", http.StatusFound)
}
