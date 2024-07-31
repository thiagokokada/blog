# Generating YAML files with Nix

I hate YAML. Instead of writing an essay on why I hate YAML, I can just link to
[noyaml.com](https://noyaml.com/). In my personal projects I will never use it,
preferring either JSON, [TOML](https://toml.io/en/) or even plain old
[INI](https://en.wikipedia.org/wiki/INI_file) files depending on the use case.
However the ship has sailed already, there are tons of projects everywhere that
uses YAML: from most CI systems ([GitHub
Actions](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions),
[CircleCI](https://circleci.com/docs/introduction-to-yaml-configurations/),
[Travis](https://docs.travis-ci.com/user/build-config-yaml), et tu
[builds.sr.ht](https://man.sr.ht/builds.sr.ht/), to
[Kubernetes](https://kubernetes.io/docs/concepts/overview/working-with-objects/),
or in almost every
[Rails](https://guides.rubyonrails.org/configuring.html#configuring-a-database)
application.

One way to avoid at least some issues with the language is to write YAML in
another language. I will show my solution in one of my [personal
repositories](https://github.com/thiagokokada/nix-configs/), writing Nix to
generate GitHub Actions configuration files. Bonus points for validating the
result against the schema of GitHub Actions, so the famous "this is supposed to
be string instead of a list of strings" is gone.

Let's start with the basics: YAML is supposed to be a [superset of
JSON](https://stackoverflow.com/a/1729545). What that means is that a JSON file
[can be parsed](https://yaml.org/spec/1.2-old/spec.html#id2759572) by a YAML
parser. And Nix itself generates JSON natively, after all, Nix can be imagined
as ["JSON with functions"](https://nix.dev/tutorials/nix-language.html).

To make things easier, I will assume that you have the `nix-commands` and
`flakes` enabled as `experimental-features` in your Nix configuration. If not,
go [here](https://wiki.nixos.org/wiki/Flakes).

Using the `nix eval` command, we can generate a JSON expression from Nix by:

```console
$ nix eval --expr '{ foo = "bar"; }' --json
{"foo":"bar"}
```

However, typing long excerpts of Nix code inside the console would be
impractical. We can write the following code inside a `foo.nix` file instead:

```nix
{
  foo = "bar";
}
```

And:

```console
$ nix eval --file foo.nix --json
{"foo":"bar"}
```

While you can use a JSON output as an input for YAML parsers, it is probably
not the [best idea](https://metacpan.org/pod/JSON::XS#JSON-and-YAML). Sadly (or
maybe not), Nix has no native functionality to export data to YAML. However,
since we are using Nix, it is trivial to use `nixpkgs` to use some program to
convert from JSON to YAML.

To start, let's create a new directory, move our `foo.nix` file to it, create a
new `flake.nix` file and put the following contents:

```nix
{
  description = "Generate YAML files with Nix";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { nixpkgs, ... }:
    {
      packages.x86_64-linux =
        let
          inherit (nixpkgs) lib;
          pkgs = import nixpkgs { system = "x86_64-linux"; };
        in
        {
          toYAML =
            let
              file = import ./foo.nix;
              json = builtins.toJSON file;
            in
            pkgs.runCommand "toYAML" {
              buildInputs = with pkgs; [ yj ];
            } ''
              mkdir -p $out
              echo ${lib.escapeShellArg json} | yj -jy > $out/foo.yaml
            '';
        };
    };
}
```

We are loading the `./foo.nix` as a Nix file, converting it to JSON with
`builtins.toJSON` function, and finally, using `pkgs.runCommand` to pipe the
output of the JSON file to [yj](https://github.com/sclevine/yj), that allow
convertion between serialisation formats. `-jy` flag means "JSON to YAML". The
reason I choose `yj` is mostly because it is a single binary Go program, but
you can use whatever you prefer.

By the way, there is a
[`lib.generators.toYAML`](https://github.com/NixOS/nixpkgs/blob/9f918d616c5321ad374ae6cb5ea89c9e04bf3e58/lib/generators.nix#L805)
inside `nixpkgs.lib`, but as of 2024-07-31 it only calls `lib.strings.toJSON`
(that in turn, calls `builtins.toJSON`). So it doesn't really help here.

If we run the following commands, we can see the result:

```console
$ nix build .#packages.x86_64-linux.toYAML
$ cat result/foo.yaml
foo: bar
```

That is the basic idea. To have a more realistic example, let's convert the
[`go.yml`](https://github.com/thiagokokada/blog/blob/4e3f25485c6682f3e066b219df2290934bc0d256/.github/workflows/go.yml),
that builds this blog, to Nix:

```nix
{
  name = "Go";
  on.push.branches = [ "main" ];

  jobs = {
    build = {
      runs-on = "ubuntu-latest";
      permissions.contents = "write";
      steps = [
        { uses = "actions/checkout@v4"; }
        {
          name = "Set up Go";
          uses = "actions/checkout@v4";
          "with".go-version = "1.21";
        }
        {
          name = "Update";
          run = "make";
        }
        {
          name = "Publish";
          run = "make publish";
          env.MATAROA_TOKEN = ''''${{ secrets.MATAROA_TOKEN }}'';
        }
        {
          name = "Commit";
          uses = "stefanzweifel/git-auto-commit-action@v5";
          "with".commit_message = "README/rss:update";
        }
      ];
    };
  };
}
```

Some interesting things to highlight: `with` is a reserved word in Nix, so we
need to quote it. Not a problem, but something to be aware. And the template
string in GitHub Actions uses the same `${}` that Nix uses, so we need to
escape.

And after running the following commands:

```
$ nix build .#packages.x86_64-linux.toYAML
$ cat result/go.yaml
jobs:
  build:
    permissions:
      contents: write
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/checkout@v4
        with:
          go-version: "1.21"
      - name: Update
        run: make
      - env:
          MATAROA_TOKEN: ${{ secrets.MATAROA_TOKEN }}
        name: Publish
        run: make publish
      - name: Commit
        uses: stefanzweifel/git-auto-commit-action@v5
        with:
          commit_message: README/rss:update
name: Go
"on":
  push:
    branches:
      - main
```

Yes, the keys are not in the same order as we defined, since Nix, like most
programming languages (with the exception of
[Python](https://mail.python.org/pipermail/python-dev/2017-December/151283.html)),
do not guarantee the insertion order in maps/dicts/attrsets/whatever. But I
really hope whatever is consuming your YAML is not relying in the order the
keys are defined (this would be more cursed than YAML already is).

So that is basically it. For the bonus points that I talked at the start of the
post, we can modify `pkgs.runCommand` to run some kind of validator. I use
[`action-validator`](https://github.com/mpalmer/action-validator), one that I
particularly packaged in
[nixpkgs](https://github.com/NixOS/nixpkgs/pull/260217) to use in those cases.
But you could use e.g.: a validator of Kubernetes YAML. Or a generic YAML lint
like this [one](https://github.com/adrienverge/yamllint). The possibilities are
endless.

Let's modify our `flake.nix` to add the validation:

```nix
{
  # ...
  outputs = { nixpkgs, ... }:
    {
      packages.x86_64-linux =
        let
          inherit (nixpkgs) lib;
          pkgs = import nixpkgs { system = "x86_64-linux"; };
        in
        {
          toYAML =
            let
              file = import ./go.nix;
              json = builtins.toJSON file;
            in
            pkgs.runCommand "toYAML" {
              buildInputs = with pkgs; [ action-validator yj ];
            } ''
              mkdir -p $out
              echo ${lib.escapeShellArg json} | yj -jy > $out/go.yaml
              action-validator -v $out/go.yaml
            '';
        };
    };
}
```

And let's add an error in our `go.nix` file:

```diff
diff --git a/go.nix b/go.nix
index 25e0596..8c00033 100644
--- a/go.nix
+++ b/go.nix
@@ -5,7 +5,7 @@
   jobs = {
     build = {
       runs-on = "ubuntu-latest";
-      permissions.contents = "write";
+      permissions.contents = [ "write" ];
       steps = [
         { uses = "actions/checkout@v4"; }
         {
```

Finally, let's try to build our YAML file again:

```
$ nix build .#packages.x86_64-linux.toYAML
error: builder for '/nix/store/j8wr6j1pvyf986sf74hqw8k31lvlzac5-toYAML.drv' failed with exit code 1;
       last 25 log lines:
       >                                 "Additional property 'runs-on' is not allowed",
       >                             ),
       >                             path: "/jobs/build",
       >                             title: "Property conditions are not met",
       >                         },
       >                         Properties {
       >                             code: "properties",
       >                             detail: Some(
       >                                 "Additional property 'steps' is not allowed",
       >                             ),
       >                             path: "/jobs/build",
       >                             title: "Property conditions are not met",
       >                         },
       >                         Required {
       >                             code: "required",
       >                             detail: None,
       >                             path: "/jobs/build/uses",
       >                             title: "This property is required",
       >                         },
       >                     ],
       >                 },
       >             ],
       >         },
       >     ],
       > }
       For full logs, run 'nix log /nix/store/j8wr6j1pvyf986sf74hqw8k31lvlzac5-toYAML.drv'.
```

Yes, the output of `action-validator` is awfully verbose, but it is still
better than making ["8 commits/push in one
hour"](https://x.com/eric_sink/status/1430954572848287744).
