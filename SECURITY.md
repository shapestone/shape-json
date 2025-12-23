# Security Policy

## Supported Versions

We actively support Shape with security updates:

- **Latest Release:** Fully supported with security updates
- **Previous Minor Version:** Supported for critical security fixes (30 days)
- **Older Versions:** No longer supported

We strongly recommend always using the latest release for the best security posture and newest features.

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue in Shape, please report it responsibly.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please use one of the following methods:

1. **GitHub Private Vulnerability Reporting** (Preferred)
   - Navigate to the [Security tab](https://github.com/shapestone/shape/security) of this repository
   - Click "Report a vulnerability"
   - Fill out the vulnerability report form

2. **Email**
   - Send details to: security@shapestone.com
   - Include "Shape Security" in the subject line

### What to Include

Please provide the following information in your report:

- **Description**: Clear description of the vulnerability
- **Impact**: Potential impact and attack scenario
- **Reproduction**: Step-by-step instructions to reproduce the issue
- **Affected Versions**: Which versions are affected
- **Proof of Concept**: Code snippets or test cases (if applicable)
- **Suggested Fix**: If you have ideas for remediation

### Response Timeline

We are committed to responding promptly to security reports:

- **Acknowledgment**: Within 48 hours of receiving your report
- **Initial Assessment**: Within 5 business days
- **Regular Updates**: Every 7 days until resolved
- **Patch Release**: Within 30 days for high/critical severity issues

### Disclosure Policy

- We follow a **coordinated disclosure** process
- We will work with you to understand and validate the issue
- Once a fix is ready, we will:
  1. Release a patch version
  2. Publish a security advisory (GitHub Security Advisory)
  3. Credit you in the advisory (unless you prefer to remain anonymous)
- We request that you do not publicly disclose the vulnerability until we have released a fix

## Security Considerations

### Known Attack Surfaces

Shape is a parser library that processes untrusted input. Be aware of these potential security considerations:

1. **Denial of Service (DoS)**
   - Malformed input with deeply nested structures
   - Extremely large files or payloads
   - Recursive or circular references

2. **Resource Exhaustion**
   - Memory consumption from large ASTs
   - CPU consumption from complex validation rules

3. **Input Validation**
   - Always validate and sanitize data after parsing
   - Set appropriate limits on input size and complexity

### Best Practices

When using Shape in production:

- **Set Timeouts**: Implement timeouts for parsing operations
- **Limit Input Size**: Enforce maximum file size limits
- **Validate Output**: Don't trust parsed data without validation
- **Monitor Resources**: Track memory and CPU usage
- **Keep Updated**: Regularly update to the latest version

## Security Updates

Security updates are released as patch versions and announced via:

- GitHub Security Advisories
- Release notes in CHANGELOG.md
- GitHub Releases page

Subscribe to releases or watch this repository to stay informed.

## Recognition

We appreciate the security research community's efforts. Researchers who responsibly disclose vulnerabilities will be credited in:

- The security advisory (with their permission)
- A SECURITY_ACKNOWLEDGMENTS.md file (if applicable)
- Release notes for the security fix

## Questions?

If you have questions about this security policy or Shape's security posture, please open a public discussion in GitHub Discussions or contact security@shapestone.com.

---

Thank you for helping keep Shape and its users safe!
