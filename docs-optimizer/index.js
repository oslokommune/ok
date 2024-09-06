import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkStringify from 'remark-stringify';
import { glob } from 'glob';
import path from 'path';
import fs from 'fs/promises';
import { fileURLToPath } from 'url';
import prettier from 'prettier';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const processor = unified()
  .use(remarkParse)
  .use(() => (tree) => {
    const visitor = (node) => {
      // Ensure that headings are at the correct level
      if (node.type === 'heading' && node.depth > 1) {
        node.depth -= 1;
      }

      // Change "Linux:" to "Linux" and "macOS:" to "macOS" in level 3 headings
      if (node.type === 'heading' && node.depth === 3) {
        const textNode = node.children[0];
        if (textNode && textNode.type === 'text') {
          if (textNode.value === 'Linux:') {
            textNode.value = 'Linux';
          } else if (textNode.value === 'macOS:') {
            textNode.value = 'macOS';
          }
        }
      }

      // Add 'sh' language to code blocks without a specified language
      if (node.type === 'code' && !node.lang) {
        node.lang = 'sh';
      }
    };

    const visitNodes = (node) => {
      visitor(node);
      if (node.children) {
        node.children.forEach(visitNodes);
      }
    };

    visitNodes(tree);
  })
  .use(remarkStringify);

async function processFile(filePath) {
  const content = await fs.readFile(filePath, 'utf8');
  const result = await processor.process(content);

  // Format the result with Prettier
  const formattedResult = await prettier.format(String(result), {
    parser: 'markdown',
    proseWrap: 'always',
  });

  await fs.writeFile(filePath, formattedResult);
  console.log(`Processed and formatted: ${filePath}`);
}

async function main() {
  const docsDir = path.resolve(__dirname, '..', 'docs');
  const files = await glob('**/*.md', { cwd: docsDir });

  for (const file of files) {
    await processFile(path.join(docsDir, file));
  }
}

main().catch(console.error);
