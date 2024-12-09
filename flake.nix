{
  description = "core";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    crane = {
      url = "github:ipetkov/crane";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, crane, rust-overlay }:
    let
      overlays = [
        (import rust-overlay)
        (self: super: {
          rustToolchain = super.rust-bin.fromRustupToolchainFile ./rust-toolchain.toml;
        })
      ];

      # supported systems
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];

      # system-specific attributes
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit overlays system; };
        craneLib = crane.lib.${system};
      });
    in
    {
      # package
      packages = forAllSystems ({ pkgs, craneLib }: {
        default = craneLib.buildPackage {
            src = craneLib.cleanCargoSource (craneLib.path ./.);
            strictDeps = true;
            buildInputs = (with pkgs; [
            ]) ++ pkgs.lib.optionals pkgs.stdenv.isLinux (with pkgs; [
              openssl
            ]) ++ pkgs.lib.optionals pkgs.stdenv.isDarwin (with pkgs; [
              libiconv
              darwin.apple_sdk.frameworks.CFNetwork
              darwin.apple_sdk.frameworks.SystemConfiguration
            ]);
            nativeBuildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux (with pkgs; [
              pkg-config
            ]);
        };
      });

      # development environment
      devShells = forAllSystems ({ pkgs, craneLib }: {
        default = pkgs.mkShell {
          packages = (with pkgs; [
	    git
	    just
            rustToolchain
          ]) ++ pkgs.lib.optionals pkgs.stdenv.isLinux (with pkgs; [
            openssl
            pkg-config
          ]) ++ pkgs.lib.optionals pkgs.stdenv.isDarwin (with pkgs; [
            libiconv
            darwin.apple_sdk.frameworks.CFNetwork
            darwin.apple_sdk.frameworks.SystemConfiguration
          ]);
        };
      });
    };
}
