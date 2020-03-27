# Comment on PR via GitHub Action

A GitHub action to comment on the open PR when a commit is pushed.

## Usage

- Requires the `GITHUB_TOKEN` secret.
- Requires the comment's message in the `msg` parameter.
- Supports `push` and `pull_request` event types.

### Example

```
name: cmtpr example
on: pull_request
jobs:
  example:
    name: example
    runs-on: ubuntu-latest
    steps:
      - name: comment PR
        uses: khezen/cmtpr@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          msg: "Looks good to me!"
```
