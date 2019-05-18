const AppTemplate = `
<div class="appblock">
			<div class="col30">

				<div style="display: inline-block;">
					<img src="<%= info.app.icon %>" class="appimg" alt="">
				</div>
				<div class="appinfo">
				<h1><%= info.app.name %></h1>
        	<% if (info.installed_version !== null) { %>
					<b>Installed version:</b> <%= info.installed_version %><br>
        	<% } %>
        	<% if (info.installed_version !== info.current_version) { %>
					<b>Available version:</b> <%= info.current_version %><br>
        	<% } %>
					<!--<b>Size:</b> 17.0 MB-->
				</div>
			</div>
			<div class="col70">
				<div class="buttonblock">
        	<% if (info.installed_version !== null) { %>
					<button id="btn_open" data-url="<%= info.app.url %>" class="buttonblue bwidth smbutton">Open</button>
        	<% } %>
        	<% if (info.installed_version === null) { %>
					<button id="btn_install" class="buttonblue bwidth smbutton" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Installing...">Install</button>
        	<% } %>
        	<% if (info.installed_version !== null && info.installed_version !== info.current_version) { %>
					<button id="btn_upgrade" class="buttongreen bwidth smbutton" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Upgrading...">Upgrade</button>
        	<% } %>
        	<% if (info.installed_version !== null) { %>
					<button id="btn_remove"  class="buttongrey bwidth smbutton" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Removing...">Remove</button>
        	<% } %>
        <% if (info.installed_version !== null) { %>
					<button id="btn_backup"  class="buttonblue bwidth smbutton" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Creating backup...">Backup</button>
        	<% } %>
				</div>
				<div class="btext"><%= info.app.description %></div>
			</div>
		</div>
`

module.exports = {
	AppTemplate
};
