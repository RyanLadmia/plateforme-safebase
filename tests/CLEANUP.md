# Nettoyage des Utilisateurs de Test

## Problème

Lors de l'exécution des tests E2E avec Cypress, de nombreux utilisateurs de test sont créés avec des emails se terminant par `@e2e.com`. Ces utilisateurs ne sont pas automatiquement supprimés et peuvent s'accumuler dans la base de données.

## Solution

Trois méthodes sont disponibles pour nettoyer les utilisateurs de test :

### 1. Nettoyage Automatique après les Tests (Recommandé)

Les tests Cypress sont configurés pour nettoyer automatiquement les utilisateurs de test après l'exécution de chaque fichier de test grâce au hook `after()` dans `e2E/support/e2e.ts`.

**Avantages :**
- Automatique, aucune action manuelle requise
- S'exécute après chaque fichier de test
- Maintient la base de données propre

**Comment ça fonctionne :**
```typescript
// Dans e2E/support/e2e.ts
after(() => {
  cy.cleanupTestUsers()
})
```

### 2. Script Shell Manuel

Utilisez le script `cleanup-test-users.sh` à la racine du projet pour nettoyer manuellement tous les utilisateurs de test.

**Usage :**
```bash
# À la racine du projet
./cleanup-test-users.sh
```

**Prérequis :**
- Le backend doit être en cours d'exécution (Docker ou local)
- Le backend doit être accessible sur `http://localhost:8080`

**Exemple de sortie :**
```
==================================
  Nettoyage des utilisateurs E2E
==================================
[OK] Backend accessible
[INFO] Suppression des utilisateurs de test (@e2e.com)...
[OK] 25 utilisateur(s) de test supprimé(s)
==================================
  Nettoyage terminé avec succès
==================================
```

### 3. Commande Cypress Manuelle

Dans vos tests ou dans la console Cypress, vous pouvez appeler manuellement :

```typescript
cy.cleanupTestUsers()
```

**Usage dans un test :**
```typescript
it('should cleanup test users', () => {
  cy.cleanupTestUsers()
})
```

## Endpoint API

Le nettoyage utilise l'endpoint backend suivant :

**Endpoint :** `POST /api/test/cleanup-users`

**Réponse :**
```json
{
  "message": "Test users cleaned up successfully",
  "deleted_count": 25
}
```

**Sécurité :**
- Cet endpoint n'est disponible qu'en environnement non-production
- En production (`GO_ENV=production`), l'endpoint retourne une erreur 403

**Critères de suppression :**
- Tous les utilisateurs dont l'email se termine par `@e2e.com` sont supprimés
- La suppression est définitive (hard delete)

## Configuration

### Backend

Le handler de nettoyage est défini dans :
- `backend/internal/handlers/test_handler.go`
- `backend/internal/routes/test_route.go`

### Frontend (Cypress)

Les commandes personnalisées sont définies dans :
- `tests/e2E/support/commands.ts` (commande `cleanupTestUsers`)
- `tests/e2E/support/index.d.ts` (types TypeScript)
- `tests/e2E/support/e2e.ts` (hook global `after()`)

## Troubleshooting

### Le script shell échoue

**Problème :** `Le backend ne répond pas`

**Solution :**
```bash
# Démarrer le backend avec Docker
docker-compose up -d backend

# Vérifier que le backend est accessible
curl http://localhost:8080/health
```

### L'endpoint retourne 403 Forbidden

**Problème :** `Test cleanup not allowed in production`

**Solution :** L'endpoint de nettoyage n'est pas disponible en production. Assurez-vous que `GO_ENV` n'est pas défini sur `production`.

### Cypress ne nettoie pas automatiquement

**Vérifiez :**
1. Le hook `after()` est bien présent dans `e2E/support/e2e.ts`
2. Le backend est accessible pendant l'exécution des tests
3. Les tests s'exécutent jusqu'à la fin (le hook `after()` ne s'exécute qu'après tous les tests d'un fichier)

## Bonnes Pratiques

1. **Exécuter le script de nettoyage** après chaque session de tests prolongée
2. **Utiliser des emails uniques** avec timestamp pour éviter les conflits :
   ```typescript
   const uniqueEmail = `test.${Date.now()}@e2e.com`
   ```
3. **Ne pas désactiver le nettoyage automatique** sauf pour le débogage
4. **Vérifier régulièrement** le nombre d'utilisateurs de test dans la base de données

## Commandes Utiles

```bash
# Nettoyer les utilisateurs de test
./cleanup-test-users.sh

# Vérifier le statut du backend
curl http://localhost:8080/health

# Lancer les tests avec nettoyage automatique
npm run test:e2e
```

## Intégration CI/CD

Le nettoyage automatique est intégré dans les tests Cypress, donc aucune configuration supplémentaire n'est nécessaire pour la CI/CD. Les tests nettoieront automatiquement après leur exécution.

Si vous souhaitez forcer un nettoyage avant les tests, ajoutez à votre pipeline CI/CD :

```yaml
# Exemple GitHub Actions
- name: Cleanup test users before tests
  run: ./cleanup-test-users.sh
  continue-on-error: true  # Ne pas échouer si aucun utilisateur à nettoyer
```

