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

// Article structure pour contenir les métadonnées des articles
type Article struct {
	Slug        string
	Title       string
	Description string
	ImageURL    string
}

// Fonction pour analyser les fichiers Markdown et extraire les métadonnées
func parseMarkdown(filePath string) (Article, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Article{}, err
	}

	lines := strings.Split(string(content), "\n")
	title := strings.TrimPrefix(lines[0], "# ")

	var description, imageURL string
	for _, line := range lines {
		if strings.HasPrefix(line, "description:") {
			description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		}
		if strings.HasPrefix(line, "image:") {
			imageURL = strings.TrimSpace(strings.TrimPrefix(line, "image:"))
		}
	}

	slug := strings.TrimSuffix(filepath.Base(filePath), ".md")
	return Article{
		Slug:        slug,
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
	}, nil
}

// Gestionnaire principal
func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/contact" || r.URL.Path == "/about" {
		http.NotFound(w, r)
		return
	}
	if r.URL.Path == "/" {
		// Lire les fichiers dans le dossier "posts"
		files, err := ioutil.ReadDir("posts")
		if err != nil {
			http.Error(w, "Impossible de lire les articles", http.StatusInternalServerError)
			return
		}

		// Analyser les fichiers Markdown pour récupérer les articles
		var articles []Article
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".md") {
				article, err := parseMarkdown(filepath.Join("posts", file.Name()))
				if err == nil {
					articles = append(articles, article)
				}
			}
		}

		// Charger et exécuter le template pour la page d'accueil
		// Définir l'article à la une (par exemple, le premier article)
		var featuredArticle *Article
		if len(articles) > 0 {
			featuredArticle = &articles[0]
			articles = articles[1:] // Supprimez l'article à la une de la liste générale
		}

		// Charger et exécuter le template pour la page d'accueil
		tmpl, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			log.Fatal("Erreur lors du chargement du template d'accueil : ", err)
		}

		homeData := struct {
			Title           string
			Articles        []Article
			FeaturedArticle *Article
		}{
			Title:           "Nos Articles : ",
			Articles:        articles,
			FeaturedArticle: featuredArticle,
		}

		err = tmpl.Execute(w, homeData)
		if err != nil {
			log.Fatal("Erreur lors de l'exécution du template : ", err)
		}

		return
	}

	// Gestion des articles individuels
	filePath := filepath.Join("posts", r.URL.Path+".md")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Convertir le Markdown en HTML avec blackfriday
	htmlContent := blackfriday.Run(content)

	// Extraire le titre du Markdown
	lines := strings.Split(string(content), "\n")
	title := strings.TrimPrefix(lines[0], "# ")

	tmpl, err := template.ParseFiles("./templates/article.html")
	if err != nil {
		log.Fatal("Erreur lors du chargement du template d'article : ", err)
	}

	articleData := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: template.HTML(htmlContent),
	}

	err = tmpl.Execute(w, articleData)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution du template : ", err)
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/contact.html")
	if err != nil {
		log.Fatal("Erreur lors du chargement du template de contact : ", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution du template : ", err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/about.html")
	if err != nil {
		log.Fatal("Erreur lors du chargement du template à propos : ", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution du template : ", err)
	}
}

func main() {
	// Servir les fichiers statiques (CSS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Gérer les requêtes
	http.HandleFunc("/", handler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/about", aboutHandler)

	// Lancer le serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
