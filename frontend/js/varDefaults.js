const fillEditor = (type, title, text) => {
  let editorName = (type==="md")? "markdown" : "asciidoc";
  
  let titleElement = document.querySelector(`#${editorName}Title`);
  let textElement = document.querySelector(`#${editorName}Editor > textarea`);

  document.querySelector(`#${editorName}Title`).value = title;  
  document.querySelector(`#${editorName}Editor > textarea`).value = text;
  document.querySelector(`#${editorName}Editor > textarea`).dispatchEvent( new Event("input"));
};

const fillMarkdownDefaultText = () => {
  let title = "default markdown title here";
  let text = `# Heading level 1

First line with two spaces after.  
And the next line.

First line with the HTML tag after.<br>
And the next line.

### Heading level 3

sone items here
- First item
- Second item
- Third item
- Fourth item


---

##### Heading level 5
who do links in markdown

[markdown guide](markdownguide.org/basic-syntax/)  
[gitlab markdown](https://docs.gitlab.com/ee/user/markdown.html)

`;

  fillEditor("md", title, text);
};

const fillAsciidocDefaultText = () => {
  let title = "default markdown title here";
  let text = `== First Headline
Text just for the hight stuff
some day build some to fill it at startup
and move it out of the markup way

=== Secound Headline

some more lines
more here

out now !

image::https://asciidoctor.org/images/octocat.jpg[GitHub mascot]

        `;

  fillEditor("adoc", title, text);
};