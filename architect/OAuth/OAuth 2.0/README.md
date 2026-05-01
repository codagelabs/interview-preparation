# OAuth (Open Authorization)

## What is OAuth?
OAuth (Open Authorization) is an open-standard authorization framework or protocol that provides applications the ability for "secure designated access." For example, it allows you to tell a service like Facebook that it's okay for a different service like ESPN.com to access your profile or post updates to your timeline without having to give ESPN your Facebook password. This minimizes risk in a major way: in the event ESPN suffers a breach, your Facebook password remains safe.

OAuth doesn't share password data, but instead uses authorization tokens to prove an identity between consumers and service providers. 

## What is OAuth 1.0?
OAuth 1.0 is the original version of the OAuth protocol, released in 2007. It provided a secure way to authorize access to APIs without sharing passwords.

**Key Characteristics:**
- Highly secure but complex to implement.
- Required cryptographic signatures for every single API request.
- Based on symmetric cryptography (shared secrets).
- Mostly designed for traditional server-side web applications, making it difficult to scale and use with mobile apps or single-page applications (SPAs).
- Required the client to perform complex signature generation (e.g., HMAC-SHA1).

## What is OAuth 2.0?
OAuth 2.0 is the complete rewrite of OAuth 1.0, released in 2012. It is currently the industry-standard protocol for authorization. It was designed to address the shortcomings of OAuth 1.0 by simplifying the developer experience.

**Key Differences from 1.0:**
- **No cryptographic signatures required:** Relies on HTTPS/TLS encryption for security instead of complex request signatures.
- **Support for more use cases:** Provides different "grant types" (flows) tailored for different types of applications (web, mobile, SPAs, Smart TVs/IoT).
- **Bearer Tokens:** Uses "Bearer tokens", which are significantly easier to use (just pass them in the HTTP Authorizaton header).
- **Separation of roles:** Clearly separates the roles of the Authorization Server (which authenticates and issues tokens) and Resource Server (which hosts the APIs).

## Terminology
To understand OAuth 2.0, you need to understand these key terms (roles & components):
- **Resource Owner:** The user who authorizes an application to access their account/data.
- **Client:** The application (web, mobile, SPA, etc.) that wants to access the user's account.
- **Authorization Server:** The server that authenticates the Resource Owner and issues access tokens to the Client.
- **Resource Server:** The server hosting the user's protected data/APIs. It accepts access tokens to grant access.
- **Access Token:** A credential (usually a string, often a JWT) that the Client uses to access the Resource Server. It has a limited lifespan (usually short).
- **Refresh Token:** A credential issued alongside the Access Token. When the Access Token expires, the Client can use the Refresh Token to securely get a new Access Token without asking the user to log in again.
- **Scope:** Used to specify the level or boundaries of access that the application is requesting. For example, `read:profile` or `write:documents`. This limits what the application can do with the token.
- **Authorization Code (Code):** A temporary, short-lived code that the Client exchanges for an Access Token securely behind the scenes (used specifically in the Authorization Code Flow).

## Architecture
OAuth 2.0 involves four main entities:
1. **Resource Owner** (User)
2. **Client** (Application)
3. **Authorization Server** (e.g., Auth0, Okta, Google Login, Keycloak)
4. **Resource Server** (e.g., Your Backend API)

```mermaid
graph LR
    User(Resource Owner) <--> Client[Client Application]
    Client <--> AuthServer[Authorization Server]
    Client <--> ResServer[Resource Server]
    AuthServer -.Trust.-> ResServer
```

## Flow: Authorization Code Grant (Sequence)
The Authorization Code Flow is the most common and secure OAuth 2.0 flow. It is heavily used by server-side web applications and mobile/SPA apps (often with PKCE extension).

```mermaid
sequenceDiagram
    participant User as Resource Owner (Browser)
    participant Client as Client Application
    participant Auth as Authorization Server
    participant Resource as Resource Server

    Note over User,Auth: 1. Authorization Request
    User->>Client: Clicks "Login/Connect"
    Client->>User: Redirects to Authorization Server
    User->>Auth: User authenticates and grants permission (Consent)
    
    Note over User,Client: 2. Authorization Grant
    Auth->>User: Redirects back to Client with an Authorization Code
    User->>Client: Delivers Authorization Code to Client callback URI
    
    Note over Client,Auth: 3. Token Exchange (Backend)
    Client->>Auth: Sends Auth Code + Client ID + Client Secret
    Auth->>Client: Validates and returns Access Token (+ Optional Refresh Token)
    
    Note over Client,Resource: 4. Resource Access
    Client->>Resource: API Request with Access Token in Header
    Resource->>Client: Returns Protected Data
```

## APIs -- Request and Response Example

Below is a pure HTTP request/response illustration for the Authorization Code Flow.

### 1. Authorization Request (Browser Redirect)
The client redirects the user's browser to the Authorization Server.
**Request:**
```http
GET /authorize?response_type=code&client_id=YOUR_CLIENT_ID&redirect_uri=https://client.example.com/callback&scope=read&state=xyz123 HTTP/1.1
Host: server.example.com
```

### 2. Authorization Response (Browser Redirect)
If the user grants access, the Authorization Server redirects the browser back to the client's `redirect_uri` with the short-lived `code`.
**Response:**
```http
HTTP/1.1 302 Found
Location: https://client.example.com/callback?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz123
```

### 3. Token Request (Backend to Backend)
The Client exchanges the `code` for an `access_token`. This is an out-of-band POST request from the Client's backend server to the Authorization server (invisible to the browser).
**Request:**
```http
POST /token HTTP/1.1
Host: server.example.com
Content-Type: application/x-www-form-urlencoded
Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW

grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA&redirect_uri=https://client.example.com/callback
```
*(Note: `Authorization: Basic` is base64 encoded `client_id:client_secret`)*

### 4. Token Response
The Authorization Server responds with the tokens.
**Response:**
```http
HTTP/1.1 200 OK
Content-Type: application/json;charset=UTF-8
Cache-Control: no-store
Pragma: no-cache

{
  "access_token": "2YotnFZFEjr1zCsicMWpAA",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "tGzv3JOkF0XG5Qx2TlKWIA",
  "scope": "read"
}
```

### 5. API Request Using Access Token
The Client makes a request to the Resource Server (API) passing the `access_token` it received.
**Request:**
```http
GET /api/user/profile HTTP/1.1
Host: resource.example.com
Authorization: Bearer 2YotnFZFEjr1zCsicMWpAA
```

**Response:**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": 12345,
  "name": "John Doe",
  "email": "johndoe@example.com"
}
```
