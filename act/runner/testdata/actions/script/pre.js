import {execFileSync} from 'node:child_process';

const script = process.env.INPUT_PRE;
if (script) {
  execFileSync('bash', ['-eo', 'pipefail', '-c', script], {stdio: 'inherit'});
}
