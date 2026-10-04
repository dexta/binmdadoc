const preLoadJSON = (pathObjs, callback) => {
	let obj = null;
	if(pathObjs.length>0) {
		// [{path:"1st path",target:"key of db"},{path:"1st path",target:"key of db"}]
		obj = pathObjs.pop();
	} else {
		return callback.call();
	}
	if(!(obj||false)) return 0;
	// console.dir(obj);
	// console.dir(pathObjs);
	superagent
	 .get(obj.path)
	 .then(res => {
 		db[obj.target] = res.body;
 		preLoadJSON(pathObjs, callback);
	 })
	 .catch(err => {
	 	console.log(err.message);
      console.dir(err.response);
	 });
};

const loadDataToTemplate = (path, tempId, target) => {
	superagent
   .get(path)
   .then(res => {
      // res.body, res.headers, res.status
      // console.dir(res.body);

      mechs.master = res.body;
      renderTemplate(tempId, mechs.master, target);
   })
   .catch(err => {
      // err.message, err.response
      console.log(err.message);
      console.dir(err.response);
   });
}

const loadTemplate = (tempId, callback) => {
	superagent
	 .get('tpl/'+tempId+'.handlebars')
	 .then(res => {
	 		// console.log(res.body);
	 		templateCache[tempId] = res.text;
	 		callback(res.text);
	 })
	 .catch(err => {
	 		cacheErrCount++
	 		if(cacheErrCount <= 10) {
		 		console.log(err.message);
	      console.dir(err.response);
	    }
	 });
};

const renderTemplate = (tempId, dataSet, target) => {
	if(templateCache[tempId]||false) {
		const template = templateCache[tempId];
		const comTempl = Handlebars.compile(template);

		const html = comTempl(dataSet);
		document.getElementById(target).innerHTML = html;
	} else {
		loadTemplate(tempId, () => { renderTemplate(tempId, dataSet, target); });
	}
};

const saveState = () => {
	let strState = JSON.stringify(state);
	localStorage.setItem('userState', strState);
};

const loadState = () => {
	let loState = localStorage.getItem('userState');
	if(loState||false) {
		let tmpState = Object.assign(state, JSON.parse(loState));
		state = tmpState;
	} else {
		console.log('error while loading the state from');
		saveState();
	}
};