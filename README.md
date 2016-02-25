# Sausage

*A silly HTTP server*

__Warning:__ This project is under development and not ready for public use.

Sausage is a simple open sauce HTTP/S test server which responds to any HTTP
request with some valid, configurable data (usually containing the keyword
'sausages').

It was built for the purpose of testing the throughput of HTTP proxies, without
hammering real end point web servers. A HTTP client or stress tester can query
a proxy which will forward all requests to a fast and safe silly sausage server.

Why __sausage__?

When I was a kid we played a game. Everyone sits in a circle with one volunteer
standing in the middle. The kids in the circle take turns asking the volunteer
any questions they deem appropriate. No matter what, the kid in the middle must
always answer "sausages" and try not to laugh. Whoever can ask a questions that
makes the kid in the middle laugh wins a sausage... or something.

No matter what request you throw at sausage server, it will always answer
"sausages" (though it's pretty difficult to make it laugh). Hopefully this means
it can outperform your proxy server to highlight where its capacity limits are.


## License

Copyright (c) 2016 Ryan Armstrong

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
