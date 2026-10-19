# Shell-Go

A Unix-style command-line shell built from scratch in **Go**, inspired by the CodeCrafters *Build Your Own Shell* challenge.

Shell-Go implements the core behavior of an interactive shell, including command parsing, process execution, I/O redirection, pipelines, background jobs, command history, autocompletion, and shell variable expansion.

The goal of this project was to understand how shells work internally — from interpreting user input to coordinating processes and managing their input and output.

## Features

### Core Shell Functionality
- Interactive REPL with a command prompt.
- Command parsing and invalid-command handling.
- Built-in commands: `exit`, `echo`, `type`, `pwd`, and `cd`.
- Execution of external programs through executable lookup.
- Directory navigation using absolute paths, relative paths, and the home directory.

### Command Parsing and Quoting
- Single-quoted strings.
- Double-quoted strings.
- Backslash escaping outside quotes and within double quotes.
- Literal handling of single-quoted content.
- Execution of commands with quoted executable names.

### I/O Redirection
- Redirect standard output (`>` and `1>`).
- Redirect standard error (`2>`).
- Append standard output (`>>` and `1>>`).
- Append standard error (`2>>`).

### Pipelines
- Connect commands using pipes (`|`).
- Support pipelines involving built-in commands.
- Execute multi-command pipelines.
- Pass the output of one command directly to the input of another.

### Background Jobs
- Start commands in the background using `&`.
- Track and list background jobs with the `jobs` built-in.
- Display background job status.
- Reap completed jobs and manage job lifecycle.
- Recycle job numbers after completed jobs are removed.

### Command Autocompletion
- Complete built-in commands and executable names.
- Complete filenames, nested paths, and directories.
- Handle partial input, multiple matches, and missing completions.
- Support command completion with arguments.
- Display multiple candidates and calculate the longest common prefix.

### Programmable Completion
- Register and unregister completion specifications with the `complete` built-in.
- Support command-specific completion scripts.
- Pass command-line arguments and environment variables to completion scripts.
- Handle multiple completion candidates and commands with no matching completions.

### Command History
- Store and list previously entered commands with `history`.
- Limit the number of displayed history entries.
- Navigate history using the Up and Down arrow keys.
- Execute commands from history.
- Persist history across sessions by reading from and writing to a history file.

### Shell Variables and Parameter Expansion
- Declare and store shell variables.
- Validate variable names.
- Expand variables using `$VAR` syntax.
- Support braced expansion using `${VAR}`.
- Handle unset and empty variables.

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/)
- Git
- A Unix-like development environment is recommended.

### Clone the Repository

```bash
git clone https://github.com/7odsa/shell-go.git
cd shell-go
```

### Run the Shell

If the repository includes the CodeCrafters launcher script, make it executable if necessary:

```bash
chmod +x your_program.sh
./your_program.sh
```

The shell should start an interactive session. You can try commands such as:

```bash
echo "Hello, World!"
pwd
cd ..
type echo
ls -la
```

To exit the shell, use:

```bash
exit
```

## Examples

### Output Redirection

```bash
echo "Hello, Shell-Go!" > output.txt
cat output.txt
```

### Pipelines

```bash
echo "pear-mango" | wc
cat output.txt | wc
```

### Background Jobs

```bash
sleep 5 &
jobs
```

### Command History

```bash
history
```

### Variable Expansion

```bash
declare NAME=Shell-Go
echo $NAME
echo ${NAME}
```

## Technical Concepts

This project explores several important systems-programming concepts:

- **REPL design:** Reading, interpreting, and executing commands in an interactive loop.
- **Process management:** Launching external programs and tracking their execution.
- **Standard streams:** Working with standard input, standard output, and standard error.
- **File descriptors and redirection:** Connecting command output to files and directing errors separately.
- **Pipelines:** Connecting processes through streams of data.
- **Concurrency:** Managing background jobs while keeping the shell responsive.
- **Parsing and escaping:** Interpreting shell syntax, quoting rules, and special characters.
- **Filesystem interaction:** Resolving executable paths, navigating directories, and completing file paths.
- **State and persistence:** Maintaining shell variables, command history, and job information.

## Built With

- **Go** — core implementation and process management.
- **Git** — version control.
- **CodeCrafters** — guided, test-driven project challenge.

## Project Goals

Shell-Go was built to go beyond simply using command-line tools and explore how they work internally. Implementing these features required dealing with process lifecycles, stream handling, parsing edge cases, persistent state, and interactions between commands.

## Acknowledgements

Built as part of the [CodeCrafters Build Your Own Shell](https://codecrafters.io/challenges/shell) challenge.

---

**Author:** [7odsa](https://github.com/7odsa)

