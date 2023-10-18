# Usage:
# GITHUB_TOKEN=<token> python main.py --old <old commit> --new <new commit> --path <path>
#
# The above command will emit markdown release notes down to stdout for changes between the <old commit>
# and the <new commit>. Only commits that change files in <path> will be included (e.g. use path=protocol
# for protocol releases). A github token is required to avoid rate limits.

import argparse
import os
import requests

GITHUB_TOKEN_ENV_VAR = "GITHUB_TOKEN"
COMMITS_ENDPOINT = "https://api.github.com/repos/dydxprotocol/v4-chain/commits"
GET_COMMIT_ENDPOINT = "https://api.github.com/repos/dydxprotocol/v4-chain/commits/%s"
LIST_COMMIT_PULLS_ENDPOINT = "https://api.github.com/repos/dydxprotocol/v4-chain/commits/%s/pulls"


def get_commit(session, commit_sha):
    r = session.get(GET_COMMIT_ENDPOINT % commit_sha)
    r.raise_for_status()
    return r.json()

# Return a string that is a markdown list entry. It should summarize the changes in this commit.
def commit_to_entry(session, commit_json):
    r = session.get(LIST_COMMIT_PULLS_ENDPOINT % commit_json["sha"])
    r.raise_for_status()
    pulls = r.json()

    if len(pulls) > 1:
        # I don't know when a commit would have more than one PR. If this happens, report a bug.
        raise NotImplementedError("Commit unexpectedly had more than one associated PR. Please file a bug.")
    elif len(pulls) == 1:
        # If there is an associated PR, use PR title and link to PR.
        pull = pulls[0]
        return "- %s ([#%s](%s))" % (pull["title"], pull["number"], pull["html_url"])
    else:
        # If there is no associated PR, use first line of commit message and link to commit.
        # This may happen if cherry-picks are pushed directly to a release branch.
        return "- %s ([%s](%s))" % (
            commit_json["commit"]["message"].partition('\n')[0],
            commit_json["sha"][:7],
            commit_json["html_url"]
        )


def get_release_notes(session, new, old, path):
    old_json = get_commit(session, old)
    r = session.get(COMMITS_ENDPOINT, params={
        "sha": new,
        "since": old_json["commit"]["author"]["date"],
        "path": path,
    })
    r.raise_for_status()

    # old commit will be the last commit iff it changes files in path. If it's there, remove it.
    commits = r.json()
    if commits[-1]["sha"] == old_json["sha"]:
        commits = commits[:-1]

    ret = []
    for commit in commits:
        ret.append(commit_to_entry(session, commit))
    return "\n".join(ret)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--old", help="Oldest commit, exclusive")
    parser.add_argument("--new", help="Newest commit, inclusive")
    parser.add_argument("--path", help="Path to filter commits by")
    args = parser.parse_args()

    if not os.getenv(GITHUB_TOKEN_ENV_VAR):
        raise ValueError("Environment variable %s is not set." % GITHUB_TOKEN_ENV_VAR)

    session = requests.Session()
    session.headers.update({'Authorization': 'Bearer %s' % os.getenv(GITHUB_TOKEN_ENV_VAR)})

    print(get_release_notes(session, args.new, args.old, args.path))
