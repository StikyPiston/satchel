{
  description = "satchel devshell and package";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          name = "satchel-devshell";

          packages = with pkgs; [
            go
            gopls
            gotools
            delve
          ];
        };

        packages.satchel = pkgs.buildGoModule {
          pname = "satchel";
          version = "2026.03.01-a";

          src = self;

          vendorHash = "sha256-ZxugMmzZG6DdOED4fDHSbDjwjuqpqvcXgO55nNcAPyY=";

          subPackages = [ "." ];
          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "A simple command-line file clipboard to make file management easier";
            license = licenses.mit;
            platforms = platforms.all;
          };
        };

        apps.satchel = {
          type = "app";
          program = "${self.packages.${system}.satchel}/bin/satchel";
        };
      });
}
