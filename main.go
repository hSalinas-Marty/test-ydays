package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

// Fonction pour afficher la liste des articles disponibles
func handler(w http.ResponseWriter, r *http.Request) {
	// Si l'URL est "/"
	if r.URL.Path == "/" {
		// Lire les fichiers dans le dossier "posts"
		files, err := ioutil.ReadDir("posts")
		if err != nil {
			http.Error(w, "Impossible de lire les articles", http.StatusInternalServerError)
			return
		}

		// Filtrer les fichiers Markdown
		var articles []string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".md") {
				articles = append(articles, strings.TrimSuffix(file.Name(), ".md"))
			}
		}

		// Charger le template de la page d'accueil
		tmpl, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			log.Fatal("Erreur lors du chargement du template d'accueil : ", err)
		}

		// Préparer les données pour le template
		homeData := struct {
			Title    string
			Articles []string
		}{
			Title:    "Accueil",
			Articles: articles, // Liste des articles disponibles
		}

		// Exécuter le template avec les données
		err = tmpl.Execute(w, homeData)
		if err != nil {
			log.Fatal("Erreur lors de l'exécution du template : ", err)
		}
		return
	}

	// Pour les autres URL, traiter les articles
	filePath := filepath.Join("posts", r.URL.Path+".md")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Convertir le Markdown en HTML
	htmlContent := blackfriday.Run(content)

	// Extraire le titre de l'article (première ligne)
	lines := strings.Split(string(content), "\n")
	title := strings.TrimSpace(lines[0])

	// Charger le template pour les articles
	tmpl, err := template.ParseFiles("./templates/article.html")
	if err != nil {
		log.Fatal("Erreur lors du chargement du template d'article : ", err)
	}

	// Préparer les données pour le template de l'article
	article := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: template.HTML(htmlContent),
	}

	// Exécuter le template pour afficher l'article
	err = tmpl.Execute(w, article)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution du template : ", err)
	}
}

func main() {
	// Servir les fichiers statiques (CSS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Gérer les requêtes pour l'accueil et les articles
	http.HandleFunc("/", handler)

	// Lancer le serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
