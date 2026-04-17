# 📥 How to use the Binary

If you choose the binary route to save those extra seconds of container startup time, use the following snippets for your respective platforms.

## 🐙 GitHub Actions
In GitHub, you typically run this after a step like **Release Drafter** to catch the output variables.

```yaml
- name: 📥 Download Changelog Updater
  run: |
    curl -sSL https://github.com/bugs5382/changelog-updater-action/releases/latest/download/updater-linux-amd64 -o updater 
    chmod +x updater

- name: 🖊️ Run Update
  run: ./updater --version="${{ steps.drafter.outputs.tag_name }}" --notes="${{ steps.drafter.outputs.body }}"
```

## 🦊 GitLab CI/CD
In GitLab, you’ll usually pull the binary in the `before_script` or a specific job step. Since GitLab uses standard environment variables, the syntax is even cleaner.

```yaml
- name: 📥 Download Changelog Updater
  before_script:
    - curl -sSL https://github.com/bugs5382/changelog-updater-action/releases/latest/download/updater-linux-amd64 -o updater
    - chmod +x updater

- name: 🖊️ Run Update
  script:
    - ./updater --version="$TAG_NAME" --notes="$RELEASE_NOTES"
```

## 🏎️ Why the Binary?
* **Instant Execution:** Zero overhead from Docker daemon startup.
* **Minimal Footprint:** No need to pull heavy layers; just a single, small **Go** executable.
* **Standardized:** The exact same binary runs across any Linux-based runner, ensuring consistent behavior whether you're on GitHub, GitLab, or a local Jenkins node.