import os
from shutil import which

ROOT_DIR = os.path.dirname(__file__)
KITSUNEC2_DIR = "kitsuneC2"
KITSUNEC2_SERVER_DIR = "kitsuneC2/server"
GO_SERVER_ENV = "kitsuneC2/server/.example-env"


def main():
    if os.name != "nt" and os.name != "posix":
        print("Unsupported operating system, quitting!")
        exit(1)

    if not coreDependenciesInstalled():
        print("go is not installed, or not in $PATH. Please install go and add it to $PATH")
        exit(1)

    i = getYNInput("Install with web-interface? (y/n)")
    if i=="y":
        if not webDependenciesInstalled():
            print("Docker is not installed, or not in $PATH. Please install docker and add it to $PATH")
            exit(1)
        
        if os.name == "posix":
            linuxWebInstall()
        else:
            windowsWebInstall()
    else:
        if os.name == "posix":
            linuxCliInstall()
        else:
            windowsCliInstall()

    

def windowsWebInstall():
    pass

def windowsCliInstall():
    pass

def linuxCliInstall():
    os.chdir(os.path.join(ROOT_DIR, KITSUNEC2_DIR))
    
    print("Attemping to download go.mod dependencies...\n")
    if os.system("go mod download") != 0:
        print("Error while downloading go.mod dependencies, quiting!")
        exit(1)

    print("Writing .env file...\n")
    envVars = ""
    with open(os.path.join(ROOT_DIR, GO_SERVER_ENV), 'r') as file:
        envVars = file.read()
        envVars = envVars.replace('''ENABLE_WEB_API = "true"''', '''ENABLE_WEB_API = "false"''')
        file.close()
    with open(os.path.join(ROOT_DIR, GO_SERVER_ENV.replace(".example-env", ".env")), 'w') as file:
        file.write(envVars)
        file.close()

    print("Attempting to build server...\n")
    os.chdir(os.path.join(ROOT_DIR, KITSUNEC2_SERVER_DIR))
    if os.system("go build .") != 0:
        print("Error while building server, quitting!")
        exit(1)

    print("Successfully built server at {}".format(os.path.join(ROOT_DIR, KITSUNEC2_SERVER_DIR, "server")))
    print("Navigate to {} and run ./server. Don't move the server binary as it relies on relative paths.".format(os.path.join(ROOT_DIR, KITSUNEC2_SERVER_DIR)))

    


def linuxWebInstall():
    pass

    

    

def coreDependenciesInstalled():
    return which("go") is not None 

def webDependenciesInstalled():
    return which("docker") is not None

def getInput(prompt):
    i = input(prompt + "\n")
    return i.lower()

def getYNInput(prompt):
    i = ""
    while i != "n" and i != "N" and i != "y" and i != "Y":
        i = input(prompt + "\n")

    return i.lower()

if __name__ == "__main__":
    main()










