<h1>Simple telnet client</h1>

Try to connect to remote address (host:port) with connection timeout (default is 10s)<br>
All stdin input will be sent to remote address<br>
All remote machine output will reached stdout<br>
<strong>Ctrl+D</strong> on terminal (input EOF) or <strong>Ctrl+C</strong> (int/term signal) will stops telnet client and close connection with remote<br>
If remote close connection by itself telnet will also stoped and exit<br><br>


<h2>Usage</strong></h2>
telnet -address=host:port [-timeout=10s]<br><br>

<h3>Run tests on source code</h3>
go test -v -race ./...<br>
