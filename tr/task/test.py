#! /bin/python
import sys
import os

file = os.environ["MEDIA"]

print("/home/"+file+";"+sys.argv[1])
