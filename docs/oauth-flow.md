# MercadoLibre OAuth 2.0 Flow Documentation

## Overview

The Jopit ETL service uses OAuth 2.0 to securely connect with MercadoLibre's API on behalf of users. This document explains how the authentication flow works, how tokens are managed, and the security measures in place.

---

## Table of Contents

1. [What is OAuth 2.0?](#what-is-oauth-20)
2. [Why OAuth for MercadoLibre?](#why-oauth-for-mercadolibre)
3. [The Complete Flow](#the-complete-flow)
4. [Token Management](#token-management)
5. [Security Considerations](#security-considerations)
6. [Troubleshooting](#troubleshooting)
7. [Implementation Reference](#implementation-reference)

---

## What is OAuth 2.0?

OAuth 2.0 is an industry-standard authorization protocol that allows third-party applications (like Jopit) to access a user's data from another service (like MercadoLibre) without exposing the user's password.

### Key Concepts

**Authorization vs Authentication**:
- **Authentication**: Verifying who the user is (login)
- **Authorization**: Granting permission to access specific resources

OAuth 2.0 focuses on authorization - allowing Jopit to access a seller's MercadoLibre catalog without needing their MercadoLibre password.

**Key Players**:
1. **Resource Owner**: The MercadoLibre seller (user)
2. **Client**: The Jopit ETL service
3. **Authorization Server**: MercadoLibre's OAuth server
4. **Resource Server**: MercadoLibre's API (items, size charts, etc.)

**Token Types**:
- **Authorization Code**: Short-lived code exchanged for tokens
- **Access Token**: Used to make API requests (expires in ~6 hours)
- **Refresh Token**: Used to get new access tokens (long-lived)

---

## Why OAuth for MercadoLibre?

MercadoLibre requires OAuth 2.0 for several reasons:

1. **Security**: Users never share their MercadoLibre password with Jopit
2. **Scoped Access**: Jopit only gets permission for specific actions (read items)
3. **Revocable**: Users can revoke Jopit's access anytime via MercadoLibre settings
4. **Standard Protocol**: Industry-standard, well-documented, secure
5. **Token Expiration**: Automatic security through time-limited access tokens

---

## The Complete Flow

### Phase 1: Initial Connection

**User Action**: User clicks "Connect MercadoLibre" in Jopit app

```
┌────────────┐         ┌────────────┐         ┌──────────────┐
│   Jopit    │         │    ETL     │         │ MercadoLibre │
│  Frontend  │         │   Service  │         │              │
└─────┬──────┘         └─────┬──────┘         └──────┬───────┘
      │                      │                       │
      │  1. Request Auth URL │                       │
      │─────────────────────>│                       │
      │                      │                       │
      │                      │  2. Build URL with    │
      │                      │     - client_id       │
      │                      │     - redirect_uri    │
      │                      │     - response_type   │
      │                      │                       │
      │  3. Return Auth URL  │                       │
      │<─────────────────────│                       │
      │                      │                       │
      │  4. Redirect user to MercadoLibre           │
      │────────────────────────────────────────────>│
      │                      │                       │
```

**What Happens**:

1. **Frontend calls ETL Service**: `GET /mercadolibre/auth-url?user_id=X&shop_id=Y`

2. **ETL Service builds authorization URL** with parameters:
   - `client_id`: Jopit's MercadoLibre app ID
   - `response_type`: "code" (we want an authorization code)
   - `redirect_uri`: Where to send user after they approve
   - `state`: Optional parameter for CSRF protection

3. **Frontend receives URL** like:
   ```
   https://auth.mercadolibre.com/authorization?
     response_type=code&
     client_id=1234567890&
     redirect_uri=https://jopit.com/callback
   ```

4. **User is redirected to MercadoLibre** login page

**Implementation**: `src/main/domain/handlers/company_layout.go` - `GetAuthURL()`

---

### Phase 2: User Authorization

**User Action**: User logs into MercadoLibre and approves permissions

```
┌────────────┐                           ┌──────────────┐
│    User    │                           │ MercadoLibre │
└─────┬──────┘                           └──────┬───────┘
      │                                         │
      │  1. Enter MercadoLibre credentials     │
      │───────────────────────────────────────>│
      │                                         │
      │  2. View permission request             │
      │    "Jopit wants to:"                    │
      │    - Read your product catalog          │
      │    - Read size guide information        │
      │<───────────────────────────────────────│
      │                                         │
      │  3. User clicks "Allow"                 │
      │───────────────────────────────────────>│
      │                                         │
```

**What Happens**:

1. **User sees MercadoLibre login screen** (or is already logged in)

2. **MercadoLibre shows permission screen** listing what Jopit wants to access

3. **User approves** by clicking "Permitir" or "Allow"

4. **MercadoLibre generates authorization code** - a temporary code valid for ~10 minutes

**Security Note**: Jopit never sees the user's MercadoLibre password. The entire login happens on MercadoLibre's domain.

---

### Phase 3: Authorization Code Exchange

**Automatic Process**: MercadoLibre redirects user back to Jopit with code

```
┌────────────┐         ┌────────────┐         ┌──────────────┐
│   Jopit    │         │    ETL     │         │ MercadoLibre │
│  Frontend  │         │   Service  │         │    OAuth     │
└─────┬──────┘         └─────┬──────┘         └──────┬───────┘
      │                      │                       │
      │  1. Redirect with code                      │
      │<────────────────────────────────────────────│
      │    callback?code=TG-xxxxx                   │
      │                      │                       │
      │  2. Send code to ETL │                       │
      │─────────────────────>│                       │
      │    + user_id         │                       │
      │    + shop_id         │                       │
      │                      │                       │
      │                      │  3. Exchange code     │
      │                      │     POST /oauth/token │
      │                      │     - code            │
      │                      │     - client_id       │
      │                      │     - client_secret   │
      │                      │     - grant_type      │
      │                      │     - redirect_uri    │
      │                      │──────────────────────>│
      │                      │                       │
      │                      │  4. Return tokens     │
      │                      │<──────────────────────│
      │                      │     {                  │
      │                      │       access_token,    │
      │                      │       refresh_token,   │
      │                      │       expires_in       │
      │                      │     }                  │
      │                      │                       │
      │                      │  5. Save to MongoDB   │
      │                      │     (encrypted)       │
      │                      │                       │
      │  6. Success response │                       │
      │<─────────────────────│                       │
      │                      │                       │
```

**What Happens**:

1. **MercadoLibre redirects** user to Jopit's redirect URI with authorization code:
   ```
   https://jopit.com/callback?code=TG-abcdef123456
   ```

2. **Frontend captures code** from URL and sends to ETL Service:
   ```
   POST /mercadolibre/callback
   {
     "code": "TG-abcdef123456",
     "user_id": "firebase_uid",
     "shop_id": "jopit_shop_123"
   }
   ```

3. **ETL Service exchanges code for tokens** by calling MercadoLibre:
   ```
   POST https://api.mercadolibre.com/oauth/token
   {
     "grant_type": "authorization_code",
     "client_id": "your_app_id",
     "client_secret": "your_app_secret",
     "code": "TG-abcdef123456",
     "redirect_uri": "https://jopit.com/callback"
   }
   ```

4. **MercadoLibre returns tokens**:
   ```
   {
     "access_token": "APP_USR-123456-long-token",
     "token_type": "Bearer",
     "expires_in": 21600,
     "refresh_token": "TG-7890-refresh-token",
     "user_id": 789012345,
     "scope": "offline_access read write"
   }
   ```

5. **ETL Service stores credentials** in MongoDB `company_layout` collection:
   - User ID (Firebase)
   - Shop ID (Jopit)
   - Access token
   - Refresh token
   - Expiration timestamp
   - MercadoLibre user ID

6. **Success confirmation** returned to frontend

**Implementation**:
- Handler: `src/main/domain/handlers/company_layout.go` - `Callback()`
- Service: `src/main/domain/services/company_layout.go` - `SaveCredentials()`
- Repository: `src/main/domain/repositories/company_layout.go` - `SaveCompanyLayout()`

**Security Note**: The client secret is never exposed to the frontend - it stays secure on the backend.

---

### Phase 4: Using Access Tokens

**When**: Every time ETL service needs to call MercadoLibre API

```
┌────────────┐         ┌──────────────┐         ┌──────────────┐
│    ETL     │         │   MongoDB    │         │ MercadoLibre │
│  Service   │         │              │         │     API      │
└─────┬──────┘         └──────┬───────┘         └──────┬───────┘
      │                       │                        │
      │  1. Fetch credentials │                        │
      │──────────────────────>│                        │
      │                       │                        │
      │  2. Return tokens     │                        │
      │<──────────────────────│                        │
      │                       │                        │
      │  3. Check if token expiring soon (< 1 hour)   │
      │                       │                        │
      │  4. If yes, refresh token first (see Phase 5) │
      │                       │                        │
      │  5. Make API request with access token        │
      │───────────────────────────────────────────────>│
      │      Authorization: Bearer APP_USR-123456...  │
      │                       │                        │
      │  6. Return data       │                        │
      │<───────────────────────────────────────────────│
      │                       │                        │
```

**What Happens**:

1. **ETL Service retrieves credentials** from MongoDB using user_id

2. **Checks token expiration**:
   - If expires in >1 hour: Use existing access token
   - If expires in <1 hour: Refresh token first (Phase 5)

3. **Makes API request** with access token in Authorization header:
   ```
   GET https://api.mercadolibre.com/users/789012345/items/search
   Authorization: Bearer APP_USR-123456-long-token
   ```

4. **MercadoLibre validates token** and returns data

**Implementation**:
- Service: `src/main/domain/services/mercadolibre.go` - Various methods
- Client: `src/main/domain/clients/http.go` - Handles token injection and refresh
- Repository: `src/main/domain/repositories/company_layout.go` - `GetCredentialsByUserID()`

---

### Phase 5: Token Refresh (Automatic)

**When**: Access token expires or is about to expire

```
┌────────────┐         ┌────────────┐         ┌──────────────┐
│    ETL     │         │   MongoDB  │         │ MercadoLibre │
│  Service   │         │            │         │    OAuth     │
└─────┬──────┘         └─────┬──────┘         └──────┬───────┘
      │                      │                       │
      │  1. Check expiration │                       │
      │     (< 1 hour left)  │                       │
      │                      │                       │
      │  2. POST /oauth/token                        │
      │     grant_type: refresh_token                │
      │     client_id: xxx                           │
      │     client_secret: yyy                       │
      │     refresh_token: TG-7890...                │
      │─────────────────────────────────────────────>│
      │                      │                       │
      │  3. Return new tokens                        │
      │<─────────────────────────────────────────────│
      │     {                │                       │
      │       access_token: NEW_TOKEN,               │
      │       refresh_token: NEW_REFRESH,            │
      │       expires_in: 21600                      │
      │     }                │                       │
      │                      │                       │
      │  4. Update MongoDB   │                       │
      │─────────────────────>│                       │
      │     with new tokens  │                       │
      │                      │                       │
      │  5. Continue with original API call          │
      │                      │                       │
```

**What Happens**:

1. **ETL Service detects token expiring soon**:
   - Calculated: `expires_at - current_time < 1 hour`

2. **Calls MercadoLibre token endpoint**:
   ```
   POST https://api.mercadolibre.com/oauth/token
   {
     "grant_type": "refresh_token",
     "client_id": "your_app_id",
     "client_secret": "your_app_secret",
     "refresh_token": "TG-7890-refresh-token"
   }
   ```

3. **MercadoLibre returns new tokens**:
   - New access token (6 hours validity)
   - New refresh token (replaces old one)
   - New expiration time

4. **ETL Service updates database** with new credentials

5. **Original API request proceeds** with fresh access token

**Key Points**:
- **Proactive refresh**: Happens 1 hour before expiration (not after)
- **Automatic**: No user intervention required
- **Transparent**: User never knows tokens are being refreshed
- **New refresh token**: Each refresh gives a new refresh token (old one becomes invalid)

**Implementation**: `src/main/domain/clients/http.go` - Auto-refresh logic before each API call

---

## Token Management

### Storage Strategy

**Where**: MongoDB database, `company_layout` collection

**Document Structure** (simplified):
```
{
  _id: ObjectId,
  user_id: "firebase_uid",
  shop_id: "jopit_shop_id",
  user_id_meli: 789012345,
  access_token: "APP_USR-123456...",
  refresh_token: "TG-7890...",
  expires_at: ISODate("2026-02-10T21:00:00Z"),
  updated_at: ISODate("2026-02-10T15:00:00Z")
}
```

**Security**:
- MongoDB database uses encryption at rest
- Access tokens never exposed in frontend
- Credentials scoped by user_id (multi-tenant isolation)
- Only backend services can access credentials collection

### Token Lifecycle

**Access Token**:
- **Lifetime**: ~6 hours (21600 seconds)
- **Purpose**: Authenticate API requests
- **Renewal**: Automatically refreshed 1 hour before expiration
- **Invalidation**: Expires after time limit or when user revokes access

**Refresh Token**:
- **Lifetime**: Long-lived (months/years until revoked)
- **Purpose**: Get new access tokens without user re-authentication
- **Renewal**: Gets new refresh token each time it's used
- **Invalidation**: User revokes access via MercadoLibre settings

### Expiration Handling

**Proactive Refresh** (recommended approach - currently implemented):
- Check expiration before each API call
- If expires in <1 hour, refresh immediately
- User never experiences expired token errors

**Reactive Refresh** (alternative, not implemented):
- Try API call with current token
- If gets 401 Unauthorized, refresh and retry
- More API calls but simpler logic

### Token Revocation

**User-Initiated**:
1. User goes to MercadoLibre account settings
2. Finds "Connected Applications"
3. Revokes Jopit's access
4. All tokens immediately invalidated

**Effect on Jopit**:
- Next ETL attempt gets 401 Unauthorized
- User must reconnect MercadoLibre (full OAuth flow again)
- Previous ETL data remains in Jopit (not deleted)

**Developer-Initiated** (future feature):
- Delete credentials from MongoDB
- User must reconnect to use ETL again

---

## Security Considerations

### What's Secure

✅ **Password Protection**: User's MercadoLibre password never touches Jopit servers

✅ **Token Storage**: Tokens stored in backend database, never in frontend localStorage

✅ **Client Secret Protection**: Secret never exposed to frontend or users

✅ **HTTPS**: All OAuth communication over encrypted connections

✅ **Automatic Refresh**: Tokens refreshed proactively to minimize exposure

✅ **User Revocation**: Users can revoke access anytime via MercadoLibre

✅ **Scoped Access**: Jopit only gets permissions user explicitly approves

✅ **Multi-Tenant Isolation**: Each user's tokens isolated by user_id

### Potential Improvements

⚠️ **Token Encryption**: Currently MongoDB handles encryption at rest, could add application-level encryption

⚠️ **Audit Logging**: Log all token usage for security audits

⚠️ **Rate Limiting**: Prevent abuse by limiting OAuth attempts

⚠️ **CSRF Protection**: Add state parameter validation in OAuth callback

⚠️ **Token Rotation**: Implement periodic forced token rotation

---

## Troubleshooting

### Common Issues

#### Issue: "Invalid authorization code"

**Cause**: Authorization code expired or already used

**Solution**: 
- Codes expire after ~10 minutes
- Each code can only be used once
- User needs to restart OAuth flow

#### Issue: "Invalid refresh token"

**Causes**:
- User revoked access via MercadoLibre
- Refresh token expired (rare)
- Client secret changed in MercadoLibre app settings

**Solution**:
- User must reconnect MercadoLibre (complete OAuth flow again)
- Check if client secret is correct

#### Issue: "Redirect URI mismatch"

**Cause**: Redirect URI in request doesn't match registered URI in MercadoLibre app settings

**Solution**:
- Verify `MERCADO_LIBRE_REDIRECT_URI` environment variable
- Check MercadoLibre app settings
- URLs must match exactly (including https vs http, trailing slash)

#### Issue: "Token not found in database"

**Cause**: Credentials never saved or user_id/shop_id mismatch

**Solution**:
- Check that callback endpoint saved credentials successfully
- Verify user_id matches between OAuth flow and ETL request
- Check MongoDB connection and collection name

#### Issue: "Access token expired" (despite auto-refresh)

**Causes**:
- Clock skew between servers
- Token expiration calculation error
- Refresh failed silently

**Solution**:
- Check server time synchronization (NTP)
- Review logs for refresh attempts
- Verify refresh token is valid

---

## Implementation Reference

### Key Files

**Handlers** (HTTP endpoints):
- `src/main/domain/handlers/company_layout.go`
  - `GetAuthURL()` - Phase 1: Generate authorization URL
  - `Callback()` - Phase 3: Exchange code for tokens

**Services** (Business logic):
- `src/main/domain/services/company_layout.go`
  - `SaveCredentials()` - Store tokens in database
- `src/main/domain/services/mercadolibre.go`
  - Methods that use tokens to call MercadoLibre API

**Clients** (External communication):
- `src/main/domain/clients/http.go`
  - Auto-refresh logic
  - Token injection in requests
  - HTTP communication with MercadoLibre

**Repositories** (Data access):
- `src/main/domain/repositories/company_layout.go`
  - `SaveCompanyLayout()` - Persist credentials
  - `GetCredentialsByUserID()` - Retrieve credentials
  - `UpdateCredentials()` - Update tokens after refresh

**Models**:
- `src/main/domain/models/company_layout.go`
  - `CompanyLayout` struct with token fields

### Configuration

**Environment Variables Required**:
```bash
MERCADO_LIBRE_CLIENT_ID=1234567890
MERCADO_LIBRE_CLIENT_SECRET=your_secret_here
MERCADO_LIBRE_REDIRECT_URI=https://your-domain.com/callback
```

**MercadoLibre App Setup**:
1. Create app at https://developers.mercadolibre.com
2. Configure redirect URIs
3. Request necessary scopes (offline_access, read)
4. Get client ID and secret

---

## Flow Summary (Plain English)

Here's the complete OAuth flow in simple terms:

1. **User clicks "Connect MercadoLibre"** in Jopit app

2. **Jopit asks MercadoLibre for permission** by redirecting user to MercadoLibre's login page

3. **User logs into MercadoLibre** and sees what Jopit wants to access

4. **User clicks "Allow"** to give Jopit permission

5. **MercadoLibre sends a temporary code** back to Jopit

6. **Jopit exchanges the code for long-term tokens** (access token and refresh token)

7. **Jopit stores tokens securely** in its database

8. **When running ETL**, Jopit uses the access token to fetch products from MercadoLibre

9. **When token expires**, Jopit automatically gets a new one using the refresh token

10. **User never deals with tokens** - it all happens in the background

**Bottom line**: The user authorizes once, and Jopit handles everything else automatically.

---

**Last Updated**: February 10, 2026  
**OAuth Version**: OAuth 2.0 (Authorization Code Grant)  
**MercadoLibre API Documentation**: https://developers.mercadolibre.com/en_us/authentication-and-authorization
