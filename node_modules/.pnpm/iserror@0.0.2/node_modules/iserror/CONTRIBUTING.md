# Contributing
Please take a moment to review this document in order to make the contribution
process easy and effective for everyone involved.

Following these guidelines helps to communicate that you respect the time of
the developers managing and developing this project. In return, they should
reciprocate that respect in addressing your issue or assessing patches and
features.

## Bug reports
A bug is a demonstrable problem that is caused by the code in the repository.

A good bug report shouldn't leave others needing to chase you up for more
information. Please try to be as detailed as possible in your report. What is
your environment? What steps will reproduce the issue? What browser(s) and OS
experience the problem? What would you expect to be the outcome? All these
details will help people to fix any potential bugs.

Example:

> Bug issue type
>
> Short and descriptive example bug report title
>
> A summary of the issue and the browser/OS environment in which it occurs. If
> suitable, include the steps required to reproduce the bug.
>
> 1. This is the first step
> 2. This is the second step
> 3. Further steps, etc.
>
> Any other information you want to share that is relevant to the issue being
> reported. This might include the lines of code that you have identified as
> causing the bug, and potential solutions (and your opinions on their
> merits).

## Feature requests
Feature or improvement requests are welcome. But take a moment to find out
whether your idea fits with the scope and aims of the project. It's up to you
to make a strong case to convince the project's developers of the merits of
this feature or improvement. Please provide as much detail and context as
possible.

Example:

> New Feature or Improvement issue type
>
> Short and descriptive example feature request title
>
> A summary of the issue and the impact of new feature or improvement.
>
> Any other information you want to share that is relevant to the new feature
> or improvement.


## Pull requests
Good pull requests - patches, improvements, new features. They should remain
focused in scope and avoid containing unrelated commits.

Please adhere to the coding conventions used throughout a project (whitespace,
accurate comments, etc.) and any other requirements (such as lint and test
coverage).

Follow this process to make your work considered for inclusion in the project:

Follow this process to make your work considered for inclusion in the project:

1. [Fork](http://help.github.com/fork-a-repo/) the project, clone your fork,
   and configure the remotes:

   ```bash
   # Clone your fork of the repo into the current directory
   git clone https://github.com/<your-username>/iserror
   # Navigate to the newly cloned directory
   cd iserror
   # Assign the original repo to a remote called "upstream"
   git remote add upstream https://github.com/yefremov/iserror
   ```

2. If you cloned a while ago, get the latest changes from upstream (and other
   remotes):

   ```bash
   git fetch --all
   git checkout master
   git pull upstream master
   ```

3. Never work directly on `master`. Create a new topic branch (off the latest
   version of `master`) to contain your feature, improvement or a bug fix:

   ```bash
   git checkout -b <topic-branch-name>
   ```

4. Commit your changes in logical chunks. Please adhere to these [git commit
   message conventions](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
   or your code is unlikely be merged into the main project. Use Git's
   [interactive rebase](https://help.github.com/articles/interactive-rebase)
   feature to tidy up your commits before making them public.

5. Locally rebase the upstream development branch into your topic branch:

   ```bash
   git pull --rebase upstream master
   ```

6. Push your topic branch up to your fork:

   ```bash
   git push origin <topic-branch-name>
   ```

7. [Open a Pull Request](https://help.github.com/articles/using-pull-requests/)
   with a clear title and description.
