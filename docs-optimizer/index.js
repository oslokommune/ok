import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkStringify from "remark-stringify";
import { glob } from "glob";
import path from "path";
import { readFile, writeFile } from "fs/promises";
import { fileURLToPath } from "url";
import prettier from "prettier";

const currentScriptPath = fileURLToPath(import.meta.url);
const currentScriptDir = path.dirname(currentScriptPath);

const processMarkdownNode = (node) => {
  switch (node.type) {
    case "heading": {
      // Ensure headings start at level 1 and maintain hierarchy
      if (node.depth > 1) {
        node.depth = Math.max(1, node.depth - 1);
      }

      if (node.children?.[0]?.type === "text") {
        const { value } = node.children[0];
        switch (node.depth) {
          case 2:
            if (value === "SEE ALSO") {
              node.children[0].value = "See also";
            }
            break;
          case 3: {
            const updates = {
              "Linux:": "Linux",
              "macOS:": "macOS",
            };
            if (value in updates) {
              node.children[0].value = updates[value];
            }
            break;
          }
        }
      }
      break;
    }
    case "code":
      if (!node.lang) {
        node.lang = "sh";
      }
      break;
  }
};

const markdownProcessor = unified()
  .use(remarkParse)
  .use(() => tree => {
    const visitAndProcessNodes = node => {
      processMarkdownNode(node);
      if (node.children) {
        node.children = node.children.map(visitAndProcessNodes);
      }
      return node;
    };

    return visitAndProcessNodes(tree);
  })
  .use(remarkStringify);

  const processMarkdownFile = async (markdownFilePath) => {
    const markdownContent = await readFile(markdownFilePath, "utf8");
    const processedMarkdown = await markdownProcessor.process(markdownContent);

    const formattedMarkdown = await prettier.format(String(processedMarkdown), {
      parser: "markdown",
      proseWrap: "always",
    });

    await writeFile(markdownFilePath, formattedMarkdown);
    console.log(`Processed and formatted: ${markdownFilePath}`);
  };

  const processAllMarkdownFiles = async () => {
    const docsDirectoryPath = path.resolve(currentScriptDir, "..", "docs");
    const markdownFilePaths = await glob("**/*.md", { cwd: docsDirectoryPath });

    await Promise.all(
      markdownFilePaths.map((relativeFilePath) =>
        processMarkdownFile(path.join(docsDirectoryPath, relativeFilePath))
      )
    );
  };

  processAllMarkdownFiles().catch(console.error);
