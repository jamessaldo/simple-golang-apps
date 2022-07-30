# Authorization Service

## Objectives

1. Implement **Access** Control to authorize **User** access to Conversa
**Applications** and their **Endpoints** based on the user's **Role** within a
**Team**.
2. Implement **Applications** within the Authorization context, their **Team**
ownership, and the surrounding Console interaction.
3. Manage **Team** **Membership** and **Invitations**.

## Domain Model

```mermaid
classDiagram
    APIKey --> Application
    Application <-- Ownership
    Team <-- Ownership
    Membership --> Team
    Membership --> User
    Membership --> Role
    Access --> Endpoint
    Access --> Role
    Invitation --> Team
    Invitation --> Role

    class APIKey
    APIKey: +int id
    APIKey: +str key
    APIKey: +Application application

    class Application
    Application: +int id
    Application: +str name
    Application: +str description
    Application: +int console_app
    Application: +create()
    Application: +setOwner(Team)

    class Ownership
    Ownership: +int id
    Ownership: +Application application
    Ownership: +Team team
    Ownership: +create()

    class Team
    Team: +int id
    Team: +str name
    Team: +str description
    Team: +str avatar_path
    Team: +bool is_personal
    Team: +create()
    Team: +inviteMember(str, Role)
    Team: +addMember(User, Role)

    class User
    User: +int id
    User: +str email
    User: +str full_name
    User: +bool is_active
    User: +str avatar_path
    User: +deactivate()

    class Role
    Role: +int id
    Role: +str name
    Role: +create()

    class Membership
    Membership: +int id
    Membership: +User user
    Membership: +Team team
    Membership: +Role role
    Membership: +create()

    class Endpoint
    Endpoint: +id int
    Endpoint: +str path
    Endpoint: +str method
    Endpoint: +int version
    Endpoint: +create()

    class Access
    Access: +id int
    Access: +Endpoint endpoint
    Access: +Role role
    Access: +create()
    Access: +authorize()

    class Invitation
    Invitation: +int id
    Invitation: +Team team
    Invitation: +str email
    Invitation: +Role role
    Invitation: +datetime expires_at
    Invitation: +str status
    Invitation: +createJob()
```

## Use Cases

### Access Control

**Access** is defined as a mapping between a **Role** and an **Endpoint**. The
funcionality to create access mappings is not provided to users directly (yet).
Instead, these mappings should be defined when the service is initialized. Keep
in mind that functionality is distributed accross several services, and that the
Endpoints registered here reflect all these services.

```mermaid
sequenceDiagram
    actor Engineer
    Engineer ->> BE Client: Submit Mapping File
    BE Client ->> BE Server: Create Endpoint(s)
    BE Client ->> BE Server: Create Role(s)
    BE Client ->> BE Server: Create Access(')
```

Requests made to any backend services need to be authorized. Authorization can
occur at the boundary between frontend and backend using an API Gateway, *or* at
the frontend (server-side) by providing a complete list of access mapping for
each user for the frontend to check. All request *must* contain the
Application's APIKey and the User's access token in the header.

```mermaid
sequenceDiagram
    FE Server ->> Gateway: Any request
    Gateway ->> BE Auth: Authorize request
    BE Auth -->> Gateway: Return OK
    Gateway ->> BE Service: Pass request
    BE Service -->> FE Server: Return OK
```

### SSO Authentication

All user-facing operations within the system are preceeded by authentication
via SSO. This sequence is irrelevant to backend services, but is placed here for
convenience.

```mermaid
sequenceDiagram
    actor User
    User->>FE Client: Open application
    FE Client->>FE Server: Check session token
    FE Client->>Prosa SSO: Redirect to login page
    User->>Prosa SSO: Input credentials
    Prosa SSO-->>FE Client: Return access code
    Prosa SSO->>FE Client: Redirect to Application
    FE Client->>FE Server: Pass access code
    FE Server->>Prosa SSO: Get access/refresh token
    Prosa SSO-->>FE Server: Return access/refresh token
    FE Server->>FE Server: Generate session token
    FE Server-->>FE Client: Return session token
```

