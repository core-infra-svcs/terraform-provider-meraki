# Contributing

## Welcome
We are thrilled that you're interested in contributing to this project. This is a unique space where individuals who are entry-level in network engineering or application development can make a significant impact. The Meraki Terraform Provider is the first declarative tool in the ecosystem, designed to empower teams and individuals to accelerate their automation strategy with confidence and simplicity.

By contributing to this project, you have the opportunity to embark on an exciting learning journey. You'll gain in-depth knowledge about the Cisco Meraki Platform and Product line, enabling you to understand the intricacies of network engineering and application development in a practical setting. Additionally, you'll have the chance to hone your software engineering skills with Golang, a powerful and popular programming language, and practice proper git hygiene, which is crucial for collaborative development.

Your contributions will not only shape the future of the Meraki Terraform Provider but also contribute to the growth and success of the entire ecosystem. We value your unique perspective and ideas, and we believe that together, we can create a robust and user-friendly tool that revolutionizes automation in network engineering and application development.

Thank you for considering contributing to the Meraki Terraform Provider.

## Who is a Contributor?

A contributor is anyone who is willing to dedicate their time and skills to make a positive impact on the OpenSource community.

## What Counts as a Contribution?

Anything of material benefit to our community efforts is considered a contribution. While development talent is highly sought after, there are many ways to make a meaningful impact.
Providing feedback, opening issues, suggesting enhancements, and advocating for new features are all excellent ways to contribute.

We are eager to see the innovative solutions and insights you will bring. Let's build something extraordinary together!

## Getting Started

Anyone interested in the progress of this provider is encouraged to follow the [project status board](https://github.com/orgs/core-infra-svcs/projects/1) for the latest updates.

## Onboarding

As a developer, your contributions are crucial to the growth of the project. Below is a list of resources designed to get you started:

- **[Development Workflow](./WORKFLOW.md):** This document outlines the steps to get your development environment up and running on your local machine.

- **[Local Development Environment Setup Guide](./DEVELOPMENT.md):** Here you'll find detailed instructions for configuring your machine for Meraki Terraform provider development.

- **[Integration Testing](./TESTING.md):** Learn how to set up a test environment in the Meraki dashboard to ensure your contributions work as expected in a live setting.

- **[Troubleshooting](https://github.com/core-infra-svcs/terraform-provider-meraki/wiki):** Access a collection of knowledge base articles in our Wiki to assist in resolving common issues you might encounter during provider development.


Should you have any questions or require assistance throughout your journey, do not hesitate to reach out. Our community is supportive, and we aim to foster a welcoming environment where everyone feels encouraged to contribute.

Happy contributing!

## How to Contribute

For information on how to set up your environment, write code, and submit pull requests, check out our [Development Guide](./DEVELOPMENT.md).

### Reporting Bugs or Issues
- Use the issue tracker to report bugs.
- Describe the bug and include additional details to help maintainers reproduce the problem.
- Follow the template provided for bug reports.

### Suggesting Enhancements
- Open an issue with a tag suggesting an enhancement.
- Clearly describe the feature and its benefits.
- Be ready to discuss its implementation and impact.

### Pull Requests

1. Fork and Clone the Repository
   Firstly, fork the main repository into your own GitHub account and then clone it locally to your development environment.

```sh
git clone https://github.com/core-infra-svcs/terraform-provider-meraki.git
cd terraform-provider-meraki
```

2. Create a New Branch
   For each new feature or bug fix, create a new branch from the `main` branch. Use a descriptive name for your branch that reflects the change.

```sh
git checkout -b feature/add-new-resource
```

3. Code Your Changes
   Implement your changes locally, adhering to the coding standards and guidelines provided in the repository's `README.md`. Ensure you write or update unit tests to cover the new functionality.

4. Test Your Changes
   Make sure to test your changes thoroughly. This could involve:

- Running existing tests to ensure they pass with your changes.
- Testing the functionality in a live Meraki environment, if possible.

5. Commit and Push Your Changes
   Once you are satisfied with your changes and all tests pass, commit your changes to your local branch and push the branch to your GitHub fork.

```sh
git add .
git commit -m "Add new resource for XYZ feature"
git push origin feature/add-new-resource
```

6. Open a Pull Request
   From your fork on GitHub, open a pull request to the `main` branch of the main repository. Fill out the pull request template, clearly describing the changes and any pertinent details.

7. Code Review
   Once your pull request is open, other contributors and maintainers can review your code. They may provide feedback or request changes.

8. Incorporate Feedback
   If changes are requested, make them in your branch, then push the updates. Your pull request will automatically update with the new commits.

9. Final Review and Merge
   After your pull request has been approved by the reviewers, a project maintainer will merge your changes into the main codebase.

10. Clean Up
    After your changes have been merged, you can delete your local and remote feature branches.



## Code of Conduct

In the interest of fostering an open and welcoming environment, we expect contributors to adhere to our [Code of Conduct](CODE_OF_CONDUCT.md) at all times.

## Questions or Need Help?

Connecting with the community is a part of the collaborative process. If you're looking for more information, news, and updates join the [meraki community forum](https://community.meraki.com)

## Resources

To effectively contribute to the terraform-provider-meraki, you may need a solid foundation in both Cisco Meraki solutions and Terraform. Below are curated resources to help build that knowledge base:

### Meraki

For those new to the Cisco Meraki ecosystem:

1. **Cisco Meraki Product Documentation:** Start with the [official documentation](https://documentation.meraki.com/) to understand the Meraki product suite.

2. **Cisco Meraki Dashboard API Documentation:** Review the [API documentation](https://developer.cisco.com/meraki/api-v1/) for details on interfacing with Meraki via the Dashboard API.

3. **Cisco Meraki Dashboard OpenAPI Specification:** Familiarize yourself with the [OpenAPI Spec](https://github.com/meraki/openapi) to understand the full range of API functionalities.

4. **Postman Collection:** Utilize the [Meraki API Postman collection](https://documenter.getpostman.com/view/897512/SzYXYfmJ#6a9be9f0-49e9-4644-b93d-5a82f09d9899) for testing and troubleshooting API requests.

5. **Dashboard-api-go:** Explore the [dashboard-api-go](https://github.com/meraki/dashboard-api-go) client, a tool for making API requests to Meraki in Go, which our provider uses for resource and data source interactions.

### Terraform (Power User)

For network engineers new to Terraform:

1. **Terraform Documentation:** The [official Terraform documentation](https://developer.hashicorp.com/terraform/docs) is the best place to start learning about Terraform.

2. **Terraform Tutorials:** Engage with hands-on [tutorials](https://developer.hashicorp.com/terraform/tutorials?product_intent=terraform) to get practical experience.

3. **Examples:** Review the [examples folder](../examples) in this repository for real-world usage of the Meraki provider.

4. **Terraform Registry:** Explore other providers and modules in the [Terraform Registry](https://registry.terraform.io/browse/providers).

### Terraform (Developer)

For developers new to Terraform provider development:

1. **Developer Documentation:** Dive into the [Plugin Development Documentation](https://developer.hashicorp.com/terraform/plugin) to understand the ins and outs of Terraform provider development.

2. **HandsOn Tutorial:** The [HashiCup tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) is an interactive resource to familiarize yourself with the Terraform Plugin Framework.

Armed with these resources, you're now ready to deepen your knowledge and start contributing to the terraform-provider-meraki. Happy learning, and we look forward to your innovative contributions!
