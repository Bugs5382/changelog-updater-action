# 📝 Changelog Updater Action

> The **Changelog Updater Action** ensures your `CHANGELOG.md` is always in sync with your upcoming release. By updating the file *during* the CI process, your documentation matches the tag version and source code exactly at the moment of release. 🚀

## 💡 The Philosophy
Most tools update the changelog **after** the tag is created, which means the version you just released contains a changelog that doesn't actually mention itself! 🔄

This action flips the script:
1. 🏗️ **Anticipate:** Use a tool like **Release Drafter** to determine the next version and compile notes.
2. ✍️ **Update:** This action injects those notes into a new version block in your `CHANGELOG.md`.
3. 💾 **Commit:** The updated file is committed back to the repo *before* the final tag.
4. 🏷️ **Tag:** Your release tag now points to a commit that already includes its own history.

## ⚡ Why this Action?
This project is a high-performance successor to `stefanzweifel/changelog-updater-action`.

* **Powered by Go:** While the original implementation is written in PHP, this version is written in **Go**. 🏎️
* **Efficiency:** By using Go, the action benefits from faster execution times and more efficient memory usage when processing large Markdown files.
* **Streamlined CI:** Reduced execution time means faster feedback loops in your GitHub Actions pipelines.

## 🛠️ Implementation Example

This workflow integrates with **Release Drafter** — the gold standard for drafting releases — to pull anticipated version names and bodies directly into your file.

```yaml
jobs:
  Update_On_Main:
    name: 🚀 Update Changelog & Prep Release
    runs-on: ubuntu-latest
    steps:
      - name: 📂 Checkout Code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: 📝 Release Drafter (Anticipate Version)
        id: drafter
        uses: release-drafter/release-drafter@v6
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      # This is our script action!
      - name: 📤 Changelog Action
        uses: Bugs5382/changelog-updater-action@v0.3.2 # This might not be the latest version!
        env:
          RELEASE_NOTES: >-
            ${{ steps.drafter.outputs.body }}
          RELEASE_VERSION: v${{ steps.drafter.outputs.resolved_version }}
        with:
          tag: ${{ env.RELEASE_VERSION }}
          notes: ${{ env.RELEASE_NOTES }}

      - name: 📤 Commit and Push Version Update
        uses: stefanzweifel/git-auto-commit-action@v4
        with:
          commit_message: "chore(pre-release): v${{ steps.drafter.outputs.resolved_version }} [skip ci]"
```

View this projects ``job-release-and-version-example.yaml`` inside the [examples](examples) folder.

## 🌟 Key Benefits
* **No Stale Logs:** Your `CHANGELOG.md` is never one version behind. 📉
* **Automation:** Eliminates the manual "forgot to update the changelog" commit. 🤖
* **Clean History:** Keeps the rest of your history intact while only modifying the header block. 🧹
* **CI Friendly:** Uses `[skip ci]` in the commit message to prevent recursive workflow loops. 🔃

## 🚀 Versatile Distribution

Whether you want ease of use or raw execution speed, we've got you covered. This action is distributed in two formats:

### 🐳 Docker Container (Recommended)

**Best for: GitHub Actions & GitLab CI/CD**

* **Zero Setup:** No need to install Go or manage dependencies.
* **Plug & Play:** Works instantly with the `uses:` syntax in GitHub or as a `services/image` in GitLab.
* **Isolated:** Environment-agnostic and won't conflict with other tools in your pipeline.

### ⚡ Pre-compiled Binary

**Best for: High-performance pipelines & Local CLI use**

* **Blazing Fast:** Skip the container overhead. Ideal for large-scale monorepos with massive `CHANGELOG.md` files.
* **Portable:** Download the binary directly from our [Releases page](https://www.google.com/search?q=%23) for Linux, macOS, or Windows.
* **Scriptable:** Perfect for local development hooks or custom CI runners where you want to call the tool directly.

#### 🎛️ Flags

Review [flags](examples/FLAGS.md).

### Examples

Review the [examples](examples) folder for more information.

## 🤝 Contributing

We welcome Pull Requests! Please follow these steps:

* **✅ Validation:** Run `make lint` to verify code quality.
* **🧪 Testing:** New features must include unit tests.
* **✍️ Security:** All commits must be **signed** (GPG/SSH).

## ❤️ Acknowledgments

* **[stefanzweifel](https://github.com/stefanzweifel):** For the original foundation of this project.
* **Family:** A special thanks to my wife, daughter, and son for their patience while I work in "geek mode."

## 📄 License

This project is licensed under the **ISC**. See the [LICENSE](LICENSE) file for details.