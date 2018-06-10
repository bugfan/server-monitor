# 此代码用来模拟数据源，不停的写入到log.txt
import os
import time
while True:
    time.sleep(3)
    os.system("echo zxy >> ./log.txt")
