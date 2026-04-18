# 🦊 GitLab CI/CD Integration

The **Changelog Updater Action** is highly portable. Since GitLab doesn't use a "Marketplace" system, you can integrate this tool either as a native Docker image or by calling the Go binary directly.

## ⚡ Option 1: Using the Binary (Maximum Performance)
This approach is the fastest. It downloads the pre-compiled **Go** binary and executes it within your existing environment. Perfect for lightweight `alpine` or `ubuntu` runners.

```yaml
Update_ChangeLog:
  stage: deploy
  image: alpine:latest
  before_script:
    - 📦 apk add --no-cache curl git
    - 📥 curl -sSL https://github.com/bugs5382/changelog-updater-action/releases/latest/download/changelog-updater-action-linux-amd64 -o changelog-updater-action
    - 🔑 chmod +x changelog-updater-action
  script:
    # Determines your version and notes via environment variables
    - 🖊️ ./changelog-updater-action --tag="$TAG_NAME" --notes="$RELEASE_NOTES"
    
    # Commit and push back to your GitLab repo
    - git config --global user.email "ci@gitlab.com"
    - git config --global user.name "GitLab CI"
    - git add CHANGELOG.md
    - git commit -m "chore(pre-release): $TAG_NAME [skip ci]"
    - git push https://oauth2:${PROJECT_ACCESS_TOKEN}@gitlab.com/${CI_PROJECT_PATH}.git HEAD:${CI_COMMIT_REF_NAME}
```

## 🐳 Option 2: Using the Docker Image (Cleanest Setup)
Since the tool is packaged as a container, you can run your job directly inside the image. This eliminates the need for manual downloads or environment setup.

```yaml
Update_ChangeLog:
  stage: deploy
  image: bugs5382/changelog-updater-action:latest
  script:
    # The entrypoint is the Go binary—simply pass your flags!
    - 🚀 /changelog-updater-action --tag="$TAG_NAME" --notes="$RELEASE_NOTES"
```

## 🎛️ Flags

| Flag        | Short | Description                                                      | Default |
|-------------|-------|------------------------------------------------------------------|---------|
| `--tag`     | `-t`  | Release tag name, e.g. `v1.2.0`. **Required.**                   | —       |
| `--notes`   | `-n`  | Release notes body (markdown). **Required.**                     | —       |
| `--path`    | `-p`  | Directory (relative to the repo root) containing `CHANGELOG.md`. | `.`     |
| `--date`    |       | Release date injected into the version header (`YYYY-MM-DD`).    | today   |
| `--diff`    |       | Show the diff (if any) of changes.                               | `false` |
| `--dry`     |       | Dry run — parse and log without modifying `CHANGELOG.md`.        | `false` |
| `--verbose` | `-v`  | Enable debug level logging.                                      | `false` |


## 💡 Key Differences in GitLab

* **🔐 Permissions:** Unlike GitHub’s automatic `GITHUB_TOKEN`, GitLab requires a **Project Access Token** or **Deploy Token** to push changes back to the repository.
* **🏎️ Speed:** Because this is built in **Go**, the actual processing happens in milliseconds. Your only overhead is the brief moment it takes to pull the image or `curl` the binary.
* **🔧 Flexibility:** This is a pure CLI tool at heart. It accepts standard inputs and flags, meaning you aren't locked into specific platform syntax—it just works.