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
        packages.default = pkgs.buildGoModule {
          pname = "gio-term";
          version = "0.1.0";
          src = ./.;

          tags = [
            "nox11"
          ];

          vendorHash = "sha256-WpYO5aQcvIcrM0g5q3ibsIJzZHKLKs6YLoy1v6D6jlM=";
          
          nativeBuildInputs = with pkgs; [ pkg-config ];
          buildInputs = with pkgs; [
            go 
            libxkbcommon
            libx11
            libxcursor
            libxfixes
            vulkan-headers
            wayland
            libGL
            vulkan-loader
          ];

        };
        devShells = 
        {
          default = with pkgs;
            mkShell ({
              packages = [ clang ]
                ++ (if stdenv.isLinux then [
                  go
                  vulkan-headers
                  libxkbcommon
                  wayland
                  libx11
                  libxcursor
                  libxfixes
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
