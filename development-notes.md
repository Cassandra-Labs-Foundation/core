## Mar 18th 2025

- ok let's do this!
    - I also should diagram the schema at the end

- let's add Business entities
    - we added `/api/business/handler.go` and `/service/business/service.go` and `/repository/business.go` 
        - added routes to `main.go`
    - now let's add business entities on Supabase
        - done
    - next step is to write a test 
        - `auth-and-business.sh` passes

- next step is KYC verification 
    - ok so we added `government_id`, `nationality`, and `kyc_document_url`

### Person Entity Schema

| Field                | Type        | Nullable | Description                                                |
|----------------------|-------------|----------|------------------------------------------------------------|
| **id**               | uuid        | No       | Primary key, auto-generated (default: `gen_random_uuid()`) |
| **first_name**       | text        | No       | First name of the person                                   |
| **last_name**        | text        | No       | Last name of the person                                    |
| **date_of_birth**    | date        | No       | Date of birth in `YYYY-MM-DD` format                       |
| **ssn**              | text        | Yes      | Social Security Number (optional)                          |
| **email**            | text        | Yes      | Email address (optional)                                   |
| **phone_number**     | text        | Yes      | Contact phone number (optional)                            |
| **street1**          | text        | Yes      | Primary street address (optional)                          |
| **street2**          | text        | Yes      | Secondary street address (optional)                        |
| **city**             | text        | Yes      | City of residence (optional)                               |
| **state**            | text        | Yes      | State or province (optional)                               |
| **postal_code**      | text        | Yes      | ZIP or postal code (optional)                              |
| **country**          | text        | Yes      | Country name or code (optional)                            |
| **kyc_status**       | text        | No       | KYC status (default: `"pending"`)                          |
| **kyc_verified_at**  | timestamptz | Yes      | Timestamp when KYC was verified (optional)                 |
| **government_id**    | text        | Yes      | Government-issued ID number (optional)                     |
| **nationality**      | text        | Yes      | Nationality of the person (optional)                       |
| **kyc_document_url** | text        | Yes      | URL for the uploaded KYC document (optional)               |
| **created_at**       | timestamptz | No       | Record creation timestamp                                  |
| **updated_at**       | timestamptz | No       | Record last update timestamp                               |


### Business Entity Schema

| Field                 | Type        | Nullable | Description                                                |
|-----------------------|-------------|----------|------------------------------------------------------------|
| **id**                | uuid        | No       | Primary key, auto-generated (default: `gen_random_uuid()`) |
| **name**              | text        | No       | Business name                                              |
| **registration_number** | text     | No       | Unique registration number of the business                 |
| **address**           | text        | No       | Business address                                           |
| **country**           | text        | No       | Country where the business is registered                   |
| **kyc_status**        | text        | No       | KYC status (default: `"pending"`)                          |
| **kyc_verified_at**   | timestamptz | Yes      | Timestamp when KYC was verified (optional)                 |
| **tax_id**            | text        | Yes      | Tax identification number (optional)                     |
| **kyc_document_url**  | text        | Yes      | URL for the uploaded KYC document (optional)               |
| **created_at**        | timestamptz | No       | Record creation timestamp                                  |
| **updated_at**        | timestamptz | No       | Record last update timestamp                               |



## Mar 17th 2025

- ok let's do this
    - the first step is to break this down further
        - ok done, it's now all time-boxed in my calendar

- next step is to figure out where exactly I left off
    - in a way, this is an opportunity to start over
        - I'm not going to use SuperGrok just yet, but I will transitioning into explaining the project to GPT
    - the next step here is to fit the documentation into o3
    - I used Repomix on the current state of the repo to have o3 understand where we left off
        - You’ve already decided on a Bearer key–based approach and implemented the authentication middleware.
        - The endpoints for token issuance (/auth/login), token refresh (/auth/refresh), and validation (/auth/validate) are all in place and are being exercised by your test script.
        - The endpoints for the Person entity (POST /entities/person, GET /entities/person, PATCH /entities/person, and list functionality) have been implemented.
        - The Business endpoints and additional KYC details (like comprehensive KYC fields and document uploads) are still pending.
        - In cmd/server/main.go, you've added logic to create a Supabase client using configuration values.
        - You've set up a person repository (using the Supabase REST API) and integrated the corresponding service and HTTP handler.
        - New routes for person entities have been added under the /entities/person endpoint (supporting POST, GET, GET by ID, and PATCH). This indicates that the Entity Onboarding piece for the Person entity is now functional, albeit likely with basic fields and behavior.
        - The configuration file (internal/config/config.go) was updated to include Database and Supabase settings, ensuring that your application now pulls these values from environment variables. This lays the groundwork for both the Person onboarding and future features that may require database access.

