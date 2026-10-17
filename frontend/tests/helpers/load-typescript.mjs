import { readFileSync } from 'node:fs';
import { createRequire } from 'node:module';
import { runInNewContext } from 'node:vm';
import ts from 'typescript';

const require = createRequire(import.meta.url);

export const loadTypeScript = (url, globals, modules = {}) => {
  const source = readFileSync(url, 'utf8');
  const { outputText } = ts.transpileModule(source, {
    compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2022 },
  });
  const exports = {};
  runInNewContext(
    outputText,
    {
      exports,
      require: (name) => modules[name] ?? require(name),
      ...globals,
    },
    { filename: url.pathname },
  );
  return exports;
};
