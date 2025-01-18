# Ydays 

## Projet : Blog sur l'inovation

**Actuellement, pour le lancer il faut :**


- installer git ainsi que golang sur votre machine.
- ouvrir un powershell et se placer où vous le souhaitez pour mettre le dossier.
- ecrire la commande suivante : 
```powershell
git clone https://github.com/hSalinas-Marty/test-ydays.git
```
- Pour finir, il vous suffit d'aller dans le dossier précédement créé et de lancer les serveur :
```powershell
 cd .\test-ydays\ | go run main.go
```

Vous pouvez maintenant ouvrir votre navigateur et y rentrer l'url suivante : "http://localhost:8080"

**------------------------------------------------------------------------------**

### Il reste maintenant à faire correspondre les pages HTML, le css et rentrer les articles.

### Les articles doivent être ajouter dans le dossier "posts" en format markdown. 

Avec les infos qui s'afficheront sur la page d'accueil : le titre en premiere ligne, la description en 2eme précédée de "description: " et le lien de l'image précédé de "image: ". Vous trouverez un exemple dans le premier article : "article1.md"