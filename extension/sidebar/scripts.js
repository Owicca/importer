let windowId;

browser.windows.getCurrent({populate: true}).then((windowInfo) => {
  windowId = windowInfo.id;
});

function run(type, url) {
  data = {
    url: url,
  };
  options = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  };
  url = `https://importer.local/${type}/`;
  fetch(url, options)
    .then(t => t.json())
    .then(t => {
      msg = t.message, ' ', t.error;
      alert(msg);
    });
}

["r", "v", "ai"].map((v) => {
  el = document.querySelector(`#${v}`);
  el.addEventListener("click", (e) => {
    browser.tabs.query({windowId: windowId, active: true})
      .then((tabs) => {
        run(v, tabs[0].url);
      })
  })
})
