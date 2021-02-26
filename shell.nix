let
  unstable = import ./unstable.nix;
in
{ pkgs ? unstable }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go_1_15
  ];
}
