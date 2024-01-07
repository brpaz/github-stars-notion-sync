[![github-workflow-shield]][github-workflow-url]
[![Contributors][contributors-shield]][contributors-url]
[![MIT License][license-shield]][license-url]
[![Go Version][gomod]][gomod-url]

<br />

<div>
<h3 align="center">GitHub Stars Notion sync</h3>

  <p align="center">
    A command line tool to sync your GitHub starred repositories with a <a href="https://notion.com">Notion</a> database.
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>


## ℹ️ About The Project

This project allows to syncronize your starred repositories from GitHub to a Notion database, allowing easily filtering your favorite GitHub using the advanced features of a Notion database.

### Built With

* Golang and [Cobra](https://cobra.dev)

## 🚀 Getting Started

### Setup your notion database.

For this integration to work properly, you must create a Notion database with at least the following collumns:

| Field            | Type            | Description                                                                                     |
|------------------|-----------------|-------------------------------------------------------------------------------------------------|
| Name             | text            | This field will store the name of the GitHub repository.                                        |
| Description      | text            | This field will store the description of the GitHub repository.                                 |
| Language         | select          | This field will store the main language of the GitHub repository.                                |
| Topics           | multi-select    | This field will store the topics of the GitHub repository.                                       |
| Repository URL   | url             | This field will store the URL of the GitHub repository.                                          |
| Repository ID    | number | Internal field that will store the repository ID. You can hide it from the table but must not change. It is used to keep track of the already synced repository. |
| Created Time     | time            | Will keep track of the date this repository was synced. You can also hide it from the table but must not be removed. |

> [!TIP]
> You can have any other collumns in your database. They won´t be touched by this command.

You can use [this template](https://brpaz-dev.notion.site/75dd9254235f4577a9d4d259df6a2b64?v=a2ecaa84752c4699b02a982fbb8872a6&pvs=4) to get started.

### Configure notion integration

Next you need to create a Notion API Token and give it access to your database.

1. **Create a Notion integration:** Go to https://www.notion.so/my-integrations and create a new integration. Give it a name like "GitHub stars Syncer" and associate in to the workspace where your database is. Make sure to save the generated token in a Safe place.
2. **Enable integration for your database**: Open your Notion database page, and on the `...` menu at top right, click on "Connections" -> "Add connection" and select the integration you created on 1. This will ensure the integration have access to your database.
3. **Find your database id** - Open your database page in Notion. You should see in your browser an url similar to `https://www.notion.com/fer6ff3d5fcs3dff1d2134349192cc?v=4rf43545..`. Grab the first id. This is your database id and you will need it when running the command.

### Create a GitHub access token.

You will also need a GitHub access token, in order to retrieve your starred repos from GitHub

1. Login in GitHub and open `https://github.com/settings/tokens`
2. Create a new "General token" with `read:user` permission.


### Installation

1. Download the latest release from [GitHub](https://github.com/brpaz/github-stars-notion-sync/releases) for your operating system.

## ▶️ Usage

To sync your GitHub starred repos with your Notion database, run the following command:

```shell
github-stars-notion-sync sync --github-token=<token> --notation-token=<notion-token> --notion-database-id=<database-id>
```

Instead of using flags to set the command options, you can also use envrionment variables.

Ex:

```shell
GITHUB_TOKEN=<github-token> NOTION_TOKEN=<notion-token> NOTION_DATABASE_ID=<database-id> github-stars-notion-sync sync
```

### Run with docker

If you prefer, you can also use Docker.

```shell
docker run --rm \
    -e GITHUB_TOKEN=<token> \
    -e NOTION_TOKEN=<notion-token> \
    -e NOTION_DATABASE_ID=<database-id> \
    ghcr.io/brpaz/github-stars-notion-sync:latest sync
```


## 🤝 Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🫶 Support

If you find this project helpful and would like to support its development, there are a few ways you can contribute:

[![Sponsor me on GitHub](https://img.shields.io/badge/Sponsor-%E2%9D%A4-%23db61a2.svg?&logo=github&logoColor=red&&style=for-the-badge&labelColor=white)](https://github.com/sponsors/brpaz)

<a href="https://www.buymeacoffee.com/Z1Bu6asGV" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>

## 📃 License

Distributed under the MIT License. See [LICENSE](LICENSE.md) file for details.

## 📩 Contact

- Bruno Paz - [https://brunopaz.dev](https://brunopaz.dev) - oss@brunopaz.dev

## 🏅 Acknowledgments

* [Anatoly Nosov](https://github.com/jomei) for creating the [Notion API golang client](https://github.com/jomei/notionapi), which helped a lot building this integration.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/brpaz/github-stars-notion-sync.svg?style=for-the-badge
[contributors-url]: https://github.com/brpaz/github-stars-notion-sync/graphs/contributors
[license-shield]: https://img.shields.io/github/license/brpaz/github-stars-notion-sync.svg?style=for-the-badge
[license-url]: https://github.com/brpaz/github-stars-notion-sync/blob/main/LICENSE.md
[github-workflow-shield]: https://img.shields.io/github/actions/workflow/status/brpaz/github-stars-notion-sync/CI.yml?style=for-the-badge
[github-workflow-url]: https://github.com/brpaz/github-stars-notion-sync/actions
[gomod]: https://github.com/github-stars-notion-sync
[gomod-url]: https://img.shields.io/github/go-mod/go-version/brpaz/github-stars-notion-sync



