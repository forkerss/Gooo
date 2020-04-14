#! /usr/bin/python3
import os, sys

def clean(rootpath, cleardirs=["__pycache__"], clearfiles=[".DS_Store"]):
    def remove(root, name):
        cur_path = os.path.join(root, name)
        print("* remove %s" % cur_path)
        os.system("rm -rf %s" % cur_path)
    
    for root, dirs, files in os.walk(rootpath):
        for name in files:
            if name in clearfiles:
                remove(root, name)
        for name in dirs:
            if name in cleardirs:
                remove(root, name)

if __name__ == "__main__":
    print("start clean %s" %sys.argv[1])
    clean(sys.argv[1])