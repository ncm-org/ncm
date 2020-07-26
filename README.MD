<!--<link rel="stylesheet" type="text/css" href="assets/asciinema-player.css" />-->
<!--<script src="assets/asciinema-player.js"></script>-->

#### NCM是什么？
- **N**ice **C**ommit **M**essage
- 用于校验 `git commit message` 是否合符[`约定式提交`](https://www.conventionalcommits.org/zh-hans/v1.0.0-beta.4/)规范。
- 生成符合[`约定式提交`](https://www.conventionalcommits.org/zh-hans/v1.0.0-beta.4/)规范的`git commit message`。

#### 什么是约定式提交

> 一种用于给提交信息增加人机可读含义的规范！

约定式提交规范是一种基于提交消息的轻量级约定。 它提供了一组用于创建清晰的提交历史的简单规则； 这使得编写基于规范的自动化工具变得更容易。 这个约定与 [`SemVer`](https://semver.org/lang/zh-CN/) 相吻合， 在提交信息中描述新特性、bug 修复和破坏性变更。

提交说明的结构如下所示：
```
<类型>[可选的作用域]: <描述>

[可选的正文]

[可选的脚注]
```
约定式提交规范受到了[`angular.js`](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#-git-commit-guidelines) 提交准则的启发，并在很大程度上以其为依据。关于类型、作用域、描述、正文、脚注的详细规范信息可参考[`这里`](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#-git-commit-guidelines)。

#### 为什么要使用[`约定式提交`](https://www.conventionalcommits.org/zh-hans/v1.0.0-beta.4/)
> 在多人协作项目中，代码风格统一、代码提交信息的说明等重要性不言而喻。规范的 `LOG` 不仅有助于他人 `REVIEW` ，还可以有效的输出 `CHANGELOG` ，甚至对于项目的研发质量都有很大的提升。
- 基于提交的类型，自动决定语义化的版本变更。
- 向同事、公众与其他利益关系者传达变化的性质。
- 触发构建和部署流程。
- 让人们探索一个更加结构化的提交历史，以便降低对你的项目做出贡献的难度。

#### 现有的开源工具
- [**commitizen/cz-cli**](https://github.com/commitizen/cz-cli)：**`star: 9.4K`**，提供 `git cz` 命令替代 `git commit` 命令, 帮助开发者生成符合规范的 `commit message`。
- [**conventional-changelog/commitlint**](https://github.com/conventional-changelog/commitlint)：**`star: 6.1K`**，检查输入的 `commit message` 是否满足固定格式的工具。
- [**typicode/husky**](https://github.com/typicode/husky)：**`star: 19.7K`**，为了可以在每次 `commit` 时执行 `commitlint` 检查输入的 message，还需要用到一个工具 [`husky`]((https://github.com/typicode/husky))。[`husky`]((https://github.com/typicode/husky)) 是一个增强的 `git hook` 工具。可以在 `git hook` 的各个阶段执行 `package.json` 中配置好的 `npm script`。例如：<br><br>
    ```
    {
      "scripts": {
        "commitmsg": "commitlint -E GIT_PARAMS"
      }
    }
    ```

#### 现有开源工具的缺点
- 安装、配置繁琐，要实现一整套规范流程，首先需要安装 [`nodejs`](https://nodejs.org/)、[`npm`](https://www.npmjs.com/) 环境。然后依次安装上述的3种工具，再进行配置，以实现工具在提交过程中的衔接。
- `git cz` 命令在运行前没有对暂存区进行审查。导致输入完提交信息后却被告知 `changes not staged for commit`。
- 配置可选项过于繁多，容易造成团队内不同项目之间规范存在差异。

#### NCM特点
- 跨平台，Go语言天然优势。支持Window、MacOS、Linux。
- 安装方便，只需要将可执行文件配置在环境变量即可。
- 在线升级，方便及时升级到最新版本。
- 功能强大，一个程序即可实现提交、校验的功能。
- 兼容性好，兼容 `git cz`，兼容 [`jetbrains`](https://www.jetbrains.com/) 的 [`Git Commit Template`](https://plugins.jetbrains.com/plugin/9861-git-commit-template) 插件。

#### NCM安装
- `MacOS` `Linux` 系统运行以下命令，等待执行完成后运行 `ncm info` 验证是否安装成功。

    > 安装命令的运行依赖与 `curl`  以及 `unzip`  , 请先确认系统是否安装了依赖工具。

    ```shell
    source <(curl -LkSs sh.ncm.ink)
    ```

    ![ncm install](assets/ncm_install.gif)
    <!--<asciinema-player width=100% poster="npt:0:20" src="assets/ncm_install.json" />-->

- `Windows` 系统参考如下步骤安装

   1. 进入项目的发布页面，选择合适的zip文件进行下载，并解压到合适的目录。

   2. 将解压目录配置到环境变量中，以便于可以不带路径运行 `ncm` 。

   3. 运行 `ncm info` 验证是否安装成功。

#### NCM使用
1. ##### NCM支持的所有命令

   ![ncm help](assets/ncm_help.gif)
   <!--<asciinema-player width=100% poster="npt:0:4" src="assets/ncm_help.json" />-->

2. ##### 使用NCM校验commit message
   在git项目中运行 `ncm init` 命令。此命令会在项目中的 `.git/hooks` 文件夹中生成内容为`ncm check`的`commit-msg` hook文件。有关`git commit-msg hook`的详细信息可点击[这里](https://git-scm.com/docs/githooks#_commit_msg)查阅。

   > 注意：如果项目中已经通过其他方式添加了 `commit-msg` 文件。则在运行命令后会提示是否覆盖文件，输入 `y/N` 进行选择。回车默认选择`N`，则在 `commit message` 时不会触发 `ncm` 的校验。

   ![ncm init](assets/ncm_init.gif)
   <!--<asciinema-player width=100% poster="npt:0:26" src="assets/ncm_init.json" />-->

   之后，每次进行 `git commit` 时都会触发 `commit-msg` hook，执行 `ncm check` 命令。如果提交信息不符合[`约定式提交`](https://www.conventionalcommits.org/zh-hans/v1.0.0-beta.4/)规范，则会中断提交，并输出帮助提示。

   ![ncm check](assets/ncm_check.gif)
   <!--<asciinema-player width=100% poster="npt:0:19" src="assets/ncm_check.json" />-->

   > 如果之后不想使用此功能，可运行 `ncm uninit` 命令对 `commit-msg` hook进行移除。

3. ##### 生成commit message信息并提交
   在需要 `commit` 时，请先使用 `git add` 命令，将需要提交的文件添加到暂存区，然后运行 `ncm` 命令 。此时会进入交互式模式。根据实际情况以及输出的提示，依次选择或填写相关信息。

   ![ncm](assets/ncm.gif)
   <!--<asciinema-player width=100% poster="npt:0:34" src="assets/ncm.json" />-->

   如果只想提交指定的文件，可以运行 `ncm add file ...`

   ![ncm add file ...](assets/ncm_add.gif)
   <!--<asciinema-player width=100% poster="npt:0:21" src="assets/ncm_add.json" />-->

   如果想提交所有文件，可以运行 `ncm add .`

   ![ncm add .](assets/ncm_add_all.gif)
   <!--<asciinema-player width=100% poster="npt:0:28" src="assets/ncm_add_all.json" />-->

4. ##### 版本升级
   运行 `ncm upgrade`

   ![ncm upgrade](assets/ncm_upgrade.gif)
   <!--<asciinema-player width=100% poster="npt:0:7" src="assets/ncm_upgrade.json" />-->

5. ##### 意见反馈

   POPO群号：1574363

   运行 `ncm info` 命令，复制输出的信息，粘贴到群里，并描述具体的问题。

   ![ncm info](assets/ncm_info.gif)
   <!--<asciinema-player width=100% poster="npt:0:3" src="assets/ncm_info.json" />-->

#### 小提示
- 对于已经在合理合法使用可视化工具生成符合[`约定式提交`](https://www.conventionalcommits.org/zh-hans/v1.0.0-beta.4/)规范的个人、团队或者项目仍可继续使用该工具。为了进一步确保提交信息符合规范，可在项目中运行 `ncm init` 命令注入hook即可。
- 运行 `ncm -h` 可查看支持的所有命令。
- 运行 `ncm -v` 可查看当前版本。
- 运行 `ncm info` 可查看当前版本的详细信息。
- 除非特殊情况，否则尽量不要运行 `ncm uninit`。

#### 感谢
- [**go**](https://go.dev)
- [**git**](https://git-scm.com)
- [**GoLand**](https://www.jetbrains.com/go)
- [**SemVer**](https://semver.org)
- [**asciinema**](https://asciinema.org/)
- [**spf13/cobra**](https://github.com/spf13/cobra)
- [**gookit/color**](https://github.com/gookit/color)
- [**typicode/husky**](https://github.com/typicode/husky)
- [**commitizen/cz-cli**](https://github.com/commitizen/cz-cli)
- [**angular/angular.js**](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#-git-commit-guidelines)
- [**schollz/progressbar**](https://github.com/schollz/progressbar)
- [**conventionalcommits**](https://www.conventionalcommits.org)
- [**goreleaser/goreleaser**](https://github.com/goreleaser/goreleaser)
- [**conventional-changelog/commitlint**](https://github.com/conventional-changelog/commitlint)