---
Title: GitHub Pages workflow
Description:
  A basic GitHub Pages build and deploy workflow using GitHub Actions.
Created: 2025-03-22
Modified: 2025-03-22
Tags:
  - GitHub Actions
  - GitHub Pages
Styles:
  - /static/styles/article.css
  - /static/styles/code.css
Scripts:
  - /static/scripts/copy-code.js
---

# {{.Title}}

This is a basic GitHub Pages workflow that builds and deploys a static site using GitHub Actions. This workflow is used on this site and is a good starting point for a simple static site.

## [Triggers](#triggers) {#triggers}

This workflow is triggered on push to the `main` branch and through a manual trigger/dispatch.

```yml
on:
  workflow_dispatch:
  push:
    branches:
      - main
```

## [Permissions](#permissions) {#permissions}

This workflow requires write permissions to the repository via an oidc token as well as GitHub Pages write permissions to request a GitHub Pages build.

```yaml
permissions:
  id-token: write
  pages: write
```

## [Build job](#build) {#build}

The build job does what it says on the tin, it checks out the repository and installs the required dependencies, then builds a binary and runs it to generate the static site. Finally, it uploads the static files as an artifact.

By using the [actions/setup-go](https://github.com/actions/setup-go) action, we can install the required version of Go for the project and get caching of the Go modules used in the project for free.

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    env:
      GITHUB_USERNAME: $\{\{ github.actor \}\}
      GITHUB_TOKEN: $\{\{ secrets.GITHUB_TOKEN \}\}
    steps:
      - uses: actions/checkout@v4.2.2
      - uses: actions/setup-go@v5.3.0
        with:
          go-version: "1.23.4"
      - name: Build Go binary
        run: go build -o ./bin/main .
      - name: Generate static site with Go binary
        run: ./bin/main
      - name: Upload static files as artifact
        id: deployment
        uses: actions/upload-pages-artifact@v3.0.1
        with:
          path: build/

  # ...
```

### Aside

I do plan on refactoring/extending this at some point to cache the binary and reuse it in a scheduled job to re-generate the static site instead of discarding it after each build. This would allow for a faster build and deploy process and would allow the fetched data ingested by the program to be up to date without needing to rebuild the entire binary. But that's a problem for another day.

## [Deploy job](#deploy) {#deploy}

This job deploys the static site to GitHub Pages using the artifact generated in the build job.

This step uses the concurrency and environment features of GitHub Actions to ensure that only one deployment is waiting for approval at a time and to provide a link to the deployed site in the GitHub UI.

```yaml
jobs:
  # ...

  deploy:
    runs-on: ubuntu-latest
    needs: build
    environment:
      name: github-pages
      url: $\{\{ steps.deployment.outputs.page_url \}\}
    concurrency:
      group: $\{\{ github.workflow \}\}-$\{\{ github.ref \}\}-deploy
      cancel-in-progress: true
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Full Workflow

Full raw workflow can be found [here](./deploy.yml).
