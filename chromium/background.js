chrome.browserAction.onClicked.addListener(function(tab) {
  chrome.tabs.query({active: true, currentWindow: true}, function(tabs) {
    var activeTab = tabs[0];
    var baseUrl = "http://0.0.0.0:8000";
    var title = encodeURIComponent(activeTab.title);
    var url = encodeURIComponent(activeTab.url);
    var redirectUrl = baseUrl+"/bookmarks/add?url="+url+"&title="+title;

    chrome.tabs.update(activeTab.id, {url: redirectUrl});
  });
});
