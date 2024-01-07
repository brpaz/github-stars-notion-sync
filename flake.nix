{
  description = "Command Line tool to sync GitHub starred repositories to a Notion database";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let

       # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      version = self.version or "0.0.0"; # TODO find a better way to handle this

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          rev = self.rev or "unknown";
        in
        {
          github-stars-notion-sync = pkgs.buildGoModule {
            pname = "github-stars-notion-sync";
            inherit version;
            # In 'nix develop', we don't need a copy of the source tree
            # in the Nix store.
            src = ./.;
            CGO_ENABLED = 0;
            subPackages = [ "cmd" ];
            ldflags = [
                "-s"
                "-w"
                "-X main.version=${version}"
                "-X main.gitCommit=${rev}"
                "-X main.buildDate=${lastModifiedDate}"
                "-extldflags=-static"
            ];

            postInstall = ''
              mv $out/bin/cmd $out/bin/github-stars-notion-sync
            '';

            # This hash locks the dependencies of this package. It is
            # necessary because of how Go requires network access to resolve
            # VCS.  See https://www.tweag.io/blog/2021-03-04-gomod2nix/ for
            # details. Normally one can build with a fake sha256 and rely on native Go
            # mechanisms to tell you what the hash should be or determine what
            # it should be "out-of-band" with other tooling (eg. gomod2nix).
            # To begin with it is recommended to set this, but one must
            # remeber to bump this hash when your dependencies change.
            vendorHash = "sha256-wBESb9F8Edl2+RXBMXT0SC6n9ww/PJDnVJtUZuKWJ+s=";

            meta = with pkgs.lib; {
              homepage = "https://github.com/brpaz/${pname}";
              description = "Command Line tool to sync GitHub starred repositories to a Notion database";
              platforms = platforms.linux ++ platforms.darwin;
              license = licenses.mit;
              maintainers = [];
            };
          };
        });

      # The default package for 'nix build'. This makes sense if the
      # flake provides only one package or there is a clear "main"
      # package.
      defaultPackage = forAllSystems (system: self.packages.${system}.github-stars-notion-sync);
    };
}
