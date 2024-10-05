
# HTTP Server
This is a static site server from a Boot.dev course. It implements a fake website called 'Chirpy', meant to be a clone of X nee twitter.

# Setup
## Database
You should download and install [PostGreSQL](https://www.postgresql.org/download/).
Create a 'chirpy' database.

## Environment
Most of the configuration values are stored in the `.env` file in the root of the project. 
Example .env file
```
# Connection string to your chirpy database 'sslmode-disable' is required
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" 

# dev allows for calls to the reset api endpoint to clear the database
PLATFORM="dev" 

# a random string used to sign the JWTs
JWT_KEY="W1VTHsNO76huHODrd3laYzPjK+k11Bc/XJGgQkXRYv7VujrfYjFokNTFFo4DSHF5xzc6VcOtKpPAUgbggxSqQosXmFPhEN15M+313eNseQfr76LXezZrteOIwHtOgZ3cQCcSvNSodRQoeJToeSiRicBXhHD6vLlcg9WzO8uQvK3I+g1sJlc64ViLMBOKtULok6JS6y9GDOxbC9H1+rIxwRj470pbfCF4mhekShqXApDQYuzVeUel332h1/kFrgzQ gVX4+LOvQkP6sJA0hlIt97Wy9QUQq2HvOSzUVqgNVUCZiAgS6E5tKS1XgzJUG7cO/aXGtURCZ7gRxi/1nk5zVw=="

# A random string used to authenticate with fictional Polka payment processor
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
```

# API

	POST /api/chirps 
    Create a new Chirp.
    Headers: Authorization - The JWT as a bearer token
    Query Params: None
    Body: application/json
    ``` 
    {
        "body": "Chirp content"
    }
    ```

	GET /api/chirps
    Get all chirps from the database.
    Headers: None
    Query Params: sort "asc"|"desc", author_id uuid (optional)
    Body: None

	GET /api/chirps/{chirpID}
    Get a single chirp from the database.
    Headers: None
    Query Params: None
    Body: None

	DELETE /api/chirps/{chirpID}
    Delete a single chirp from the DB. Only the chirp creater can delete a chirp
    Headers: Authorization - The JWT as a bearer token
    Query Params: None
    Body: None

	POST /api/users
    Create a new User
    Headers: None
    Query Params: None
    Body: application/json
    ``` 
    {
        "email": "user@sample.com",
        "password": "abc123"
    }
    ```

	PUT /api/users
    Update the user's email and password, only the user themselves can update their information
    Headers: Authorization - The JWT as a bearer token
    Query Params: None
    Body: application/json
    ``` 
    {
        "email": "user@sample.com",
        "password": "abc123"
    }
    ```

	POST /api/login
    Login a user and get a JWT and Refresh Token
    Headers: None
    Query Params: None
    Body: application/json
    ``` 
    {
        "email": "user@sample.com",
        "password": "abc123"
    }
    ```
    Response Body: application/json
    ```
    {
        "id": <uuid>
        "created_at": timestamp
        "updated_at": timestamp
        "is_chirpy_red" bool
        "token": string
        "refresh_token" : string
    }
    ```

	POST /api/refresh
    Get a new JWT token
    Headers: Authorization - The refresh token as an ApiKey token
    Query Params: None
    Body: None
    Response Body: application/json
    ```
    {
        token: string
    }
    ```

	POST /api/revoke
    Revoke the refresh token
    Headers: Authorization - The refresh token as an Bearer token
    Query Params: None
    Body: None
    Response Body: None

	GET /api/healthz
    Get a 200 ok message

	GET /admin/metrics
    Webpage that displays the number of file hits

	POST /admin/reset
    Reset the database and metrics. Only available when PLATFORM environment variable is set to "dev"


