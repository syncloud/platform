let template = `
		<% if (apps.length == 0) { %>
		<h2 class="bh2">You don't have any installed apps yet. You can install one from App Center</h2>
		<a href="appcenter.html" class="appcenterh">App Center</a>
    <% } %>

    <%for (s=0; s < apps.length; s++) {
        var app = apps[s]; %>

				<% if (app.id != "store" && app.id != "settings") { %>
				<a href="app.html?app_id=<%= app.id %>" class="colapp app">
					<img src="<%= app.icon %>" class="appminimg">
					<div class="appname"><%= app.name %></div>
					<div class="appdesc"></div>
				</a>
				<% } %>

    <% } %>`;

module.exports = template;
