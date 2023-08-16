{
  perSystem = {pkgs, ...}: {
    packages = {
      godfs = pkgs.callPackage ./godfs.nix {};
    };
  };
}
