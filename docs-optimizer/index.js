import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkStringify from 'remark-stringify';
import { glob } from 'glob';
import path from 'path';
import { readFile, writeFile } from 'fs/promises';
import { fileURLToPath } from 'url';
import prettier from 'prettier';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

function processNode(node) {
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
}

const processor = unified()
  .use(remarkParse)
  .use(() => (tree) => {
    const visitNodes = (node) => {
      processNode(node);
      if (node.children) {
        node.children.forEach(visitNodes);
      }
    };

    visitNodes(tree);
  })
  .use(remarkStringify);

async function processFile(filePath) {
  const content = await readFile(filePath, 'utf8');
  const result = await processor.process(content);

  const formattedResult = await prettier.format(String(result), {
    parser: 'markdown',
    proseWrap: 'always',
  });

  await writeFile(filePath, formattedResult);
  console.log(`Processed and formatted: ${filePath}`);
}

async function main() {
  const docsDir = path.resolve(__dirname, '..', 'docs');
  const files = await glob('**/*.md', { cwd: docsDir });

  await Promise.all(files.map(file => processFile(path.join(docsDir, file))));
}

main().catch(console.error);
