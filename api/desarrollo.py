from subprocess import call
call(["cmd", "/c", "cls&&go", "build", "-tags", "desarrollo full", "&&api"])
#call(["cmd", "/c", "dir&&hostname"])