{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        packages.default = pkgs.stdenv.mkDerivation {
          pname = "gio-term";
          version = "0.1.0";
          src = ./.;
          buildInputs = with pkgs; [
            go 
            vulkan-headers
            wayland
            libGL
            pkg-config
            vulkan-loader
          ];

          installPhase = ''
            mkdir -p $out/bin
            cp bin/gioterm $out/bin/
          '';

        };
        devShells = 
        {
          default = with pkgs;
            mkShell ({
              packages = [ clang ]
                ++ (if stdenv.isLinux then [
                  go
                  vulkan-headers
                  wayland
                  libGL
                  pkg-config
                ] else
                  [ ]);
            } // (if stdenv.isLinux then {
              LD_LIBRARY_PATH = "${vulkan-loader}/lib";
            } else
              { }));
        };
      });
}
