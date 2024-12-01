{
  outputs = { flake-utils, nixpkgs, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        setup-internal = pkgs.writeShellApplication {
          name = "setup-internal";
          text = let files = "{age,encrypted_keys,parse,tui,wordlist}"; in ''
            rm -fv internal/${files}.go
            cp -v ${pkgs.age.src}/cmd/age/${files}.go internal
            sed -Ei 's|package main|package internal|' internal/*.go
          '';
        };
      in
      rec {
        packages.default = pkgs.buildGoModule {
          pname = "sillysecrets";
          version = "0";
          src = ./.;
          vendorHash = "sha256-OruMcvH46u7Ms62o8w2cH4ptqoLIrpgmg2ER93yZrd4=";

          preBuild = pkgs.lib.getExe setup-internal;

          CGO_ENABLED = "0";
          ldflags = [ "-s" "-w" ];

          meta.mainProgram = "sesi";
        };

        devShells.default = pkgs.mkShell {
          inputsFrom = [ packages.default ];
          packages = with pkgs; [
            age
            age-plugin-fido2-hmac
            air
            delve
            setup-internal
            ssh-to-age
          ];
        };
      });
}
