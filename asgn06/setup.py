#!/bin/python3
from os import path, system
import platform
import subprocess

def setup_script():
    user = input("What is your ssh username?: ")
    hostname = platform.node()
    print("Hostname found:" + hostname)
    ssh_path = path.expanduser("~/.ssh/")
    if not path.exists(ssh_path + "id_rsa.pub"):
        print("ERROR: public ssh key at" + ssh_path + "does not exist, please generate one")
        print("Run:ssh-keygen")
        exit(1)

    print("Setting up connections with Lab 127 Nodes")
    alive_nodes = list()
    # Get health of Lab 127 and find all active nodes
    for i in range(1, 39):
        alt = ""
        if i < 10:
            alt = "127x0" + str(i) 
            cmd = alt + ".csc.calpoly.edu"
        else:
            alt = "127x" + str(i) 
            cmd = alt + ".csc.calpoly.edu"
        val = 0
        try:
            val = subprocess.check_output(['ping', '-c 1', cmd])
        except subprocess.CalledProcessError as e:
            continue
        if val is not 0:
            alive_nodes.append(cmd)
            alive_nodes.append(alt)

    # add nodes to known_hosts file
    system("ssh-keyscan -H " + hostname + " >> ~/.ssh/known_hosts")
    system("ssh-copy-id -i ~/.ssh/id_rsa.pub " + user + "@" + hostname)
    for node in alive_nodes:
        system("ssh-keygen -R " + node)
        system("ssh-keyscan -H " + node + " >> ~/.ssh/known_hosts")
        # setup public ssh key login

    print("Successful setup ssh keys with:")
    print(alive_nodes)


    print("Setup Complete")

setup_script()