- all the action is happening `internal/repository/person.go` 
    - `./auth-and-person.sh` runs the test 
    - most thigns seem to be working, except for the updating the person entity with KYC status
        - the key turned out to be making the type on Supabase not just a Timestamp but also a Timezone

- ok now that all the tests have passed, let's take some time properly document all of this
    - `main` package contains `main.go` and it initializes the app, loads configs, sets up routes and runs the server
        - Divided into authentication (/auth) and protected routes (/), which include Person entity routes.
    - `internal` package is for code that used internally by the project
        - `api` contains the endpoints for how the service interacts with externals
            - `auth/handler.go`: Handles login, token refresh, and validation endpoints.
            - `middleware/auth.go`: Implements a Gin middleware to validate JWT tokens.
            - `person/handler.go`: Defines CRUD-like endpoints (POST, GET, PATCH, etc.) for Person entities.
        - `service` holds business logic, core operations and validation rules
            - `auth/service.go` is authentication logic (login, token refresh, token validation).
            - `person/service.go` is person onboarding logic (e.g., validating input, parsing dates, setting default KYC status).
        - `repository` is for interactions with external data sources
            - `person.go` contains PersonEntity struct and the personRestRepository that talks to Supabase for CRUD operations.
            - this is where we would put TigerBeetle interactions
        - `config` manages config loading, env variables, and settings
            - `config.go` loads environment variables into a Config struct.
                - Defines server settings (port, timeouts), JWT settings, database config, and Supabase config.
        - our `internal/clients` is where we have `supabase/client.go` I.E. the Supabase client to interact with the Supabase API
        - our `internal/database` is where we connect with our PostGres DB
    - `pkg` holds can code that can be reused outside of the project
        - `jwt/jwt.go` handles JWT creation and validation (signing, parsing, verifying claims).
    - `go.mod` and `go.sum` handle external packages and dependencies
- client sends a request to the `api` endpoint which calls a method in `service` that validates/transforms the input to then call a `repository` method to persist the data 
    - A client (e.g., auth-and-person-test.sh) sends a request to one of the endpoints (e.g., POST /api/v1/entities/person).
    - The request goes through the Gin router (in main.go), which calls the appropriate handler function (in internal/api/person/handler.go).
    - The handler validates and parses input, then calls the service layer (in internal/service/person/service.go).
    - The service layer applies business logic (e.g., verifying data, setting defaults) and calls the repository (in internal/repository/person.go) to actually interact with Supabase.
    - The repository uses the Supabase client (in internal/clients/supabase/client.go) to send HTTP requests to the Supabase REST API.
    - The result is passed back up through the layers to the handler, which returns a response to the client.


## Feb 26th 2025

- the focus is on Supabase, let's get it working

- ok so I added the Supabase logic and wrote a test that authenaticates and uses the auth-key to create a person Entity
    - this test is currently failing for reasons to be determined
    - it might be related to the .env file not being loaded properly 

- noooo, Claude3.7 is done until 9pm...
    - I need to figure out a way to easily transition between models so I stop being bottlenecked like this 

## Feb 25th 2025

