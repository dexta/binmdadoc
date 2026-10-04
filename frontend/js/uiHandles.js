const getValues = () => {
  let titleSelector = "";
  let textSelector = "";
  if(state.type==="md") {
    titleSelector = "#markdownTitle";
    textSelector = "#markdownEditor > textarea";
  } else if(state.type==="adoc") {
    titleSelector = "#asciidocTitle";
    textSelector = "#asciidocEditor > textarea";
  } else { return null;}
  let title = document.querySelector(titleSelector).value;
  let text = document.querySelector(textSelector).value;
  return {title, text};
};

const renderDocumentIndexList = (documentIndex) => {
  renderTemplate('documentIndexList', {list: documentIndex}, 'documentLoadList');
};

const switchEditType = (toSwitchTo) => {
  state.type = toSwitchTo;
};

const resizeEditors = () => {
  document.querySelector("#asciidocEditor > pre").style.maxHeight = (window.innerHeight-364)+"px";
  document.querySelector("#markdownEditor > pre").style.maxHeight = (window.innerHeight-364)+"px";
};

const toggleSpellcheck = (codeBedderElm) => {
  let tElment = document.querySelector(`#${codeBedderElm} > textarea`);
  if(tElment.spellcheck) {
    tElment.spellcheck = false;
  } else {
    tElment.spellcheck = true;
  }
};

const callTheDefaults = () => {
  fillMarkdownDefaultText();
  fillAsciidocDefaultText();
  resizeEditors();
  window.addEventListener("resize", resizeEditors);
};

const handleViewURL = () => {
  let nUrlValue = "";
  if(state.md5||false) {
    let origin = window.location.origin;
    nUrlValue = `${origin}/view/${state.md5}`;
  }
  document.querySelector("#viewUrlTarget").value = nUrlValue;
};

const themeSelector = (themeName) => {
  let cssList = {"dark":"prism","default":"prism_default","tomorrowNight":"prism_torwNight","solarLight":"prism_solarLight","twilight":"prism_twilight","okaidia":"prism_okaidia","coy":"coy"};
  let themeFi = `/css/${cssList[themeName]}.css`;
  document.querySelector("#themeSwitcher").href = themeFi;
  state.theme = themeName;
  document.querySelector("#themeBtn").innerText = themeName;
  // console.log(themeName, cssList[themeName]);
};

const copyUrlToClipboard = () => {
  let urltocopy = document.querySelector("#viewUrlTarget").value;
  navigator.clipboard.writeText(urltocopy).then(function() {
    // console.log('Async: Copying to clipboard was successful!');
    // some ui feedback todo - 
  }, function(err) {
    console.error('Async: Could not copy text: ', err);
  });
};

const loadDocumentById = (toLoadId) => {
  let idFromState = state.docIndex.filter( e => {return e.md5 === toLoadId} );
  if( !(idFromState[0].title||false) ) {
    console.log("load id not match local index");
    return;
  }

  superagent
  .get(`/api/doc/by/id/${toLoadId}`)
  .responseType('text')
  .then(res => {
    // console.dir(res.text);
    fillEditor(idFromState[0].type, idFromState[0].title, res.text);
    switchEditType(idFromState[0].type);
    let toSwitchToTabName = (idFromState[0].type==="md")? "markdown" : "asciidoc";
    document.querySelector(`#${toSwitchToTabName}Edit-tab`).click();
    state.mode = "update";
    state.md5 = toLoadId;
    saveUpdateBtnToggler();
    handleViewURL();
  })
  .catch(err => {
    console.log(err.message);
    console.dir(err.response);
  });
};

const newDocument = () => {
  state.mode = "new";
  state.md5 = "";
  saveUpdateBtnToggler();
  fillEditor(state.type, "noch kein Titel", "\n\n\n\n\n\n\n\n");
  handleViewURL();
};

const saveUpdateBtnToggler = () => {
  let usButton = document.querySelector("#saveUpdateBtn");
  if(state.mode==="new") {
    usButton.text = "Save";
  } else if(state.mode==="update") {
    usButton.text = "Update";
  }
};