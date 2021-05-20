FROM ubuntu:latest

ENV TIMEZONE=Europe/Berlin

RUN AREA=$(echo $TIMEZONE | cut -d'/' -f1) && \
    REGION=$(echo $TIMEZONE | cut -d'/' -f2) && \
    rm -f /etc/localtime && \
    ln -s /usr/share/zoneinfo/$AREA/$REGION /etc/localtime

RUN apt-get -y update
# RUN apt-get install -y --no-install-recommends docker.io
RUN apt-get install -y --no-install-recommends curl ca-certificates git tini
RUN apt-get install -y --no-install-recommends gcc g++ libxrender1 libxext-dev pkg-config

# COPY Anaconda3-2020.11-Linux-x86_64.sh /root/conda.sh
# COPY Anaconda2-2019.10-Linux-x86_64.sh /root/conda.sh
COPY Miniconda3-py39_4.9.2-Linux-x86_64.sh /root/conda.sh
RUN bash /root/conda.sh -p /anaconda3 -b
RUN echo "PATH=/anaconda3/condabin/:$PATH" >> ~/.profile
SHELL ["/bin/bash", "--login", "-c"]
RUN rm /root/conda.sh

WORKDIR /root
RUN /anaconda3/condabin/conda update -n base -c defaults conda

COPY environment.yml /root/environment.yml
RUN conda env create -f /root/environment.yml

RUN mkdir -p /app
WORKDIR /app

COPY ./src/requirements.txt /app/requirements.txt

# Make RUN commands use the new environment:
SHELL ["/anaconda3/condabin/conda", "run", "--no-capture-output", "-n", "chem-spectra", "/bin/bash", "-c"]
# SHELL ["/anaconda3/condabin/conda", "run", "-n", "chem-spectra", "/bin/bash", "-c"]

RUN cd /app && pip install -r requirements.txt

COPY ./scripts/fake-docker.py /bin/docker
RUN chmod +x /bin/docker

COPY ./src /app
RUN mkdir -p /app/instance/
RUN ln -s /shared /app/chem_spectra/tmp

EXPOSE 3007

ENV FLASK_ENV=development
ENTRYPOINT ["/bin/tini", "--", "/anaconda3/condabin/conda", "run", "--no-capture-output", "-n", "chem-spectra"]
CMD gunicorn -w 4 -b 0.0.0.0:3007 server:app
# CMD python -u server.py
