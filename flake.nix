{
  description = "A basic flake with a shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11-small";
    flake-utils.url = "flake:flake-utils";
    overlays.url = "github:thelonelyghost/blank-overlay-nix";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, overlays, flake-compat }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ overlays.overlays.default ];
        };

        rig = pkgs.buildGoModule {
          pname = "rig";
          version = "1.0.0";
          src = ./.;

          # vendorHash = pkgs.lib.fakeHash;
          vendorHash = "sha256-skYMlL9SbBoC89tFCTIzyRViEJaviXENASEqr6zSvoo=";

          meta = {
            description = "Instantly jump to your ripgrep matches";
            homepage = "https://github.com/thelonelyghost/rig";
          };
        };
      in
      {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = [
            pkgs.bashInteractive
            pkgs.go
            pkgs.statix
          ];
          buildInputs = [
            rig
          ];
        };

        packages = {
          inherit rig;

          default = rig;
        };
      });
}
