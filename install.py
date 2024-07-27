import os
from shutil import which
from secrets import token_bytes
from base64 import b64encode

ROOT_DIR = os.path.dirname(__file__)
KITSUNEC2_DIR = "kitsuneC2"
KITSUNEC2_SERVER_DIR = "kitsuneC2/server"
GO_SERVER_ENV = "kitsuneC2/server/.template-env"
KITSUNEC2_FRONTEND_ENV = "kitsune-frontend/.template-env"

def main():
    if os.name != "nt" and os.name != "posix":
        print("Unsupported operating system, quitting!")
        exit(1)

    i = getYNInput("Install web-interface? Answering no will install CLI instead (y/n)")
    if i=="y":
        if not webDependenciesInstalled():
            print("Docker is not installed, or not in $PATH. Please install docker and add it to $PATH")
            exit(1)
        
        if os.name == "posix":
            linuxWebInstall()
        else:
            windowsWebInstall()
    else:
        if not cliDependenciesInstalled():
            print("Go is not installed, or not in $PATH. Please install go and add it to $PATH")
            exit(1)

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
        envVars = envVars.format(webEnabled="false")
        file.close()
    with open(os.path.join(ROOT_DIR, GO_SERVER_ENV.replace(".template-env", ".env")), 'w') as file:
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
    username = getInput("Username for web-login: ")
    password = getInput("Password for web-login: ")
    listenerPort = getInput("What port should web-interface listen on?")
    if not listenerPort.isdigit() or int(listenerPort) <= 0 or int(listenerPort) >= 65535:
        print("Port number should be valid, quitting!")
        exit(1)
        

    redirectUrl = getInput("Domain (or ip address) that server will be running on (e.g. \"kitsunec2.com\" or \"123.456.789.123\"): ")
    redirectUrl = "http://" + redirectUrl + ":{}/".format(listenerPort)

    print("Generating authentication secret...\n")
    nextAuthSecret = b64encode(token_bytes(32)).decode()

    print("Writing frontend .env file...\n")
    envVars = ""
    with open(os.path.join(ROOT_DIR, KITSUNEC2_FRONTEND_ENV), 'r') as file:
        envVars = file.read()
        envVars = envVars.format(_nextAuthUrl=redirectUrl, _nextAuthSecret=nextAuthSecret, _username=username, _password=password)
        file.close()
    with open(os.path.join(ROOT_DIR, KITSUNEC2_FRONTEND_ENV.replace(".template-env", ".env")), 'w') as file:
        file.write(envVars)
        file.close()
    
    print("Writing backend .env file...\n")
    with open(os.path.join(ROOT_DIR, GO_SERVER_ENV), 'r') as file:
        envVars = file.read()
        envVars = envVars.format(webEnabled="true")
        file.close()
    with open(os.path.join(ROOT_DIR, GO_SERVER_ENV.replace(".template-env", ".env")), 'w') as file:
        file.write(envVars)
        file.close()


    print("Environment setup done. Navigate to {} and run command \"sudo docker compose up\"".format(ROOT_DIR))
    

    

def cliDependenciesInstalled():
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










