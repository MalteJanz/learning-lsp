# Learning LSP

Just a learning project to get familiar with the [LSP Specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.18/specification/)
by building a very primitive LSP server from scratch in Golang.

## Goal

To build towards a real use case this project tries to implement LSP functionality for administration snippets in [Shopware](https://github.com/shopware/shopware) projects.
For now the following is implemented (with some ideas to explore further):
- [x] when you open `snippet/en.json` / `snippet/de.json` file in your IDE these snippets get indexed (stored in a lookup table)
- [x] when you have a `.twig` / `.js` / `.ts` file open and hover open a snippet key string, and if that is found in the index, it shows all known definitions and their values
- [ ] index all snippet files in workspace root directory, not only opened / changed ones
- [ ] autocomplete indexed snippet keys
- [ ] instead of just json unmarshal build a lexer + parser for snippet definitions, properly extracting line + utf-16 character information for each snippet key
- [ ] goto definitions based on the exact file + position information of each snippet key
- [ ] goto references
- [ ] diagnostics on missing keys for other languages (e.g. defined in 'en' but not in 'de')
- [ ] diagnostics on key overrides for the same language (e.g. leaf defined in one file, overridden in another)

## Used learning ressources

- Youtube: [LSP: Building a Language Server From Scratch](https://youtu.be/Xo5VXTRoL6Q) (in Typescript)
  - GitHub: [Minimum Viable VS Code Language Server Extension](https://github.com/semanticart/minimum-viable-vscode-language-server-extension) where this Repo got its scaffolding code from, especially the VSCode client / extension and some further details below in this README on how to build / debug with VSCode.
- Youtube: [Learn By Building: Language Server Protocol](https://youtu.be/YsdlcQoHqPY) (in Golang)
  - GitHub: [educationalsp](https://github.com/tjdevries/educationalsp)
- Blog: [Why building a Rust LSP is hard](https://rust-glancer.github.io/blog/why-lsp-is-hard/)
  - Interesting blog post which contains many useful links
- Also worth exploring all the great blog posts from [matklad](https://matklad.github.io/) / [rust-analyzer](https://rust-analyzer.github.io/blog) about parsing and LSP topics

## Getting Started

1. Run `npm install` from the repo root.

To make it easy to get started, this language server will run on _every_ file type by default. To target specific languages, change

`package.json`'s `activationEvents` to something like

```
"activationEvents": [
  "onLanguage:plaintext"
],
```

And change the `documentSelector` in `client/src/extension.ts` to replace the `*` (e.g.)

```
documentSelector: [{ scheme: "file", language: "plaintext" }],
```

2. Inside the `./server` directory, run `go build .`

## Developing your extension

From the root directory of this project, run `code .` Then in VS Code

1. Build the extension (both client and server) with `⌘+shift+B` (or `ctrl+shift+B` on windows)
2. Open the Run and Debug view and press "Launch Client" (or press `F5`). This will open a `[Extension Development Host]` VS Code window.
3. Opening or editing a file in that window should show an information message in VS Code like you see below.

   ![example information message](https://semanticart.com/misc-images/minimum-viable-vscode-language-server-extension-info-message.png)

4. Inspect the LSP server log under `/tmp/learning-lsp.log` by opening that file

[Debugging instructions can be found here][debug]

## Distributing your extension

Read the full [Publishing Extensions doc][publish] for the full details.

Note that you can package and distribute a standalone `.vsix` file without publishing it to the marketplace by following [these instructions][vsix].

## Anatomy

```
.
├── .vscode
│   ├── launch.json         // Tells VS Code how to launch our extension
│   └── tasks.json          // Tells VS Code how to build our extension
├── LICENSE
├── README.md
├── client
│   ├── package-lock.json   // Client dependencies lock file
│   ├── package.json        // Client manifest
│   ├── src
│   │   └── extension.ts    // Code to tell VS Code how to run our language server
│   └── tsconfig.json       // TypeScript config for the client
├── package-lock.json       // Top-level Dependencies lock file
├── package.json            // Top-level manifest
├── server                  // Golang LSP written from scratch, the heart of this project
└── tsconfig.json           // Top-level TypeScript config
```

[debug]: https://code.visualstudio.com/api/language-extensions/language-server-extension-guide#debugging-both-client-and-server
[sample]: https://github.com/microsoft/vscode-extension-samples/tree/main/lsp-sample
[publish]: https://code.visualstudio.com/api/working-with-extensions/publishing-extension
[vsix]: https://code.visualstudio.com/api/working-with-extensions/publishing-extension#packaging-extensions
