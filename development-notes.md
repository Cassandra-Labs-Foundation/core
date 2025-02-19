## Feb 19th 2025

o3-mini says the first step is to actually build the authentication module, so I setup server.js
    - idk how I feel about using Node for this
    - yea, dropping Node to use Rust so that @Jacob would be proud

I might be regretting this Rust thing. I still haven't been able to run the server because the dependencies don't compile...

## Feb 16th 2025

Ok we are going to set up a TigerBeetle mock 

I setup an EC2 instance on AWS. I am using the mock-server-test.pem for key-pair auth (added to .gitignore)
    - ended up dropping this because o3-mini-high could not figure out why none of the instances were ssh accessible

Went on GCP and easily setup a VM with Docker, using [TigerBeetle's image](https://docs.tigerbeetle.com/operating/docker/)

The development process will be the following
    1. Develop banking logic integrated with TigerBeetle ledger locally 
    2. Deploy on GCP with Docker alongside an API server to expose our endpoints
    3. Provide documentation for the API endpoints so that the fintech partners can simulate transactions

