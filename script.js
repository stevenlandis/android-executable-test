async function main(args) {
  const cmd = args[0];
  const otherArgs = args.slice(1);
  ({
    build: runBuild,
    upload: runUpload,
    run: runRun,
    shell: runShell,
    test: runTest,
  }[cmd](otherArgs));
}

const EXECUTABLES = ["serve", "kill_active_server", "top", "start_server"];

async function runBuild() {
  await cmd("mkdir", "-p", "dist").run();
  process.env.GOOS = "linux";
  process.env.GOARCH = "arm";
  process.env.GOARM = "7";
  function buildExe(name) {
    return cmd("go", "build", "-o", `dist/${name}`, `${name}/main.go`).run();
  }
  await Promise.all([
    ...EXECUTABLES.map((name) => buildExe(name)),
    shell("yarn", "build").inDir("site-js").runShell(),
    cmd("cp", "-r", "public", "dist").run(),
  ]);
}

async function runUpload() {
  await runBuild();
  await cmd("adb", "push", "--sync", "dist", `/data/local/tmp`).run();
  await cmd(
    "adb",
    "shell",
    [
      "cd /data/local/tmp/dist",
      ...EXECUTABLES.map((name) => `chmod 777 ${name}`),
    ].join(" && ")
  ).run();
}

async function runRun() {
  await runUpload();
  const env = JSON.parse(await fs.promises.readFile("env.json", "utf8"));
  const input =
    `
cd /data/local/tmp/dist
./kill_active_server
./start_server ${JSON.stringify(JSON.stringify(env))}
exit
  `.trim() + "\n";
  console.log(input);
  await cmd("echo", input).pipe("adb", "shell").get();
}

async function runShell() {
  await cmd(
    "adb",
    "shell",
    ["cd /data/local/tmp/dist", "PATH=/data/local/tmp/dist:$PATH", "sh"].join(
      " && "
    )
  ).run();
}

async function runTest() {
  // await cmd("ls").pipe("grep", "d").run();
  // const t = await cmd("ls").pipe("grep", "d").get();
  // console.log(t);
  // await cmd("bash").run();
  // await cmd("ls").toFile("tmp.txt");
  // await cat(".gitignore").pipe("wc", "-l").run();
  // await cat(".gitignore").pipe("wc", "-l").run();
}

function runPrintIp(args) {
  let [s] = args;
  let [ip, port] = s.split(":");
  const parts = [];
  for (let i = 0; i + 1 < ip.length; i += 2) {
    const t = ip.substr(i, 2);
    parts.push(parseInt(t, 16));
  }
  parts.reverse();
  console.log(parts.join(".") + ":" + parseInt(port, 16).toString());
}

const child_process = require("child_process");
const fs = require("fs");
function cmd(...parts) {
  const c = new _CMD();
  c.pipe(...parts);
  return c;
}
function shell(...parts) {
  const c = new _CMD();
  c.pipeShell(...parts);
  return c;
}

function cat(filePath) {
  const c = new _CMD().pipe("cat", filePath);
  return c;
}
class _CMD {
  constructor() {
    this.cmds = [];
    this.cwd = undefined;
    this.stdin = process.stdin;
    this.stdout = process.stdout;
  }
  pipe(...parts) {
    this.cmds.push({ parts, shell: false });
    return this;
  }
  pipeShell(...parts) {
    this.cmds.push({ parts, shell: true });
    return this;
  }
  inDir(dir) {
    this.cwd = dir;
    return this;
  }
  async __exec(stdin, stdout, stderr) {
    return new Promise((resolve, reject) => {
      let proc = undefined;
      for (let i = 0; i < this.cmds.length; i++) {
        const parts = this.cmds[i].parts;
        const _stdin = i === 0 ? stdin : proc.stdout;
        const _stdout = i === this.cmds.length - 1 ? stdout : "pipe";
        const _stderr = stderr;
        proc = child_process.spawn(parts[0], parts.slice(1), {
          stdio: [_stdin, _stdout, _stderr],
          shell: this.cmds[i].shell,
          cwd: this.cwd,
        });
      }
      let stdoutTxt = undefined;
      if (stdout === "pipe") {
        stdoutTxt = "";
        proc.stdout.on("data", (data) => {
          stdoutTxt += data;
        });
      }
      proc.on("close", (code) => {
        if (code === 0) {
          resolve({ stdout: stdoutTxt });
        } else {
          reject(code);
        }
      });
    });
  }
  async run() {
    await this.__exec(this.stdin, process.stdout, process.stderr);
  }
  async runShell() {
    this.shell = true;
    await this.run();
  }
  async get() {
    const { stdout } = await this.__exec(this.stdin, "pipe", process.stderr);
    return stdout;
  }
  async toFile(file) {
    await this.__exec(this.stdin, fs.createWriteStream(file), process.stderr);
  }
}

function randomHexString(N) {
  let s = "";
  for (let i = 0; i < N; i++) {
    s += Math.floor(16 * Math.random()).toString(16);
  }
  return s;
}

main(process.argv.slice(2));
