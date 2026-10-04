const express = require('express');
const bodyParser = require('body-parser')

const asciidoctor = require('asciidoctor');
const markdownit = require('markdown-it');

const host = "0.0.0.0";
const port = 8080;

const app = express();

// app.use(bodyParser.text({ type: '*' }));
app.use(bodyParser.json());


app.get('/api/test/asciidoctor', async (req, res) => {
  const Asciidoctor = asciidoctor();
  const content = 'http://asciidoctor.org[*Asciidoctor*] ' +
    'running on https://opalrb.com[_Opal_] ' +
    'brings AsciiDoc to Node.js!';
  const html = Asciidoctor.convert(content);
  res.status(200).send(html);
});

app.get('/api/test/markdown', async (req, res) => {
  const md = markdownit();
  const html = md.render('# markdown-it rulezz!\nget doc for [markdown-it](https://markdown-it.github.io/markdown-it/)');

  res.status(200).send(html);
});


app.post('/api/echo/asciidoc', async (req, res) => {
  const Asciidoctor = asciidoctor();
  const content = req.body.text;

  const html = await Asciidoctor.convert(content);
  res.status(200).send(html);
});


app.post('/api/echo/markdown', async (req, res) => {
  const md = markdownit();
  const content = req.body.text;

  const html = md.render(content);
  res.status(200).send(html);
});

app.listen(port, host);
console.log(`Running on http://${host}:${port}`);