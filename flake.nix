{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
    zig.url = "github:mitchellh/zig-overlay";
  };

  outputs = { self, nixpkgs, flake-compat, zig }:
    let pkgs = import nixpkgs {
        system = "x86_64-linux";
    };
    in
    {
      devShell.x86_64-linux =
        pkgs.mkShell {
          buildInputs = with pkgs;[
            go_1_21
            gopls
            go-task
            golangci-lint
            goose
            zig.packages.${system}.master
            goreleaser
            syft
          ];
          shellHook = ''
            export PATH=$PWD/tools:$PATH
          '';
        };
    };
}
