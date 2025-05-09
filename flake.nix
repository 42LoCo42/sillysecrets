{
  outputs = { flake-utils, nixpkgs, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; }; in
      rec {
        packages.default = pkgs.buildGoModule rec {
          pname = "sillysecrets";
          version = "2.1.4";
          src = ./.;

          ldflags = [ "-s" "-w" ];
          vendorHash = "sha256-XIh0+k/CJ4RrRw9+Lvl/BKPh4TM5g+gY2jnrWwFvfTs=";

          nativeBuildInputs = with pkgs; [
            installShellFiles
            makeBinaryWrapper
          ];

          postInstall = ''
            mv $out/bin/{${pname},sesi}

            wrapProgram $out/bin/sesi \
              --prefix PATH : ${pkgs.moreutils}/bin

            $out/bin/sesi man
            installManPage man/*

            for i in bash fish zsh; do
              $out/bin/sesi completion $i > sesi.$i
              installShellCompletion sesi.$i
            done
          '';

          meta.mainProgram = "sesi";
        };

        devShells.default = pkgs.mkShell {
          shellHook = ''
            PATH="$PWD:$PATH"
          '';

          inputsFrom = [ packages.default ];
          packages = with pkgs; [
            cobra-cli
            delve
          ];
        };
      });
}