- After some review, I'm dropping Rust in favor of Go at this stage
    - Rust is essentially [a better version of C++](https://www.youtube.com/watch?v=5C_HPTJg5ek), where the focus is on checking for issues at compile time as opposed to runtime
    - Rust has a [notoriously steep learning curve](https://www.youtube.com/watch?v=2hXNd6x9sZs), which by itself isn't a problem, but at the current stage it would definitely slow us down 
    - Most of the benefits of memory efficiency would only make sense at large scales, but we are not operating a large scale system with lots of parallelization 

- Next, given that TigerBeetle takes care of the reliability of the ledger, development speed is the second major factor in this decision
    - TigerBeetle does NOT have a Rust client, but it does have [Go, Node, and Python clients](https://github.com/tigerbeetle/tigerbeetle/blob/main/src/clients/README.md)
    - As much as I'd love to use JavaScript for this as well, it makes sense to use Go and o3-mini agrees 
    - Using these languages also gives us a boost with LLMs because they are better trained on older programming language than Rust which is so new 

- Overall, we are not dropping Rust forever, but it doesn't make sense to pay such high development costs for memory management when it's less impactfulgo mod init
    - It makes more sense to iterate in Go and once we have stabilized the MVP we can re-write it in Rust if we want to 

- Ok now, I've integrated Claude with our Repo so that we can start implementing Go
    - just like GPT, it suggests with starting off with auth, which is painful, but ok 
    - Go code is organized into packages. Each directory represents a package, and files in the same directory are part of the same package.
    - In the Gin framework, middleware is a function that processes requests before they reach the handlers. We created an authentication middleware that validates JWT tokens.

- Here is the current structure
    - pkg/jwt: Reusable JWT package (could be used in other projects)
    - internal/config: Configuration management
    - internal/service/auth: Authentication business logic
    - internal/api/auth: HTTP handlers for authentication
    - internal/api/middleware: HTTP middleware for authentication
    - cmd/server: Application entry point

- The authentication module now provides these endpoints:
    - `POST /api/v1/auth/login` - For user login and token generation
    - `POST /api/v1/auth/refresh` - For refreshing expired tokens
    - `GET /api/v1/auth/validate` - For validating tokens (protected endpoint)

- Turns out, TigerBeetle is not a general-purpose database and has intentional limitations:
    - Limited data model - TigerBeetle only supports Accounts (for storing balances) and Transfers (for moving money between accounts)
    - No support for complex queries - It doesn't have SQL-like query capabilities or support for complex joins, filtering, etc.
    - No document storage - There's no way to store unstructured data like KYC documents, images, etc.
    - Limited field types - TigerBeetle has a fixed schema with specific field types for its account and transfer objects.

- The standard approach is to use TigerBeetle for what it excels at (the financial ledger) and pair it with a more flexible database like PostgreSQL for everything else. This is the architecture recommended by TigerBeetle's own documentation (apparently).
    - What if we use Supabase? According to Clause "This gives us the best of both worlds: A modern, developer-friendly PostgreSQL solution with Supabase. The high-performance, reliable financial ledger with TigerBeetle"
    - The strategy has two tiers
        1. Entity Management using Supabase
            - Store entity data (people, businesses, KYC info) in Supabase PostgreSQL tables
            - Use Supabase Storage for document uploads and retrieval
            - Leverage Supabase Auth for authentication if desired
        2. Financial Ledger using TigerBeetle
            - Keep TigerBeetle for the core accounting functions
            - Link accounts in TigerBeetle to entities in Supabase via IDs
            - Use TigerBeetle for all financial transactions and balance tracking

- Ok we are trying to start by setting up Persons entities on Supabase
    - Gotta check the .env file when you get back from the break (Claude is frozen)


## Feb 19th 2025

- o3-mini says the first step is to actually build the authentication module, so I setup server.js
    - idk how I feel about using Node for this
    - yea, dropping Node to use Rust so that @Jacob would be proud

- I might be regretting this Rust thing. I still haven't been able to run the server because the dependencies don't compile...

## Feb 16th 2025

- Ok we are going to set up a TigerBeetle mock 

- I setup an EC2 instance on AWS. I am using the mock-server-test.pem for key-pair auth (added to .gitignore)
    - ended up dropping this because o3-mini-high could not figure out why none of the instances were ssh accessible

- Went on GCP and easily setup a VM with Docker, using [TigerBeetle's image](https://docs.tigerbeetle.com/operating/docker/)

- The development process will be the following
    1. Develop banking logic integrated with TigerBeetle ledger locally 
    2. Deploy on GCP with Docker alongside an API server to expose our endpoints
    3. Provide documentation for the API endpoints so that the fintech partners can simulate transactions

