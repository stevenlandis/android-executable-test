const http = require("http");
const fs = require("fs");

const server = http.createServer((req, res) => {
  let txt = "";
  req.on("data", (chunk) => {
    txt += chunk;
  });
  req.on("end", () => {
    fs.appendFile(
      "local/log.txt",
      JSON.stringify({ txt, t: Date.now() }) + "\n",
      () => {}
    );
    console.log("recieved:", txt);
    res.writeHead(200);
    res.end();
  });
});
server.listen(4562);
