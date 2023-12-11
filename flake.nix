{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
    zig.url = "github:mitchellh/zig-overlay";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-compat, flake-utils, zig }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs;[
            go_1_21
            gopls
            go-task
            golangci-lint
            goose
            # https://github.com/ziglang/zig/issues/17130
            zig.packages.${system}.master-2023-08-15
            goreleaser
            syft
            # snyk
          ];
          shellHook = ''
            export PATH=$PWD/tools:$PATH
          '';
        };
      }
    );
}
