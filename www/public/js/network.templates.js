const NetworkTemplate = `
    <%
		if (interfaces.length == 0) {
    %>

		<span class="span">No networks found</span>

    <% } %>

    <%
    for (i=0; i < interfaces.length; i++) {
        var interface = interfaces[i];
    %>

        <div class="setline">
				<span class="span">IPv4 (<%= interface.name %>): </span>

            <%
                if (interface.ipv4) {
            %>

                <%
                    addresses = interface.ipv4;
                    for (j=0; j < addresses.length; j++) {
                        var address = addresses[j];
                %>

                    <%
                        if (j > 0) {
                    %>

                        <span class="span">, </span>

                    <% } %>


                    <span class="span"><%=address.addr %></span>

                <% } %>

            <% } else {  %>

                <span class="span">No IPv4</span>

            <% } %>

        </div>
        <div class="setline">
				<span class="span">IPv6 (<%= interface.name %>): </span>

            <%
                if (interface.ipv6) {
            %>

                <%
                    addresses = interface.ipv6;
                    for (j=0; j < addresses.length; j++) {
                        var address = addresses[j];
                %>

                    <%
                        if (j > 0) {
                    %>

                        <span class="span">, </span>

                    <% } %>


                    <span class="span"><%=address.addr %></span>

                <% } %>

            <% } else {  %>

                <span class="span">No IPv6</span>

            <% } %>

        </div>

    <% } %>

`

modules.export = NetworkTemplate