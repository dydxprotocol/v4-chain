### Changelist
[Describe or list the changes made in this PR]

### Test Plan
[Describe how this PR was tested (if applicable)]

### Author/Reviewer Checklist
- [ ] If this PR has changes that result in a different app state given the same prior state and transaction list, manually add the `state-breaking` label.
- [ ] If the PR has breaking postgres changes to the indexer add the `indexer-postgres-breaking` label.
- [ ] If this PR isn't state-breaking but has changes that modify behavior in `PrepareProposal` or `ProcessProposal`, manually add the label `proposal-breaking`.
- [ ] If this PR is one of many that implement a specific feature, manually label them all `feature:[feature-name]`.
- [ ] If you wish to for mergify-bot to automatically create a PR to backport your change to a release branch, manually add the label `backport/[branch-name]`.
- [ ] Manually add any of the following labels: `refactor`, `chore`, `bug`.
