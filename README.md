# cmt

**cmt** is a command-line utility that generates [Conventional Commit](https://www.conventionalcommits.org/) messages using OpenAI's GPT models based on staged Git changes.

It automates the process of writing clear and structured commit messages, enhancing your Git workflow and ensuring consistency across projects.

## Features

- **Automated Commit Messages**: Generates commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification.
- **Interactive Approval**: Allows you to review and approve the generated commit message before committing.
- **Interactive Edit**: Supports editing the commit message interactively before committing.
- **Custom Prefixes**: Supports adding custom prefixes to commit messages for better traceability (e.g., task IDs, issue numbers).
- **Changelog Generation**: Automatically creates changelogs based on your commit history.
- **Integration with OpenAI GPT**: Utilizes GPT to analyze your staged changes and produce meaningful commit messages.

## Prerequisites

Before installing and using **cmt**, ensure you have the following:

- **Go**: Version 1.16 or higher is recommended.
- **Git**: Ensure Git is installed and initialized in your project.
- **OpenAI API Key**: Obtain an API key from [OpenAI](https://platform.openai.com/account/api-keys) to use GPT models.

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

_For permanent setup, add the above line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.)._

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

## Configuration

**cmt** can be configured using a YAML configuration file named `cmt.yaml` placed in the current directory. Below is an example configuration file with the default settings:

```yaml
api:
  retry_count: 3  # Number of retry attempts for API requests
  timeout: 60s    # Timeout duration for API requests

editor: vim       # Editor to use for interactive editing

model:
  name: gpt-4.1-nano # OpenAI model to use
  max_tokens: 500    # Maximum tokens for the model response
  temperature: 0.7   # Controls randomness of the model output

logging:
  format: console    # Logging format (console or json)
  level: info        # Logging level (debug, info, warn, error)
```

You can customize any of these settings to fit your preferences.
If no configuration file is found, **cmt** will use these default values.

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

Accept, edit, or cancel? (y/e/n):
```

Type **y** to accept and commit the changes, **e** to edit message or **n** to abort.

```sh
🚀 Changes committed:
[feature/jwt 29ca12d] feat(core): Add user authentication
 2 files changed, 106 insertions(+), 68 deletions(-)
 ...
```

### Optional Prefix

Optional prefix for the commit message can be set with the `--prefix` flag:

```sh
cmt --prefix "TASK-1234"
```

Resulting commit message:

```sh
💬 Message: TASK-1234 feat(core): Add user authentication

Implemented JWT-based authentication for API endpoints. Users can now log in and receive a token for subsequent requests.

Accept, edit, or cancel? (y/e/n):
```

### Changelog generation

Run the `cmt changelog` to generate a changelog based on your commit history:

```sh
cmt changelog SHA1..SHA2
```

```sh
cmt changelog v1.0.0..v1.1.0
```

The command will output the changelog in the following format:

```sh
# CHANGELOG

## [1.1.0]

### Features

- **feat(jwt):** Add user authentication
- **feat(api):** Implement rate limiting for API endpoints

### Bug Fixes

- **fix(auth):** Resolve token expiration issue

...
```

## FAQ

**Q:** How do I obtain an OpenAI API key?

**A:** You can obtain an API key by signing up at [OpenAI's website](https://platform.openai.com/account/api-keys). After signing in, navigate to the API keys section to generate a new key.

---

**Q:** How can I ensure that private information isn't shared with OpenAI?

**A:** Here are some best practices to prevent sharing private information with OpenAI:

1. **Review Staged Changes**: Before running the `cmt` command, carefully review the changes you have staged using `git diff --staged`. Ensure that no sensitive information (like passwords, API keys, or personal data) is included.
2. **Exclude Sensitive Files**: Use `.gitignore` to exclude files that contain sensitive information from being tracked and staged. For example:

   ```gitignore
   .env
   secrets/
   ```

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

- [OpenAI](https://openai.com/) for providing the GPT models.
- [Conventional Commits](https://www.conventionalcommits.org/) for the commit message specification.
