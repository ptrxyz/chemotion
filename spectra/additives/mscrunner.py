from flask import Flask, request
import json
import os
import subprocess as sbp
import logging
from gevent.pywsgi import WSGIServer

app = Flask(__name__)
logging.basicConfig(level=logging.DEBUG)

# Hack to make Flask's initial banner go away.
# This only works with werkzeug<2.1.0
# os.environ["WERKZEUG_RUN_MAIN"] = "true"


@app.route("/ping")
def ping():
    return ("pong\n", 200, {"Content-Type": "text/plain"})


@app.route("/run", methods=["POST"])
def run():
    cmd = request.json
    app.logger.info("Received: %s" % cmd)
    app.logger.info("%s" % type(cmd))
    try:
        assert type(cmd) is list
        assert cmd[0] == "wine"
        app.logger.info("Running: %s" % " ".join(cmd))
        sbp.run(cmd, timeout=10, check=True)
    except Exception as ex:
        return (str(ex), 400)
    return (json.dumps(cmd), 200, {"Content-Type": "application/json"})


def main():
    port = int(os.getenv("MSC_PORT", "8088"))
    logging.info("Starting server on port %d" % port)
    http_server = WSGIServer(('0.0.0.0', int(port)), app)
    http_server.serve_forever()
    print("Good Bye.")


if __name__ == "__main__":
    main()
