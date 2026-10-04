const sendEchoTest = () => {
  console.log("send test to ")
  superagent
  .post('/api/echomd/')
  .set('Content-Type', 'application/json')
  .send({"type":"md","text":"# World\n\n## Hello\n\n### dexta\n\n"})
  .responseType('text')
  .then(res => {
    console.log(res.text);
  }).catch(err => {
    console.log(err.message);
    console.dir(err.response);
  });
};

const renderPreview = () => {
  let values = getValues();
  let taUrl = "";
  let taCSS = "";
  if (state.type==="md") {
    taUrl = "markdown";
    taCSS = "retroMarkdownStyle.css";
  } else  {
    taUrl = "asciidoc";
    taCSS = "defaultAsciidoctor.css";
  }
  superagent
  .post(`/api/echo/${taUrl}`)
  .set('Content-Type', 'application/json')
  .send({"type":state.type,"md5":"","title":values.title,"text":values.text})
  .responseType('text')
  .then(res => {
    // console.log(res.text);    
    let htBody = `<!DOCTYPE html>\n<html><head><link rel="stylesheet" href="css/${taCSS}"></head>`;
    htBody += res.text;
    htBody += `</body>\n</html>`;
    renderTemplate('previewRenderer', {text: htBody}, 'previewContainerTarget');
    document.querySelector("#previewTitleTarget").innerText = values.title;
  }).catch(err => {
    console.log(err.message);
    console.dir(err.response);
  });
};

const saveUpdateDocument = () => {
  console.log("click on preview");
  let values = getValues();
  let url = "/api/";
  let sendObj = {type:state.type, title:values.title, text:values.text};
  if(state.mode==="update" && (state.md5||false)) {
    url += "update/doc";
    sendObj.md5 = state.md5;
    state.mode = "update";
    saveUpdateBtnToggler();
  } else {
    url += "add/doc";
  }
  superagent
  .post(url)
  .set('Content-Type', 'application/json')
  .send(sendObj)
  .then(res => {
    state.md5 = res.body.md5;
    handleViewURL();
    // console.dir(res.body);
  }).catch(err => {
    console.log(err.message);
    console.dir(err.response);
  });
};

const getDocumentIndex = () => {
  superagent
  .get('/api/doc/index')
  .then(res => {
    // console.dir(res.body)
    state.docIndex = res.body;
    state.mode = "update";
    saveUpdateBtnToggler();
    renderDocumentIndexList(res.body);
  })
  .catch(err => {
    console.log(err.message);
    console.dir(err.response);
  });
};