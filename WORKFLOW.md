## Resource Development Workflow

### Getting Started

Before you begin contributing to the development of the Meraki Terraform provider, it's important to set up a proper workflow. This document provides a detailed explanation of the steps you should follow to contribute effectively.

### Workflow Steps

#### 1. Fork and Clone the Repository
Firstly, fork the main repository into your own GitHub account and then clone it locally to your development environment.

```sh
git clone https://github.com/core-infra-svcs/terraform-provider-meraki.git
cd terraform-provider-meraki
```

#### 2. Create a New Branch
For each new feature or bug fix, create a new branch from the `main` branch. Use a descriptive name for your branch that reflects the change.

```sh
git checkout -b feature/add-new-resource
```

#### 3. Code Your Changes
Implement your changes locally, adhering to the coding standards and guidelines provided in the repository's `README.md`. Ensure you write or update unit tests to cover the new functionality.

#### 4. Test Your Changes
Make sure to test your changes thoroughly. This could involve:

- Running existing tests to ensure they pass with your changes.
- Testing the functionality in a live Meraki environment, if possible.

#### 5. Commit and Push Your Changes
Once you are satisfied with your changes and all tests pass, commit your changes to your local branch and push the branch to your GitHub fork.

```sh
git add .
git commit -m "Add new resource for XYZ feature"
git push origin feature/add-new-resource
```

#### 6. Open a Pull Request
From your fork on GitHub, open a pull request to the `main` branch of the main repository. Fill out the pull request template, clearly describing the changes and any pertinent details.

#### 7. Code Review
Once your pull request is open, other contributors and maintainers can review your code. They may provide feedback or request changes.

#### 8. Incorporate Feedback
If changes are requested, make them in your branch, then push the updates. Your pull request will automatically update with the new commits.

#### 9. Final Review and Merge
After your pull request has been approved by the reviewers, a project maintainer will merge your changes into the main codebase.

#### 10. Clean Up
After your changes have been merged, you can delete your local and remote feature branches.

### Project Board and Selecting a Resource to Develop

#### 1. Examine the Project Board
Go to the [project board](https://github.com/orgs/core-infra-svcs/projects/1/views/1) to find out what's currently being worked on and what's in the pipeline.

#### 2. Pick a Resource
Look for resources or tasks that are not assigned and align with your skills and interests.

#### 3. Comment on the Task
Leave a comment on the task you're interested in to let others know you are working on it.

#### 4. Start Contributing
Follow the development workflow to start contributing.

### Remember
- Always keep your fork's `main` branch in sync with the upstream `main` branch.
- Do not merge your own pull requests. A second pair of eyes on your code is invaluable, even for seasoned developers.

With this workflow, you can now start contributing to the Meraki Terraform provider with confidence. Your contributions are invaluable to the project's success, and we appreciate your involvement!