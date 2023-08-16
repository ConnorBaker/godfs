{
  buildGoModule,
  lib,
}:
buildGoModule {
  pname = "godfs";
  version = "0.0.1";
  src = lib.sources.sourceByRegex ../.. [
    "^go.(:?mod|sum)$"
    "main\\.go$"
    "(:?cmd|database|decode|encode|utils)(:?/.*)?"
    # NOTE: Include a directory and all of its contents with "directory_name(:?/.*)?".
  ];

  vendorHash = "sha256-wN/c3MdN9U98EIKb2m3EAIfttM0JCF8W+ZPfb26M2DQ=";

  meta = with lib; {
    description = "Erasure codes test";
    homepage = "https://github.com/ConnorBaker/godfs";
    maintainers = with maintainers; [connorbaker];
  };
}
