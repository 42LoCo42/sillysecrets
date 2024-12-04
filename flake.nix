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
        packages.default = pkgs.buildGoModule rec {
          pname = "sillysecrets";
          version = "1.3.0";
          src = ./.;
          vendorHash = "sha256-BdPKXOwWPGQGCvdXVKC7kax86UtUmLjD9wco8Gay2pE=";

          preBuild = pkgs.lib.getExe setup-internal;

          CGO_ENABLED = "0";
          ldflags = [ "-s" "-w" ];

          nativeBuildInputs = with pkgs; [
            installShellFiles
            makeBinaryWrapper
          ];

          postInstall = ''
            mv $out/bin/{${pname},sesi}

            wrapProgram $out/bin/sesi \
              --prefix PATH : ${pkgs.moreutils}/bin

            for i in bash fish zsh; do
              $out/bin/sesi completion $i > sesi.$i
              installShellCompletion sesi.$i
            done
          '';

          meta.mainProgram = "sesi";
        };

        devShells.default = pkgs.mkShell {
          inputsFrom = [ packages.default ];
          packages = with pkgs; [
            age
            age-plugin-fido2-hmac
            cobra-cli
            delve
            setup-internal
            ssh-to-age
          ];
        };
      });
}
