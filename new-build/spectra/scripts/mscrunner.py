from flask import Flask, request
import json
import os
import subprocess as sbp
import logging

app = Flask(__name__)
logging.basicConfig(level=logging.DEBUG)
os.environ["WERKZEUG_RUN_MAIN"] = "true"
port = os.getenv("MSC_PORT", "8088")


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
    return (json.dumps(cmd), 200)


def main():
    app.run(host="0.0.0.0", port=port, debug=False)
    print("Good Bye.")


if __name__ == "__main__":
    main()
