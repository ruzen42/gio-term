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
                  xorg.libX11
                  xorg.libXcursor
                  xorg.libXfixes
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
