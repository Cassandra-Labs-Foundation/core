fmt:
    cargo fmt --all

lint:
    cargo clippy --all --all-targets --all-features

build:
    cargo build --all --all-targets --all-features

test:
    cargo test --all --all-targets --all-features

ci: lint build test
