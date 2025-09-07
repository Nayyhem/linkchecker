# LinkChecker

**LinkChecker** est un outil web pour détecter facilement les liens morts sur n’importe quel site. Il analyse une page web, vérifie tous les liens et affiche les résultats avec une interface moderne basée sur Tailwind CSS.

---

## 🚀 Fonctionnalités

* Vérifie les liens morts (404, 500…) sur une page web.
* Gère les liens relatifs et absolus.
* Interface web responsive et moderne.
* Indique clairement le nombre de liens morts et les URLs concernées.
* Prise en charge de la limitation de requêtes simultanées pour éviter la surcharge.

---

## 🧰 Technologies utilisées

* **Go** pour le backend et la logique de vérification.
* **GoQuery** pour parser le HTML.
* **net/http** pour les requêtes HTTP.
* **HTML / Tailwind CSS** pour le frontend.
* **Font Awesome** pour les icônes.

---

## 📂 Structure du projet

```text
LinkChecker/
├─ main.go                  # Point d'entrée du serveur web
├─ checker/                 # Package pour la vérification des liens
│   └─ checker.go           # Fonctions: isLinkAlive, checkAllLinks, makeAbsolute, findLinks, CheckLinkPage
├─ templates/
│   └─ index.html           # Template HTML avec Tailwind CSS
└─ go.mod                   # Fichier de module Go
```

---

## ⚡ Installation

1. **Cloner le projet**

```bash
git clone https://github.com/votre-utilisateur/linkchecker.git
cd linkchecker
```

2. **Installer les dépendances**

```bash
go mod tidy
```

3. **Lancer le serveur**

```bash
go run main.go
```

4. **Accéder à l’interface**

Ouvrez votre navigateur et allez sur [http://localhost:8080](http://localhost:8080)

---

## 🧪 Liens de test

Pour vérifier le fonctionnement du checker, vous pouvez utiliser ces liens :

* `https://example.com` → Accessible
* `https://httpbin.org/status/200` → Accessible
* `https://httpbin.org/status/404` → Lien mort
* `https://httpbin.org/status/500` → Lien mort
* `https://www.wikipedia.org` → Peut poser problème selon restrictions du site

---

## 📝 Utilisation

1. Entrez l’URL complète de votre site web dans la barre de recherche (ex : `https://example.com`).
2. Cliquez sur **Exécuter** ou appuyez sur **Entrée**.
3. Les liens morts détectés apparaîtront dans la section « Résultats de l’analyse ».
4. Si aucun lien mort n’est trouvé, un message de succès s’affiche.

---

## 📦 Personnalisation

Vous pouvez modifier :

* Le nombre de workers simultanés pour limiter la charge (`rateLimiter` dans `checker.go`).
* Le délai d’attente pour les requêtes HTTP (`Timeout` dans `main.go`).
* Le style et l’interface via le template HTML et Tailwind CSS.

---

## ⚠️ Limitations

* Certains sites (Wikipedia, Cloudflare, etc.) peuvent bloquer les requêtes automatisées.
* Les liens nécessitant une authentification ne sont pas pris en charge.
* Les liens JavaScript (`javascript:`) et mailto sont ignorés.

---

## 📝 Licence

MIT License © 2025 Nayyhem
