// Test de la validation frontend pour les URLs
console.log("Test de validation frontend des URLs de base de données");

// Simuler la fonction validateDatabaseData
function validateDatabaseData(data) {
    if (!data.name || data.name.trim() === '') {
        throw new Error('Le nom de la base de données est requis');
    }

    if (!data.type || !['mysql', 'postgresql'].includes(data.type)) {
        throw new Error('Le type de base de données doit être mysql ou postgresql');
    }

    // Si une URL est fournie, les champs individuels ne sont pas requis
    if (data.url && data.url.trim() !== '') {
        // Validation basique de l'URL
        if (!data.url.includes('://')) {
            throw new Error('L\'URL doit être au format: mysql://user:pass@host:port/db ou postgresql://user:pass@host:port/db');
        }
        console.log("✅ URL valide, champs individuels non requis");
        return; // Pas besoin de valider les champs individuels
    }

    // Si pas d'URL, tous les champs individuels sont requis
    console.log("❌ Pas d'URL, vérification des champs individuels");
    if (!data.host || data.host.trim() === '') {
        throw new Error('L\'hôte est requis');
    }

    if (!data.port || data.port.trim() === '') {
        throw new Error('Le port est requis');
    }

    if (!data.username || data.username.trim() === '') {
        throw new Error('Le nom d\'utilisateur est requis');
    }

    if (!data.db_name || data.db_name.trim() === '') {
        throw new Error('Le nom de la base de données est requis');
    }

    // Validation du port (doit être un nombre)
    const portNum = parseInt(data.port);
    if (isNaN(portNum) || portNum < 1 || portNum > 65535) {
        throw new Error('Le port doit être un nombre entre 1 et 65535');
    }
}

// Test cases
const testCases = [
    {
        name: "URL valide MySQL",
        data: {
            name: "Test DB",
            type: "mysql",
            url: "mysql://root:password@localhost:8889/mydb"
        },
        expected: "success"
    },
    {
        name: "URL valide PostgreSQL",
        data: {
            name: "Test DB",
            type: "postgresql",
            url: "postgresql://user:pass@remotehost:5432/prod_db"
        },
        expected: "success"
    },
    {
        name: "URL invalide (pas de protocol)",
        data: {
            name: "Test DB",
            type: "mysql",
            url: "root:password@localhost:8889/mydb"
        },
        expected: "error"
    },
    {
        name: "Champs individuels sans URL",
        data: {
            name: "Test DB",
            type: "mysql",
            host: "localhost",
            port: "8889",
            username: "root",
            password: "secret",
            db_name: "mydb"
        },
        expected: "success"
    },
    {
        name: "URL vide + champs manquants",
        data: {
            name: "Test DB",
            type: "mysql",
            url: "",
            host: "",
            port: "8889"
        },
        expected: "error"
    }
];

console.log("\n=== TESTS DE VALIDATION ===\n");

testCases.forEach((testCase, index) => {
    console.log(`Test ${index + 1}: ${testCase.name}`);
    console.log(`Données:`, testCase.data);

    try {
        validateDatabaseData(testCase.data);
        if (testCase.expected === "success") {
            console.log("✅ Test RÉUSSI\n");
        } else {
            console.log("❌ Test ÉCHOUÉ - Devrait échouer mais a réussi\n");
        }
    } catch (error) {
        if (testCase.expected === "error") {
            console.log(`✅ Test RÉUSSI - Erreur attendue: ${error.message}\n`);
        } else {
            console.log(`❌ Test ÉCHOUÉ - Erreur inattendue: ${error.message}\n`);
        }
    }
});

console.log("=== FIN DES TESTS ===");
