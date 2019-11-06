<h1>Simple echo server</h1>

Run server on address (host:port) with timeout (default is 0ms means without timeout)<br>
If server receive any message it will be cloned and send back (echo)<br>
Each timeout tick server also send current time<br>
Purpose on these simple server to test simple [telnet client](https://github.com/MitrickX/otus-golang-2019/tree/master/19/telnet)

<h2>Usage</h2>
echo -address=host:port [-timeout=0ms]<br>

