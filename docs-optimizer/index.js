import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkStringify from 'remark-stringify';
import { glob } from 'glob';
import path from 'path';
import { readFile, writeFile } from 'fs/promises';
import { fileURLToPath } from 'url';
import prettier from 'prettier';

const currentFilePath = fileURLToPath(import.meta.url);
const currentDirPath = path.dirname(currentFilePath);

const processMarkdownNode = (node) => {
  switch (node.type) {
    case 'heading':
      // Ensure headings start at level 1 and maintain hierarchy
      if (node.depth > 1) {
        node.depth = Math.max(1, node.depth - 1);
      }

      if (node.children?.[0]?.type === 'text') {
        switch (node.depth) {
          case 2:
            if (node.children[0].value === 'SEE ALSO') {
              node.children[0].value = 'See also';
            }
            break;
          case 3:
            switch (node.children[0].value) {
              case 'Linux:':
                node.children[0].value = 'Linux';
                break;
              case 'macOS:':
                node.children[0].value = 'macOS';
                break;
            }
            break;
        }
      }
      break;
    case 'code':
      if (!node.lang) {
        node.lang = 'sh';
      }
      break;
  }
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
  .use(remarkStringify);

const processMarkdownFile = async (filePath) => {
  const markdownContent = await readFile(filePath, 'utf8');
  const processedMarkdown = await markdownProcessor.process(markdownContent);

  const formattedMarkdown = await prettier.format(String(processedMarkdown), {
    parser: 'markdown',
    proseWrap: 'always',
  });

  await writeFile(filePath, formattedMarkdown);
  console.log(`Processed and formatted: ${filePath}`);
};

const processAllMarkdownFiles = async () => {
  const docsDirectoryPath = path.resolve(currentDirPath, '..', 'docs');
  const markdownFilePaths = await glob('**/*.md', { cwd: docsDirectoryPath });

  await Promise.all(markdownFilePaths.map(filePath => processMarkdownFile(path.join(docsDirectoryPath, filePath))));
};

processAllMarkdownFiles().catch(console.error);
