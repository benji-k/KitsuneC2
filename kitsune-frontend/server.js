const { createServer } = require('https');
const { parse } = require('url');
const next = require('next');
const fs = require('fs');
const path = require('path');


certDir = path.resolve("./certificates/") //relative path
if (!checkSSLCerts(certDir)){
    console.log("Could not find certificates in " + certDir)
    process.exit(1)
}

// Load your SSL certificate and key
const httpsOptions = {
  key: fs.readFileSync(path.join(certDir, 'kc2SSL.key')), 
  cert: fs.readFileSync(path.join(certDir, 'kc2SSL.pem')), 
};

const dev = process.env.NODE_ENV !== 'production';
const app = next({ dev });
const handle = app.getRequestHandler();

app.prepare().then(() => {
  createServer(httpsOptions, (req, res) => {
    const parsedUrl = parse(req.url, true);
    handle(req, res, parsedUrl);
  }).listen(8080, "0.0.0.0", err => {
    if (err) console.log(err);
    console.log('> Ready on https://0.0.0.0:8080');
  });
});

function checkSSLCerts(dir){
    fileExists = true
    fs.stat(path.join(certDir, 'kc2SSL.key'), function(err, stat) {
        if (err != null) {
            fileExists  = false
        } 
    });

    fs.stat(path.join(certDir, 'kc2SSL.pem'), function(err, stat) {
        if (err != null) {
            fileExists = false
        } 
    });
    return fileExists
}