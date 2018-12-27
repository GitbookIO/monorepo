# MonoRepo

> A big home for small repos

`monorepo` groups multiple (sub)repos into one big repo and supports bidirectional sync (`push` and `pull`).

# Setup

1. Start with a new (empty) repo to be your monorepo
2. Add a `monofile.yml` (see [_examples/simple](https://github.com/GitbookIO/monorepo/tree/master/_examples/simple))
3. Run `monorepo pull`
    - This will pull the files from the subrepos
    - Creaate `monofile.lock` with the specific SHAs pulled in
4. Commit the monorepo (e.g: `git commit -am "Initial pull"`)

# CI

You can use a CI system like Travis to automate syncing between your monorepo and sub repos.

## 1. Setup each subrepo to pull in new updates:
```
git clone https://github.com/org/bigrepo
cd bigrepo
monorepo pull
git add .
git commit -am "Update"
git push
```

## 2. Setup monorepo to dispatch changes

```
cd bigrepo
monorepo push
```

# Usage

```
❯ monorepo --help
NAME:
   monorepo - A big home for small repos

USAGE:
   monorepo [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR(S):
   Aaron O'Mullan <aaron@gitbook.com>

COMMANDS:
     list, ls
     pull
     add
     rm

GLOBAL OPTIONS:
   --force        Force action, may result in git force-pushes [$MONOREPO_FORCE]
   --root value   Path to the root of the monorepo [$MONOREPO_ROOT]
   --help, -h     show help
   --version, -v  print the version
```

## `monofile.yml`

```yaml
- repos:
    - path: folder1
      url: https://github.com/org/repo1
      ref: master
    - path: folder2
      url: https://github.com/org/repo2
      ref: v2
```

# Notes

- `monorepo` should never force push to your original repos (you can use the `--force` flag if you chose to do so)
