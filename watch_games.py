import subprocess
import os
import sys

d = sys.argv[1]
print d

for root, dirs, files in os.walk('./' + d):
	for f in files:
		print f
		if not f.startswith('.'):
			subprocess.check_output("electron /Users/bovardtiberi/Code/me/chlorine -o ./{}/{}".format(d,f), shell=True)
