import os

command = 'node write.js {0} {1}'.format(123,"456")
with os.popen(command) as nodejs:
    result = nodejs.read().replace('\n','')
print(result)
