from http.server import BaseHTTPRequestHandler, HTTPServer
from threading import Thread


class HTTPHandler(BaseHTTPRequestHandler):

    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/plain; chart=utf-8')
        self.end_headers()
        with open('resources/pac.js', 'rb') as fp:
            pac = fp.read()
        self.wfile.write(pac)


def start_server(host: str, port: int) -> HTTPServer:
    s = HTTPServer((host, port), HTTPHandler)
    Thread(target=s.serve_forever, daemon=True).start()
    return s
