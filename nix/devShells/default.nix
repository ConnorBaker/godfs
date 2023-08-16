{
  perSystem = {
    config,
    lib,
    pkgs,
    ...
  }: {
    devShells = {
      godfs = pkgs.mkShell {
        strictDeps = true;
        inputsFrom = [config.packages.godfs];
        packages = with pkgs; [
          go
          gopls
          delve
          go-tools
          revive
        ];
      };
    };
  };
}
