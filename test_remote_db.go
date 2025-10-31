package main

import (
	"fmt"
	"net"
	"time"
)

// Test de connectivité réseau pour diagnostiquer les problèmes de bases distantes
func testNetworkConnectivity(host, port string) error {
	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return fmt.Errorf("impossible de se connecter à %s:%s: %v", host, port, err)
	}
	conn.Close()
	return nil
}

func main() {
	fmt.Println("🩺 Outil de diagnostic pour bases de données distantes")
	fmt.Println("==================================================")

	// Test MySQL typique
	fmt.Println("\n🔍 Test MySQL distant (port 3306):")
	testHost := "votre-serveur-mysql.com"  // Remplacez par votre serveur
	testPort := "3306"

	if err := testNetworkConnectivity(testHost, testPort); err != nil {
		fmt.Printf("❌ Échec: %v\n", err)
		fmt.Println("💡 Causes possibles:")
		fmt.Println("   - Serveur non accessible (firewall, IP bloquée)")
		fmt.Println("   - Port incorrect ou fermé")
		fmt.Println("   - Nom d'hôte incorrect")
	} else {
		fmt.Printf("✅ Connexion réseau réussie à %s:%s\n", testHost, testPort)
		fmt.Println("💡 Si mysqldump échoue encore, vérifier:")
		fmt.Println("   - Identifiants de connexion (user/password)")
		fmt.Println("   - Droits utilisateur sur la base de données")
		fmt.Println("   - Base de données existe")
	}

	// Test PostgreSQL typique
	fmt.Println("\n🔍 Test PostgreSQL distant (port 5432):")
	testHost = "votre-serveur-postgres.com"  // Remplacez par votre serveur
	testPort = "5432"

	if err := testNetworkConnectivity(testHost, testPort); err != nil {
		fmt.Printf("❌ Échec: %v\n", err)
		fmt.Println("💡 Causes possibles:")
		fmt.Println("   - Serveur non accessible (firewall, IP bloquée)")
		fmt.Println("   - Port incorrect ou fermé")
		fmt.Println("   - Nom d'hôte incorrect")
	} else {
		fmt.Printf("✅ Connexion réseau réussie à %s:%s\n", testHost, testPort)
		fmt.Println("💡 Si pg_dump échoue encore, vérifier:")
		fmt.Println("   - Identifiants de connexion (user/password)")
		fmt.Println("   - Droits utilisateur sur la base de données")
		fmt.Println("   - Base de données existe")
		fmt.Println("   - SSL requis par le serveur")
	}

	fmt.Println("\n📋 Checklist de dépannage:")
	fmt.Println("1. Vérifier que le serveur distant accepte les connexions depuis votre IP")
	fmt.Println("2. Confirmer que les identifiants (user/password) sont corrects")
	fmt.Println("3. S'assurer que l'utilisateur a les droits nécessaires sur la base")
	fmt.Println("4. Vérifier que la base de données existe")
	fmt.Println("5. Pour MySQL: vérifier si SSL est requis")
	fmt.Println("6. Pour PostgreSQL: vérifier la configuration pg_hba.conf")

	fmt.Println("\n🔧 Commandes de test manuelles:")

	fmt.Println("\nMySQL:")
	fmt.Println("mysql -h VOTRE_SERVEUR -P 3306 -u VOTRE_USER -p")
	fmt.Println("mysqldump -h VOTRE_SERVEUR -P 3306 -u VOTRE_USER -p VOTRE_DB --single-transaction > test.sql")

	fmt.Println("\nPostgreSQL:")
	fmt.Println("psql -h VOTRE_SERVEUR -p 5432 -U VOTRE_USER -d VOTRE_DB")
	fmt.Println("pg_dump -h VOTRE_SERVEUR -p 5432 -U VOTRE_USER -d VOTRE_DB > test.sql")
}
