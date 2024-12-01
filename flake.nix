{
  outputs = { flake-utils, nixpkgs, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; }; in rec {
        packages.default = pkgs.buildGoModule {
          pname = "sillysecrets";
          version = "0";
          src = ./.;
          vendorHash = null;

          CGO_ENABLED = "0";
          ldflags = [ "-s" "-w" ];

          meta.mainProgram = "sesi";
        };

        devShells.default = pkgs.mkShell {
          inputsFrom = [ packages.default ];
          packages = with pkgs; [ air ];
        };
      });
}
