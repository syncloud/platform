const AppsTemplate = `
<%for (s=0; s < apps.length; s++) {
        var app = apps[s]; %>
				<a href="app.html?app_id=<%= app.id %>" class="colapp app">
					<img src="<%= app.icon %>" class="appminimg">
					<div class="appname withline"><%= app.name %></div>
					<div class="appdesc"></div>
				</a>
<% } %>`

module.exports = {
  AppsTemplate
};
