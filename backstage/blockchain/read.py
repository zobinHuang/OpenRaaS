import os

command = 'node read.js {0}'.format(123)
with os.popen(command) as nodejs:
    result = nodejs.read().replace('\n','')
print(result)
