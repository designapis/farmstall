{
  nixpkgs ? <nixpkgs>
}:

# This is all for Strest, although will be adding more for the go server

let pkgs = import nixpkgs { overlays = [(self: super: {
      nodejs = super.nodejs-14_x;
      yarn = super.yarn.override {
        nodejs = self.nodejs;
      };
    })]; };
in

pkgs.mkShell {
  buildInputs = [ pkgs.nodejs pkgs.yarn ];
  shellHook = ''
      mkdir -p .nix-node
      export NODE_PATH=$PWD/.nix-node
      export NPM_CONFIG_PREFIX=$PWD/.nix-node
      export PATH=$NODE_PATH/bin:$PATH
  '';
}
