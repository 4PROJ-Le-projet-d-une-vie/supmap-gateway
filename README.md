# supmap-gateway

Une gateway API légère et performante écrite en Go, permettant de router et de rediriger les requêtes vers différents microservices.

## Fonctionnalités

- Reverse proxy pour plusieurs microservices :
    - Service Utilisateurs (authentification, gestion des comptes)
    - Service Incidents (gestion et suivi des incidents)
    - Service GIS (géocodage et services géographiques)
    - Service Navigation (WebSocket pour la navigation en temps réel)
- Configuration simple via variables d'environnement
- Support des WebSockets
- Gestion des routes avec réécriture d'URL

## Architecture

La gateway agit comme un point d'entrée unique pour différents microservices :

```mermaid
flowchart TD
    Client[Client] --> Gateway[Gateway API]
    
    subgraph Services
        Gateway --> Users[Service Users]
        Gateway --> Incidents[Service Incidents]
        Gateway --> GIS[Service GIS]
        Gateway --> Navigation[Service Navigation]
    end

    subgraph Endpoints Users
        Users --> |/login| Auth[Authentification]
        Users --> |/logout| Logout[Déconnexion]
        Users --> |/register| Register[Inscription]
        Users --> |/refresh| Token[Refresh Token]
        Users --> |/users| UsersMgmt[Gestion Utilisateurs]
    end

    subgraph Endpoints Incidents
        Incidents --> |/incidents| IncidentsCRUD[CRUD Incidents]
        Incidents --> |/incidents/types| Types[Types d'incidents]
        Incidents --> |/incidents/interactions| Inter[Interactions]
        Incidents --> |/incidents/me/history| History[Historique]
    end

    subgraph Endpoints GIS
        GIS --> |/address| Address[Recherche Adresses]
        GIS --> |/geocode| Geocode[Géocodage]
        GIS --> |/route| Route[Calcul Routes]
    end

    subgraph Endpoints Navigation
        Navigation --> |/navigation/ws| WS[WebSocket Navigation]
    end
```

## Configuration

La configuration se fait via des variables d'environnement :

| Variable               | Description                  |
|------------------------|------------------------------|
| SUPMAP_GATEWAY_PORT    | Port d'écoute de la gateway  |
| SUPMAP_USERS_HOST      | Hôte du service utilisateurs |
| SUPMAP_USERS_PORT      | Port du service utilisateurs |
| SUPMAP_INCIDENTS_HOST  | Hôte du service incidents    |
| SUPMAP_INCIDENTS_PORT  | Port du service incidents    |
| SUPMAP_GIS_HOST        | Hôte du service GIS          |
| SUPMAP_GIS_PORT        | Port du service GIS          |
| SUPMAP_NAVIGATION_HOST | Hôte du service navigation   |
| SUPMAP_NAVIGATION_PORT | Port du service navigation   |

Points d'accès (Endpoints)
- Service Utilisateurs
  - /login - Authentification
  - /logout - Déconnexion
  - /register - Création de compte
  - /refresh - Rafraîchissement du token
  - /users - Gestion des utilisateurs
- Service Incidents
  - /incidents - CRUD des incidents
  - /incidents/types - Types d'incidents
  - /incidents/interactions - Interactions sur les incidents
  - /incidents/me/history - Historique personnel
- Service GIS
  - /address - Recherche d'adresses
  - /geocode - Géocodage
  - /route - Calcul d'itinéraires
- Service Navigation
  - /navigation/ws - WebSocket pour la navigation en temps réel

## Démarrage
Configurez les variables d'environnement nécessaires
Lancez l'application :

```sh
go run .
```

## Développement

Le projet utilise :
- net/http/httputil pour le reverse proxy
- github.com/caarlos0/env/v11 pour la gestion de la configuration
- Go modules pour la gestion des dépendances