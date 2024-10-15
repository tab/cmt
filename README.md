# cmt

**cmt** is a command-line utility that generates Conventional Commit messages using OpenAI's GPT models based on your staged git changes. It automates the process of writing clear and structured commit messages, improving your git workflow.

## Features

- **Automated Commit Messages**: Generates commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification.
- **Integration with OpenAI GPT**: Uses GPT to analyze your staged changes and produce meaningful commit messages.
- **Interactive Approval**: Allows you to review and approve the generated commit message before committing.

## Prerequisites

- **Go**: Version 1.16 or higher is recommended.
- **Git**: Ensure git is installed and initialized in your project.
- **OpenAI API Key**: You need an API key from OpenAI to use GPT models.

## Installation

1. **Clone the Repository**

   ```sh
   git clone https://github.com/tab/cmt.git
   ```

2. **Navigate to the Project Directory**

   ```sh
   cd cmt
   ```

3. **Set Up Environment Variables**


   ```sh
   export OPENAI_API_KEY=your-api-key-here
   ```

4. **Build the Binary**

   ```sh
   go build -o cmd/cmt cmd/main.go
    ```

5. **Make the Binary Executable**

   ```sh
   chmod +x cmd/cmt
   ```

6. **Move the Binary to Your PATH**

   ```sh
   sudo ln -s $(pwd)/cmd/cmt /usr/local/bin/cmt
   ```

7. **Verify Installation**

   ```sh
    cmt --version
    ```

## Usage

Navigate to your git repository and stage the changes you want to commit:

```sh
git add .
```

Run the `cmt` command to generate a commit message:

```sh
cmt
```

Review the generated commit message and choose whether to commit or not.

```sh
💬 Message: feat(core): Add user authentication

Implemented JWT-based authentication for API endpoints. Users can now log in and receive a token for subsequent requests.

Accept? (y/n):
```

Type **y** to accept and commit the changes, or **n** to abort.

```sh
🚀 Changes committed:
[feature/jwt 29ca12d] feat(core): Add user authentication
 2 files changed, 106 insertions(+), 68 deletions(-)
 ...
```

### Configuration

Optional prefix for the commit message can be set with the `--prefix` flag:

```sh
cmt --prefix "TASK-1234"
```

or

```sh
cmt -p "TASK-1234"
```

Resulting commit message:

```sh
💬 Message: TASK-1234 feat(core): Add user authentication

Implemented JWT-based authentication for API endpoints. Users can now log in and receive a token for subsequent requests.

Accept? (y/n):
```

### Changelog generation

Run the `cmt changelog` to generate a changelog based on your commit history:

```sh
cmt changelog sha1..sha2
```

```sh
cmt changelog v1.0.0..v1.1.0
```

The command will output the changelog in the following format:

```sh
# CHANGELOG

[1.1.0]

### Features

- **feat(core):** Add user authentication
...
```

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

- [OpenAI](https://openai.com/) for providing the GPT models.
- [Conventional Commits](https://www.conventionalcommits.org/) for the commit message specification.