### First Login

During the first login ie. accessing the application for the first time after
registration, after logging in via SSO the User needs to be added to the authorization
service. A private Team is created with the new User as the owner.

```mermaid
sequenceDiagram
    actor User
    User->>FE Client: SSO Authentication
    FE Client->>FE Client: Redirect to Home
    FE Client->>FE Server: Get Home data
    FE Server->>FE Server: Decrypt access token
    FE Server->>BE Server: Get User, Teams, Applications
    BE Server-->>FE Server: 404 User not found
    FE Server->>BE Server: Create User
    BE Server->>BE Server: Create User
    BE Server->>BE Server: Create Team
    BE Server->>BE Server: Add User as Owner
    FE Server->>BE Server: Get User, Teams, Applications
    BE Server-->>FE Server: Return User, Teams, Applications
    FE Server->>FE Client: Return Home data
    FE Client-->>User: Render Home
```

Subsequent logins will return User, Teams, and Applications data without
encountering the 404. Subequent diagrams will assume the User is logged in and
at Home.

### Team Management

Teams are used to define the relationship between Applications and Users.
Applications are owned by Teams, not by Users. It is the User's Role within the
Team that determines which functionality (of the Team's Applications) the User
has Access to.

Team creation should be straightforward. Teams *must* have at least one Member,
and exactly one Member *must* have the Role of `owner`. A Team is said to be
`personal` when it can only have one member at most.

```mermaid
sequenceDiagram
    actor User
    User->>FE Client: Open Teams page
    FE Client->>FE Server: Submit form data
    FE Server->>BE Server: Create Team
    BE Server->>BE Server: Create Team
    BE Server-->>FE Server: Return OK
```

Members are added to Teams with a specified Role. If the User is already
registered within the system, they can be searched when inviting new members.
Regardless of whether the User to be added is registered within the system
(Prosa SSO), an email invitation is sent to the User. Following the invitation
link given in the email, after authentication, the User is explicitly prompted
to accept/deny the invitation. Invitations can be accepted, declined, expired,
and pending. Unacknowledged invitations expire after a specified amount of time.

```mermaid
sequenceDiagram
    actor User
    User->>FE Client: Open Team page
    FE Client->>FE Server: Search member
    FE Server->>BE Server: Invite member
    BE Server->>Worker: Create invitation job
    Worker->>Worker: Generate invitation
    Worker->>SMTP: Send email
    Worker->>Worker: Update invitation status
    BE Server-->>FE Server: Return OK
```

### Application Management

Applications are the primary method of storing Chatbots. The functionality of
Applications are distributed across multiple services. The Authorization Service
is responsible for interacting with the Console for the quota and payment-related
functionality of Applications.

In order to create Applications, the Team Owner must have at least one Payment
Profile. This is outside of the scope of this service, but is placed here for
convenience.

```mermaid
sequenceDiagram
    actor Team Owner
    Team Owner->>FE: Manage profiles
    FE->>Console: Add payment profile
    Console-->>FE: Return OK
    FE->>Console: Set default profile
    Console-->>FE: Return OK
```

Applications are created in the context of a Team, and are created using the
Console credentials of the Team Owner.

```mermaid
sequenceDiagram
    actor Team Admin
    Team Admin->>FE: Create Application
    FE->>BE Auth: Get profiles
    BE Auth->>BE Auth: Get Team Owner
    BE Auth->>SSO: Authenticate Owner
    SSO-->>BE Auth: Access Token
    BE Auth->>Console: Get profiles
    Console-->>BE Auth: Return profiles
    BE Auth-->>FE: Return profiles
    FE->>BE Auth: Create Application
    BE Auth->>Console: Create app
    Console-->>BE Auth: Return app
    BE Auth->>BE Auth: Create Application
    BE Auth->>BE Auth: Set Team Ownership
    BE Auth->>Queue: Application Creation Event
    BE Auth-->>FE: Return Application
    FE->>FE: Redirect to Home
    FE-->>Team Admin: Render Home
```
