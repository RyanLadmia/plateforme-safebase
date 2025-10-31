package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// Simule le parsing d'URL MySQL
func parseMySQLURL(dbURL string) (host, port, user, pass, dbname string, err error) {
	if !strings.HasPrefix(dbURL, "mysql://") {
		return "", "", "", "", "", fmt.Errorf("URL doit commencer par mysql://")
	}

	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("URL invalide: %v", err)
	}

	host, port, err = net.SplitHostPort(parsedURL.Host)
	if err != nil {
		// Si pas de port, utiliser le défaut
		host = parsedURL.Host
		port = "3306"
	}

	user = parsedURL.User.Username()
	pass, _ = parsedURL.User.Password()
	dbname = strings.TrimPrefix(parsedURL.Path, "/")

	return host, port, user, pass, dbname, nil
}

func main() {
	fmt.Println("🔍 Débogage du parsing d'URL MySQL")
	fmt.Println("===================================")

	// Test avec l'URL qui pose problème
	testURL := "mysql://localhost:3306"

	fmt.Printf("URL testée: %s\n", testURL)

	host, port, user, pass, dbname, err := parseMySQLURL(testURL)
	if err != nil {
		fmt.Printf("❌ Erreur de parsing: %v\n", err)
	} else {
		fmt.Printf("✅ Parsing réussi:\n")
		fmt.Printf("   Host: '%s'\n", host)
		fmt.Printf("   Port: '%s'\n", port)
		fmt.Printf("   User: '%s'\n", user)
		fmt.Printf("   Pass: '%s'\n", pass)
		fmt.Printf("   DB: '%s'\n", dbname)
	}

	// Test avec une URL complète
	fullURL := "mysql://ryan-ladmia:motdepasse@localhost:3306/ryan-ladmia_cinetech"
	fmt.Printf("\nURL complète testée: %s\n", fullURL)

	host2, port2, user2, pass2, dbname2, err2 := parseMySQLURL(fullURL)
	if err2 != nil {
		fmt.Printf("❌ Erreur de parsing: %v\n", err2)
	} else {
		fmt.Printf("✅ Parsing réussi:\n")
		fmt.Printf("   Host: '%s'\n", host2)
		fmt.Printf("   Port: '%s'\n", port2)
		fmt.Printf("   User: '%s'\n", user2)
		fmt.Printf("   Pass: '%s'\n", pass2)
		fmt.Printf("   DB: '%s'\n", dbname2)
	}

	fmt.Println("\n📋 Diagnostic:")
	fmt.Println("Si mysqldump voit 'mysql://localhost:3306' comme hostname,")
	fmt.Println("cela signifie que l'URL complète est stockée dans le champ Host")
	fmt.Println("au lieu d'être parsée et distribuée dans les bons champs.")
}
