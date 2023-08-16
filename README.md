TODO:

- [ ] File is added to the sqlite database
- [x] sqlite database primary key is the hash
- [x] Move to SRI hashes
- [ ] think about how best to model relationships between files
  - [ ] A file has exactly one sri hash
  - [ ] A file has at least one shard
  - [ ] A shard can be used by multiple files
    - This is because a shard is a block of data, and some files will share blocks of data, and we are using content addressable storage
- [ ] How to handle upserts/updates?
- [ ] Constraints on keys (see Nix database?)