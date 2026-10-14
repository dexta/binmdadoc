const express = require('express');
const bodyParser = require('body-parser')

const asciidoctor = require('asciidoctor');
const markdownit = require('markdown-it');

const host = process.env.HOST || "0.0.0.0";
const port = parseInt(process.env.PORT, 10) || 8080;

const app = express();

app.use(bodyParser.json({ limit: '10mb' }));


app.get('/api/test/asciidoctor', async (req, res) => {
  try {
    const Asciidoctor = asciidoctor();
    const content = 'http://asciidoctor.org[*Asciidoctor*] ' +
      'running on https://opalrb.com[_Opal_] ' +
      'brings AsciiDoc to Node.js!';
    const html = Asciidoctor.convert(content);
    res.status(200).send(html);
  } catch (error) {
    console.error('Error in /api/test/asciidoctor:', error);
    res.status(500).json({ error: 'Internal Server Error' });
  }
});

app.get('/api/test/markdown', async (req, res) => {
  try {
    const md = markdownit();
    const html = md.render('# markdown-it rulezz!\nget doc for [markdown-it](https://markdown-it.github.io/markdown-it/)');
    res.status(200).send(html);
  } catch (error) {
    console.error('Error in /api/test/markdown:', error);
    res.status(500).json({ error: 'Internal Server Error' });
  }
});


app.post('/api/echo/asciidoc', async (req, res) => {
  try {
    const Asciidoctor = asciidoctor();
    const content = req.body.text;
    const html = await Asciidoctor.convert(content);
    res.status(200).send(html);
  } catch (error) {
    console.error('Error in /api/echo/asciidoc:', error);
    res.status(500).json({ error: 'Internal Server Error' });
  }
});


app.post('/api/echo/markdown', async (req, res) => {
  try {
    const md = markdownit();
    const content = req.body.text;
    const html = md.render(content);
    res.status(200).send(html);
  } catch (error) {
    console.error('Error in /api/echo/markdown:', error);
    res.status(500).json({ error: 'Internal Server Error' });
  }
});

app.listen(port, host);
console.log(`Running on http://${host}:${port}`);