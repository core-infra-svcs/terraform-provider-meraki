# Contributing


## Welcome
We are thrilled that you're interested in contributing to this project. This is a unique space where individuals who are entry-level in network engineering or application development can make a significant impact. The Meraki Terraform Provider is the first declarative tool in the ecosystem, designed to empower teams and individuals to accelerate their automation strategy with confidence and simplicity.

By contributing to this project, you have the opportunity to embark on an exciting learning journey. You'll gain in-depth knowledge about the Cisco Meraki Platform and Product line, enabling you to understand the intricacies of network engineering and application development in a practical setting. Additionally, you'll have the chance to hone your software engineering skills with Golang, a powerful and popular programming language, and practice proper git hygiene, which is crucial for collaborative development.

Your contributions will not only shape the future of the Meraki Terraform Provider but also contribute to the growth and success of the entire ecosystem. We value your unique perspective and ideas, and we believe that together, we can create a robust and user-friendly tool that revolutionizes automation in network engineering and application development.

Thank you for considering contributing to the Meraki Terraform Provider. We can't wait to see the innovative solutions and insights you bring to the table. Let's build something extraordinary together!

If you have any questions or need assistance along the way, please don't hesitate to reach out. Happy contributing!


## Workflow Overview

To ensure the quality and reliability of the codebase, we have established a pull request workflow that involves peer review and change approval. The following guidelines apply:

1. Fork the repository and create a new branch for your changes.

2. Make your desired changes in the new branch, ensuring that you follow the coding conventions and guidelines of the project.

3. Once you are ready to submit your changes, open a pull request from your branch to the main repository's main branch.

4. Peer reviewers will review your code, provide feedback, and suggest any necessary improvements or changes.

5. Make the necessary revisions based on the feedback received during the peer review.

6. Once all the review comments have been addressed, one of the project maintainers will perform the final review and merge your changes into the development branch.

7. From that point on, your change will follow the process described in the README.md in the `.github/workflows` directory.


## Getting Started with Terraform Development

If this is your first time developing a Terraform provider, or leveraging the new [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) we recommend the following as prerequisite resources to get you up to speed:

1. **Developer Documentation:** Explore the official [Plugin Development Documentation](https://developer.hashicorp.com/terraform/plugin) to gain in-depth knowledge of Terraform and its ecosystem. This documentation provides detailed guides, tutorials, and references that cover various aspects of Terraform usage and development.

2. **HandsOn Tutorial:** A fun and interactive [HashiCup tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) that showcases the capabilities of Terraform in a practical and hands-on manner.

By leveraging these resources, you can deepen your understanding of Terraform, explore best practices, and enhance your contributions to the Meraki Terraform provider. Happy learning and coding!


## Setup your Local Development Environment

See the [Local Development Environment Setup](.github/workflow-docs/local-development-setup.md) instructions.

### Project Board and Choosing your first Resource

To contribute effectively to the development of the Meraki Terraform provider, it's essential to have visibility into the project's progress and identify areas where your contributions can make a significant impact.

The project board provides an overview of the tasks, issues, and feature requests related to the Meraki Terraform provider. It helps in tracking the project's development, identifying ongoing work, and finding areas where you can contribute. Here are a few steps to get started:

1. Visit the [project board](https://github.com/orgs/core-infra-svcs/projects/1/views/1) to view the current status of the project and the tasks being worked on.

2. Review the different columns on the project board. These columns represent different stages or categories of work, such as "To Do," "In Progress," "Blocked," and "Completed." Each task or issue is typically represented as a card within the respective column.

3. Explore the cards to find specific tasks or issues that align with your interests or expertise. These cards may include feature requests, bug fixes, documentation improvements, or other areas where contributions are needed.

4. Once you identify a card of interest, click on it to view more details. You'll find additional information about the task, including any associated discussions, requirements, or dependencies.

5. If you decide to work on a particular task, leave a comment on the card expressing your intention to contribute. This helps avoid duplication of efforts and allows the project maintainers to provide guidance or feedback.


## Meraki Learning Resources

If you come from a development background or are largely unfamiliar with the Cisco Meraki solution, the following resources will help you familiarize yourself with the platform:

1. **Cisco Meraki Product Documentation:** Explore the official [Cisco Meraki documentation](https://documentation.meraki.com/) to gain an understanding of the various Meraki products, features, and configurations. This comprehensive documentation provides detailed guides, tutorials, and best practices for deploying and managing Meraki networks.

2. **Cisco Meraki Dashboard API Documentation:** The Meraki Terraform provider leverages the [Cisco Meraki Dashboard API](https://developer.cisco.com/meraki/api-v1/) for interacting with Meraki devices and configurations. 

3. **Cisco Meraki Dashboard OpenAPI Specification:** Refer to the [OpenAPI Spec](https://github.com/meraki/openapi) to learn about the available API endpoints, request/response formats, and authentication methods. Understanding the API will enable you to leverage the provider effectively.

4. **Postman Collection:** The [Meraki API Postman collection](https://documenter.getpostman.com/view/897512/SzYXYfmJ#6a9be9f0-49e9-4644-b93d-5a82f09d9899) is an excellent resource for ad-hoc troubleshooting of API calls. 

5. ** Dashboard-api-go:** The [dashboard-api-go](https://github.com/meraki/dashboard-api-go) HTTP client is our first open source contribution to the Cisco Meraki ecosystem hosted in the official meraki repo. It is used in our resources and data source to perform the underlying API calls abstracted by the provider.
