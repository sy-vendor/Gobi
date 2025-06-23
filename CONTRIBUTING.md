# Contributing to Gobi

We welcome contributions from the community! Thank you for your interest in making Gobi better.

## How to Contribute

- **Report bugs**: If you find a bug, please create a bug report.
- **Suggest features**: If you have an idea for a new feature, please create a feature request.
- **Write code**: If you want to contribute code, please follow the development setup below.
- **Improve documentation**: If you see an error or an area for improvement in the documentation, please submit a pull request.

## Development Setup

1.  **Fork the repository**: Click the "Fork" button on the top right of the GitHub page.
2.  **Clone your fork**:
    ```bash
    git clone https://github.com/sy-vendor/gobi.git
    cd gobi
    ```
3.  **Create a branch**:
    ```bash
    git checkout -b feature/your-new-feature
    ```
4.  **Make your changes**: Write your code and add tests.
5.  **Run tests**:
    ```bash
    go test -v ./...
    ```
6.  **Commit your changes**:
    ```bash
    git commit -m "feat: Add your new feature"
    ```
7.  **Push to your branch**:
    ```bash
    git push origin feature/your-new-feature
    ```
8.  **Create a Pull Request**: Go to the original Gobi repository and create a pull request.

## Code Style

- Follow standard Go formatting (`gofmt`).
- Use meaningful variable and function names.
- Write clear and concise comments for complex logic.
- Ensure your code is well-tested.

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification. Your commit messages should be in the following format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.

**Example**:
```
feat: Add support for 3D bubble charts
```

## Pull Request Process

1.  Ensure all tests are passing.
2.  Update the `README.md` and other relevant documentation if your changes require it.
3.  Your pull request will be reviewed by the maintainers. We may ask for changes.
4.  Once your pull request is approved, it will be merged into the `main` branch.

Thank you for your contribution! 