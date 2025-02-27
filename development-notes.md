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

