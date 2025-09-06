package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{
	Timeout: 5 * time.Second, // par ex. 5 secondes
}

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LinkChecker - Vérificateur de Liens</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <script>
        tailwind.config = {
            darkMode: 'class',
            theme: {
                extend: {
                    colors: {
                        dark: {
                            100: '#161b22',
                            200: '#0d1117',
                            300: '#010409',
                        },
                        github: {
                            gray: '#21262d',
                            border: '#30363d',
                            text: '#c9d1d9',
                            blue: '#58a6ff',
                            green: '#3fb950',
                            red: '#f85149'
                        }
                    }
                }
            }
        }
    </script>
    <style>
        body {
            font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
        }
        .terminal-window {
            border: 1px solid #30363d;
            border-radius: 6px;
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
        }
        .terminal-header {
            background-color: #21262d;
            border-bottom: 1px solid #30363d;
            padding: 8px 12px;
            display: flex;
            align-items: center;
        }
        .terminal-button {
            width: 12px;
            height: 12px;
            border-radius: 50%;
            margin-right: 6px;
        }
        .terminal-red { background-color: #f85149; }
        .terminal-yellow { background-color: #d29922; }
        .terminal-green { background-color: #3fb950; }
        .glow {
            box-shadow: 0 0 10px rgba(88, 166, 255, 0.7);
        }
        .code-block {
            background-color: #0d1117;
            border: 1px solid #30363d;
            border-radius: 6px;
            padding: 16px;
            overflow-x: auto;
        }
        .link-item {
            position: relative;
            transition: all 0.2s ease;
        }
        .link-item:hover {
            background-color: #1c2129;
        }
        .link-item:before {
            content: ">";
            position: absolute;
            left: -20px;
            color: #58a6ff;
            opacity: 0;
            transition: opacity 0.2s ease;
        }
        .link-item:hover:before {
            opacity: 1;
        }
    </style>
</head>
<body class="bg-dark-300 text-github-text min-h-screen">
    <div class="container mx-auto px-4 py-8 max-w-4xl">
        <header class="flex items-center justify-between mb-10">
            <div class="flex items-center">
                <div class="w-10 h-10 rounded-md bg-github-blue flex items-center justify-center mr-3 glow">
                    <i class="fas fa-link text-white"></i>
                </div>
                <h1 class="text-2xl font-semibold">LinkChecker<span class="text-github-blue">.dev</span></h1>
            </div>
            <div class="flex space-x-2">
                <div class="w-3 h-3 rounded-full bg-github-red"></div>
                <div class="w-3 h-3 rounded-full bg-github-yellow"></div>
                <div class="w-3 h-3 rounded-full bg-github-green"></div>
            </div>
        </header>

        <div class="mb-8">
            <h2 class="text-xl mb-4">$ Vérifiez les liens morts sur votre site</h2>
            <p class="text-gray-400 mb-6">Entrez l'URL de votre site web pour analyser tous les liens et détecter ceux qui sont brisés.</p>
            
            <div class="terminal-window bg-dark-200">
                <div class="terminal-header">
                    <div class="terminal-button terminal-red"></div>
                    <div class="terminal-button terminal-yellow"></div>
                    <div class="terminal-button terminal-green"></div>
                    <span class="text-sm ml-2">bash - vérification de liens</span>
                </div>
                <div class="p-6">
                    <form method="POST" action="/">
                        <div class="flex mb-4">
                            <span class="inline-flex items-center px-3 rounded-l-md border border-r-0 border-github-border bg-dark-100 text-gray-400">
                                https://
                            </span>
                            <input 
                                type="text" 
                                name="url" 
                                placeholder="votresite.com" 
                                class="flex-1 min-w-0 block w-full px-3 py-2 rounded-none border-github-border bg-dark-100 text-white focus:ring-github-blue focus:border-github-blue"
                                required
                            >
                            <button 
                                type="submit" 
                                class="inline-flex items-center px-4 rounded-r-md border border-l-0 border-github-border bg-github-blue text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-github-blue"
                            >
                                <i class="fas fa-terminal mr-2"></i> Exécuter
                            </button>
                        </div>
                    </form>
                    <div class="text-sm text-gray-400 flex items-center">
                        <i class="fas fa-info-circle mr-2"></i> Appuyez sur Entrée ou cliquez sur Exécuter pour lancer l'analyse
                    </div>
                </div>
            </div>
        </div>

        {{if .}}
        <div class="terminal-window bg-dark-200 mt-10">
            <div class="terminal-header">
                <div class="terminal-button terminal-red"></div>
                <div class="terminal-button terminal-yellow"></div>
                <div class="terminal-button terminal-green"></div>
                <span class="text-sm ml-2">résultats de l'analyse</span>
            </div>
            
            <div class="p-6">
                <div class="flex justify-between items-center mb-4">
                    <h3 class="text-lg font-medium">
                        <i class="fas fa-file-code text-github-blue mr-2"></i>Liens morts détectés
                    </h3>
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-github-red text-white">
                        {{len .}} erreur(s)
                    </span>
                </div>
                
                <div class="code-block mt-4">
                    <ul class="space-y-3">
                        {{range .}}
                        <li class="link-item pl-6 py-2">
                            <div class="flex items-start">
                                <span class="text-github-red mr-2">→</span>
                                <span class="text-github-red break-all">{{.}}</span>
                            </div>
                            <div class="text-xs text-gray-500 mt-1 ml-4">
                                <i class="fas fa-clock mr-1"></i>Status: 404 - Lien non trouvé
                            </div>
                        </li>
                        {{end}}
                    </ul>
                </div>
                
                <div class="mt-6 pt-4 border-t border-github-border flex justify-between">
                    <div class="text-sm text-gray-400">
                        <i class="fas fa-history mr-1"></i> Analyse terminée
                    </div>
                    <button onclick="window.location.href='/'" class="inline-flex items-center px-3 py-1 border border-github-border rounded text-sm text-white bg-dark-100 hover:bg-dark-300">
                        <i class="fas fa-redo mr-1"></i> Nouvelle analyse
                    </button>
                </div>
            </div>
        </div>
        {{else}}
        <div class="terminal-window bg-dark-200 mt-10">
            <div class="terminal-header">
                <div class="terminal-button terminal-red"></div>
                <div class="terminal-button terminal-yellow"></div>
                <div class="terminal-button terminal-green"></div>
                <span class="text-sm ml-2">comment utiliser</span>
            </div>
            <div class="p-6">
                <div class="code-block">
                    <pre class="text-sm text-gray-300">
<span class="text-github-green"># Exemple d'utilisation:</span>
<span class="text-github-blue">$</span> curl -X POST -d "url=https://votresite.com" https://linkchecker.dev/

<span class="text-github-green"># Format de sortie:</span>
<span class="text-github-blue">></span> [ERROR] https://votresite.com/lien-cassé - 404 Not Found
<span class="text-github-blue">></span> [SUCCESS] https://votresite.com/page-valide - 200 OK

<span class="text-github-green"># Options disponibles:</span>
<span class="text-github-blue">--depth</span>       Niveau de profondeur de l'analyse (par défaut: 2)
<span class="text-github-blue">--timeout</span>     Délai d'attente pour chaque requête (par défaut: 5s)
<span class="text-github-blue">--workers</span>     Nombre de workers simultanés (par défaut: 10)</pre>
                </div>
            </div>
        </div>
        {{end}}

        <footer class="mt-12 pt-8 border-t border-github-border text-center text-gray-500 text-sm">
            <p>Développé avec <i class="fas fa-heart text-github-red"></i> pour les développeurs - LinkChecker v1.0</p>
            <p class="mt-2">
                <a href="#" class="text-github-blue hover:underline"><i class="fab fa-github mr-1"></i>Code source</a> • 
                <a href="#" class="text-github-blue hover:underline">Documentation API</a> • 
                <a href="#" class="text-github-blue hover:underline">Statut du service</a>
            </p>
        </footer>
    </div>

    <script>
        // Animation simple pour le texte de terminal
        document.addEventListener('DOMContentLoaded', function() {
            const terminalText = document.querySelector('h2.text-xl');
            if (terminalText && !{{if .}}true{{else}}false{{end}}) {
                const originalText = terminalText.textContent;
                terminalText.textContent = '';
                
                let i = 0;
                function typeWriter() {
                    if (i < originalText.length) {
                        terminalText.textContent += originalText.charAt(i);
                        i++;
                        setTimeout(typeWriter, 50);
                    }
                }
                
                setTimeout(typeWriter, 500);
            }
        });
    </script>
</body>
</html>
`))

// Handler HTTP
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		link := r.FormValue("url")
		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			tmpl.Execute(w, []string{"URL invalide - doit commencer par http:// ou https://"})
			return
		}

		resp, ok := linkTest(link)
		if !ok {
			tmpl.Execute(w, []string{"Lien inaccessible"})
			return
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			tmpl.Execute(w, []string{"Impossible de parser la page"})
			return
		}

		links := findLinks(doc)
		deadLinks := checkAllLinks(link, links)
		if len(deadLinks) == 0 {
			deadLinks = []string{"Aucun lien mort trouvé ✅"}
		}

		err = tmpl.Execute(w, deadLinks)
		if err != nil {
			return
		}
	} else {
		tmpl.Execute(w, nil)
	}
}

// Lancement du serveur web
func main() {
	http.HandleFunc("/", handler)
	log.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func linkTest(link string) (*http.Response, bool) {
	resp, err := client.Get(link)
	if err != nil || resp.StatusCode != 200 {
		return nil, false
	}
	return resp, true
}

var rateLimiter = make(chan struct{}, 10)

func isLinkAlive(link string) bool {

	if link == "" || link[0] == '#' ||
		strings.HasPrefix(link, "mailto:") ||
		strings.HasPrefix(link, "javascript:") ||
		strings.HasPrefix(link, "tel:") {
		return true
	}

	rateLimiter <- struct{}{}
	defer func() { <-rateLimiter }()

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "LinkChecker/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 400
}

func findLinks(doc *goquery.Document) []string {
	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})
	return links
}

func makeAbsolute(base, href string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return href
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		return href
	}
	return baseURL.ResolveReference(hrefURL).String()
}

func checkAllLinks(baseLink string, links []string) []string {
	var deadLinks []string
	deadChan := make(chan string)
	var wg sync.WaitGroup

	for _, l := range links {
		fullLink := makeAbsolute(baseLink, l)
		wg.Add(1)
		go func(fl string) {
			defer wg.Done()
			if !isLinkAlive(fl) {
				deadChan <- fl
			}
		}(fullLink)
	}

	go func() {
		wg.Wait()
		close(deadChan)
	}()

	for dl := range deadChan {
		deadLinks = append(deadLinks, dl)
	}

	return deadLinks
}
