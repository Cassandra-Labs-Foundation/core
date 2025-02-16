## Feb 16th 2025

Ok we are going to set up a TigerBeetle mock 

I setup an EC2 instance on AWS. I am using the mock-server-test.pem for key-pair auth (added to .gitignore)
    - ended up dropping this because o3-mini-high could not figure out why none of the instances were ssh accessible

Went on GCP and easily setup a VM with Docker, using [TigerBeetle's image](https://docs.tigerbeetle.com/operating/docker/)

The development process will be the following
    1. Develop banking logic integrated with TigerBeetle ledger locally 
    2. Deploy on GCP with Docker alongside an API server to expose our endpoints
    3. Provide documentation for the API endpoints so that the fintech partners can simulate transactions

