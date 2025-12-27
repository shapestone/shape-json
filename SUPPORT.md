# Support

Thank you for using shape-json! This document will help you find the right resources for getting help.

## Getting Help

### Documentation

Before opening an issue, please check our comprehensive documentation:

- **[README.md](README.md)** - Quick start, installation, and feature overview
- **[USER_GUIDE.md](USER_GUIDE.md)** - Complete API documentation with examples
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - How shape-json works internally
- **[examples/](examples/)** - Working code samples for common use cases
- **[API Reference](https://pkg.go.dev/github.com/shapestone/shape-json)** - Generated Go package documentation
- **[JSONPath Documentation](pkg/jsonpath/README.md)** - JSONPath query engine guide
- **[CHANGELOG.md](CHANGELOG.md)** - Recent changes and migration guides

**Most questions are answered in these resources.** Please check them first!

---

## Questions & Discussions

**Have a question about using shape-json?**

### Where to Ask

- **GitHub Discussions** - [Ask questions here](https://github.com/shapestone/shape-json/discussions)
  - Best for: General questions, usage help, design discussions
  - Community-driven with maintainer participation
  - Searchable by others with similar issues

- **Search Closed Issues** - [Check if it's been asked before](https://github.com/shapestone/shape-json/issues?q=is%3Aissue+is%3Aclosed)
  - Many questions have already been answered

### Common Questions

Before asking, check if your question is covered in the guides:

- **"How do I parse JSON?"** → See [README Quick Start](README.md#usage) and [USER_GUIDE.md](USER_GUIDE.md#quick-start)
- **"How do I build JSON without type assertions?"** → See [DOM API Guide](USER_GUIDE.md#fluent-dom-api-recommended)
- **"How do I handle large files?"** → See [Streaming Parser](USER_GUIDE.md#streaming-parser)
- **"How do I query JSON with JSONPath?"** → See [JSONPath Guide](pkg/jsonpath/README.md)
- **"How is this different from encoding/json?"** → See [ARCHITECTURE.md](ARCHITECTURE.md#key-design-decisions)

### Please Do NOT Open Issues for Questions

**GitHub Issues are for bug reports and feature requests only.**

Opening issues for usage questions:
- Clutters the issue tracker
- Makes it harder to track real bugs
- Takes longer to get answered (issues aren't a discussion forum)

Use GitHub Discussions instead - you'll get faster, better answers!

---

## Reporting Bugs

Found a bug? We want to fix it! Please use our **[bug report template](.github/ISSUE_TEMPLATE/bug_report.yml)**.

### Before Reporting

1. **Search existing issues** - Check if it's already reported: [Open Issues](https://github.com/shapestone/shape-json/issues)
2. **Verify your version** - Make sure you're using the latest release
3. **Check the CHANGELOG** - It may be a known issue that's already fixed
4. **Create minimal reproduction** - Reduce to the smallest code that shows the problem

### What to Include

Your bug report should include:

- **shape-json version** - Check with `go list -m github.com/shapestone/shape-json`
- **Go version** - Run `go version`
- **Operating system** - macOS, Linux, Windows, etc.
- **Minimal code sample** - Smallest code that reproduces the issue
- **Expected behavior** - What should happen
- **Actual behavior** - What actually happens
- **Error messages** - Full error output if applicable

**Good bug reports save everyone time!** The template will guide you through this.

### What Happens Next

- We aim to acknowledge bugs within **2-3 business days**
- Critical bugs (crashes, data loss) are prioritized
- You may be asked for additional information
- Once confirmed, we'll add appropriate labels and milestones

---

## Feature Requests

Have an idea for a new feature? Use our **[feature request template](.github/ISSUE_TEMPLATE/feature_request.yml)**.

### Before Requesting

1. **Check existing requests** - Search [enhancement issues](https://github.com/shapestone/shape-json/issues?q=is%3Aissue+label%3Aenhancement)
2. **Review our scope** - See [Format Addition Policy](CONTRIBUTING.md#format-addition-policy)
3. **Consider fit** - Does it align with shape-json's mission (JSON parsing/manipulation)?

### What We Accept

✅ **JSON parsing improvements** - Better error messages, performance, RFC compliance
✅ **API enhancements** - New methods, better ergonomics, type safety
✅ **Performance optimizations** - Speed, memory usage improvements
✅ **Bug fixes** - Always welcome
✅ **Documentation improvements** - Examples, guides, clarifications
✅ **Test coverage** - Additional tests for edge cases

### What We Generally Don't Accept

❌ **New format parsers** - shape-json is JSON-only (see [Shape ecosystem](https://github.com/shapestone/shape) for other formats)
❌ **Breaking API changes** - We follow semantic versioning
❌ **Features outside scope** - Non-JSON-related functionality
❌ **Custom DSLs** - Use Shape's tokenizer framework directly in your own project

See our [Contributing Guide](CONTRIBUTING.md) for details on scope and what we're looking for.

### Response Timeline

- Feature requests are reviewed during planning cycles
- We may ask clarifying questions about use cases
- Not all requests will be accepted (scope, maintenance burden, etc.)
- Rejected requests will have clear explanations

---

## Security Vulnerabilities

**🚨 Do NOT open public issues for security vulnerabilities**

Security issues require private disclosure to protect users.

### How to Report

**Preferred method:**
1. Go to the [Security tab](https://github.com/shapestone/shape-json/security)
2. Click "Report a vulnerability"
3. Fill out the private vulnerability report form

**Alternative method:**
- Email: security@shapestone.com
- Subject: "shape-json Security Issue"

### What to Include

- Clear description of the vulnerability
- Potential impact and attack scenario
- Steps to reproduce
- Affected versions (if known)
- Proof of concept code (if applicable)
- Suggested fix (if you have ideas)

### Our Commitment

- **Acknowledgment**: Within 48 hours
- **Initial assessment**: Within 5 business days
- **Regular updates**: Every 7 days until resolved
- **Patch release**: Within 30 days for high/critical issues

See our complete [Security Policy](SECURITY.md) for details on coordinated disclosure and supported versions.

---

## Response Times

shape-json is an open source project maintained by Shapestone. We aim to respond within these timeframes:

| Type | Response Time | Notes |
|------|---------------|-------|
| **Security issues** | 48 hours | See [SECURITY.md](SECURITY.md) for full policy |
| **Bug reports** | 2-3 business days | Critical bugs prioritized |
| **Feature requests** | Reviewed during planning | May take longer for complex requests |
| **Questions on Discussions** | Best effort | Community-driven, maintainers participate when available |
| **Pull requests** | 3-5 business days | Initial review; may require iterations |

**Note:** These are goals, not guarantees. Response times may vary based on maintainer availability, holidays, and issue complexity.

---

## Contributing

Want to contribute code, documentation, or tests? Fantastic!

See our **[Contributing Guide](CONTRIBUTING.md)** for:
- Development setup instructions
- Code style and testing requirements
- Pull request process
- Branching workflow
- What kinds of contributions we're looking for

Quick links:
- **[Local Development Setup](CONTRIBUTING.md#development-setup)**
- **[Testing Guidelines](CONTRIBUTING.md#testing-guidelines)**
- **[Pull Request Process](CONTRIBUTING.md#pull-request-process)**

---

## Community Guidelines

All interactions in the shape-json community (issues, discussions, PRs) are governed by our **[Code of Conduct](CODE_OF_CONDUCT.md)**.

**Expected behavior:**
- Be respectful and inclusive
- Provide constructive feedback
- Focus on what's best for the project and community
- Show empathy toward other community members

**Unacceptable behavior:**
- Harassment, trolling, or personal attacks
- Public or private harassment
- Publishing others' private information
- Other conduct inappropriate in a professional setting

To report violations, contact: conduct@shapestone.com

---

## Additional Resources

### Learning Resources

- **[Go Documentation](https://go.dev/doc/)** - General Go programming help
- **[JSON RFC 8259](https://datatracker.ietf.org/doc/html/rfc8259)** - Official JSON specification
- **[JSONPath RFC 9535](https://datatracker.ietf.org/doc/html/rfc9535)** - JSONPath specification
- **[Shape Core](https://github.com/shapestone/shape-core)** - Universal AST and tokenizer framework
- **[Shape Ecosystem](https://github.com/shapestone/shape)** - Multi-format parser ecosystem

### Related Projects

- **[shape-core](https://github.com/shapestone/shape-core)** - Core infrastructure for parsers
- **[shape](https://github.com/shapestone/shape)** - Multi-format parser library
- Other Shape parsers: shape-xml, shape-yaml, shape-csv, etc.

---

## What We Don't Support

To keep the project focused and maintainable, we generally don't provide:

❌ **Support for EOL Go versions** - We support Go 1.21+ (see [go.mod](go.mod))
❌ **Custom implementation help** - We can't debug your specific project code
❌ **Third-party integration debugging** - Issues with other libraries that use shape-json
❌ **Format additions** - See our [scope policy](CONTRIBUTING.md#format-addition-policy)
❌ **Non-JSON features** - shape-json is focused on JSON parsing/manipulation

For custom needs:
- Hire a Go consultant for your specific project
- Fork the project if you need custom functionality

---

## Quick Reference

**I want to...**

- ❓ **Ask a question** → [GitHub Discussions](https://github.com/shapestone/shape-json/discussions)
- 🐛 **Report a bug** → [Bug Report Template](.github/ISSUE_TEMPLATE/bug_report.yml)
- 💡 **Request a feature** → [Feature Request Template](.github/ISSUE_TEMPLATE/feature_request.yml)
- 🔒 **Report a security issue** → [Private Vulnerability Reporting](https://github.com/shapestone/shape-json/security)
- 📖 **Learn how to use shape-json** → [USER_GUIDE.md](USER_GUIDE.md)
- 🔧 **Contribute code** → [CONTRIBUTING.md](CONTRIBUTING.md)
- 🏗️ **Understand the architecture** → [ARCHITECTURE.md](ARCHITECTURE.md)
- 📋 **See recent changes** → [CHANGELOG.md](CHANGELOG.md)

---

## Still Need Help?

If you've:
- ✅ Checked the documentation
- ✅ Searched existing issues and discussions
- ✅ Asked on GitHub Discussions
- ✅ Still can't find an answer

Then please open a **[discussion](https://github.com/shapestone/shape-json/discussions)** with:
- What you're trying to accomplish
- What you've already tried
- Specific error messages or unexpected behavior
- Minimal code example

The community and maintainers will do their best to help!

---

**Thank you for being part of the shape-json community!** 🎉

---

*Last Updated: December 21, 2025*
