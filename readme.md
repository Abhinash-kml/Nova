#Nova
A multiplayer backend solution designed from the ground up.
Features -
* Users
* Authentication
  * OAuth 2.0
  * Social Provider
  * Device ID
* Posts
* Comments
* Clans
* Realtime
  * Lobby
  * Chat
  * Channels
  * Notifications
  * Presence
  
**Stack**
* Golang
* Gin
* Postgres
* Redis
* Websocket

```mermaid
graph LR
    A[Nova] --> B[Modules]
    B --> C[Users]
    C --> AA[Auth]
    AA --> Oauth2
    AA --> Social
    AA --> Device
    C --> Entitlements
    C --> Posts
    C --> Comments
    C --> Clans
    B --> D[Realtime]
    D --> Lobby
    D --> Channels
    D --> Notifications
    D --> Presence
    D --> Chat
    B --> E[Economy]
    E --> Store
    E --> Currency
    E --> Wallet
    B --> I[Inventory]
    I --> Items
    I --> Storage
```

**User login flow**

```mermaid
sequenceDiagram
autonumber
    Client->>Auth: POST /login(credentials)
    Auth->>Ban: Check if banned
    Ban-->>Auth: Banned response
    Auth-->>Client: Banned Response
    Auth->>Database: Fetch userinfo(not banned)
    Database-->>Auth: userinfo response
    Auth-->Client: 200 OK (JWT Token)
    Note over Client, Database: One single request
```

**User ban flow**
```mermaid
sequenceDiagram
autonumber
    Control Pane->>Database: Ban initiated, persist in db
    Database-->>Control Pane: Ban details persisted
    Control Pane->>Redis Session Store: Fetch user edge server info
    Redis Session Store-->>Control Pane: connected edge server info
    Control Pane->>Edge Server: Send ban details message
    Edge Server->>Client: Send ban details and disconnect websocket session
```
**Realtime Message Routing Flow**

```mermaid
sequenceDiagram
autonumber
    Client->>Edge Server: Send message to edge server
    Edge Server->>Kafka/RabbitMQ: Send to message queue
    Kafka/RabbitMQ->>Database: Persist message in database
    Kafka/RabbitMQ-->>Edge Server: Send Acknowledgement
    Edge Server->>Target Client: Target exists on same edge server
    Edge Server->>Pub/Sub: Target doesn't exist on same edge server
    Pub/Sub->>Edge Server: Route to Target's edge server
    Edge Server->>Target Client: Target client receives message
```

**Inventory creation flow**
```mermaid
sequenceDiagram
autonumber
    Client->>Inventory Service: Fetch inventory
    Inventory Service->>Entitlement Service: Fetch active entitlements
    Entitlement Service-->>Inventory Service: active entitlements response
    Inventory Service-->>Client: 200 OK (response)
```

**Presence update flow**
```mermaid
sequenceDiagram
autonumber
    Client->>Edge Server: Update presence status
    Edge Server-->>Client: Response OK
    Edge Server->>Presence Service: Update state
    Edge Server->>Redis Presence Store: Fetch subscribed users
    Redis Presence Store-->>Presence Service: response list
    Presence Service->>Edge server: Fan out presence update message to subscribed users
```
