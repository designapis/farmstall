{
  nixpkgs ? <nixpkgs>
}:

# This is all for Strest, although will be adding more for the go server

let pkgs = import nixpkgs { overlays = [(self: super: {
      nodejs = super.nodejs-14_x;
    })]; };
in

pkgs.mkShell {
  buildInputs = [ pkgs.nodejs pkgs.yarn pkgs.bash ];
  shellHook = ''
      mkdir -p .nix-node
      export NPM_CONFIG_PREFIX=$PWD/.nix-node
      export NODE_PATH=$PWD/.nix-node

      export PATH=$NODE_PATH/bin:$PATH
      export PATH="$PWD/node_modules/.bin/:$PATH"
  '';
}
