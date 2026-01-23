# Contributing to MyTrack

Thank you for your interest in contributing to MyTrack! We welcome contributions from the community to help make this personal finance tool better for everyone.

## Getting Started

1.  **Fork the repository** on GitHub.
2.  **Clone your fork** locally:
    ```bash
    git clone https://github.com/YOUR_USERNAME/go-eila.git
    cd go-eila
    ```
3.  **Create a new branch** for your feature or bug fix:
    ```bash
    git checkout -b feature/amazing-feature
    ```

## Code Style

*   **Go Formatter**: Always run `go fmt ./...` before committing.
*   **Linting**: We recommend using `golangci-lint` to catch common issues.
*   **Comments**: proper documentation comments for public functions and types are encouraged.

## Architecture Guidelines

*   **Repository Pattern**: All database access goes through `internal/repository`. Do not interact with `sql.DB` directly in the UI code.
*   **Validation**: Validate all user inputs in the UI layer before calling repository methods.
*   **Error Handling**: Show meaningful error messages to the user via dialogs. Log technical details if necessary.

## Pull Request Process

1.  Ensure all code compiles and runs.
2.  Test your changes manually (since UI testing is manual for now).
3.  Update documentation (README.md) if you are changing features or installation steps.
4.  Submit a Pull Request to the `main` branch.
5.  Provide a clear description of what your changes do and why.

## License

By contributing, you agree that your contributions will be licensed under the MIT License of this project.
