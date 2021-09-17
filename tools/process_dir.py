import os
import sys

print("当前路径:",sys.path[0])
print("目标目录:",sys.argv[1])

print("===========================")

dirpath=sys.argv[1]

def process_file(root, filename):
    cmd = "D:\\软件\\ffmpeg-gpl-4.4\\bin\\ffmpeg.exe -y -i \"{}\" -filter:v \"crop=iw*1550/1920:ih\" \"{}\"".format(
        os.path.join(root,filename), "new-"+filename
    )
    ret = os.popen(cmd)
    print(filename, ret.readlines())


for root,dirs,names in os.walk(dirpath):
    for filename in names:
        process_file(root, filename)


def scandir(path):
    for name in os.scandir(path):
        if name.is_dir():
            print("dirpath:",name.path)
            scandir(name.path)
        else:
            print("filename:",name.name) 
#函数调用
#scandir(dirpath)
