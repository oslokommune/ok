# OK Tool MCP Server Roadmap

This document outlines planned MCP tools for the ok CLI, ranked by implementation priority and complexity.

## Current Status

✅ **Implemented:**
- `hello_world` - Simple greeting tool (demo/testing)
- `version` - Get ok tool version information

## Planned Tools (Ranked)

### Tier 1: High Value, Low Complexity
*Simple read-only operations with immediate utility*

1. **`list_latest_releases`** ⭐⭐⭐
   - **Purpose**: Get latest versions of all boilerplate packages
   - **Implementation**: Use `githubreleases.GetLatestReleases()`
   - **Input**: Optional organization filter
   - **Output**: JSON list of {component, version} pairs
   - **Value**: Quick version checking for dependency management

2. **`read_package_manifest`** ⭐⭐⭐
   - **Purpose**: Parse and display package manifest contents
   - **Implementation**: Use `common.LoadPackageManifest()`
   - **Input**: Path to packages.yml (defaults to current directory)
   - **Output**: Structured package configuration
   - **Value**: Essential for understanding project structure

3. **`validate_package_structure`** ⭐⭐
   - **Purpose**: Check if directory uses valid package structure
   - **Implementation**: Use existing validation logic
   - **Input**: Directory path (defaults to current directory)
   - **Output**: Validation results and recommendations
   - **Value**: Project health checking

### Tier 2: Good Value, Medium Complexity
*More complex operations requiring parameters or business logic*

4. **`get_package_info`** ⭐⭐
   - **Purpose**: Get detailed information about specific packages
   - **Implementation**: Combine manifest loading with GitHub API calls
   - **Input**: Package name or path
   - **Output**: Package details, available versions, dependencies
   - **Value**: Deep-dive package analysis

5. **`generate_aws_config`** ⭐⭐
   - **Purpose**: Generate AWS CLI configuration for IAM Identity Center
   - **Implementation**: Use `aws.config.Generate()`
   - **Input**: Profile configuration parameters
   - **Output**: AWS CLI config snippet
   - **Value**: DevOps automation

6. **`download_github_file`** ⭐
   - **Purpose**: Download files from boilerplate repositories
   - **Implementation**: Use `githubreleases.DownloadGithubFile()`
   - **Input**: Repository, file path, destination
   - **Output**: Download status and location
   - **Value**: Template file retrieval

### Tier 3: Specialized Use Cases
*Complex or niche functionality for advanced users*

7. **`resolve_template_paths`** ⭐
   - **Purpose**: Resolve template paths within package structure
   - **Implementation**: Use existing path resolution logic
   - **Input**: Template reference, context
   - **Output**: Resolved absolute paths
   - **Value**: Advanced templating support

8. **`parse_package_versions`** ⭐
   - **Purpose**: Parse and compare semantic versions from package refs
   - **Implementation**: Use semver parsing utilities
   - **Input**: Version strings or package references
   - **Output**: Parsed version objects and comparisons
   - **Value**: Version management utilities

## Implementation Guidelines

### Design Principles
- **Read-only first**: Prioritize tools that don't modify system state
- **Simple inputs**: Minimize required parameters, provide sensible defaults
- **Structured output**: Return JSON for programmatic consumption
- **Error handling**: Graceful failure with helpful error messages
- **Documentation**: Clear descriptions and examples

### Technical Considerations
- Reuse existing ok tool functionality where possible
- Handle GitHub API rate limiting gracefully
- Support both authenticated and unauthenticated GitHub access
- Provide consistent error response format
- Include validation for all inputs

### Testing Strategy
- Unit tests for each tool handler
- Integration tests with real GitHub API (rate-limited)
- Mock GitHub responses for CI/CD
- Test error conditions and edge cases

## Future Enhancements

### Potential Advanced Features
- **Batch operations**: Process multiple packages simultaneously
- **Caching**: Local cache for GitHub API responses
- **Webhooks**: Real-time notifications for package updates
- **Diffing**: Compare package versions and configurations
- **Dependency analysis**: Identify package dependency chains

### Integration Opportunities
- **IDE extensions**: Enhanced package management in editors
- **CI/CD pipelines**: Automated package validation and updates
- **Documentation generation**: Auto-generate package documentation
- **Monitoring**: Track package usage and health across projects

---

*This roadmap prioritizes practical, immediately useful tools while building toward more sophisticated package management capabilities.*