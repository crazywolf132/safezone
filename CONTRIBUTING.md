# Contributing to SafeZone

First off, thank you for considering contributing to SafeZone! It's people like you that make SafeZone such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by the [SafeZone Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [maintainer's email].

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for SafeZone. Following these guidelines helps maintainers and the community understand your report, reproduce the behavior, and find related reports.

- Use the bug report template when creating an issue.
- Be as detailed as possible in your report.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for SafeZone, including completely new features and minor improvements to existing functionality.

- Use the feature request template when creating an issue for feature requests.
- Provide as much detail and context as possible.

### Your First Code Contribution

Unsure where to begin contributing to SafeZone? You can start by looking through these `beginner` and `help-wanted` issues:

- [Beginner issues](https://github.com/crazywolf132/safezone/labels/beginner) - issues which should only require a few lines of code, and a test or two.
- [Help wanted issues](https://github.com/crazywolf132/safezone/labels/help%20wanted) - issues which should be a bit more involved than `beginner` issues.

### Pull Requests

- Fill in the required template
- Do not include issue numbers in the PR title
- Include screenshots and animated GIFs in your pull request whenever possible
- Follow the Go styleguides
- Include thoughtfully-worded, well-structured tests
- Document new code
- End all files with a newline

## Styleguides

### Git Commit Messages

We use conventional commit messages. This means each commit message consists of a header, a body and a footer. The header has a special format that includes a type, a scope and a subject:

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

The header is mandatory and the scope of the header is optional.

Examples:

```
feat(zone): add new recovery strategy

This new strategy allows for ...

Closes #123
```

```
docs: update README with new API methods
```

```
fix(core): resolve race condition in error handling

This patch fixes a race condition that could occur when ...

Fixes #456
```

More examples:

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools and libraries such as documentation generation

### Go Styleguide

All Go code should adhere to the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and be formatted using `gofmt`.

## Additional Notes

### Issue and Pull Request Labels

This section lists the labels we use to help us track and manage issues and pull requests.

- `bug`: Confirmed bugs or reports that are very likely to be bugs.
- `enhancement`: Feature requests.
- `documentation`: Improvements or additions to documentation.
- `good first issue`: Good for newcomers.
- `help wanted`: Extra attention is needed.
- `question`: Further information is requested.

## Thank You!

Your contributions to open source, large or small, make great projects like this possible. Thank you for taking the time to contribute.