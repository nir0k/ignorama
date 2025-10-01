Thought for 10 seconds


# **Proposed name:** **Ignorama**

*A tiny “ignore‑aware” file lister that sweeps a directory tree, honours all your `.gitignore` rules (plus an optional second ignore file) and writes the matched files—together with their contents—into a single text file.*

---

```text
Ignorama
├── Reads <dir>/.gitignore  ➜ understands full git‑ignore syntax
├── Reads an extra ignore   ➜ same syntax, supplied with -e/--exclusions
├── Filters by extensions   ➜ -x "go,md,txt"  (empty ⇒ all files)
└── Writes list + contents  ➜ default: go_files.txt
```

---

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Command‑line flags](#command-line-flags)
5. [Examples](#examples)
6. [Exit codes](#exit-codes)
7. [License](#license)

---

## Features

* **Ignore aware** – merges patterns from a project‑local `.gitignore` **and** a user‑supplied ignore file. Supports `*`, `**/`, `!negate`, comments, etc.
* **Selective export** – optional extension filter lets you grab just `*.go`, `*.md`, whatever.
* **Self‑contained output** – each file path followed by its raw source between triple back‑ticks; perfect for quick greps or diff‑friendly code reviews.
* **Cross‑platform** – builds on Linux, macOS, Windows; Go 1.20+.
* **One‑shot viewer** – `--open` pops the resulting file in your system’s default editor.

---

## Installation

```bash
go install github.com/your‑org/ignorama@latest
```

The binary `ignorama` will appear in `$(go env GOPATH)/bin`.

---

## Usage

```bash
ignorama [flags]
```

### Command‑line flags

| Short | Long           | Default          | Description                                                                      |
| ----- | -------------- | ---------------- | -------------------------------------------------------------------------------- |
| `-d`  | `--directory`  | `.`              | Root directory to scan.                                                          |
| `-e`  | `--exclusions` | `exclusions.txt` | Extra ignore file (full git‑ignore syntax). Ignored if the file is missing.      |
| `-x`  | `--extensions` | *(empty)*        | Comma‑separated list of extensions **without dots**. Empty ⇒ include all.        |
| `-o`  | `--open`       | `false`          | After generation, open the output file in the system’s default editor/viewer.    |
|       |                |                  | *(implicit)* Output file name is `go_files.txt`; override by renaming after run. |

---

## Examples

1. **Export every non‑ignored file under the current repo**

   ```bash
   ignorama
   ```

2. **Only Go + Markdown sources, honouring extra patterns**

   ```bash
   ignorama -x "go,md" -e .scanignore
   ```

3. **Scan `~/projects/demo`, then open the list right away**

   ```bash
   ignorama -d ~/projects/demo --open
   ```

---

## Exit codes

| Code | Meaning                                 |
| ---- | --------------------------------------- |
|  0   | Success                                 |
|  1   | Any runtime error (IO, walk, open, etc) |

---

## License

MIT – see **LICENSE** file for details.

Enjoy clean, predictable file exports with **Ignorama**!
