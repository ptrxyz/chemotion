#!/bin/bash
function curlProbe() {
    # since we deal with return codes, we remember: 0=true, 1=false...
    curl -f --connect-timeout 2 "${ELNHOST}":4000 &>/dev/null && probe=0 || probe=1
    return $probe
}

function pidProbe() {
    [[ -e "${PIDFILE}" ]] && probe=0 || probe=1
    return $probe
}

cntr=0
ELNHOST=${CONFIG_ELN_HOST-eln}
until pidProbe || curlProbe; do
    if [[ $((cntr % 30)) -eq 0 ]]; then
        echo "Waiting for ELN [expecting PIDFILE ${PIDFILE} or ELN responding on ${ELNHOST}:4000]...";
    fi
    cntr=$((cntr + 1))
    sleep 1
done
echo "Found ELN!"
