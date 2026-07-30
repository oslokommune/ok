import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkStringify from "remark-stringify";
import { glob } from "glob";
import path from "path";
import { readFile, writeFile } from "fs/promises";
import { fileURLToPath } from "url";
import prettier from "prettier";
import { visit } from "unist-util-visit";

const currentScriptPath = fileURLToPath(import.meta.url);
const currentScriptDir = path.dirname(currentScriptPath);

const removeHeadingAndSubsectionsPlugin = ({ targetHeading = "" } = {}) => {
  return (tree) => {
    visit(tree, "heading", (currentHeading, currentIndex, parentNode) => {
      if (currentHeading.children[0]?.value !== targetHeading) return;

      const nextSiblingHeadingIndex = findNextSiblingHeadingIndex(
        parentNode.children,
        currentIndex,
        currentHeading.depth
      );
      const sectionLength = calculateSectionLength(
        nextSiblingHeadingIndex,
        parentNode.children.length,
        currentIndex
      );

      parentNode.children.splice(currentIndex, sectionLength);
      return [visit.SKIP, currentIndex];
    });
  };
};

const findNextSiblingHeadingIndex = (siblings, startIndex, currentDepth) => {
  return siblings.findIndex(
    (sibling, index) =>
      index > startIndex &&
      sibling.type === "heading" &&
      sibling.depth <= currentDepth
  );
};

const calculateSectionLength = (endIndex, totalSiblings, startIndex) => {
  return endIndex === -1 ? totalSiblings - startIndex : endIndex - startIndex;
};

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
      node.lang = node.lang ?? "sh";
      break;
  }
};

// Literals that cobra emits as bare or single-quoted words mid-sentence. Outside
// code font they read as misspellings, so each match becomes an inlineCode node.
const inlineLiteralPattern = new RegExp(
  [
    // Package names that cobra's help text wraps in single quotes.
    String.raw`'(?<quoted>bash-completion)'`,
    // The CLI's own name used as a word in a sentence. The lookarounds keep it
    // from matching inside a longer word, a filename or a possessive.
    String.raw`(?<![\w'./-])(?<command>ok)(?![\w'./-])`,
  ].join("|"),
  "g"
);

const splitTextNodeOnLiterals = (node, parentNode, nodeIndex) => {
  const literalPattern = new RegExp(inlineLiteralPattern);
  if (!literalPattern.test(node.value)) return;

  literalPattern.lastIndex = 0;
  const replacementNodes = [];
  let lastMatchEnd = 0;

  for (const match of node.value.matchAll(literalPattern)) {
    if (match.index > lastMatchEnd) {
      replacementNodes.push({
        type: "text",
        value: node.value.slice(lastMatchEnd, match.index),
      });
    }
    const { quoted, command } = match.groups;
    replacementNodes.push({ type: "inlineCode", value: quoted ?? command });
    lastMatchEnd = match.index + match[0].length;
  }

  if (lastMatchEnd < node.value.length) {
    replacementNodes.push({
      type: "text",
      value: node.value.slice(lastMatchEnd),
    });
  }

  parentNode.children.splice(nodeIndex, 1, ...replacementNodes);
  return nodeIndex + replacementNodes.length;
};

const codeFontLiteralsPlugin = () => (tree) => {
  visit(tree, "text", (node, nodeIndex, parentNode) => {
    if (!parentNode) return;
    const nextIndex = splitTextNodeOnLiterals(node, parentNode, nodeIndex);
    if (nextIndex !== undefined) return nextIndex;
  });
};

// A bare command name in a sentence reads as a misspelling, because that is what
// `fmt` and `aws` are outside code font. Cobra emits the command path bare in
// both the page title and the "See also" link labels.
const isCommandPath = (value) => /^ok(\s|$)/.test(value);

const wrapChildInCodeFont = (node) => {
  const [firstChild] = node.children;
  if (firstChild?.type !== "text" || !isCommandPath(firstChild.value)) return;
  node.children = [{ type: "inlineCode", value: firstChild.value }];
};

const codeFontCommandNamesPlugin = () => (tree) => {
  visit(tree, "heading", (node) => {
    // Only the page title holds a command path; the rest are section headings.
    if (node.depth !== 1) return;
    wrapChildInCodeFont(node);
  });

  visit(tree, "link", (node) => {
    // "See also" links point at sibling pages; leave external links alone.
    if (!node.url.endsWith(".md")) return;
    wrapChildInCodeFont(node);
  });
};

const markdownProcessor = unified()
  .use(remarkParse)
  .use(() => (tree) => {
    const visitAndProcessNodes = (node) => {
      processMarkdownNode(node);
      if (node.children) {
        node.children = node.children.map(visitAndProcessNodes);
      }
      return node;
    };

    return visitAndProcessNodes(tree);
  })
  .use(removeHeadingAndSubsectionsPlugin, {
    targetHeading: "Options inherited from parent commands",
  })
  .use(codeFontCommandNamesPlugin)
  .use(codeFontLiteralsPlugin)
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
  const docsDirectoryPath = path.join(currentScriptDir, "..", "docs");
  const markdownFilePaths = await glob("**/*.md", { cwd: docsDirectoryPath });

  await Promise.all(
    markdownFilePaths.map((relativeFilePath) =>
      processMarkdownFile(path.join(docsDirectoryPath, relativeFilePath))
    )
  );
};

processAllMarkdownFiles().catch(console.error);
